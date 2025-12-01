# 04 – Display Passenger Current Location on Map

## Overview
Enable passengers to see their real-time GPS location on the Home screen map. This feature helps users request nearby rides easily by showing their current position with high accuracy and updating the map reactively as their location changes.

## Source
- Issue: Internal Story (no Jira URL)
- Type: Story
- Priority: High
- Labels: UI, Passenger, Location, Map
- Reporter: Product Team
- Linked Dependencies: None

## Scope
### In Scope
- Request location permission from user
- Fetch live GPS coordinates using device location services
- Display current location as a blue dot/marker on the map
- Auto-center map on user’s location on initial load
- Update map reactively as location changes
- Handle permission denied and GPS loss gracefully

### Out of Scope
- Backend location storage or streaming
- Ride request logic
- Historical location tracking
- Map styling beyond marker and centering
- Location sharing with other users

## Requirements
- Functional:
  - Request and handle location permission
  - Fetch and monitor live GPS coordinates
  - Display current location marker (blue dot) on map
  - Auto-center map on user’s location at startup
  - Update marker and map position as location changes
  - Show alert and fallback state if permission denied
  - Retain last known location if GPS is lost
- Non-Functional:
  - GPS accuracy < 50m
  - Location updates must be reflected in UI within 1 second
  - Map must load within 2 seconds on Home screen
  - Handle permission and GPS errors gracefully

## Data & Contracts
- No backend API changes required
- Uses device geolocation API (foreground)
- Data shape: `{ latitude: number, longitude: number, accuracy: number }`
- UI state: permission status, current coordinates, last known location

## UI / Components
- Home screen map (React Native Maps or Google Maps SDK)
- Blue dot/marker for passenger’s current location
- Permission request dialog
- Alert for permission denied
- Fallback UI for GPS loss (show last known location)

## State & Logic
- State: permission status, current coordinates, last known location
- Logic:
  - On Home screen load, request location permission
  - If granted, fetch and subscribe to location updates
  - Center map on current location
  - Update marker and map position as location changes
  - If permission denied, show alert and fallback state
  - If GPS lost, keep last known location visible

## Acceptance Criteria
- [ ] Requests location permission on Home screen load
- [ ] Fetches live GPS coordinates with accuracy < 50m
- [ ] Displays blue dot/marker for current location
- [ ] Auto-centers map on user’s location at startup
- [ ] Updates marker and map position in real-time
- [ ] Shows alert and fallback UI if permission denied
- [ ] Keeps last known location if GPS is lost
- [ ] Map loads within 2 seconds
- [ ] Location updates reflected in UI within 1 second

## Edge Cases
- Permission denied: show alert, fallback UI
- GPS temporarily lost: retain last known location
- Location accuracy > 50m: show warning or fallback
- App backgrounded: pause location updates
- Rapid location changes: debounce updates to avoid UI jitter
- Device location services disabled: show error

## Risks & Mitigations
| Risk | Impact | Mitigation |
|------|--------|------------|
| Permission denied | Feature unusable | Show clear alert, fallback UI |
| GPS signal weak | Inaccurate location | Show warning, fallback to last known |
| Location update lag | Poor UX | Use high-frequency updates, debounce UI |
| Battery drain | User dissatisfaction | Use foreground geolocation only, optimize update interval |
| Map fails to load | Broken UI | Add loading state, error handling |
| Device API differences | Inconsistent behavior | Test on multiple devices, use standard libraries |

## Open Questions / Assumptions
- Q: Should we show a warning if accuracy > 50m? / A: Yes, show fallback or warning
- Q: Is map required to update in background? / A: No, only foreground updates
- Q: What map provider? / A: React Native Maps or Google Maps SDK
- Q: Should we support custom marker icons? / A: No, use standard blue dot for MVP
- Q: Is historical location needed? / A: No, only current location

## Story Points
3 – Reason: Simple UI feature using standard libraries, requires permission handling, real-time updates, and error states. No backend changes or complex logic.

## Implementation Plan
1. Integrate React Native Maps or Google Maps SDK in Home screen
2. Implement location permission request logic
3. Fetch and subscribe to live GPS coordinates
4. Display blue dot/marker for current location
5. Auto-center map on user’s location at startup
6. Update marker and map position reactively
7. Handle permission denied and GPS loss (alert, fallback)
8. Debounce rapid location updates for smooth UI
9. Test on multiple devices for accuracy and performance
10. Add loading/error states for map and location

## TODO Checklist
- [ ] Integrate map component in Home screen
- [ ] Implement location permission request
- [ ] Fetch live GPS coordinates
- [ ] Display blue dot/marker for current location
- [ ] Auto-center map on initial load
- [ ] Subscribe to location updates
- [ ] Update marker/map position in real-time
- [ ] Handle permission denied (alert, fallback)
- [ ] Retain last known location if GPS lost
- [ ] Debounce rapid location changes
- [ ] Add loading/error states
- [ ] Test on multiple devices
- [ ] Code review and merge
- [ ] Update documentation

## Testing Strategy
- Unit: Permission logic, location update handler
- Component: Map rendering, marker display, alert/fallback UI
- Integration: End-to-end Home screen flow, permission denied, GPS loss
- Manual QA: Test on real devices (iOS/Android), check accuracy, UI responsiveness

## Performance Considerations
- Use high-frequency location updates only in foreground
- Debounce rapid updates to avoid UI jitter
- Optimize map rendering for low-end devices
- Minimize battery usage by limiting update interval
- Target map load < 2 seconds, location update < 1 second

## Definition of Done
- All acceptance criteria met
- Checklist complete
- No errors/warnings
- Code clean and reviewed
- Documentation updated if needed
