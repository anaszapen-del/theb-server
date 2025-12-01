# THEB Ride App – Database Structure

This document defines the database schema for THEB's MVP.
The structure is optimized for scalability, simplicity, and real-time operations.

---

## 1. USERS TABLE
Stores all users: passengers, captains, or both.

| Field           | Type                               | Notes                              |
|----------------|-------------------------------------|------------------------------------|
| user_id        | UUID (PK)                           | Unique user identifier             |
| name           | string                              | User name                          |
| email          | string (unique)                     | User email address                 |
| phone          | string (unique)                     | Login identifier                   |
| date_of_birth  | date                                | User's date of birth               |
| gender         | enum(male, female, other)           | User's gender                      |
| role           | enum(passenger, captain, both)      | User's role                        |
| password_hash  | string                              | Secure password storage            |
| created_at     | timestamp                           |                                    |
| updated_at     | timestamp                           |                                    |

---

## 2. CAPTAINS TABLE
Details for users who work as captains.

| Field           | Type       | Notes                                 |
|----------------|------------|---------------------------------------|
| captain_id     | UUID (PK)  | Unique captain identifier             |
| user_id        | UUID (FK)  | References users.user_id              |
| vehicle_type   | string     | Sedan, SUV, Pickup, etc.              |
| vehicle_model  | string     | Car model name                        |
| vehicle_year   | string     | e.g. 2018                             |
| plate_number   | string     | Car license plate                     |
| license_verified | boolean  | Admin verification status             |
| is_online      | boolean    | Captain availability                  |
| current_lat    | float      | Live location                         |
| current_lng    | float      | Live location                         |
| last_updated   | timestamp  | Last location update                  |

---

## 3. RIDES TABLE
Main operational table for ride requests.

| Field         | Type       | Notes                                  |
|--------------|------------|----------------------------------------|
| ride_id      | UUID (PK)  | Ride identifier                        |
| passenger_id | UUID (FK)  | References users.user_id               |
| captain_id   | UUID (FK)  | References captains.captain_id         |
| pickup_lat   | float      |                                        |
| pickup_lng   | float      |                                        |
| dropoff_lat  | float      |                                        |
| dropoff_lng  | float      |                                        |
| status       | enum(requested, matched, on_the_way, in_progress, completed, canceled) | Ride lifecycle |
| fare_estimate| float      | Pre-calculated estimate                |
| fare_final   | float      | Final fare                             |
| created_at   | timestamp  |                                        |

---

## 4. PAYMENTS TABLE

| Field        | Type       | Notes                                  |
|-------------|------------|----------------------------------------|
| payment_id  | UUID (PK)  | Payment identifier                     |
| ride_id     | UUID (FK)  | References rides.ride_id               |
| passenger_id| UUID (FK)  | References users.user_id               |
| captain_id  | UUID (FK)  | References captains.captain_id         |
| amount      | float      | Payment amount                         |
| method      | enum(cash, wallet, card) | Payment type             |
| status      | enum(pending, paid, failed) | Payment status        |
| created_at  | timestamp  |                                        |

---

## 5. LOCATIONS_HISTORY TABLE

| Field      | Type       | Notes                                  |
|-----------|------------|----------------------------------------|
| record_id | UUID (PK)  |                                        |
| user_id   | UUID (FK)  |                                        |
| lat       | float      |                                        |
| lng       | float      |                                        |
| timestamp | timestamp  | Tracking history                       |

---

## 6. RATINGS TABLE

| Field        | Type       | Notes                                  |
|-------------|------------|----------------------------------------|
| rating_id   | UUID (PK)  |                                        |
| ride_id     | UUID (FK)  |                                        |
| passenger_id| UUID (FK)  |                                        |
| captain_id  | UUID (FK)  |                                        |
| rating_value| int (1–5)  |                                        |
| review      | string     | Optional                               |
| timestamp   | timestamp  |                                        |

---

## 7. NOTIFICATIONS TABLE

| Field           | Type       | Notes                                  |
|----------------|------------|----------------------------------------|
| notification_id| UUID (PK)  |                                        |
| user_id        | UUID (FK)  |                                        |
| title          | string     | Notification title                      |
| body           | string     | Notification message                    |
| is_read        | boolean    |                                        |
| created_at     | timestamp  |                                        |

---

## 8. APP_SETTINGS TABLE

| Field      | Type       | Notes                                  |
|-----------|------------|----------------------------------------|
| setting_id| UUID (PK)  |                                        |
| key       | string     | e.g. "base_fare", "per_km_rate"         |
| value     | string     | Stored value                           |

---

# End of Schema
