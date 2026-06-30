package skillkit

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

func TestAgentfieldSkillBrandingAndCompatibilitySurfaces(t *testing.T) {
	repoRoot := skillkitRepoRoot(t)

	t.Run("canonical skill prose uses Silmari", func(t *testing.T) {
		skillPath := filepath.Join(repoRoot, "skills", "agentfield", "SKILL.md")
		skillText := readTextFile(t, skillPath)

		for _, needle := range []string{
			"# Silmari",
			"Design and ship a multi-agent system on Silmari.",
			"runnable Silmari project",
		} {
			if !strings.Contains(skillText, needle) {
				t.Fatalf("%s missing %q", skillPath, needle)
			}
		}
		for _, legacy := range []string{"# AgentField", "on AgentField", "AgentField gives you"} {
			if strings.Contains(skillText, legacy) {
				t.Fatalf("%s still contains legacy product guidance %q", skillPath, legacy)
			}
		}
	})

	t.Run("command help uses Silmari", func(t *testing.T) {
		commandPath := filepath.Join(repoRoot, "skills", "agentfield", "commands", "agentfield.md")
		commandText := readTextFile(t, commandPath)

		for _, needle := range []string{
			"Design and ship a multi-agent system on Silmari",
			"multi-agent system on Silmari",
		} {
			if !strings.Contains(commandText, needle) {
				t.Fatalf("%s missing %q", commandPath, needle)
			}
		}
		if strings.Contains(commandText, "AgentField") {
			t.Fatalf("%s still contains legacy brand text", commandPath)
		}
	})

	t.Run("reference docs use Silmari without legacy product prose", func(t *testing.T) {
		for _, rel := range []string{
			filepath.Join("references", "anti-patterns.md"),
			filepath.Join("references", "capability-playbook.md"),
			filepath.Join("references", "examples-map.md"),
			filepath.Join("references", "live-docs.md"),
			filepath.Join("references", "memory-events.md"),
			filepath.Join("references", "model-selection.md"),
			filepath.Join("references", "patterns-emerge.md"),
			filepath.Join("references", "project-claude-template.md"),
			filepath.Join("references", "scaffold-recipe.md"),
			filepath.Join("references", "verification.md"),
		} {
			path := filepath.Join(repoRoot, "skills", "agentfield", rel)
			text := readTextFile(t, path)
			if !strings.Contains(text, "Silmari") {
				t.Fatalf("%s missing Silmari-visible guidance", path)
			}
			if strings.Contains(text, "AgentField") {
				t.Fatalf("%s still contains legacy product guidance", path)
			}
		}
	})

	t.Run("canonical skill markdown has no legacy AgentField product name", func(t *testing.T) {
		sourceRoot := filepath.Join(repoRoot, "skills", "agentfield")
		var matches []string

		err := filepath.WalkDir(sourceRoot, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() || filepath.Ext(path) != ".md" {
				return nil
			}

			text := readTextFile(t, path)
			if strings.Contains(text, "AgentField") {
				rel, err := filepath.Rel(sourceRoot, path)
				if err != nil {
					return err
				}
				matches = append(matches, filepath.ToSlash(rel))
			}
			return nil
		})
		if err != nil {
			t.Fatalf("walk %s: %v", sourceRoot, err)
		}
		if len(matches) > 0 {
			sort.Strings(matches)
			t.Fatalf("legacy AgentField product text remains in %v", matches)
		}
	})

	t.Run("compatibility surfaces stay unchanged and manifested", func(t *testing.T) {
		skillPath := filepath.Join(repoRoot, "skills", "agentfield", "SKILL.md")
		commandPath := filepath.Join(repoRoot, "skills", "agentfield", "commands", "agentfield.md")
		manifestPath := filepath.Join(repoRoot, "docs", "silmari-rebrand-manifest.md")

		skillText := readTextFile(t, skillPath)
		commandText := readTextFile(t, commandPath)
		manifestText := readTextFile(t, manifestPath)

		for _, needle := range []string{
			"name: agentfield",
			"aliases: [agentfield-multi-reasoner-builder]",
			"`agentfield` skill",
		} {
			if !strings.Contains(skillText, needle) && !strings.Contains(commandText, needle) {
				t.Fatalf("expected compatibility surface %q in skill docs", needle)
			}
		}
		if filepath.Base(commandPath) != "agentfield.md" {
			t.Fatalf("command file = %q want agentfield.md", filepath.Base(commandPath))
		}

		for _, manifestRow := range []string{
			"| skills/agentfield/SKILL.md | name: agentfield | skill-or-repo-slug |",
			"| skills/agentfield/SKILL.md | aliases: [agentfield-multi-reasoner-builder] | skill-or-repo-slug |",
			"| skills/agentfield/SKILL.md | agentfield.ai | published-link-target |",
			"| skills/agentfield/commands/agentfield.md | /agentfield | cli-command |",
			"| skills/agentfield/commands/agentfield.md | https://agentfield.ai/llms.txt | published-link-target |",
			"| skills/agentfield/references/live-docs.md | agentfield.ai | published-link-target |",
			"| skills/agentfield/references/live-docs.md | AGENTFIELD_HOME | env-var |",
			"| skills/agentfield/references/primitives-snapshot.md | agentfield.ai | published-link-target |",
			"| skills/agentfield/references/primitives-snapshot.md | AGENTFIELD_SERVER | env-var |",
			"| skills/agentfield/references/primitives-snapshot.md | sdk/python/agentfield/agent.py | import-module-path |",
			"| skills/agentfield/references/project-claude-template.md | agentfield/control-plane:latest | docker-image-or-volume |",
			"| skills/agentfield/references/scaffold-recipe.md | from agentfield import | import-module-path |",
			"| skills/agentfield/references/scaffold-recipe.md | AGENTFIELD_STORAGE_MODE | env-var |",
			"| skills/agentfield/references/scaffold-recipe.md | AGENTFIELD_HTTP_ADDR | env-var |",
			"| skills/agentfield/references/scaffold-recipe.md | AGENTFIELD_HTTP_PORT | env-var |",
			"| skills/agentfield/references/scaffold-recipe.md | AGENTFIELD_SERVER | env-var |",
			"| skills/agentfield/references/scaffold-recipe.md | agentfield/control-plane:latest | docker-image-or-volume |",
			"| skills/agentfield/references/scaffold-recipe.md | agentfield-data | docker-image-or-volume |",
			"| skills/agentfield/references/scaffold-recipe.md | agentfield | package-name |",
			"| skills/agentfield/references/triggers.md | agentfield.ai/docs | published-link-target |",
			"| skills/agentfield/references/triggers.md | from agentfield import | import-module-path |",
		} {
			if !strings.Contains(manifestText, manifestRow) {
				t.Fatalf("%s missing manifest row %q", manifestPath, manifestRow)
			}
		}
	})
}

