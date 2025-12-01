# 05 – Search With Auto-Complete for Destination

## Overview
Enable passengers to search for destinations using an auto-complete search bar powered by Google Places API (or Mapbox). This feature improves trip planning by providing real-time suggestions and allowing users to select a destination with full place details.

## Source
- Issue: Internal Story (no Jira URL)
- Type: Story
- Priority: High
- Labels: UI, Passenger, Search, Map, Destination
- Reporter: Product Team
- Linked Dependencies: None

## Scope
### In Scope
- Auto-complete search bar for destination input
- Integration with Google Places Autocomplete API (or Mapbox)
- Display suggestions as user types
- Select suggestion to set destination (lat/lng + name)
- Debounce input (300 ms)
- Secure API key storage
- Clear button to reset search
- No API call if input < 2 characters

### Out of Scope
- Backend place storage or history
- Multi-destination support
- Custom place ranking or filtering
- Offline search
- Map rendering beyond destination marker

## Requirements
- Functional:
  - Show auto-complete suggestions as user types
  - Display suggestions popup below input
  - Selecting a suggestion sets destination (lat/lng + name)
  - Each result includes: name, place_id, latitude, longitude
  - Clear button resets search input and suggestions
  - No API call if input < 2 characters
  - Debounce input by 300 ms before API call
  - Securely store API key (not exposed in client code)
- Non-Functional:
  - Suggestions must appear within 500 ms of typing
  - API key must not be hardcoded or exposed
  - Handle API errors gracefully (show fallback UI)
  - Support for international place names

## Data & Contracts
- Uses Google Places Autocomplete API (or Mapbox)
- Request: `input` (string, min 2 chars)
- Response: Array of suggestions:
  - `name` (string)
  - `place_id` (string)
  - `latitude` (float)
  - `longitude` (float)
- UI state: input value, suggestions array, selected destination

## UI / Components
- Search bar with auto-complete
- Suggestions popup below input
- Clear button
- Destination marker on map (after selection)
- Loading indicator for API calls
- Error/fallback UI for API failures

## State & Logic
- State: input value, suggestions, selected destination
- Logic:
  - On input change, debounce and call API if input >= 2 chars
  - Display suggestions popup
  - On suggestion select, set destination (lat/lng + name)
  - Clear button resets input and suggestions
  - Handle API errors (show fallback)
  - Securely load API key (env/config)

## Acceptance Criteria
- [ ] Auto-complete search bar appears on destination input
- [ ] Suggestions popup below input as user types
- [ ] Each suggestion includes name, place_id, latitude, longitude
- [ ] Selecting a suggestion sets destination in UI and backend
- [ ] No API call if input < 2 characters
- [ ] Input is debounced by 300 ms before API call
- [ ] Clear button resets search
- [ ] API key is securely stored
- [ ] Suggestions appear within 500 ms
- [ ] Error/fallback UI shown for API failures

## Edge Cases
- Input < 2 characters: no API call, no suggestions
- API error: show fallback UI
- No suggestions returned: show "No results" message
- Rapid typing: debounce prevents excessive API calls
- International characters: support Unicode input
- API rate limit exceeded: show error
- Place with missing lat/lng: skip or show warning
- Selecting same place twice: update destination

## Risks & Mitigations
| Risk | Impact | Mitigation |
|------|--------|------------|
| API key exposed | Security breach | Store key in env/config, never in client code |
| API rate limit | Feature unusable | Debounce input, handle errors gracefully |
| Slow API response | Poor UX | Show loading indicator, fallback UI |
| Incorrect place details | Wrong destination | Validate response, show details before confirm |
| No suggestions | User confusion | Show "No results" message |
| Internationalization issues | Incomplete search | Support Unicode, test with various locales |

## Open Questions / Assumptions
- Q: Which API provider? / A: Google Places preferred, Mapbox as fallback
- Q: Should backend store selected destination? / A: Yes, after selection
- Q: Is multi-destination needed? / A: No, single destination for MVP
- Q: Should suggestions be cached? / A: No, live API for accuracy
- Q: Is offline search required? / A: No

## Story Points
5 – Reason: Medium complexity. Requires API integration, UI state management, debouncing, secure key handling, and error states.

## Implementation Plan
1. Integrate Google Places Autocomplete API (or Mapbox)
2. Implement search bar with auto-complete
3. Add suggestions popup below input
4. Debounce input by 300 ms before API call
5. Display suggestions with name, place_id, lat/lng
6. Handle selection to set destination in UI and backend
7. Add clear button to reset search
8. Securely store and load API key
9. Handle API errors and show fallback UI
10. Test with international input and edge cases
11. Code review and merge
12. Update documentation

## TODO Checklist
- [ ] Integrate Places Autocomplete API
- [ ] Implement search bar with auto-complete
- [ ] Add suggestions popup
- [ ] Debounce input (300 ms)
- [ ] Display suggestions (name, place_id, lat/lng)
- [ ] Handle selection to set destination
- [ ] Add clear button
- [ ] Secure API key storage
- [ ] Handle API errors/fallback UI
- [ ] Test with international input
- [ ] Code review
- [ ] Update documentation

## Testing Strategy
- Unit: Input debounce, API call logic
- Component: Search bar, suggestions popup, clear button
- Integration: End-to-end search and selection flow
- Manual QA: Test with various inputs, API errors, rate limits

## Performance Considerations
- Debounce input to minimize API calls
- Optimize suggestions rendering for low-end devices
- Secure API key storage
- Target suggestions popup < 500 ms

## Definition of Done
- All acceptance criteria met
- Checklist complete
- No errors/warnings
- Code clean and reviewed
- Documentation updated if needed
