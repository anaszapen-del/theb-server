# THEB – Ride-Hailing App
## Product Requirements Document (PRD)

---

# 1. Overview
**THEB** is a local ride‑hailing application designed for **Mafraq, Jordan**, targeting a population with strong Bedouin culture and local mobility needs. The app offers a simple, fast, and culturally aligned transportation experience connecting passengers with nearby captains.

The brand name **“THEB / ذيب”** reflects strength, reliability, and Bedouin identity — making it memorable and relevant to the city.

---

# 2. Vision & Goals
### **Vision**
To create a simple, reliable, culturally rooted ride‑hailing solution for Mafraq, with potential expansion to other Jordanian regions.

### **Primary Goals**
- Provide passengers with a clean, fast, and intuitive ride‑request experience.
- Enable captains to receive nearby requests and earn efficiently.
- Ensure safe, accurate, real‑time tracking of captains.
- Build a lightweight system optimized for local infrastructure.

---

# 3. Target Users
### **Passengers**
- People living in Mafraq
- University students
- Families needing safe local rides
- Individuals without cars or wanting convenience

### **Captains**
- Local drivers with private cars
- Drivers wanting flexible income

---

# 4. Core Features
## 4.1 Passenger App
- Phone number authentication
- Live location detection
- Pickup & drop‑off selection (map or search)
- Fare estimation
- Nearest captain matching
- Real‑time captain tracking
- Ride status updates
- In‑app ride history
- Rating the captain

## 4.2 Captain App
- Phone authentication
- Online/Offline toggle
- Live location sharing
- Receive ride requests
- Accept/decline rides
- Navigation to pickup & destination
- Earnings summary

## 4.3 Admin Dashboard
- Manage captains
- View live map of active captains
- Monitor rides
- Handle cancellations
- Cash-out management for captains
- Basic reports (rides, revenue, captain performance)

---

# 5. System Architecture Summary
### **Mobile Apps:** React Native
### **Backend:** Golang
### **Real-time Communication:** WebSockets
### **Location Storage:** Redis (real‑time) + SQL (historical)
### **Main Database:** PostgreSQL
### **Maps:** Google Maps SDK, Directions API, Places API

---

# 6. Ride Flow (End-to-End)
1. Passenger opens the app → GPS detected.
2. Passenger selects pickup & drop‑off.
3. Price estimation from backend.
4. Passenger confirms request.
5. Backend finds nearest available captain.
6. Ride request sent to captain.
7. Captain accepts → Passenger notified.
8. Captain moves toward pickup → Passenger tracks.
9. Ride starts → Captain moves toward destination.
10. Ride ends → Payment + rating.

---

# 7. Design Guidelines
THEB must feel **modern, local, confident, and strong**.

## 7.1 Brand Colors
### **Primary Colors**
- **Wolf Black**: `#0D0D0D` — main backgrounds, power, confidence
- **Desert Gold**: `#D4A048` — Bedouin identity, accents, highlights
- **Pure White**: `#FFFFFF` — text, clarity, high contrast

### **Secondary Colors**
- **Sand Gray**: `#B6B0A2` — soft UI elements
- **Mafraq Blue**: `#3A6EA5` — subtle tech feel when needed
- **Success Green**: `#41A45A`
- **Alert Red**: `#D9534F`

## 7.2 Typography
- **Primary Font:** Cairo / Noto Sans Arabic (Arabic support)
- **English Equivalent:** Inter or SF Pro
- Style should be **bold, clean, and geometric**.

## 7.3 Logo Style
### Concept: **Geometric Wolf Head + Location Pin**
- Angular wolf face made with geometric shapes
- Inside or combined with a map location pin
- Desert Gold as a key highlight color
- Outline version for dark mode

## 7.4 UI Style
- Dark mode default (fits Wolf identity)
- Gold accents for active elements
- Rounded corners (12–16px)
- Simple, no clutter
- Clear CTAs like "Request Ride"

---

# 8. Screens Breakdown

## Passenger App
1. Splash screen (wolf logo + THEB)
2. Onboarding (simple 1–2 slides)
3. Login (phone + OTP)
4. Home map screen
   - Current location
   - Search bar
   - Destination pin
5. Ride estimation modal
6. Captain searching screen
7. Captain assigned screen
   - Captain details
   - Car info
   - ETA
8. Live ride tracking
9. Payment summary
10. Rating screen
11. Profile + Ride History

## Captain App
1. Splash screen
2. Login (phone)
3. Go Online toggle
4. Map screen with location
5. Incoming ride request popup
6. Navigation to pickup
7. In‑trip navigation
8. Ride complete
9. Daily earnings

---

# 9. Database Overview
### Key tables:
- users
- passengers
- captains
- vehicles
- captain_locations (Redis realtime + SQL log)
- orders
- order_events
- payments
- ratings
- cancellations

---

# 10. Future Enhancements
- Wallet system
- Promo codes
- Multi-city expansion
- Loyalty points
- Carpooling
- In-app chat

---

# 11. Success Metrics
- First 1000 passengers in Mafraq
- Average captain arrival time < 6 minutes
- Captain acceptance rate > 85%
- Ride completion rate > 90%
- App crash rate < 1%

---

# 12. Summary
This PRD defines the foundation of **THEB**, a culturally rooted, technically strong, locally optimized ride-hailing solution for Mafraq. With clear user flows, branding guidelines, and system requirements, the product is positioned for successful MVP development and future expansion.


## Single-App Role Architecture (Passenger + Captain)

### Overview
THEB will launch with **one mobile application** that supports both **Passenger** and **Captain** roles. This approach reduces development time, simplifies maintenance, and accelerates MVP release in Mafraq.

### Role Selection Flow
- At login or registration, the user chooses: **Passenger**, **Captain**, or (optional) **Both**.
- The user’s role determines which UI modules and permissions they can access.

### Backend Role Model
```json
{
  "user_id": "string",
  "name": "string",
  "phone": "string",
  "role": "passenger | captain | both",
  "captain_profile": {
    "vehicle_type": "string",
    "vehicle_model": "string",
    "vehicle_year": "string",
    "plate_number": "string",
    "license_verified": true
  }
}
```

### Mode Switching (Inside the App)
- A toggle in the sidebar lets users switch between **Passenger Mode** and **Captain Mode**.
- Each mode has its own dashboard and features.

### Passenger Mode Screens
- Home Map (Pickup selection)
- Choose Destination
- Ride Matching
- Ride Tracking
- Payment
- Trip History
- Account Settings

### Captain Mode Screens
- Go Online / Go Offline
- Requests Queue (New ride requests)
- Navigation to Passenger
- Navigation to Drop-off
- Earnings Dashboard
- Trip History
- Captain Profile

### Permissions Separation
| Feature | Passenger | Captain |
|---------|-----------|---------|
| Request ride | ✅ | ❌ |
| Accept ride | ❌ | ✅ |
| Location tracking | Send location when in ride | Always send when online |
| Payments | Pay | Receive |
| Navigation | Yes | Yes |

### Advantages of This Approach
- One codebase → faster development
- Unified branding and UX
- Lightweight backend
- Easier testing and deployment
- Perfect for local MVP

### Future Scalability
When THEB expands beyond Mafraq, the apps can be easily split:
- THEB Passenger App
- THEB Captain App

All backend logic remains compatible.

