# 06 – Display Suggested Places & Recent Destinations

## Overview
Enable passengers to quickly select destinations by displaying smart suggestions (e.g., home, work, popular areas) and a list of recent destinations based on user history. This feature streamlines trip planning and reduces typing for frequent users.

## Source
- Issue: Internal Story (no Jira URL)
- Type: Story
- Priority: High
- Labels: UI, Passenger, Places, Suggestions, History
- Reporter: Product Team
- Linked Dependencies: None

## Scope
### In Scope
- Store trip history in backend for each passenger
- Endpoint to fetch recent destinations (`GET /api/v1/passenger/places/recent`)
- Endpoint to fetch suggested places (`GET /api/v1/passenger/places/suggestions`)
- Display recent and suggested places on Home screen
- Show name and small icon for each item
- Selecting a suggestion or recent destination sets it as the destination
- Recent list updates automatically after each ride
- Suggested places based on usage patterns (e.g., Work at 9am, Home at night)

### Out of Scope
- Manual editing of suggested places
- Multi-destination support
- Custom place icons beyond standard set
- Offline suggestions/history
- Admin management of suggestions

## Requirements
- Functional:
  - Store trip history for each passenger
  - Return up to 5 recent destinations, sorted by most recent
  - Return smart suggested places (home, work, popular areas)
  - Display both lists on Home screen
  - Each item shows name and icon
  - Selecting an item sets destination instantly
  - Recent list updates after each ride
  - Suggestions update based on usage patterns
- Non-Functional:
  - Endpoints must respond within 500 ms
  - Data stored securely by passenger ID
  - Lists update reactively in UI
  - Handle empty states gracefully (no history/suggestions)

## Data & Contracts
- Backend endpoints:
  - `GET /api/v1/passenger/places/recent` → returns array of recent destinations (limit 5)
  - `GET /api/v1/passenger/places/suggestions` → returns array of suggested places
- Data shape:
  - `{ name: string, place_id: string, latitude: number, longitude: number, icon: string }`
- Stored by passenger ID
- Recent destinations updated after each ride
- Suggestions generated based on time, frequency, patterns

## UI / Components
- Home screen: recent destinations list, suggested places list
- Each item: name, small icon
- Empty state: "No recent destinations yet"
- Loading indicator for API calls

## State & Logic
- State: recent destinations, suggested places, selected destination
- Logic:
  - On Home screen load, fetch both lists from backend
  - Display lists with name and icon
  - Selecting an item sets destination in UI and backend
  - Recent list updates after ride completion
  - Suggestions update based on usage patterns (e.g., time of day)
  - Handle empty states and API errors

## Acceptance Criteria
- [ ] Backend stores trip history per passenger
- [ ] `GET /api/v1/passenger/places/recent` returns up to 5 recent destinations
- [ ] `GET /api/v1/passenger/places/suggestions` returns smart suggestions
- [ ] Home screen displays both lists with name and icon
- [ ] Selecting an item sets destination instantly
- [ ] Recent list updates after each ride
- [ ] Suggestions update based on usage patterns
- [ ] Lists update reactively in UI
- [ ] Empty states handled gracefully
- [ ] Endpoints respond within 500 ms

## Edge Cases
- No recent destinations: show empty state
- No suggestions: show empty state
- API error: show fallback UI
- Duplicate destinations: show only unique items
- Rapid ride completions: debounce updates
- Place with missing lat/lng: skip or show warning
- Passenger with no history: show suggestions only
- Selecting same place twice: update destination

## Risks & Mitigations
| Risk | Impact | Mitigation |
|------|--------|------------|
| API slow or fails | Poor UX | Add loading/fallback UI, optimize queries |
| Data not updating | Outdated lists | Use reactive updates, poll after ride completion |
| Duplicate places | Confusing UI | Deduplicate before display |
| Privacy concerns | Data misuse | Store history securely, limit retention |
| Incorrect suggestions | Poor relevance | Improve pattern detection, allow feedback |
| Empty states | User confusion | Show clear empty state messages |

## Open Questions / Assumptions
- Q: How are suggested places determined? / A: Based on time, frequency, patterns (e.g., Work at 9am)
- Q: Can passenger edit suggestions? / A: No, only auto-generated for MVP
- Q: How long is history retained? / A: Assume 30 days
- Q: Are icons customizable? / A: No, use standard set
- Q: Is multi-destination needed? / A: No

## Story Points
5 – Reason: Medium complexity. Requires backend history storage, smart suggestion logic, two endpoints, UI state management, and error handling.

## Implementation Plan
1. Design DB schema for trip history and suggested places
2. Implement backend storage for trip history (by passenger ID)
3. Create endpoint for recent destinations (`GET /api/v1/passenger/places/recent`)
4. Create endpoint for suggested places (`GET /api/v1/passenger/places/suggestions`)
5. Implement logic for smart suggestions (patterns, frequency, time)
6. Update Home screen to display both lists
7. Add icons to each item
8. Handle selection to set destination
9. Update recent list after each ride
10. Handle empty states and API errors
11. Test with various usage patterns and edge cases
12. Code review and merge
13. Update documentation

## TODO Checklist
- [ ] Design DB schema for trip history/suggestions
- [ ] Implement backend storage for trip history
- [ ] Create recent destinations endpoint
- [ ] Create suggested places endpoint
- [ ] Implement smart suggestion logic
- [ ] Update Home screen UI
- [ ] Add icons to items
- [ ] Handle selection to set destination
- [ ] Update recent list after ride
- [ ] Handle empty states/API errors
- [ ] Test with edge cases
- [ ] Code review
- [ ] Update documentation

## Testing Strategy
- Unit: Suggestion logic, history storage, endpoint responses
- Component: List rendering, icon display, selection handling
- Integration: End-to-end Home screen flow, ride completion updates
- Manual QA: Test with various passengers, empty states, rapid rides

## Performance Considerations
- Optimize DB queries for recent/suggested places
- Deduplicate destinations before display
- Target endpoint response < 500 ms
- Minimize data sent to client (limit 5 items)

## Definition of Done
- All acceptance criteria met
- Checklist complete
- No errors/warnings
- Code clean and reviewed
- Documentation updated if needed