func TestAgentfieldEmbeddedMirrorMatchesSource(t *testing.T) {
	repoRoot := skillkitRepoRoot(t)
	sourceRoot := filepath.Join(repoRoot, "skills", "agentfield")
	mirrorRoot := filepath.Join(repoRoot, "control-plane", "internal", "skillkit", "skill_data", "agentfield")

	sourceTree := snapshotTree(t, sourceRoot)
	mirrorTree := snapshotTree(t, mirrorRoot)
	if len(sourceTree) == 0 {
		t.Fatalf("canonical skill tree %s is empty", sourceRoot)
	}
	assertTreesEqual(t, sourceRoot, sourceTree, mirrorRoot, mirrorTree)
}

func TestSyncEmbeddedSkillsScriptScenarios(t *testing.T) {
	repoRoot := skillkitRepoRoot(t)

	t.Run("empty mirror drift", func(t *testing.T) {
		fixtureRoot := newSkillMirrorFixture(t, repoRoot)
		exitCode, output := runSyncScript(t, fixtureRoot, "--check")
		if exitCode == 0 {
			t.Fatalf("sync --check unexpectedly passed for empty mirror: %s", output)
		}
		for _, needle := range []string{
			"DRIFT: agentfield",
			"Run ./scripts/sync-embedded-skills.sh to fix the drift, then commit.",
		} {
			if !strings.Contains(output, needle) {
				t.Fatalf("sync --check output missing %q:\n%s", needle, output)
			}
		}
	})

	t.Run("single changed skill doc", func(t *testing.T) {
		fixtureRoot := newSkillMirrorFixture(t, repoRoot)
		seedMirror(t, fixtureRoot)

		mirrorSkillPath := filepath.Join(fixtureRoot, "control-plane", "internal", "skillkit", "skill_data", "agentfield", "SKILL.md")
		if err := os.WriteFile(mirrorSkillPath, []byte("legacy AgentField doc\n"), 0o644); err != nil {
			t.Fatalf("write mirror skill: %v", err)
		}

		exitCode, output := runSyncScript(t, fixtureRoot, "--check")
		if exitCode == 0 || !strings.Contains(output, "DRIFT: agentfield") {
			t.Fatalf("expected drift for single changed skill doc, exit=%d output=%s", exitCode, output)
		}

		exitCode, output = runSyncScript(t, fixtureRoot)
		if exitCode != 0 || !strings.Contains(output, "synced agentfield") {
			t.Fatalf("sync failed after single changed skill doc, exit=%d output=%s", exitCode, output)
		}

		exitCode, output = runSyncScript(t, fixtureRoot, "--check")
		if exitCode != 0 {
			t.Fatalf("sync --check failed after repair, exit=%d output=%s", exitCode, output)
		}

		assertTreesEqual(
			t,
			filepath.Join(fixtureRoot, "skills", "agentfield"),
			snapshotTree(t, filepath.Join(fixtureRoot, "skills", "agentfield")),
			filepath.Join(fixtureRoot, "control-plane", "internal", "skillkit", "skill_data", "agentfield"),
			snapshotTree(t, filepath.Join(fixtureRoot, "control-plane", "internal", "skillkit", "skill_data", "agentfield")),
		)
	})

	t.Run("many reference docs", func(t *testing.T) {
		fixtureRoot := newSkillMirrorFixture(t, repoRoot)
		seedMirror(t, fixtureRoot)

		for rel, contents := range map[string]string{
			filepath.Join("references", "anti-patterns.md"):       "legacy AgentField anti-pattern\n",
			filepath.Join("references", "capability-playbook.md"): "legacy AgentField capability playbook\n",
			filepath.Join("references", "verification.md"):        "legacy AgentField verification\n",
		} {
			path := filepath.Join(fixtureRoot, "control-plane", "internal", "skillkit", "skill_data", "agentfield", rel)
			if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
				t.Fatalf("write %s: %v", rel, err)
			}
		}

		exitCode, output := runSyncScript(t, fixtureRoot, "--check")
		if exitCode == 0 || !strings.Contains(output, "DRIFT: agentfield") {
			t.Fatalf("expected drift for many reference docs, exit=%d output=%s", exitCode, output)
		}

		exitCode, output = runSyncScript(t, fixtureRoot)
		if exitCode != 0 {
			t.Fatalf("sync failed for many reference docs, exit=%d output=%s", exitCode, output)
		}

		exitCode, output = runSyncScript(t, fixtureRoot, "--check")
		if exitCode != 0 {
			t.Fatalf("sync --check failed after repairing many reference docs, exit=%d output=%s", exitCode, output)
		}

		assertTreesEqual(
			t,
			filepath.Join(fixtureRoot, "skills", "agentfield"),
			snapshotTree(t, filepath.Join(fixtureRoot, "skills", "agentfield")),
			filepath.Join(fixtureRoot, "control-plane", "internal", "skillkit", "skill_data", "agentfield"),
			snapshotTree(t, filepath.Join(fixtureRoot, "control-plane", "internal", "skillkit", "skill_data", "agentfield")),
		)
	})

	t.Run("sync is idempotent", func(t *testing.T) {
		fixtureRoot := newSkillMirrorFixture(t, repoRoot)

		exitCode, output := runSyncScript(t, fixtureRoot)
		if exitCode != 0 {
			t.Fatalf("initial sync failed, exit=%d output=%s", exitCode, output)
		}
		firstSnapshot := snapshotTree(t, filepath.Join(fixtureRoot, "control-plane", "internal", "skillkit", "skill_data", "agentfield"))

		exitCode, output = runSyncScript(t, fixtureRoot, "--check")
		if exitCode != 0 {
			t.Fatalf("first sync --check failed, exit=%d output=%s", exitCode, output)
		}

		exitCode, output = runSyncScript(t, fixtureRoot)
		if exitCode != 0 {
			t.Fatalf("second sync failed, exit=%d output=%s", exitCode, output)
		}
		secondSnapshot := snapshotTree(t, filepath.Join(fixtureRoot, "control-plane", "internal", "skillkit", "skill_data", "agentfield"))

		exitCode, output = runSyncScript(t, fixtureRoot, "--check")
		if exitCode != 0 {
			t.Fatalf("second sync --check failed, exit=%d output=%s", exitCode, output)
		}

		assertTreesEqual(
			t,
			filepath.Join(fixtureRoot, "control-plane", "internal", "skillkit", "skill_data", "agentfield"),
			firstSnapshot,
			filepath.Join(fixtureRoot, "control-plane", "internal", "skillkit", "skill_data", "agentfield"),
			secondSnapshot,
		)
	})

	t.Run("extra mirror file is reported as drift and removed by sync", func(t *testing.T) {
		fixtureRoot := newSkillMirrorFixture(t, repoRoot)
		seedMirror(t, fixtureRoot)

		extraPath := filepath.Join(
			fixtureRoot,
			"control-plane",
			"internal",
			"skillkit",
			"skill_data",
			"agentfield",
			"references",
			"unexpected.md",
		)
		if err := os.WriteFile(extraPath, []byte("stale mirror file\n"), 0o644); err != nil {
			t.Fatalf("write extra mirror file: %v", err)
		}

		exitCode, output := runSyncScript(t, fixtureRoot, "--check")
		if exitCode == 0 || !strings.Contains(output, "DRIFT: agentfield") {
			t.Fatalf("expected drift for extra mirror file, exit=%d output=%s", exitCode, output)
		}

		exitCode, output = runSyncScript(t, fixtureRoot)
		if exitCode != 0 {
			t.Fatalf("sync failed for extra mirror file, exit=%d output=%s", exitCode, output)
		}

		if _, err := os.Stat(extraPath); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("extra mirror file still present after sync: err=%v", err)
		}

		exitCode, output = runSyncScript(t, fixtureRoot, "--check")
		if exitCode != 0 {
			t.Fatalf("sync --check failed after removing extra mirror file, exit=%d output=%s", exitCode, output)
		}

		assertTreesEqual(
			t,
			filepath.Join(fixtureRoot, "skills", "agentfield"),
			snapshotTree(t, filepath.Join(fixtureRoot, "skills", "agentfield")),
			filepath.Join(fixtureRoot, "control-plane", "internal", "skillkit", "skill_data", "agentfield"),
			snapshotTree(t, filepath.Join(fixtureRoot, "control-plane", "internal", "skillkit", "skill_data", "agentfield")),
		)
	})

	t.Run("missing source skill directory returns an error", func(t *testing.T) {
		fixtureRoot := newSkillMirrorFixture(t, repoRoot)
		if err := os.RemoveAll(filepath.Join(fixtureRoot, "skills", "agentfield")); err != nil {
			t.Fatalf("remove canonical skill: %v", err)
		}

		exitCode, output := runSyncScript(t, fixtureRoot, "--check")
		if exitCode == 0 {
			t.Fatalf("sync --check unexpectedly passed with missing source skill directory: %s", output)
		}
		if !strings.Contains(output, "ERROR: skill source") || !strings.Contains(output, "does not exist") {
			t.Fatalf("unexpected missing-source output:\n%s", output)
		}
	})
}

