Transform the provided feature idea and optional context into A single INVEST-compliant User Story

Inputs:

- Feature Idea: TJ
- Context / PRD Excerpt (optional): {context_text}
- Extra Constraints (optional): {constraints}

Process:

1. Parse inputs; extract persona(s), core value, key functional and non-functional hints. List any gaps as Questions & Assumptions.
2. Define minimal viable scope (avoid future phases). Exclude out-of-scope items explicitly.
3. Draft Story (no platform-specific tech):
   - Title (<100 chars) "<Persona> can <capability> to <benefit>".
   - Narrative: As a {persona} I want {capability} so that {benefit}.
   - Business Value: bullet list.
   - Functional Requirements: bullets (FR1, FR2 ...).
   - Non-Functional: performance, reliability, security/privacy, usability/accessibility, localization (if relevant), observability.
   - Acceptance Criteria: checklist; each atomic & testable (success + validation + error/edge + permission + empty state + performance threshold if relevant).
   - Definition of Done: concise bullet list (ACs pass, tests added, docs updated, no unresolved questions).
   - Open Questions & Assumptions: bullets (Q: / A: pattern).
4. Labels: derive up to 5 from domains (e.g., auth, billing, content, search, analytics, localization, accessibility, performance).
5. Output final artifacts plus JSON summary.

Output Markdown Structure:

# {story_title}

## Narrative

...

## Business Value

- ...

## Functional Requirements

- FR1 ...

## Non-Functional

- Performance: ...
- Reliability: ...
- Security & Privacy: ...
- Usability & Accessibility: ...
- Localization: ... (omit if N/A)
- Observability: ...

## Acceptance Criteria

- [ ] ...

## Definition of Done

- ...

## Open Questions & Assumptions

- Q: ... / A: Assumed ...

# TECH: {story_title}

## Overview

## Key Design Notes

## Implementation Outline

1. ...

## Story Points

{story_point_estimate} (rationale)

Guidelines:

- Keep concise; avoid platform/framework specifics.
- Use placeholders not concrete code.
- Do not create multiple stories; keep scope small.
- Proceed with reasonable assumptions if context sparse; list them.
- Use imperative, professional tone. Avoid redundant wording.
