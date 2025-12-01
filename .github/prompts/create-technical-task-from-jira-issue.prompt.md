# Create Technical Task From Jira Issue

## Purpose

Turn a Jira issue (Story or Bug) into a clear technical task file in the `tasks/` folder. The task should have an overview, acceptance criteria, plan, risks, story points, and a simple checklist.

## How to Use

1. Get the Jira issue details (issue key from user or context).
2. Use Jira Tool to read the main fields: summary, description, type (Story/Bug), priority, labels, story points (if any), acceptance criteria, attachments (summarize), linked issues (note dependencies), reporter.
3. Analyze the issue:

- Clarify the main goal and business need.
- Identify the main area (auth, UI, backend, etc.).
- List what needs to be done (functional and non-functional requirements).
- List any unclear points or assumptions.

4. Estimate story points:

- If missing, use Fibonacci (1,2,3,5,8,13) based on complexity and risk. Give a short reason.

5. Name the task file: `TASK_<issueKey>_<SNAKE_CASE_SUMMARY>.md` (use only letters, numbers, underscores, max 6 words after issue key).
6. Save the file in `/tasks` using the structure below.

## Task File Structure

```
# <Issue Key> – <Title>

## Overview
Short summary of the task and its value.

## Source
- Issue: <Jira URL>
- Type: Story | Bug
- Priority: <priority>
- Labels: <labels or ->
- Reporter: <name>
- Linked Dependencies: <keys or ->

## Scope
### In Scope
- ...
### Out of Scope
- ...

## Requirements
- Functional:
  - ...
- Non-Functional:
  - ...

## Data & Contracts
- List any new or changed APIs, data shapes, or contracts.

## UI / Components
- List new or updated UI parts, if any.

## State & Logic
- Where state lives, main logic, special considerations.

## Acceptance Criteria
- [ ] Each item should be clear and testable.

## Edge Cases
- List possible edge cases (errors, empty, slow, etc.).

## Risks & Mitigations
| Risk | Impact | Mitigation |
|------|--------|------------|

## Open Questions / Assumptions
- Q: ... / A: Assumed ...

## Story Points
<Number> – Reason: <short explanation>

## Implementation Plan
1. ...
2. ...
3. ...

## TODO Checklist
- [ ] ...
- [ ] ...

## Testing Strategy
- Unit: ...
- Component: ...
- Integration: ...
- Manual QA: ...

## Performance Considerations
- List any performance tips or needs.

## Definition of Done
- All acceptance criteria met
- Checklist complete
- No errors/warnings
- Code clean and reviewed
- Documentation updated if needed
```

## Tool Guidelines

- Always get Jira issue details first. If it fails, try again up to 2 times. If it still fails, make a partial task and mark it as BLOCKED.
- If story points field has a different name (like `customfield_10016`), find the number.

## Estimation Guide

- Very small change: 1
- Simple task: 2
- New feature with one API: 3
- Medium feature, some state or UI: 5
- Complex or cross-team: 8
- Big, risky, or unclear: 13

## Quality Reminders

- No hardcoded UI text; use translation if needed.
- Validate all data before use.
- Keep code clean and well-typed.
- Refactor if code gets too big.

## Output Rules

- Make one markdown file in `/tasks` for each Jira issue.
- Do not overwrite unless same issue key; if so, add a `## Revision <timestamp>` section.
- Use checkboxes for acceptance criteria and TODOs.
- Keep language simple and clear.
- Do not add details not in the issue or assumptions.