func skillkitRepoRoot(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
}

func readTextFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}

func snapshotTree(t *testing.T, root string) map[string][]byte {
	t.Helper()

	tree := make(map[string][]byte)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		tree[filepath.ToSlash(rel)] = data
		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}
	return tree
}

func assertTreesEqual(t *testing.T, wantRoot string, want map[string][]byte, gotRoot string, got map[string][]byte) {
	t.Helper()

	var missing []string
	var mismatched []string
	for rel, wantData := range want {
		gotData, ok := got[rel]
		if !ok {
			missing = append(missing, rel)
			continue
		}
		if !bytes.Equal(wantData, gotData) {
			mismatched = append(mismatched, rel)
		}
	}

	var extra []string
	for rel := range got {
		if _, ok := want[rel]; !ok {
			extra = append(extra, rel)
		}
	}

	sort.Strings(missing)
	sort.Strings(mismatched)
	sort.Strings(extra)
	if len(missing) == 0 && len(mismatched) == 0 && len(extra) == 0 {
		return
	}

	t.Fatalf(
		"tree mismatch between %s and %s: missing=%v mismatched=%v extra=%v",
		wantRoot,
		gotRoot,
		missing,
		mismatched,
		extra,
	)
}

func newSkillMirrorFixture(t *testing.T, repoRoot string) string {
	t.Helper()

	fixtureRoot := t.TempDir()
	for _, rel := range []string{
		"scripts",
		"skills",
		filepath.Join("control-plane", "internal", "skillkit", "skill_data"),
	} {
		if err := os.MkdirAll(filepath.Join(fixtureRoot, rel), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", rel, err)
		}
	}

	copyFile(t, filepath.Join(repoRoot, "scripts", "sync-embedded-skills.sh"), filepath.Join(fixtureRoot, "scripts", "sync-embedded-skills.sh"), 0o755)
	copyDir(t, filepath.Join(repoRoot, "skills", "agentfield"), filepath.Join(fixtureRoot, "skills", "agentfield"))

	return fixtureRoot
}

