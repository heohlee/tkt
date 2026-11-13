---
name: issue-governance
description: "Trigger: creating an issue, starting a feature, fixing a bug, labeling, planning work, or any change to the tkt repository. Issues are optional and never gate a change; when one exists, apply the canonical type/area label taxonomy best-effort."
license: MIT
metadata:
  author: "heohlee"
  version: "1.0"
---

## Activation Contract

Run when a TKT change:

- is about to start, to decide whether an issue is worth opening first;
- requires creating, labeling, or triaging a GitHub issue;
- transitions from planning to implementation with an issue in place.

Do NOT run for:

- read-only Q&A with no repository change.

## Hard Rules

- An issue is optional for every change and every contributor: a pull request may arrive without one, and no CI check requires it. Recommend opening an issue first only for large or design-changing work, so the idea can be discussed before code is written. Never block, close or delay a change for lacking an issue.
- When an issue exists, it is the intent anchor for the change (SDD artifacts, Engram keys) and the target of the PR's `Closes #<n>` line.
- No work ever waits for an approval step; `status:approved` is informational, does not gate and never blocks implementation.
- Labels come ONLY from the taxonomy below; the repository contains EXACTLY these 12 labels. Never invent new labels. Apply them best-effort with the closest fit and note the choice; a missing, ambiguous, or imperfect label never blocks or delays work.

### Label Taxonomy

| Group | Rule | Values |
| --- | --- | --- |
| Type | exactly ONE | `type:bug` — incorrect behavior; `type:feature` — visible functionality or improvement; `type:docs` — documentation only; `type:chore` — CI, tooling, tests, skills, templates, or maintenance with no functional change. |
| Area | ONE primary, at most TWO | `area:auth`, `area:tickets`, `area:users`, `area:categories-workflows`, `area:desks`, `area:settings`, `area:tooling`. |
| Status | only when applicable, never a gate | `status:approved` — optional marker. It authorizes nothing and never blocks work. |

- `area:tooling` covers CI, Playwright, skills, Engram change governance, and repository configuration.
- Cardinality rules apply to the governing issue when one exists. PR labels are optional and do not need to duplicate the issue's labels.
- If no area fits a proposed issue, pick the closest fit and note the mismatch in the issue body; never block or delay work over it.

### Issue Title Prefixes

New issue titles should start with one canonical prefix that matches the issue's `type:*` label:

| Title prefix | Label |
| --- | --- |
| `[Feature]` | `type:feature` |
| `[Bug]` | `type:bug` |
| `[Chore]` | `type:chore` |
| `[Docs]` | `type:docs` |

Example: `[Feature] Add ticket export to CSV`.

Labels remain authoritative for filtering, automation, and governance. The prefix is a human-readable aid in issue lists and search results, not a replacement for the label. A missing or wrong prefix is not a labeling violation and MUST NOT become a new decision gate; when prefix and label disagree, the label wins.

The feature request template applies the `[Feature]` prefix automatically. Apply the other prefixes manually when drafting an issue. Do not bulk-edit existing issues to add prefixes; existing titles may keep their current format.

### Why there is no issue or approval gate

Requiring an issue before every pull request discourages contributions and doubles the moderation work: an issue and a pull request for the same change. The rule is uniform across contributors and change types on purpose: any rule that has to classify the person or the change creates a fresh judgment gate. Do not reintroduce issue or approval requirements through labels, title conventions or contributor roles.

## Decision Gates

| Condition | Result |
| --- | --- |
| No issue exists for the change | GO: the issue is optional; suggest one only for large or design-changing work |
| Issue exists | GO — proceed (respect other project skills) |
| Issue lacks a `type:*` label, or has more than one | GO — apply the best fit or note the correction; never block or delay |
| Issue has zero areas or more than two | GO — apply the best fit or note the correction; never block or delay |
| Issue has no `status:approved` | GO — approval is not a gate in this repository |

## Execution Steps

1. **Check tracking**: search open issues for the requested work. If one exists, work from it. Open a new issue (title with the canonical prefix, motivation, acceptance criteria) only for large or design-changing work; otherwise start right away without one.
2. **Apply labels** (when an issue exists): pick exactly one `type:*` and one primary `area:*` (max two). Best-effort is enough — a doubtful label is never a reason to wait.
3. **Implement**: start immediately; never wait for an issue or an approval step.
4. **Keep labels truthful**: if an issue exists and scope changes mid-work so a different type/area applies, update the issue labels before continuing.
5. **Close the loop**: when an issue exists, the PR links it with a standalone `Closes #<n>` line.

## Output Contract

Return:

- `issue`: number/URL of the governing issue, or "none";
- `labels_valid`: true / false — with the exact taxonomy deviation if false;
- `labels_valid` applies only when an issue exists;
- `gate`: GO — a missing issue, missing approval or imperfect label never makes it BLOCKED.

## References

- `../../../AGENTS.md` — project skill registration (this skill is listed there).
- GitHub labels: `gh label list --repo heohlee/tkt`.
