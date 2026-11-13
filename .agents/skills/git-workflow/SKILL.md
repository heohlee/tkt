---
name: git-workflow
description: "Trigger: starting any change in tkt, creating or triaging an issue, labeling, planning work, creating a branch or worktree, or creating, preparing, reviewing or merging a pull request. Issues are optional and never a gate; each change gets its own worktree under ../tkt-worktrees with its own port and database; pull requests follow the template, the required checks and squash merge."
license: MIT
metadata:
  author: "heohlee"
  version: "1.0"
---

## Activation Contract

Use before the first write of any tkt change (feature, bug, chore, docs or skill edit), when an issue needs creating, labeling or triaging, and when a change is pushed for a pull request, needs a PR created or edited, or is checked for review readiness.

Do not use for read-only work: Q&A, reading, searching, inspecting history, running tests or building all leave the main checkout usable. Loading this skill does not authorize GitHub writes, pushing, opening a PR or merging.

Generic PR, branch or issue instructions installed in your agent describe other repositories and do not apply here. They ask for shellcheck over `scripts/*.sh`, `type:refactor` and `type:breaking-change` labels, and PR-validation workflow checks; none of those exist in tkt. Follow this skill rather than importing checks from another repository, and do not try to reconcile those generic instructions with this skill: where they differ, this skill applies and the generic instructions are ignored.

## Hard Rules

- An issue is optional for every change and every contributor: a pull request may arrive without one, and no CI check requires it. Recommend opening an issue first only for large or design-changing work, so the idea can be discussed before code is written. Never block, close or delay a change for lacking an issue, and never require an approval label or an approval round trip. Do not reintroduce issue or approval requirements through labels, title conventions or contributor roles.
- No work ever waits for an approval step: `status:approved` is informational, authorizes nothing and never blocks implementation.
- One change, one worktree, one branch. The main checkout is for reading, syncing and merging.
- One worktree, one port and one database. Two default servers both want `:8080`, and two worktrees must never share a database file.
- Every required check must pass before merge, including on documentation-only PRs. Never interpret readiness as permission to merge.
- Never put tokens, credentials or personal data in issues, PR bodies, comments, screenshots or logs. Summarize or redact logs.
- Respect the current task's authorization: if GitHub writes are prohibited, report the proposed issue, labels or PR without creating or editing anything remotely.

## Issues

An issue is optional; when one exists, it governs the change: it is the intent anchor (SDD artifacts, Engram keys) and the target of the PR's `Closes #<n>` line. Search open issues for the requested work and reuse one if present. For large or design-changing work without one, propose one (title prefix, motivation, acceptance criteria) and create it when authorized; otherwise start right away without one. If scope changes mid-work so a different type or area applies, update the issue labels before continuing.

### Issue Forms

`.github/ISSUE_TEMPLATE/` has a bug report form (title `[Bug] `, applies `type:bug`) and a feature request form (title `[Feature] `, applies `type:feature`). Docs and chore issues, and anything outside the forms, are created with `gh issue create` and labeled by hand.

### Label Taxonomy

| Group | Rule | Values |
| --- | --- | --- |
| Type | Exactly one | `type:bug`: incorrect behavior; `type:feature`: visible functionality or improvement; `type:docs`: documentation only; `type:chore`: CI, tooling, tests, skills, templates or maintenance with no functional change |
| Area | One primary, at most two | `area:auth`, `area:tickets`, `area:users`, `area:categories-workflows`, `area:desks`, `area:settings`, `area:tooling` |
| Status | Only when applicable, never a gate | `status:approved`: optional marker; it authorizes nothing and never blocks work |

- `area:tooling` covers CI, Playwright, skills, Engram change governance and repository configuration.
- The repository contains exactly these 12 labels; verify the live inventory with `gh label list --repo heohlee/tkt`. Never invent labels or import another repository's taxonomy.
- Apply labels best-effort with the closest fit and note the choice; a missing, ambiguous or imperfect label never blocks or delays work. If no area fits, pick the closest and note the mismatch in the issue body.
- Cardinality rules apply to the governing issue when one exists. PR labels are optional and need not duplicate the issue's labels.

Labels are informational, with one effect: `.github/release.yml` groups the generated changelog by PR labels (`type:feature` under Features, `type:bug` under Fixes, everything else under Other changes). An unlabeled PR still appears, under Other changes.

### Issue Title Prefixes

New issue titles start with one prefix matching the issue's `type:*` label:

| Title prefix | Label |
| --- | --- |
| `[Feature]` | `type:feature` |
| `[Bug]` | `type:bug` |
| `[Chore]` | `type:chore` |
| `[Docs]` | `type:docs` |

Example: `[Feature] Add ticket export to CSV`. The forms apply `[Bug]` and `[Feature]`; apply the others manually. The prefix is a human-readable aid, not a replacement for the label: a missing or wrong prefix is not a violation and never becomes a gate; when prefix and label disagree, the label wins. Do not bulk-edit existing titles.

## Branches and Worktrees

Branches are type/slug, lowercase and without spaces, based on an updated `main`. The worktree goes in the sibling directory next to the repository root, never inside it, under ../tkt-worktrees.

| Prefix | Use for |
| --- | --- |
| feat/ | User-facing capability |
| fix/ | Incorrect behavior |
| chore/ | Tooling, skills, repository configuration |
| test/ | Test-only change |
| docs/ | Documentation only |

These are the prefixes this repository uses; do not invent others. The issue number goes in the names only when an issue exists:

| Case | Branch | Worktree directory |
| --- | --- | --- |
| With issue 123 | `<type>/<n>-<slug>`, e.g. feat/123-metrics-link-sync | `../tkt-worktrees/issue-<n>-<slug>` |
| Without an issue | `<type>/<slug>`, e.g. feat/metrics-link-sync | `../tkt-worktrees/<type>-<slug>` |

