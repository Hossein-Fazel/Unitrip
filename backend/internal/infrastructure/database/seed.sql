
-- Cities
INSERT INTO cities (name) VALUES ('Tehran'), ('Isfahan'), ('Shiraz'), ('Mashhad'), ('Yazd'), ('Kerman'), ('Tabriz'), ('Hamedan');

-- Users
INSERT INTO users (full_name, phone, email, password_hash) VALUES
('Ali Rezaei', '09123456789', 'ali.rezaei@example.com', 'hash123'),
('Sara Mohammadi', '09129876543', 'sara.mohammadi@example.com', 'hash456'),
('Hasan Karimi', '09121112233', 'hasan.karimi@example.com', 'hash789'),
('Maryam Ghasemi', '09131234567', 'maryam.ghasemi@example.com', 'hash101112'),
('Reza Davoodi', '09147654321', 'reza.davoodi@example.com', 'hash131415'),
('Fatemeh Najafi', '09151112222', 'fatemeh.najafi@example.com', 'hash161718'),
('Babak Zandi', '09188118811', 'babak.zandi@example.com', 'hash192021');

-- Hotels
INSERT INTO hotels (name, city_id, address) VALUES
('Espinas Palace Hotel', 1, 'Saadat Abad, Behrood Sq.'), ('Abbasi Hotel', 2, 'Amadegah St.'),
('Shiraz Grand Hotel', 3, 'Quran Gate'), ('Dad Hotel', 5, '10th Farvardin St.'),
('Pars Hotel', 6, 'Jomhouri Blvd.'), ('Bu-Ali Hotel', 8, 'Bu-Ali Sina Sq.');

-- Rooms
INSERT INTO rooms (hotel_id, room_type, price, capacity, status) VALUES
(1, 'Double', 2500000, 2, 'AVAILABLE'), (1, 'Royal Suite', 5000000, 4, 'AVAILABLE'),
(2, 'Single', 1800000, 1, 'AVAILABLE'), (2, 'Double', 3000000, 2, 'AVAILABLE'),
(3, 'Suite', 4500000, 3, 'AVAILABLE'), (4, 'Single', 2200000, 1, 'AVAILABLE'),
(4, 'Double', 3500000, 2, 'AVAILABLE'), (5, 'Suite', 5500000, 4, 'AVAILABLE'),
(1, 'Single', 1900000, 1, 'AVAILABLE'), (6, 'Double', 2800000, 2, 'AVAILABLE'),
(6, 'Suite', 4000000, 3, 'AVAILABLE');

-- Buses
INSERT INTO buses (source_city_id, dest_city_id, travel_date, travel_time, price) VALUES
(1, 2, '2026-03-01', '22:30:00', 450000), (2, 3, '2026-03-05', '21:00:00', 400000),
(1, 4, '2026-03-10', '20:00:00', 550000), (2, 1, '2026-03-02', '23:00:00', 450000),
(3, 1, '2026-03-08', '21:30:00', 600000), (1, 5, '2026-03-12', '19:00:00', 500000),
(7, 1, '2026-03-15', '20:30:00', 700000), (1, 8, '2026-03-20', '23:30:00', 350000);

-- Bus Seats (25 seats per bus)
INSERT INTO bus_seats (bus_id, seat_no) SELECT 1, generate_series(1, 25);
INSERT INTO bus_seats (bus_id, seat_no) SELECT 2, generate_series(1, 25);
INSERT INTO bus_seats (bus_id, seat_no) SELECT 3, generate_series(1, 25);
INSERT INTO bus_seats (bus_id, seat_no) SELECT 4, generate_series(1, 25);
INSERT INTO bus_seats (bus_id, seat_no) SELECT 5, generate_series(1, 25);
INSERT INTO bus_seats (bus_id, seat_no) SELECT 6, generate_series(1, 25);
INSERT INTO bus_seats (bus_id, seat_no) SELECT 7, generate_series(1, 25);
INSERT INTO bus_seats (bus_id, seat_no) SELECT 8, generate_series(1, 25);

-- Assume reservations start with ID=1 and auto-increment
-- Scenario 1: Successful bus booking for 2 passengers (bus 1, seats 5 & 6)
INSERT INTO reservations (user_id, reservation_type, buyer_full_name, buyer_phone, buyer_national_code, buyer_birth_date) VALUES (1, 'BUS', 'Ali Rezaei', '09123456789', '0011223344', '1990-05-15'); -- res_id=1
INSERT INTO bus_reservations (reservation_id, bus_id) VALUES (1, 1);
INSERT INTO passengers (reservation_id, seat_id, full_name, national_code, gender) VALUES (1, 5, 'Ali Rezaei', '0011223344', 'MALE'), (1, 6, 'Maryam Ahmadi', '0099887766', 'FEMALE');
INSERT INTO payments (reservation_id, amount, status, paid_at) VALUES (1, 900000, 'PAID', CURRENT_TIMESTAMP);

