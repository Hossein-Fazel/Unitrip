# How to Run the Project

This guide explains how to run the backend project locally using **Go**.
You can either use an **existing PostgreSQL database** or run PostgreSQL easily using **Docker Compose**.

---

## Prerequisites

Make sure you have the following installed:

* Go **1.22+**
* Git
* PostgreSQL **(optional – only if you don’t use Docker)**
* Docker & Docker Compose **(optional – recommended)**

---

## 1. Clone the Repository

```bash
git clone <https://github.com/Hossein-Fazel/Unitrip.git>
cd <Unitrip/backend/>
```

---

## 2. Environment Variables

The application is configured using environment variables.
You must create a `.env` file in the backend root directory.

Start by copying the example file:

```bash
cp env.example .env
```

Then open `.env` and configure the variables as described below.

---

### Database Configuration

```env
DB_NAME=your_database_name
DB_USER=your_database_user
DB_PORT=5432
DB_HOST=localhost
DB_PASSWORD=your_database_password
```

| Variable      | Description                                                                          |
| ------------- | ------------------------------------------------------------------------------------ |
| `DB_NAME`     | Name of the PostgreSQL database the application connects to                          |
| `DB_USER`     | Database user with access to the specified database                                  |
| `DB_PORT`     | Port on which PostgreSQL is running (default: `5432`)                                |
| `DB_HOST`     | Database host (use `localhost` for local DB or `postgres` when using Docker Compose) |
| `DB_PASSWORD` | Password for the database user                                                       |

📌 **Note:**
When using **Docker Compose**, set `DB_HOST=localhost` if the backend runs on your host machine.

---

### Web Server Configuration

```env
WEB_PORT=8080
SECRET_KEY=your_secret_key
```

| Variable     | Description                                                            |
| ------------ | ---------------------------------------------------------------------- |
| `WEB_PORT`   | Port on which the HTTP server will run                                 |
| `SECRET_KEY` | Secret key used for signing JWTs and other security-related operations |

⚠️ **Security Tip:**
Do not commit the `.env` file to version control. Always keep your `SECRET_KEY` private.

---

### Admin Account Configuration

```env
ADMIN_USERNAME=admin_username
ADMIN_PASSWORD=admin_password
```

| Variable         | Description                            |
| ---------------- | -------------------------------------- |
| `ADMIN_USERNAME` | Username for the initial admin account |
| `ADMIN_PASSWORD` | Password for the initial admin account |

📌 **Behavior:**
On application startup, the system checks whether an admin account exists.
If not, it automatically creates one using these credentials.

---

## 3. Database Setup

You have **two options** for running the database.

---

### Option A: Use Docker Compose (Recommended)

If you **don’t have PostgreSQL installed**, or you want a fast and isolated setup, use Docker.

#### Steps

1. Make sure Docker is running
2. Run the following command:

```bash
docker compose up -d
```

This will:

* Start a PostgreSQL container (`postgres:16-alpine`)
* Automatically create the database
* Apply the database schema (`schema.sql`)
* Insert fake/sample data (`seed.sql`)
* Persist data using a Docker volume

📌 **Notes**

* The database will be available on the port defined in `DB_PORT`
* Database credentials are read from your `.env` file
* Data will persist between restarts

To stop the database:

```bash
docker compose down
```

---

### Option B: Use Your Own PostgreSQL Database

If you **already have PostgreSQL installed**, you can use it directly.

#### Steps

1. Create a database manually:

   ```sql
   CREATE DATABASE your_database_name;
   ```

2. Update your `.env` file with your database credentials:

   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_database_user
   DB_PASSWORD=your_database_password
   DB_NAME=your_database_name
   ```

3. Apply schema and seed data manually (optional but recommended):

   ```bash
   psql -U your_database_user -d your_database_name -f internal/infrastructure/database/schema.sql
   psql -U your_database_user -d your_database_name -f internal/infrastructure/database/seed.sql
   ```

---

## 4. Run the Backend Server

After the database is ready:

```bash
go mod tidy
go run .
```

If everything is configured correctly, you should see logs indicating that:

* Database connection is successful
* HTTP server is running

The API will be available at:

```
http://localhost:<WEB_PORT>
```

Example:

```
http://localhost:8080
```

---

## 5. Admin Account

On startup, the application uses the following environment variables to create an admin account:

```env
ADMIN_USERNAME
ADMIN_PASSWORD
```

Make sure these values are set correctly in your `.env` file.
