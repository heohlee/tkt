package references

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// TestDocumentedPathsExist fails when AGENTS.md or a project skill points at a
// path that does not exist: a renamed skill, a moved file, a deleted workflow,
// or a typo. It is the machine check that lets the skill index and the skills
// point at their sources instead of restating them, which is what issues #179
// and #182 kept deferring and what the native review kept reporting as an
// unverified reference.
//
// Blind spot, accepted and stated here on purpose: bare filenames with no slash
// (`fakes_test.go`, `go.mod`) and globs (`*_test.go`, `testdata/*.golden`) are
// skipped, because treating every bare name as a path would flag ordinary
// prose. A misspelled top-level directory inside a bare token stays undetected.
func TestDocumentedPathsExist(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}

	sources, err := Sources(root)
	if err != nil {
		t.Fatalf("sources: %v", err)
	}
	if len(sources) < 2 {
		t.Fatalf("expected the skill index and at least one skill, got %d sources", len(sources))
	}

	// Counting both sides keeps the test honest: without it, a regression in the
	// extraction rules would make this test pass while checking nothing.
	validated := 0
	skips := map[string]int{}

	for _, source := range sources {
		rel, err := filepath.Rel(root, source)
		if err != nil {
			rel = source
		}

		t.Run(rel, func(t *testing.T) {
			body, err := os.ReadFile(source)
			if err != nil {
				t.Fatalf("read %s: %v", rel, err)
			}

			for line := range strings.SplitSeq(string(body), "\n") {
				for _, match := range RefPattern.FindAllStringSubmatch(line, -1) {
					if reason := SkipReason(match[1]); reason != "" {
						skips[reason]++
					}
				}
			}

			refs := Refs(source, body)
			if len(refs) == 0 {
				t.Errorf("%s has no path references: either it lost its content or the extraction is broken", rel)
			}

			for _, ref := range refs {
				validated++
				if !Exists(root, ref) {
					t.Errorf("%s:%d: unresolved reference %q\nresolves to %s\nif this is prose and not a path, add it to NonPathRefs with a reason",
						rel, ref.Line, ref.Token, Resolve(root, ref))
				}
			}
		})
	}

	if validated == 0 {
		t.Fatal("no references were validated: the extraction is broken")
	}

	t.Logf("validated %d references across %d files; skipped %d tokens as non-paths", validated, len(sources), len(skips))
}

// TestSkillsAreSelfContained fails when a skill document names another skill
// or AGENTS.md. Each skill is loaded on its own trigger, often without the
// others, so a rule it needs from elsewhere must be stated inline rather than
// pointed at.
func TestSkillsAreSelfContained(t *testing.T) {
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}

	names, err := SkillNames(root)
	if err != nil {
		t.Fatalf("skill names: %v", err)
	}
	if len(names) < 2 {
		t.Fatalf("expected at least two skills to check against each other, got %v", names)
	}

	docs, err := SkillDocs(root)
	if err != nil {
		t.Fatalf("skill docs: %v", err)
	}
	if len(docs) < len(names) {
		t.Fatalf("expected at least one document per skill, got %d documents for %d skills", len(docs), len(names))
	}

	for _, doc := range docs {
		rel, err := filepath.Rel(root, doc)
		if err != nil {
			rel = doc
		}
		skillRel, err := filepath.Rel(SkillsDir(root), doc)
		if err != nil {
			t.Fatalf("locate %s: %v", rel, err)
		}
		own, _, _ := strings.Cut(filepath.ToSlash(skillRel), "/")

		t.Run(rel, func(t *testing.T) {
			body, err := os.ReadFile(doc)
			if err != nil {
				t.Fatalf("read %s: %v", rel, err)
			}
			for _, ref := range CrossRefs(own, body, names) {
				t.Errorf("%s:%d: names %s (%s)\nstate the needed fact inline instead of pointing at it", rel, ref.Line, ref.Target, ref.Text)
			}
		})
	}
}

// TestCrossRefsDetection pins what CrossRefs reports, so a loosened detector
// cannot make TestSkillsAreSelfContained pass while checking nothing.
func TestCrossRefsDetection(t *testing.T) {
	names := []string{"go-tests", "e2e-playwright", "design-system"}
	cases := []struct {
		name string
		line string
		want []string
	}{
		{"relative path", "- `../go-tests/SKILL.md` — Go checks", []string{"go-tests"}},
		{"nested relative path", "see ../../e2e-playwright/SKILL.md", []string{"e2e-playwright"}},
		{"repository path", "- `.agents/skills/go-tests/SKILL.md`", []string{"go-tests"}},
		{"backticked name", "(that is `e2e-playwright`)", []string{"e2e-playwright"}},
		{"skill index", "- `AGENTS.md`: skill registration", []string{"AGENTS.md"}},
		{"two skills on one line", "`go-tests` or `e2e-playwright`", []string{"go-tests", "e2e-playwright"}},
		{"own skill", "- `design-system` and `../design-system/SKILL.md`", nil},
		{"plain prose", "Playwright and Go tests and the design system", nil},
		{"unrelated path", "- `.github/workflows/e2e.yml` and `internal/adapters/http/golden_test.go`", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got []string
			for _, ref := range CrossRefs("design-system", []byte(tc.line), names) {
				got = append(got, ref.Target)
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("CrossRefs(%q) = %v, want %v", tc.line, got, tc.want)
			}
		})
	}
}
