---
applyTo: "**"
---

# Product Name

- THEB (ذيب)

## Elevator Pitch

- Local ride-hailing app for Mafraq, Jordan—culturally rooted, simple, and fast. Connecting passengers with nearby captains for reliable transportation.

## Problem Statement

- Mafraq residents lack reliable, local, culturally aligned transportation options.
- Existing ride-hailing services don't cater to local Bedouin culture and regional needs.
- Captains need flexible income opportunities with fair earnings.

## Target Users / Personas

- **Passengers**: Mafraq residents, university students, families needing safe rides, individuals without cars.
- **Captains**: Local drivers with private cars seeking flexible income, drivers wanting to serve their community.
- **Admin**: Operations team managing captains, monitoring rides, handling disputes, and cash-out management.

## Core Value Proposition

- Simple, fast ride requests with real-time captain tracking.
- Culturally aligned branding (Wolf/THEB identity, Bedouin heritage).
- Fair earnings for local captains with transparent pricing.
- Lightweight system optimized for local infrastructure.
- Phone-based authentication (no complex registration).

## Primary Use Cases / Journeys

### Passenger Journey
- Open app → GPS auto-detects location.
- Select pickup and drop-off on map.
- View fare estimation.
- Request ride → Backend finds nearest captain.
- Track captain arriving in real-time.
- Complete ride → Pay and rate captain.
- View ride history.

### Captain Journey
- Login and go online.
- Receive nearby ride requests.
- Accept/decline rides.
- Navigate to passenger pickup location.
- Pick up passenger → Navigate to destination.
- Complete ride → View earnings.
- Track daily/weekly earnings summary.

### Admin Journey
- Monitor live captain locations on map.
- Approve/manage captain registrations.
- View ride analytics and captain performance.
- Handle cancellations and disputes.
- Process captain cash-outs.

## Core Features (Scope Backbone)

### Passenger Features
- Phone number + OTP authentication.
- Live GPS location detection.
- Pickup & drop-off selection (map/search).
- Fare estimation before booking.
- Nearest captain matching.
- Real-time captain tracking.
- Ride status updates (requested → accepted → arriving → in-progress → completed).
- In-app ride history.
- Captain rating system.

### Captain Features
- Phone authentication.
- Online/Offline toggle.
- Live location sharing (WebSocket streaming).
- Receive ride requests with passenger details.
- Accept/decline rides.
- Turn-by-turn navigation to pickup and destination.
- Earnings summary (daily/weekly).

### Admin Dashboard
- Captain management (approval, suspension).
- Live map of active captains.
- Ride monitoring and analytics.
- Cancellation handling.
- Cash-out management for captains.
- Performance reports (rides, revenue, captain metrics).

## Backend API Architecture (Implementation Focus)

### Microservices
- **Auth Service**: Phone OTP authentication, JWT token generation and validation.
- **User Service**: User profiles, passenger/captain management, vehicle information.
- **Location Service**: Real-time captain location tracking, geospatial queries, captain matching.
- **Order Service**: Ride requests, captain matching, status lifecycle management.
- **Payment Service**: Fare calculation, payment processing, transaction recording.
- **Rating Service**: Post-ride ratings and reviews.
- **Notification Service**: Push notifications via Expo Push API.
- **WebSocket Server**: Real-time bidirectional communication for location and ride updates.

### Core API Endpoints
- `/api/v1/auth/login` - Send OTP to phone number
- `/api/v1/auth/verify` - Verify OTP and issue JWT tokens
- `/api/v1/auth/refresh` - Refresh access token
- `/api/v1/users/profile` - Get/update user profile
- `/api/v1/captains/online` - Toggle captain online status
- `/api/v1/captains/location` - Update captain location (deprecated in favor of WebSocket)
- `/api/v1/orders/request` - Create new ride request
- `/api/v1/orders/{id}/accept` - Captain accepts ride
- `/api/v1/orders/{id}/status` - Update ride status
- `/api/v1/orders/history` - Get ride history
- `/api/v1/payments/calculate` - Calculate fare estimation
- `/api/v1/payments/{id}` - Process payment
- `/api/v1/ratings` - Submit captain rating
- `/ws/location/stream` - WebSocket endpoint for real-time location updates
- `/ws/rides/{id}` - WebSocket endpoint for ride status updates

