Create a Jira task using the Jira tools based on the provided inputs.

Inputs (fill all applicable placeholders):

- Project Key: {jira_project_key}
- Issue Type (e.g. Task, Bug, Story, Subtask): {issue_type}
- Task Idea (raw user idea / goal): {task_idea}
- Business / Technical Context: {context}
- Priority (Blocker/Critical/High/Medium/Low): {priority}
- Impact (qualitative + any metrics): {impact}
- Dependencies (issue keys or external): {dependencies}
- Component / Area: {component}
- Labels (comma-separated): {labels}
- Due Date (YYYY-MM-DD): {due_date}
- Assignee (display name or leave blank): {assignee}
- Epic Key (if linking to an epic): {epic_key}
- Parent Issue Key (for Subtask only): {parent_issue_key}
- Acceptance Criteria (bullets or plain text): {acceptance_criteria}
- Definition of Done (bullets or plain text): {definition_of_done}
- Attachments / References (URLs, IDs): {attachments}

Instructions:

1. Parse the inputs. If a placeholder is empty or missing, use "TBD" except: skip epic_key or parent_issue_key if not provided. Do not invent false metrics.
2. Derive a concise, action-oriented Summary (max ~12 words) from {task_idea}. Prefix with an action verb. Avoid redundancy with project or epic names.
3. Construct a rich Description in Markdown with the following sections (omit a section if all data is TBD):
   **Overview**: 1–3 sentences translating {task_idea} and {context} into a clear objective.
   **Motivation / Impact**: Summarize {impact}. If metrics missing, note expected qualitative benefit.
   **Scope**: Explicit in-scope bullets inferred from idea + context.
   **Out of Scope**: List obvious exclusions to prevent scope creep (at least one when possible).
   **Dependencies**: Bullet list from {dependencies}; mark unknown as TBD.
   **Technical Notes**: Any architectural / implementation hints deduced from context (keep high-value, concise).
   **Risks & Mitigations**: 1–3 bullets if risks apparent; else omit.
   **Attachments / References**: From {attachments}.
4. Acceptance Criteria: Normalize {acceptance_criteria}. If freeform text, convert into bullet list of behavioral, testable statements using Given/When/Then or Must/Should style. Ensure each criterion is independently verifiable.
5. Definition of Done: Normalize {definition_of_done}. Append defaults if absent (e.g. "Code reviewed", "Tests updated & passing", "No critical security findings", "Documentation updated"). Avoid duplication with Acceptance Criteria.
6. Add an Epic link only if {epic_key} provided and not TBD.
7. For Subtask issue_type, require {parent_issue_key}; if missing, set summary to include "(PARENT NEEDED)".
8. Sanitize labels: lowercase, kebab-case, strip spaces. Remove empty entries. Ensure uniqueness.
9. Ensure final description length stays focused; trim filler; avoid repeating Summary.
10. Prepare data for Jira tool invocation.

Output Format:
First output a Markdown preview:

---

Summary: <generated summary>
Issue Type: {issue_type}
Project: {jira_project_key}
Priority: {priority}
Component: {component}
Labels: <normalized_labels>
Due Date: {due_date}
Assignee: {assignee}
Epic: {epic_key}
Parent: {parent_issue_key}

Description (Markdown):
<full description>

Acceptance Criteria:

- <bullet 1>
- <bullet 2>
  ...

Definition of Done:

- <bullet 1>
  ...

---

Then invoke the Jira create issue tool with (only after preview generation):
projectKey={jira_project_key}
issueTypeName={issue_type}
summary=<generated summary>
assignee_account_id=<resolved or omit if blank>
parent=<parent_issue_key if issue_type is Subtask>
epic_key=<epic_key if provided>
labels=<comma separated normalized labels>
description=<full description including Acceptance Criteria + Definition of Done sections>
due_date={due_date if not TBD}

Rules:

- Do NOT fabricate data beyond structured derivations.
- Preserve user-provided technical details verbatim when present.
- Keep everything actionable and testable.
- Use American English.
- No extraneous commentary outside specified output.