1. Create the worktree from the updated base. `--track` sets the upstream, so a later push needs no `-u`:

   ```
   git fetch origin
   git worktree add --track -b <type>/<slug> ../tkt-worktrees/<type>-<slug> origin/main
   cd ../tkt-worktrees/<type>-<slug>
   ```

   When an issue exists, use `-b <type>/<n>-<slug> ../tkt-worktrees/issue-<n>-<slug>` instead.

2. Create the database directory. The store does not create it, and the server fails with `unable to open database file (14)` without it:

   ```
   mkdir -p data
   ```

3. Run with this worktree's own port and database, so two worktrees never fight over `:8080` or share one database file:

   ```
   TKT_LISTEN=127.0.0.1:8081 TKT_DB_PATH=data/tkt.db go run ./cmd/server
   ```

   Playwright needs none of that: it picks a free port and a temporary database on its own.

4. Work, test and commit only inside that worktree; the local checks in `CONTRIBUTING.md` mirror the workflows.

5. After the PR is merged, remove the worktree from the main checkout with `git worktree remove ../tkt-worktrees/<dir>`, only with authorization.

## Pull Requests

- The PR body follows `.github/PULL_REQUEST_TEMPLATE.md`: Summary, Scope, Changes, Verification and Related issue. Describe observed verification, not intended checks; anything not tried is stated as not verified.
- When an issue exists, link it with a standalone `Closes #<n>` line. No check requires the line; a missing issue or link never blocks a PR.
- Labels on the PR are optional. Do not add a `type:*` label to satisfy a rule; when an issue exists, its labels carry the classification.
- Pending or failing required checks block merge. Check definitions live under `.github/workflows/`; point at them rather than duplicating their flags, and do not claim a check that has no workflow. Which checks are merge-blocking is set by repository protection on `main`; the workflow files alone do not establish it.
- With explicit merge authorization, use `gh pr merge --squash` and no other method.

### Required checks

Branch protection on `main` requires one check: `Go quality`, from `.github/workflows/quality.yml`. That job passes only when every Go job in the workflow passes: Go lint (gofmt, go mod tidy, vet, build) and the sharded Go race jobs.

### Other workflows

These run on pull requests and their failures are fixed or explained, but they do not block merge:

- `.github/workflows/coverage.yml`: minimum coverage threshold
- `.github/workflows/security.yml`: staticcheck and govulncheck
- `.github/workflows/e2e.yml`: frontend formatting and the Playwright suite
- `.github/workflows/container.yml`: image build and smoke test

No check validates issue linkage.

### What Proves A Pull Request Here

- `.github/PULL_REQUEST_TEMPLATE.md`: body contract
- `.github/workflows/`: the checks that actually run
- `CONTRIBUTING.md`: the Go test layers and the local equivalents of every check
- `e2e/README.md`: the browser suite and its boundary

`.github/workflows/release-container.yml` runs on pushed `vX.Y.Z` tags, not on PRs; it is not a missing PR gate. Do not report another repository's checks as missing here. Do not call a local command run proof that remote required checks passed.

### Commit Format

Use Conventional Commits with a scope, matching the history: `chore(skills): ...` and `feat(tickets): ...` are the usual shapes. Keep commits focused. Pull requests are squash-merged, so the PR title becomes the commit message on `main`: write it as a valid conventional commit and do not rely on the individual branch commits surviving.

## Decision Gates

| Condition | Result |
| --- | --- |
| No issue exists for the change | GO: the issue is optional; suggest one only for large or design-changing work |
| Issue exists | GO: no approval wait |
| Issue has no `status:approved` | GO: approval is not a gate in this repository |
| Issue lacks a `type:*` label or has more than one | GO: apply or propose the best-fit single type; never block over labels |
| Issue has zero areas or more than two | GO: apply or propose the best fit; never block over labels |
| Label is ambiguous or outside the taxonomy | GO: note or propose a correction; do not create labels |
| An issue the PR does not link | GO: suggest adding `Closes #<n>` |
| Branch name or prefix deviates | GO: report the deviation in `branch_valid` |
| A required check fails or is pending | BLOCKED for merge |
| Required checks pass | GO for readiness, not authorization to publish or merge |
| The release workflow does not run on the PR | Not a missing PR gate |

## Output Contract

Return:

- `issue`: number/URL or "none"; `labels_valid`: true/false, with the exact taxonomy deviation if false, only when an issue exists;
- `worktree`: path, branch and base commit, plus the TKT_LISTEN and TKT_DB_PATH pair in use; `branch_valid`: true/false, with the naming deviation if any;
- `pr`: number/URL or "none"; `issue_linked`: true/false/"no issue", with the exact closing line if any (informational);
- `gate`: GO or BLOCKED, naming failing or pending required checks; a missing issue, missing approval or imperfect label never makes it BLOCKED;
- verification performed, anything not verified, and proposed corrections or GitHub writes not performed because authorization was absent.

## References

Paths below are repository-relative.

- `.github/ISSUE_TEMPLATE/`: bug report and feature request forms
- `.github/PULL_REQUEST_TEMPLATE.md`: body contract
- `.github/workflows/`: actual check and release definitions
- `.github/release.yml`: changelog grouping by PR labels
- `CONTRIBUTING.md`: local development and the local equivalents of every check
- `cmd/server/main.go`: the TKT_LISTEN and TKT_DB_PATH defaults
- `internal/adapters/sqlite/sqlite.go`: the store Open that leaves the database directory to the operator
- GitHub label inventory: `gh label list --repo heohlee/tkt`