-- Scenario 2: Successful hotel room booking (hotel 2, room 4)
INSERT INTO reservations (user_id, reservation_type, buyer_full_name, buyer_phone, buyer_national_code, buyer_birth_date) VALUES (2, 'HOTEL', 'Sara Mohammadi', '09129876543', '1122334455', '1992-11-20'); -- res_id=2
INSERT INTO hotel_reservations (reservation_id, room_id, check_in, check_out) VALUES (2, 4, '2026-04-10', '2026-04-14');
INSERT INTO payments (reservation_id, amount, status, paid_at) VALUES (2, 12000000, 'PAID', CURRENT_TIMESTAMP);

-- Scenario 3: Unpaid bus reservation (bus 2, seat 10)
INSERT INTO reservations (user_id, reservation_type, buyer_full_name, buyer_phone, buyer_national_code, buyer_birth_date) VALUES (3, 'BUS', 'Hasan Karimi', '09121112233', '2233445566', '1985-01-30'); -- res_id=3
INSERT INTO bus_reservations (reservation_id, bus_id) VALUES (3, 2);
INSERT INTO passengers (reservation_id, seat_id, full_name, national_code, gender) VALUES (3, 35, 'Hasan Karimi', '2233445566', 'MALE'); -- seat_id 35 = bus 2, seat 10

-- Scenario 4: Post reviews
INSERT INTO reviews (user_id, target_type, target_id, rating, comment) VALUES (1, 'BUS', 1, 5, 'Very clean bus.'), (2, 'HOTEL', 2, 4, 'Beautiful hotel.');

-- Scenario 5: Successful bus booking (return trip, bus 4, seat 15)
INSERT INTO reservations (user_id, reservation_type, buyer_full_name, buyer_phone, buyer_national_code, buyer_birth_date) VALUES (4, 'BUS', 'Maryam Ghasemi', '09131234567', '3344556677', '1995-02-25'); -- res_id=4
INSERT INTO bus_reservations (reservation_id, bus_id) VALUES (4, 4);
INSERT INTO passengers (reservation_id, seat_id, full_name, national_code, gender) VALUES (4, 90, 'Maryam Ghasemi', '3344556677', 'FEMALE'); -- seat_id 90 = bus 4, seat 15
INSERT INTO payments (reservation_id, amount, status, paid_at) VALUES (4, 450000, 'PAID', CURRENT_TIMESTAMP);

-- Scenario 6: Successful booking for a new hotel (hotel 4, room 7)
INSERT INTO reservations (user_id, reservation_type, buyer_full_name, buyer_phone, buyer_national_code, buyer_birth_date) VALUES (1, 'HOTEL', 'Ali Rezaei', '09123456789', '0011223344', '1990-05-15'); -- res_id=5
INSERT INTO hotel_reservations (reservation_id, room_id, check_in, check_out) VALUES (5, 7, '2026-04-20', '2026-04-22');
INSERT INTO payments (reservation_id, amount, status, paid_at) VALUES (5, 7000000, 'PAID', CURRENT_TIMESTAMP);

-- Scenario 7: Add more reviews
INSERT INTO reviews (user_id, target_type, target_id, rating, comment) VALUES (4, 'BUS', 1, 3, 'The driver was too fast.'), (5, 'HOTEL', 2, 5, 'Outstanding!');

-- Scenario 8: Unpaid hotel reservation (hotel 5, room 8)
INSERT INTO reservations (user_id, reservation_type, buyer_full_name, buyer_phone, buyer_national_code, buyer_birth_date) VALUES (6, 'HOTEL', 'Fatemeh Najafi', '09151112222', '4455667788', '2000-10-01'); -- res_id=6
INSERT INTO hotel_reservations (reservation_id, room_id, check_in, check_out) VALUES (6, 8, '2026-05-01', '2026-05-05');

-- Scenario 9: New booking for the hotel in Hamedan (hotel 6, room 11)
INSERT INTO reservations (user_id, reservation_type, buyer_full_name, buyer_phone, buyer_national_code, buyer_birth_date) VALUES (7, 'HOTEL', 'Babak Zandi', '09188118811', '5566778899', '1988-08-10'); -- res_id=7
INSERT INTO hotel_reservations (reservation_id, room_id, check_in, check_out) VALUES (7, 11, '2026-05-10', '2026-05-12');
INSERT INTO payments (reservation_id, amount, status, paid_at) VALUES (7, 5600000, 'PAID', CURRENT_TIMESTAMP);
INSERT INTO reviews (user_id, target_type, target_id, rating, comment) VALUES (7, 'HOTEL', 6, 5, 'Great location.');