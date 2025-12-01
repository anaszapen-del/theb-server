Create a Jira Bug using the Jira tools based on the provided inputs.

Some of this Inputs will provided from the user:

- Project Key: {jira_project_key}
- Issue Type (must be Bug): {issue_type}
- Bug Title (raw / short): {bug_title}
- High-Level Summary (if given): {summary_hint}
- Detailed Description (freeform): {description}
- Steps To Reproduce (newline or numbered list): {steps_to_reproduce}
- Expected Result: {expected_result}
- Actual Result: {actual_result}
- Environment (OS, browser, device, build, region): {environment}
- Severity (Blocker/Critical/Major/Minor/Trivial): {severity}
- Priority (Blocker/Critical/High/Medium/Low): {priority}
- Affected Version(s): {affected_versions}
- Component / Area: {component}
- Regression? (Yes/No/Unknown): {regression}
- Introduced In (commit/issue/ref if known): {introduced_in}
- Workaround (if any): {workaround}
- Labels (comma-separated): {labels}
- Assignee (display name or blank): {assignee}
- Reporter (display name or blank): {reporter}
- Epic Key (optional): {epic_key}
- Sprint (name or id): {sprint}
- Story Points (if team estimates bugs): {story_points}
- Attachments / References (URLs, IDs): {attachments}
- Additional Notes: {additional_notes}

Instructions:

1. Parse inputs. Any missing/empty placeholder becomes "TBD" except: omit epic_key if missing. Do not fabricate stack traces or logs.
2. Derive concise action-oriented Summary (<= 11 words). Prefer pattern: "Fix <component> <issue succinctly>" or "Resolve <error/behavior> in <component>". Avoid redundant words (Bug, Issue, Jira, Project name).
3. Normalize Steps To Reproduce:
   - Split by newlines or numbering.
   - Ensure each starts with a verb.
   - If fewer than 3 steps and data sparse, keep only provided (no invention).
4. Construct Description in Markdown sections (omit any section where all fields are TBD):
   **Overview**: 1–2 sentences combining {bug_title} / {summary_hint} + impact summary.
   **Environment**: bullet list parsed from {environment}; include Affected Version(s).
   **Severity & Priority**: single line summarizing {severity} / {priority} and rationale if inferable from mismatch or impact; else just values.
   **Steps to Reproduce**: numbered list from normalized steps.
   **Expected Result**: from {expected_result}.
   **Actual Result**: from {actual_result}. Preserve user wording; do not soften issues.
   **Workaround**: from {workaround} or "None known" if explicitly blank.
   **Regression**: state {regression}; if Yes and {introduced_in} known, cite it.
   **Attachments / Logs**: list from {attachments}. Do not create fake logs.
   **Additional Notes**: from {additional_notes} if meaningful.
5. Append a **Triage Checklist** (static bullets):
   - [ ] Reproducible on latest main build
   - [ ] Minimal reproduction documented
   - [ ] Logs / console errors captured (if applicable)
   - [ ] Screenshots or video attached (if UI)
   - [ ] Security/privacy impact assessed
   - [ ] Added/updated automated test (if possible)
6. Labels: lowercase, kebab-case, trim, unique. Add derived label for severity: severity-{lowercase severity}. If regression == Yes add label: regression. Do not duplicate.
7. If {story_points} provided and numeric, append a **Story Points** section at bottom.
8. If Epic key given, include it; else skip.
9. Do not exceed necessary verbosity; keep focused.
10. Prepare for Jira tool invocation.

Output Format:
First output a Markdown preview:

---

Summary: <generated summary>
Issue Type: Bug
Project: {jira_project_key}
Severity: {severity}
Priority: {priority}
Component: {component}
Affected Versions: {affected_versions}
Labels: <normalized_labels>
Assignee: {assignee}
Reporter: {reporter}
Epic: {epic_key}
Sprint: {sprint}
Story Points: {story_points}

Description (Markdown):
<full description>

---

Then invoke the Jira create issue tool with (only after preview generation):
projectKey={jira_project_key}
issueTypeName=Bug
summary=<generated summary>
assignee_account_id=<resolved or omit if blank>
labels=<comma separated normalized labels>
epic_key=<epic_key if provided>
description=<full description including Triage Checklist + Story Points section if any>
custom_fields: affected_versions={affected_versions}; severity={severity}; regression={regression}; introduced_in={introduced_in}; sprint={sprint}; story_points={story_points}

Rules:

- Do NOT fabricate or guess technical details not provided (stack traces, HTTP codes, etc.).
- Preserve user-provided phrasing for Actual Result and Steps (aside from normalization).
- Use American English.
- No extra commentary outside specified output.
