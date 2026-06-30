package templates

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func testTemplateData() TemplateData {
	return TemplateData{
		ProjectName:       "silmari-starter",
		NodeID:            "agent-123",
		GoModule:          "example.com/silmari-starter",
		AuthorName:        "Taylor Example",
		AuthorEmail:       "taylor@example.com",
		CurrentYear:       2026,
		CreatedAt:         "2026-06-29 12:00:00 UTC",
		Language:          "python",
		ControlPlaneImage: "agentfield/control-plane:latest",
		ControlPlanePort:  8080,
		AgentPort:         8001,
		DefaultModel:      "openrouter/google/gemini-2.5-flash",
	}
}

func renderTemplate(t *testing.T, name string, data TemplateData) string {
	t.Helper()

	tmpl, err := GetTemplate(name)
	if err != nil {
		t.Fatalf("GetTemplate(%q) error = %v", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		t.Fatalf("Execute(%q) error = %v", name, err)
	}

	return buf.String()
}

func assertContainsAll(t *testing.T, text string, substrings ...string) {
	t.Helper()

	for _, substring := range substrings {
		if !strings.Contains(text, substring) {
			t.Fatalf("expected output to contain %q:\n%s", substring, text)
		}
	}
}

func assertNotContainsAny(t *testing.T, text string, substrings ...string) {
	t.Helper()

	for _, substring := range substrings {
		if strings.Contains(text, substring) {
			t.Fatalf("expected output not to contain %q:\n%s", substring, text)
		}
	}
}

func readSilmariRebrandManifest(t *testing.T) string {
	t.Helper()

	repoRoot := findRepoRoot(t)
	manifestPath := filepath.Join(repoRoot, "docs", "silmari-rebrand-manifest.md")

	content, err := os.ReadFile(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("branch-local manifest missing at %q", manifestPath)
		}
		t.Fatalf("os.ReadFile(%q) error = %v", manifestPath, err)
	}

	return string(content)
}

func findRepoRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}

	dir := wd
	for {
		gitPath := filepath.Join(dir, ".git")
		info, err := os.Stat(gitPath)
		if err == nil {
			if !info.IsDir() && !info.Mode().IsRegular() {
				t.Fatalf(".git marker at %q has unsupported mode %v", gitPath, info.Mode())
			}
			return dir
		}
		if !os.IsNotExist(err) {
			t.Fatalf("os.Stat(%q) error = %v", gitPath, err)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not locate repo root from working directory %q", wd)
		}
		dir = parent
	}
}

func TestGetTemplate(t *testing.T) {
	t.Parallel()

	data := TemplateData{
		NodeID: "agent-123",
	}

	tests := []struct {
		name        string
		template    string
		wantText    string
		wantErrText string
	}{
		{
			name:     "parses and executes python template",
			template: "python/main.py.tmpl",
			wantText: `"agent-123"`,
		},
		{
			name:     "parses and executes go template",
			template: "go/main.go.tmpl",
			wantText: `"agent-123"`,
		},
		{
			name:        "missing template returns error",
			template:    "python/does-not-exist.tmpl",
			wantErrText: "pattern matches no files",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpl, err := GetTemplate(tt.template)
			if tt.wantErrText != "" {
				if err == nil {
					t.Fatalf("GetTemplate(%q) error = nil, want %q", tt.template, tt.wantErrText)
				}
				if !strings.Contains(err.Error(), tt.wantErrText) {
					t.Fatalf("GetTemplate(%q) error = %q, want substring %q", tt.template, err.Error(), tt.wantErrText)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetTemplate(%q) error = %v", tt.template, err)
			}

			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, data); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if !strings.Contains(buf.String(), tt.wantText) {
				t.Fatalf("rendered template missing %q in output:\n%s", tt.wantText, buf.String())
			}
		})
	}
}

func TestGetTemplateFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		language    string
		want        map[string]string
		wantErrText string
	}{
		{
			name:     "python templates",
			language: "python",
			want: map[string]string{
				"python/.env.example.tmpl":     ".env.example",
				"python/.gitignore.tmpl":       ".gitignore",
				"python/README.md.tmpl":        "README.md",
				"python/main.py.tmpl":          "main.py",
				"python/reasoners.py.tmpl":     "reasoners.py",
				"python/requirements.txt.tmpl": "requirements.txt",
			},
		},
		{
			name:     "go templates",
			language: "go",
			want: map[string]string{
				"go/.env.example.tmpl": ".env.example",
				"go/.gitignore.tmpl":   ".gitignore",
				"go/README.md.tmpl":    "README.md",
				"go/go.mod.tmpl":       "go.mod",
				"go/main.go.tmpl":      "main.go",
				"go/reasoners.go.tmpl": "reasoners.go",
			},
		},
		{
			name:     "typescript templates",
			language: "typescript",
			want: map[string]string{
				"typescript/.env.example.tmpl":  ".env.example",
				"typescript/.gitignore.tmpl":    ".gitignore",
				"typescript/README.md.tmpl":     "README.md",
				"typescript/main.ts.tmpl":       "main.ts",
				"typescript/package.json.tmpl":  "package.json",
				"typescript/reasoners.ts.tmpl":  "reasoners.ts",
				"typescript/tsconfig.json.tmpl": "tsconfig.json",
			},
		},
		{
			name:        "unsupported language",
			language:    "ruby",
			wantErrText: "unsupported language: ruby",
		},
		{
			name:        "empty language",
			language:    "",
			wantErrText: "unsupported language:",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := GetTemplateFiles(tt.language)
			if tt.wantErrText != "" {
				if err == nil {
					t.Fatalf("GetTemplateFiles(%q) error = nil, want %q", tt.language, tt.wantErrText)
				}
				if !strings.Contains(err.Error(), tt.wantErrText) {
					t.Fatalf("GetTemplateFiles(%q) error = %q, want substring %q", tt.language, err.Error(), tt.wantErrText)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetTemplateFiles(%q) error = %v", tt.language, err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("GetTemplateFiles(%q) returned %d files, want %d: %#v", tt.language, len(got), len(tt.want), got)
			}
			for wantPath, wantDest := range tt.want {
				if got[wantPath] != wantDest {
					t.Fatalf("GetTemplateFiles(%q)[%q] = %q, want %q", tt.language, wantPath, got[wantPath], wantDest)
				}
			}
		})
	}
}

func TestGeneratedReadmesUseSilmariBranding(t *testing.T) {
	t.Parallel()

	t.Run("empty project name still uses Silmari in every README", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name      string
			template  string
			wantTexts []string
		}{
			{
				name:     "python README",
				template: "python/README.md.tmpl",
				wantTexts: []string{
					"A Silmari agent created with `af init`.",
					"Start the Silmari server",
					"register with Silmari.",
				},
			},
			{
				name:     "go README",
				template: "go/README.md.tmpl",
				wantTexts: []string{
					"A Silmari agent created with `af init`.",
					"Start the Silmari server",
					"register with Silmari.",
				},
			},
			{
				name:     "typescript README",
				template: "typescript/README.md.tmpl",
				wantTexts: []string{
					"A Silmari agent created with `af init`.",
					"Start the Silmari server",
					"registers with Silmari.",
				},
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				data := testTemplateData()
				data.ProjectName = ""
				rendered := renderTemplate(t, tt.template, data)

				assertContainsAll(t, rendered, tt.wantTexts...)
				assertNotContainsAny(t, rendered, "AgentField")
			})
		}
	})

	tests := []struct {
		name      string
		template  string
		wantTexts []string
	}{
		{
			name:     "python README",
			template: "python/README.md.tmpl",
			wantTexts: []string{
				"A Silmari agent created with `af init`.",
				"Start the Silmari server",
				"register with Silmari.",
				"[Silmari Documentation](https://agentfield.ai/docs/learn)",
				"[SDK Reference](https://agentfield.ai/docs/reference/sdks/python)",
				"http://localhost:8080/api/v1/execute/",
			},
		},
		{
			name:     "go README",
			template: "go/README.md.tmpl",
			wantTexts: []string{
				"A Silmari agent created with `af init`.",
				"Start the Silmari server",
				"register with Silmari.",
				"[Silmari Documentation](https://agentfield.ai/docs/learn)",
				"[SDK Reference](https://agentfield.ai/docs/reference/sdks/go)",
				"http://localhost:8080/api/v1/execute/",
			},
		},
		{
			name:     "typescript README",
			template: "typescript/README.md.tmpl",
			wantTexts: []string{
				"A Silmari agent created with `af init`.",
				"Start the Silmari server",
				"registers with Silmari.",
				"[Silmari Documentation](https://agentfield.ai/docs/learn)",
				"[SDK Reference](https://agentfield.ai/docs/reference/sdks/typescript)",
				"http://localhost:8080/api/v1/execute/",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rendered := renderTemplate(t, tt.template, testTemplateData())
			assertContainsAll(t, rendered, tt.wantTexts...)
			assertNotContainsAny(t, rendered, "AgentField")
		})
	}
}