func seedMirror(t *testing.T, fixtureRoot string) {
	t.Helper()
	copyDir(
		t,
		filepath.Join(fixtureRoot, "skills", "agentfield"),
		filepath.Join(fixtureRoot, "control-plane", "internal", "skillkit", "skill_data", "agentfield"),
	)
}

func copyDir(t *testing.T, src, dst string) {
	t.Helper()

	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		t.Fatalf("copy %s -> %s: %v", src, dst, err)
	}
}

func copyFile(t *testing.T, src, dst string, mode os.FileMode) {
	t.Helper()

	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read %s: %v", src, err)
	}
	if err := os.WriteFile(dst, data, mode); err != nil {
		t.Fatalf("write %s: %v", dst, err)
	}
}

func runSyncScript(t *testing.T, repoRoot string, args ...string) (int, string) {
	t.Helper()

	cmdArgs := append([]string{filepath.Join(repoRoot, "scripts", "sync-embedded-skills.sh")}, args...)
	cmd := exec.Command("bash", cmdArgs...)
	cmd.Dir = repoRoot
	output, err := cmd.CombinedOutput()
	if err == nil {
		return 0, string(output)
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode(), string(output)
	}
	t.Fatalf("run sync script: %v\n%s", err, output)
	return 0, ""
}
