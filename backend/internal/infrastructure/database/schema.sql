CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    full_name     VARCHAR(100),
    phone         VARCHAR(20) UNIQUE,
    email         VARCHAR(100) UNIQUE,
    password_hash VARCHAR(255),
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cities (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE hotels (
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    city_id      INT REFERENCES cities(id),
    address      TEXT,
    rating_avg   NUMERIC(2,1) DEFAULT 0,
    rating_count INT DEFAULT 0
);

CREATE TABLE rooms (
    id        SERIAL PRIMARY KEY,
    hotel_id  INT REFERENCES hotels(id),
    room_type VARCHAR(50),
    price     NUMERIC(10,2),
    capacity  INT,
    status    VARCHAR(10) NOT NULL DEFAULT 'AVAILABLE' CHECK (status IN ('AVAILABLE','PENDING','BOOKED'))
);

CREATE TABLE buses (
    id             SERIAL PRIMARY KEY,
    source_city_id INT REFERENCES cities(id),
    dest_city_id   INT REFERENCES cities(id),
    travel_date    DATE NOT NULL,
    travel_time    TIME NOT NULL,
    price          NUMERIC(10,2),
    rating_avg     NUMERIC(2,1) DEFAULT 0,
    rating_count   INT DEFAULT 0
);

CREATE TABLE bus_seats (
    id        SERIAL PRIMARY KEY,
    bus_id    INT REFERENCES buses(id),
    seat_no   INT NOT NULL,
    status    VARCHAR(10) NOT NULL DEFAULT 'AVAILABLE' CHECK (status IN ('AVAILABLE','PENDING','BOOKED'))
);

CREATE TABLE reservations (
    id                  SERIAL PRIMARY KEY,
    user_id             INT REFERENCES users(id),
    reservation_type    VARCHAR(10) NOT NULL CHECK (reservation_type IN ('BUS','HOTEL')),
    buyer_full_name     VARCHAR(100) NOT NULL,
    buyer_phone         VARCHAR(20) NOT NULL,
    buyer_national_code VARCHAR(10) NOT NULL,
    buyer_birth_date    DATE NOT NULL,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE hotel_reservations (
    id             SERIAL PRIMARY KEY,
    reservation_id INT REFERENCES reservations(id),
    room_id        INT REFERENCES rooms(id),
    check_in       DATE,
    check_out      DATE
);

CREATE TABLE bus_reservations (
    id             SERIAL PRIMARY KEY,
    reservation_id INT UNIQUE REFERENCES reservations(id),
    bus_id         INT REFERENCES buses(id)
);

CREATE TABLE passengers (
    id             SERIAL PRIMARY KEY,
    reservation_id INT REFERENCES reservations(id),
    seat_id        INT REFERENCES bus_seats(id),
    full_name      VARCHAR(100) NOT NULL,
    national_code  VARCHAR(10) NOT NULL,
    gender         VARCHAR(10) CHECK (gender IN ('MALE','FEMALE'))
);

CREATE TABLE payments (
    id             SERIAL PRIMARY KEY,
    reservation_id INT REFERENCES reservations(id),
    amount         NUMERIC(10,2),
    status         VARCHAR(20),
    paid_at        TIMESTAMP
);

CREATE TABLE reviews (
    id          SERIAL PRIMARY KEY,
    user_id     INT REFERENCES users(id),
    target_type VARCHAR(10) NOT NULL CHECK (target_type IN ('HOTEL','BUS')),
    target_id   INT NOT NULL,
    rating      INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment     TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_rating_after_review_insert()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.target_type = 'HOTEL' THEN
        UPDATE hotels
        SET rating_avg = ((rating_avg * rating_count) + NEW.rating) / (rating_count + 1),
            rating_count = rating_count + 1
        WHERE id = NEW.target_id;
    ELSIF NEW.target_type = 'BUS' THEN
        UPDATE buses
        SET rating_avg = ((rating_avg * rating_count) + NEW.rating) / (rating_count + 1),
            rating_count = rating_count + 1
        WHERE id = NEW.target_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_rating
AFTER INSERT ON reviews
FOR EACH ROW
EXECUTE FUNCTION update_rating_after_review_insert();


CREATE OR REPLACE FUNCTION lock_seat_before_passenger_insert()
RETURNS TRIGGER AS $$
DECLARE
    current_status VARCHAR(10);
BEGIN
    SELECT status INTO current_status FROM bus_seats WHERE id = NEW.seat_id FOR UPDATE;
    IF current_status <> 'AVAILABLE' THEN
        RAISE EXCEPTION 'Seat % is not available (status=%)', NEW.seat_id, current_status;
    END IF;
    UPDATE bus_seats SET status = 'PENDING' WHERE id = NEW.seat_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_lock_seat
BEFORE INSERT ON passengers
FOR EACH ROW
EXECUTE FUNCTION lock_seat_before_passenger_insert();


CREATE OR REPLACE FUNCTION finalize_seat_booking()
RETURNS TRIGGER AS $$
DECLARE
    res_type VARCHAR(10);
BEGIN
    SELECT reservation_type INTO res_type FROM reservations WHERE id = NEW.reservation_id;
    
    IF res_type = 'BUS' THEN
        UPDATE bus_seats
        SET status = 'BOOKED'
        WHERE id IN (
            SELECT seat_id
            FROM passengers
            WHERE reservation_id = NEW.reservation_id
        );
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_finalize_seat
AFTER INSERT ON payments
FOR EACH ROW
WHEN (NEW.status = 'PAID')
EXECUTE FUNCTION finalize_seat_booking();


CREATE OR REPLACE FUNCTION lock_room_before_reservation()
RETURNS TRIGGER AS $$
DECLARE
    room_status VARCHAR(10);
BEGIN
    SELECT status INTO room_status FROM rooms WHERE id = NEW.room_id FOR UPDATE;
    IF room_status <> 'AVAILABLE' THEN
        RAISE EXCEPTION 'Room % is not available', NEW.room_id;
    END IF;
    UPDATE rooms SET status = 'PENDING' WHERE id = NEW.room_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_lock_room
BEFORE INSERT ON hotel_reservations
FOR EACH ROW
EXECUTE FUNCTION lock_room_before_reservation();


CREATE OR REPLACE FUNCTION finalize_room_booking()
RETURNS TRIGGER AS $$
DECLARE
    res_type VARCHAR(10);
BEGIN
    SELECT reservation_type INTO res_type FROM reservations WHERE id = NEW.reservation_id;
    
    IF res_type = 'HOTEL' THEN
        UPDATE rooms
        SET status = 'BOOKED'
        WHERE id IN (
            SELECT room_id
            FROM hotel_reservations
            WHERE reservation_id = NEW.reservation_id
        );
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_finalize_room
AFTER INSERT ON payments
FOR EACH ROW
WHEN (NEW.status = 'PAID')
EXECUTE FUNCTION finalize_room_booking();


CREATE OR REPLACE FUNCTION release_expired_pending()
RETURNS VOID AS $$
BEGIN
    UPDATE bus_seats
    SET status = 'AVAILABLE'
    WHERE status = 'PENDING'
    AND id IN (
        SELECT p.seat_id
        FROM passengers p
        JOIN reservations r ON r.id = p.reservation_id
        WHERE r.created_at < NOW() - INTERVAL '15 minutes'
    );

    UPDATE rooms
    SET status = 'AVAILABLE'
    WHERE status = 'PENDING'
    AND id IN (
        SELECT hr.room_id
        FROM hotel_reservations hr
        JOIN reservations r ON r.id = hr.reservation_id
        WHERE r.created_at < NOW() - INTERVAL '15 minutes'
    );
END;
$$ LANGUAGE plpgsql;