func TestGeneratedSourceGuidanceUsesSilmari(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		template       string
		wantSubstrings []string
		avoid          []string
		data           TemplateData
	}{
		{
			name:     "python main docstring",
			template: "python/main.py.tmpl",
			wantSubstrings: []string{
				`"""silmari-starter - Silmari agent node.`,
				"per the Silmari scaffold recipe.",
			},
			avoid: []string{
				"AgentField agent node",
				"agentfield skill's scaffold-recipe",
			},
			data: testTemplateData(),
		},
		{
			name:     "python main docstring with empty project name",
			template: "python/main.py.tmpl",
			wantSubstrings: []string{
				"Silmari agent node.",
				"per the Silmari scaffold recipe.",
			},
			avoid: []string{
				"AgentField agent node",
				"agentfield skill's scaffold-recipe",
			},
			data: func() TemplateData {
				data := testTemplateData()
				data.ProjectName = ""
				return data
			}(),
		},
		{
			name:     "docker compose comment",
			template: "docker/docker-compose.yml.tmpl",
			wantSubstrings: []string{
				"Silmari uses the legacy-compatible image `agentfield/control-plane:latest`,",
			},
			avoid: []string{
				"NOTE: agentfield/control-plane:latest is a distroless image",
			},
			data: testTemplateData(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rendered := renderTemplate(t, tt.template, tt.data)
			assertContainsAll(t, rendered, tt.wantSubstrings...)
			assertNotContainsAny(t, rendered, tt.avoid...)
		})
	}
}

func TestReadTemplateContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		path        string
		wantText    string
		wantErrText string
	}{
		{
			name:     "reads embedded content",
			path:     "typescript/package.json.tmpl",
			wantText: `"@agentfield/sdk"`,
		},
		{
			name:        "missing content returns error",
			path:        "typescript/missing.tmpl",
			wantErrText: "file does not exist",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ReadTemplateContent(tt.path)
			if tt.wantErrText != "" {
				if err == nil {
					t.Fatalf("ReadTemplateContent(%q) error = nil, want %q", tt.path, tt.wantErrText)
				}
				if !strings.Contains(err.Error(), tt.wantErrText) {
					t.Fatalf("ReadTemplateContent(%q) error = %q, want substring %q", tt.path, err.Error(), tt.wantErrText)
				}
				return
			}

			if err != nil {
				t.Fatalf("ReadTemplateContent(%q) error = %v", tt.path, err)
			}
			if !strings.Contains(string(got), tt.wantText) {
				t.Fatalf("ReadTemplateContent(%q) missing %q in output:\n%s", tt.path, tt.wantText, string(got))
			}
		})
	}
}