### Database Structure (PostgreSQL: theb_db, password: 00962)
- **users**: user_id, name, email, phone, date_of_birth, gender, role, password_hash
- **captains**: captain_id, user_id, vehicle_type, vehicle_model, vehicle_year, plate_number, license_verified, is_online
- **orders**: order_id, passenger_id, captain_id, pickup_lat, pickup_lng, dropoff_lat, dropoff_lng, status, fare_estimate, fare_final
- **payments**: payment_id, order_id, passenger_id, captain_id, amount, method, status
- **locations_history**: record_id, user_id, lat, lng, timestamp
- **ratings**: rating_id, order_id, passenger_id, captain_id, rating_value, review, timestamp
- **notifications**: notification_id, user_id, title, body, is_read, created_at

### Real-Time Features
- **Location Streaming**: Captains send location updates via WebSocket when online.
- **Ride Matching**: Backend uses Redis GEORADIUS to find nearest captains within radius.
- **Status Updates**: Ride status changes broadcast to passenger and captain via WebSocket.
- **Connection Management**: WebSocket hub manages active connections and reconnections.
- **Event Broadcasting**: Redis Pub/Sub coordinates events across multiple server instances.

## Key Domain Entities & Relationships

- **User**: Account with phone authentication, role (passenger | captain | both).
- **Passenger**: User profile for ride requesting.
- **Captain**: User profile with vehicle details, license verification, online status.
- **Vehicle**: Belongs to Captain; has type, model, year, plate number.
- **Order/Ride**: Links Passenger, Captain, pickup/dropoff locations; has status lifecycle, fare, payment.
- **Captain Location**: Real-time GPS coordinates (Redis streaming + historical SQL log).
- **Order Event**: Status change history (requested, accepted, arriving, started, completed, cancelled).
- **Payment**: Transaction record linked to Order; fare amount, method, status.
- **Rating**: Post-ride passenger feedback for Captain.
- **Cancellation**: Records who cancelled (passenger/captain/system), reason, timestamp.

## Non-Goals / Out of Scope (MVP)

- Wallet system (future enhancement).
- Promo codes and loyalty points.
- Multi-city expansion (focus on Mafraq first).
- Carpooling/shared rides.
- In-app chat between passenger and captain.
- Multiple vehicle types (luxury, economy) – single tier initially.

## Critical Constraints & Requirements

- **Arabic + English support** with RTL layout for Arabic.
- **Phone OTP authentication** (no email/password complexity).
- **Real-time location tracking** via WebSocket (low latency).
- **Google Maps integration** for display, routing, and place search.
- **Offline resilience**: Queue requests when network unstable.
- **Low resource usage**: Optimized for mid-range Android devices common in Mafraq.
- **Privacy**: Captain location only visible to matched passenger.
- **Fast captain matching**: < 10 seconds to find nearest available captain.

## High-Level Architecture

- **Mobile App**: Single React Native (Expo) app supporting both Passenger and Captain roles with role-based navigation.
- **Backend**: Golang microservices (auth, orders, location, payments, ratings).
- **Real-time**: WebSocket server for location streaming and ride status updates.
- **Database**: PostgreSQL for persistent data (users, orders, payments, ratings).
- **Cache**: Redis for real-time captain locations and session management.
- **Maps**: Google Maps SDK (display), Directions API (routing), Places API (search).
- **Push Notifications**: Expo Push for ride alerts and captain notifications.

## Success Metrics (Directional)

- **First 1000 passengers** registered in Mafraq within 3 months.
- **Average captain arrival time** < 6 minutes.
- **Captain acceptance rate** > 85%.
- **Ride completion rate** > 90% (low cancellation).
- **App crash rate** < 1%.
- **Daily active captains** > 50 within first 2 months.
- **Passenger retention** (30-day) > 40%.

## Glossary

- **THEB (ذيب)**: Wolf in Arabic; brand name representing strength and reliability.
- **Captain**: Driver providing rides (not "driver" to emphasize respect).
- **Passenger**: Person requesting/taking a ride.
- **Order**: A ride request/transaction (used interchangeably with "Ride").
- **Online Status**: Captain availability to accept rides.
- **Fare Estimation**: Calculated price shown before ride confirmation.
- **Real-time Tracking**: Live captain location updates via WebSocket.
- **OTP**: One-Time Password for phone authentication.
- **Mafraq**: Target city in northern Jordan.
- **Redis**: In-memory data store for real-time location caching.
- **WebSocket**: Protocol for bi-directional real-time communication.
- **Bedouin**: Cultural identity central to THEB's brand and target market.
