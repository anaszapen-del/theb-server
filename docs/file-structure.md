/theb-app
│
├── app/                       # All screens, navigation & UI logic
│   ├── passenger/             # Passenger mode screens
│   │   ├── HomeScreen.tsx
│   │   ├── SelectDestinationScreen.tsx
│   │   ├── SearchingCaptainScreen.tsx
│   │   ├── RideTrackingScreen.tsx
│   │   ├── PaymentScreen.tsx
│   │   └── PassengerProfileScreen.tsx
│   │
│   ├── captain/               # Captain mode screens
│   │   ├── CaptainHomeScreen.tsx
│   │   ├── IncomingRequestScreen.tsx
│   │   ├── CaptainNavigationScreen.tsx
│   │   ├── EarningsScreen.tsx
│   │   └── CaptainProfileScreen.tsx
│   │
│   ├── auth/                  # Onboarding & auth
│   │   ├── LoginScreen.tsx
│   │   ├── RegisterScreen.tsx
│   │   └── RoleSelectionScreen.tsx
│   │
│   ├── settings/
│   │   ├── SettingsScreen.tsx
│   │   └── NotificationsScreen.tsx
│   │
│   ├── navigation/            # Navigation logic
│   │   ├── PassengerNavigator.tsx
│   │   ├── CaptainNavigator.tsx
│   │   ├── AuthNavigator.tsx
│   │   ├── MainNavigator.tsx
│   │   └── index.tsx
│   │
│   ├── components/            # Reusable UI components
│   │   ├── MapViewComponent.tsx
│   │   ├── RideCard.tsx
│   │   ├── CaptainCard.tsx
│   │   ├── Button.tsx
│   │   ├── Input.tsx
│   │   ├── Modal.tsx
│   │   ├── Loading.tsx
│   │   └── Header.tsx
│   │
│   └── hooks/                 # Logic reuse
│       ├── useLocation.ts
│       ├── useWebSocket.ts
│       ├── useAuth.ts
│       ├── useRideMatching.ts
│       └── useAppTheme.ts
│
├── assets/                    # Images, logos, fonts
│   ├── images/
│   ├── icons/
│   └── fonts/
│
├── constants/                 # App-level constants
│   ├── Colors.ts
│   ├── Fonts.ts
│   ├── Config.ts
│   └── ApiRoutes.ts
│
├── services/                  # API & external services
│   ├── api/
│   │   ├── auth.api.ts
│   │   ├── rides.api.ts
│   │   ├── payments.api.ts
│   │   ├── captains.api.ts
│   │   └── notifications.api.ts
│   │
│   ├── socket/                # WebSocket connection
│   │   └── socketClient.ts
│   │
│   └── location/
│       └── locationService.ts # Permissions + background tracking
│
├── store/                     # Global state
│   ├── auth.store.ts
│   ├── user.store.ts
│   ├── ride.store.ts
│   ├── captain.store.ts
│   ├── location.store.ts
│   └── index.ts
│
├── utils/                     # Helper functions
│   ├── formatDistance.ts
│   ├── calculateFare.ts
│   ├── toast.ts
│   └── validators.ts
│
├── types/                     # TypeScript types/interfaces
│   ├── Ride.ts
│   ├── User.ts
│   ├── Captain.ts
│   ├── Payment.ts
│   └── index.ts
│
├── app.json
├── package.json
└── index.js
