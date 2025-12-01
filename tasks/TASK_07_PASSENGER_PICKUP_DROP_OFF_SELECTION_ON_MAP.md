# 07 – Passenger Pickup & Drop-off Selection on Map

## Overview
Enable passengers to select both pickup and drop-off locations directly on the map. This feature streamlines the ride request process, allowing users to visually choose their start and end points for each trip.

## Source
- Issue: Internal Story (no Jira URL)
- Type: Story
- Priority: High
- Labels: UI, Passenger, Map, Pickup, Drop-off
- Reporter: Product Team
- Linked Dependencies: None

## Scope
### In Scope
- Interactive map for selecting pickup and drop-off locations
- Display draggable markers for pickup and drop-off
- Show address or place name for selected points
- Confirm button to finalize selection
- Update backend and UI with selected coordinates
- Auto-center map on user’s location at startup
- Support for manual address entry as fallback

### Out of Scope
- Route calculation or fare estimation
- Multi-stop trips
- Backend geocoding (handled by external API)
- Historical pickup/drop-off suggestions
- Admin override of locations

## Requirements
- Functional:
  - Display map with current location centered
  - Allow user to select pickup location by tapping or dragging marker
  - Allow user to select drop-off location by tapping or dragging marker
  - Show address or place name for each selected point
  - Confirm selection to set pickup/drop-off
  - Update backend and UI with selected coordinates
  - Support manual address entry if map selection fails
- Non-Functional:
  - Location selection must complete within 10 seconds
  - Markers must be draggable and responsive
  - UI must handle permission and GPS errors gracefully
  - Map must load within 2 seconds

## Data & Contracts
- Data shape:
  - `pickup: { latitude: number, longitude: number, address: string }`
  - `dropoff: { latitude: number, longitude: number, address: string }`
- Backend receives coordinates and address/place name
- Manual address entry supported as fallback

## UI / Components
- Map view with draggable pickup and drop-off markers
- Address display for selected points
- Confirm button
- Manual address entry field
- Loading indicator for map and geocoding
- Error/fallback UI for permission/GPS issues

## State & Logic
- State: current location, pickup marker, drop-off marker, selected addresses
- Logic:
  - On Home screen load, center map on current location
  - User taps or drags marker to select pickup/drop-off
  - Fetch address/place name via geocoding API
  - Display address for each marker
  - Confirm button sends selection to backend/UI
  - Manual address entry available if needed
  - Handle permission denied and GPS loss gracefully

## Acceptance Criteria
- [ ] Map displays with current location centered
- [ ] User can select pickup location by tapping/dragging marker
- [ ] User can select drop-off location by tapping/dragging marker
- [ ] Address or place name shown for each selected point
- [ ] Confirm button finalizes selection
- [ ] Backend and UI updated with coordinates
- [ ] Manual address entry supported as fallback
- [ ] Markers are draggable and responsive
- [ ] Map loads within 2 seconds
- [ ] Location selection completes within 10 seconds
- [ ] Error/fallback UI for permission/GPS issues

## Edge Cases
- Permission denied: show alert, fallback UI
- GPS temporarily lost: retain last known location
- Invalid address: allow manual entry
- Rapid marker movement: debounce address lookup
- Overlapping pickup/drop-off: show warning
- Map fails to load: show error
- App backgrounded: pause location selection

## Risks & Mitigations
| Risk | Impact | Mitigation |
|------|--------|------------|
| Permission denied | Feature unusable | Show alert, fallback UI |
| GPS signal weak | Inaccurate selection | Show warning, allow manual entry |
| Address lookup slow | Poor UX | Debounce requests, show loading indicator |
| Markers not responsive | User frustration | Test on multiple devices, optimize UI |
| Backend update fails | Ride request blocked | Retry logic, error feedback |
| Overlapping pickup/drop-off | Invalid trip | Show warning, require distinct points |

## Open Questions / Assumptions
- Q: Can pickup and drop-off be the same? / A: No, require distinct points
- Q: Is manual address entry always available? / A: Yes, as fallback
- Q: What geocoding API? / A: Google Maps or Mapbox
- Q: Are markers customizable? / A: Use standard icons for MVP
- Q: Is multi-stop needed? / A: No

## Story Points
5 – Reason: Medium complexity. Requires interactive map, draggable markers, geocoding, backend update, error handling, and responsive UI.

## Implementation Plan
1. Integrate map component with current location centering
2. Add draggable pickup and drop-off markers
3. Implement tap/drag logic for marker selection
4. Fetch address/place name via geocoding API
5. Display address for selected points
6. Add confirm button to finalize selection
7. Update backend and UI with coordinates
8. Support manual address entry as fallback
9. Handle permission/GPS errors and edge cases
10. Test on multiple devices for responsiveness
11. Code review and merge
12. Update documentation

## TODO Checklist
- [ ] Integrate map component
- [ ] Add draggable pickup marker
- [ ] Add draggable drop-off marker
- [ ] Implement tap/drag selection logic
- [ ] Fetch address via geocoding API
- [ ] Display address for markers
- [ ] Add confirm button
- [ ] Update backend/UI with selection
- [ ] Support manual address entry
- [ ] Handle errors and edge cases
- [ ] Test on multiple devices
- [ ] Code review
- [ ] Update documentation

## Testing Strategy
- Unit: Marker selection logic, address lookup
- Component: Map rendering, marker drag, address display
- Integration: End-to-end selection flow, backend update
- Manual QA: Test with permission denied, GPS loss, rapid marker movement

## Performance Considerations
- Debounce address lookup to minimize API calls
- Optimize map rendering for low-end devices
- Target map load < 2 seconds, selection < 10 seconds

## Definition of Done
- All acceptance criteria met
- Checklist complete
- No errors/warnings
- Code clean and reviewed
- Documentation updated if needed