func TestTemplateCompatibilityIdentifiersRemainStable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		path      string
		wantTexts []string
	}{
		{
			name: "python README compatibility identifiers",
			path: "python/README.md.tmpl",
			wantTexts: []string{
				"`af init`",
				"af server",
				"https://agentfield.ai/docs/learn",
				"https://agentfield.ai/docs/reference/sdks/python",
				"http://localhost:8080/api/v1/execute/",
			},
		},
		{
			name: "python main compatibility identifiers",
			path: "python/main.py.tmpl",
			wantTexts: []string{
				"from agentfield import Agent, AIConfig",
				`os.getenv("AGENTFIELD_SERVER", "http://localhost:8080")`,
			},
		},
		{
			name: "python reasoners import remains stable",
			path: "python/reasoners.py.tmpl",
			wantTexts: []string{
				"from agentfield import AgentRouter",
			},
		},
		{
			name: "python env example remains stable",
			path: "python/.env.example.tmpl",
			wantTexts: []string{
				"AGENTFIELD_CONTROL_PLANE_URL=http://localhost:8080",
			},
		},
		{
			name: "python requirements package remains stable",
			path: "python/requirements.txt.tmpl",
			wantTexts: []string{
				"agentfield",
			},
		},
		{
			name: "go README compatibility identifiers",
			path: "go/README.md.tmpl",
			wantTexts: []string{
				"`af init`",
				"af server",
				"https://agentfield.ai/docs/learn",
				"https://agentfield.ai/docs/reference/sdks/go",
				"http://localhost:8080/api/v1/execute/",
			},
		},
		{
			name: "go main compatibility identifiers",
			path: "go/main.go.tmpl",
			wantTexts: []string{
				"github.com/Agent-Field/agentfield/sdk/go/agent",
				"github.com/Agent-Field/agentfield/sdk/go/ai",
				"AgentFieldURL:",
			},
		},
		{
			name: "go module remains stable",
			path: "go/go.mod.tmpl",
			wantTexts: []string{
				"github.com/Agent-Field/agentfield/sdk/go",
			},
		},
		{
			name: "go reasoners imports remain stable",
			path: "go/reasoners.go.tmpl",
			wantTexts: []string{
				"github.com/Agent-Field/agentfield/sdk/go/agent",
				"github.com/Agent-Field/agentfield/sdk/go/ai",
			},
		},
		{
			name: "go env example remains stable",
			path: "go/.env.example.tmpl",
			wantTexts: []string{
				"AGENTFIELD_CONTROL_PLANE_URL=http://localhost:8080",
			},
		},
		{
			name: "typescript README compatibility identifiers",
			path: "typescript/README.md.tmpl",
			wantTexts: []string{
				"`af init`",
				"af server",
				"https://agentfield.ai/docs/learn",
				"https://agentfield.ai/docs/reference/sdks/typescript",
				"http://localhost:8080/api/v1/execute/",
			},
		},
		{
			name: "typescript main compatibility identifiers",
			path: "typescript/main.ts.tmpl",
			wantTexts: []string{
				"import { Agent } from '@agentfield/sdk';",
				"process.env.AGENTFIELD_URL",
			},
		},
		{
			name: "typescript package remains stable",
			path: "typescript/package.json.tmpl",
			wantTexts: []string{
				`"@agentfield/sdk": "^0.1.0"`,
			},
		},
		{
			name: "typescript reasoners import remains stable",
			path: "typescript/reasoners.ts.tmpl",
			wantTexts: []string{
				"import { AgentRouter } from '@agentfield/sdk';",
			},
		},
		{
			name: "typescript env example remains stable",
			path: "typescript/.env.example.tmpl",
			wantTexts: []string{
				"AGENTFIELD_URL=http://localhost:8080",
			},
		},
		{
			name: "docker compatibility identifiers",
			path: "docker/docker-compose.yml.tmpl",
			wantTexts: []string{
				"AGENTFIELD_STORAGE_MODE: local",
				"AGENTFIELD_HTTP_ADDR: 0.0.0.0:8080",
				"${AGENTFIELD_HTTP_PORT:-{{.ControlPlanePort}}}:8080",
				"AGENTFIELD_SERVER: http://control-plane:8080",
				"agentfield-data:/data",
				"agentfield/control-plane:latest",
			},
		},
		{
			name: "docker env example remains stable",
			path: "docker/.env.example.tmpl",
			wantTexts: []string{
				"AGENTFIELD_HTTP_PORT={{.ControlPlanePort}}",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			content, err := ReadTemplateContent(tt.path)
			if err != nil {
				t.Fatalf("ReadTemplateContent(%q) error = %v", tt.path, err)
			}
			assertContainsAll(t, string(content), tt.wantTexts...)
		})
	}
}

func TestTemplateManifestRecordsPreservedIdentifiers(t *testing.T) {
	t.Parallel()

	manifest := readSilmariRebrandManifest(t)

	assertContainsAll(t, manifest,
		"| control-plane/internal/templates/python/README.md.tmpl | af init | cli-command |",
		"| control-plane/internal/templates/python/main.py.tmpl | from agentfield import Agent, AIConfig | import-module-path |",
		"| control-plane/internal/templates/python/.env.example.tmpl | AGENTFIELD_CONTROL_PLANE_URL | env-var |",
		"| control-plane/internal/templates/python/requirements.txt.tmpl | agentfield | package-name |",
		"| control-plane/internal/templates/go/main.go.tmpl | github.com/Agent-Field/agentfield/sdk/go | go-module-path |",
		"| control-plane/internal/templates/typescript/package.json.tmpl | @agentfield/sdk | package-name |",
		"| control-plane/internal/templates/docker/docker-compose.yml.tmpl | agentfield/control-plane:latest | docker-image-or-volume |",
		"| control-plane/internal/templates/templates_test.go | AgentField | test-fixture |",
		"| control-plane/internal/templates/templates_test.go | https://agentfield.ai/docs/learn | published-link-target |",
		"| control-plane/internal/templates/templates_test.go | from agentfield import Agent, AIConfig | import-module-path |",
		"| control-plane/internal/templates/templates_test.go | AGENTFIELD_HTTP_PORT | env-var |",
	)
}

func TestTemplateManifestRecordsCompatibilityEdges(t *testing.T) {
	t.Parallel()

	manifest := readSilmariRebrandManifest(t)

	tests := []struct {
		name string
		row  string
	}{
		{
			name: "python README preserves af server command",
			row:  "| control-plane/internal/templates/python/README.md.tmpl | af server | cli-command |",
		},
		{
			name: "docker compose preserves volume name",
			row:  "| control-plane/internal/templates/docker/docker-compose.yml.tmpl | agentfield-data | docker-image-or-volume |",
		},
		{
			name: "python README preserves execute endpoint",
			row:  "| control-plane/internal/templates/python/README.md.tmpl | http://localhost:8080/api/v1/execute/ | api-path |",
		},
		{
			name: "go README preserves execute endpoint",
			row:  "| control-plane/internal/templates/go/README.md.tmpl | http://localhost:8080/api/v1/execute/ | api-path |",
		},
		{
			name: "typescript README preserves execute endpoint",
			row:  "| control-plane/internal/templates/typescript/README.md.tmpl | http://localhost:8080/api/v1/execute/ | api-path |",
		},
		{
			name: "template compatibility tests preserve execute endpoint",
			row:  "| control-plane/internal/templates/templates_test.go | http://localhost:8080/api/v1/execute/ | api-path |",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertContainsAll(t, manifest, tt.row)
		})
	}
}

func TestTemplateManifestUsesDirectGoTestVerification(t *testing.T) {
	t.Parallel()

	manifest := readSilmariRebrandManifest(t)

	assertContainsAll(t, manifest,
		"| control-plane/internal/templates/templates_test.go | rebranded-with-preserved-identifiers | go test ./internal/templates/... |",
		"| control-plane/internal/templates/python/README.md.tmpl | rebranded-with-preserved-identifiers | go test ./internal/templates/... |",
		"| control-plane/internal/templates/docker/docker-compose.yml.tmpl | rebranded-with-preserved-identifiers | go test ./internal/templates/... |",
		"| go test ./internal/templates/... | control-plane | 0 | Passed; generated README branding, compatibility identifiers, and template error assertions all succeeded. |",
	)
	assertNotContainsAny(t, manifest, "PATH=/tmp/go-toolchain/go/bin:$PATH go test ./internal/templates/...")
}

func TestTemplateDataCompatibilityFieldsRemainStable(t *testing.T) {
	t.Parallel()

	wantFields := []string{
		"ProjectName",
		"NodeID",
		"GoModule",
		"AuthorName",
		"AuthorEmail",
		"CurrentYear",
		"CreatedAt",
		"Language",
		"ControlPlaneImage",
		"ControlPlanePort",
		"AgentPort",
		"DefaultModel",
	}

	typ := reflect.TypeOf(TemplateData{})
	gotFields := make([]string, 0, typ.NumField())
	missingFields := make([]string, 0)
	for i := 0; i < typ.NumField(); i++ {
		gotFields = append(gotFields, typ.Field(i).Name)
	}

	for _, wantField := range wantFields {
		if _, ok := typ.FieldByName(wantField); !ok {
			missingFields = append(missingFields, wantField)
		}
	}

	if len(missingFields) > 0 {
		t.Fatalf("TemplateData missing required fields %v; got %v", missingFields, gotFields)
	}
}

func TestGetSupportedLanguages(t *testing.T) {
	t.Parallel()

	got := GetSupportedLanguages()
	want := []string{"python", "go", "typescript"}

	if !slices.Equal(got, want) {
		t.Fatalf("GetSupportedLanguages() = %v, want %v", got, want)
	}
}
