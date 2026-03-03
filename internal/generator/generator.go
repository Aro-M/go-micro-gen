package generator

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Aro-M/go-micro-gen/internal/config"
)

// templateFS is set by the main package (via SetTemplateFS) using go:embed.
// This indirection is needed because go:embed cannot reference parent directories.
var templateFS fs.FS

// SetTemplateFS injects the embedded template filesystem.
// Must be called before Generate().
func SetTemplateFS(f fs.FS) {
	templateFS = f
}

// Generator orchestrates the full service generation.
type Generator struct {
	cfg *config.ServiceConfig
}

// New creates a new Generator for the given config.
func New(cfg *config.ServiceConfig) *Generator {
	if cfg.Arch.Service == "" {
		cfg.Arch = config.GetArchFolders(cfg.Architecture)
	}
	return &Generator{cfg: cfg}
}

// Generate runs the full generation pipeline.
func (g *Generator) Generate() error {
	// 1. Create output directory
	if err := os.MkdirAll(g.cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	// 2. Walk all templates and render applicable ones
	// Note: templateFS is already sub'd to the templates root, so we walk "."
	return fs.WalkDir(templateFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Determine if this template applies to the current config
		if !g.shouldInclude(path) {
			return nil
		}

		return g.renderTemplate(path)
	})
}

// shouldInclude decides whether a template file should be rendered
// based on the current ServiceConfig (DB type, CI, Redis, etc.)
func (g *Generator) shouldInclude(tmplPath string) bool {
	// Database-specific templates
	if (strings.Contains(tmplPath, "/repository/postgres/") || strings.Contains(tmplPath, "/db/migrations/")) && g.cfg.Database != config.DBPostgres {
		return false
	}
	if strings.Contains(tmplPath, "/repository/mongo/") && g.cfg.Database != config.DBMongo {
		return false
	}

	// Message Broker templates
	if strings.Contains(tmplPath, "/broker/kafka/") && g.cfg.Broker != config.BrokerKafka {
		return false
	}
	if strings.Contains(tmplPath, "/broker/rabbitmq/") && g.cfg.Broker != config.BrokerRabbitMQ {
		return false
	}
	if strings.Contains(tmplPath, "/broker/nats/") && g.cfg.Broker != config.BrokerNATS {
		return false
	}

	// Cloud-specific templates
	if strings.Contains(tmplPath, "config/aws.go") && g.cfg.Cloud != config.CloudAWS {
		return false
	}
	if strings.Contains(tmplPath, "config/gcp.go") && g.cfg.Cloud != config.CloudGCP {
		return false
	}

	// Transport templates
	if strings.Contains(tmplPath, "/transport/grpc/") && g.cfg.Transport == config.TransportHTTP {
		return false
	}
	if strings.Contains(tmplPath, "/transport/httpx/") && g.cfg.Transport == config.TransportGRPC {
		return false
	}

	// Infra templates (Docker, K8s, Helm)
	if strings.Contains(tmplPath, "docker/") && !g.cfg.IncludeDocker {
		return false
	}
	if strings.Contains(tmplPath, "k8s/") && !g.cfg.IncludeK8s {
		return false
	}
	if strings.Contains(tmplPath, "helm/") && !g.cfg.IncludeHelm {
		return false
	}

	// CI/CD templates
	if strings.Contains(tmplPath, "github-actions") && g.cfg.CI != config.CIGitHub {
		return false
	}
	if strings.Contains(tmplPath, "gitlab-ci") && g.cfg.CI != config.CIGitLab {
		return false
	}

	return true
}

// renderTemplate reads, parses and executes a single .tmpl file,
// writing the result to the correct output path.
func (g *Generator) renderTemplate(tmplPath string) error {
	// Skip go.mod if we are running `init` directly inside an existing project that already has it.
	if strings.HasSuffix(tmplPath, "go.mod.tmpl") && g.cfg.OutputDir == "." {
		if _, err := os.Stat("go.mod"); err == nil {
			fmt.Printf("  - skipping go.mod (already exists)\n")
			return nil
		}
	}

	// Use the fs.ReadFile function (not method) since templateFS is fs.FS
	content, err := fs.ReadFile(templateFS, tmplPath)
	if err != nil {
		return fmt.Errorf("read template %s: %w", tmplPath, err)
	}

	// Parse template
	tmpl, err := template.New(filepath.Base(tmplPath)).
		Funcs(templateFuncs()).
		Parse(string(content))
	if err != nil {
		return fmt.Errorf("parse template %s: %w", tmplPath, err)
	}

	// Compute output path:
	// templates/service/internal/foo.go.tmpl → <outputDir>/internal/foo.go
	// templates/docker/Dockerfile.tmpl       → <outputDir>/docker/Dockerfile
	// templates/ci/github-actions.yml.tmpl   → <outputDir>/.github/workflows/ci.yml
	outRel := g.resolveOutputPath(tmplPath)
	outAbs := filepath.Join(g.cfg.OutputDir, outRel)

	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(outAbs), 0755); err != nil {
		return fmt.Errorf("create dir for %s: %w", outAbs, err)
	}

	out, err := os.Create(outAbs)
	if err != nil {
		return fmt.Errorf("create file %s: %w", outAbs, err)
	}
	defer out.Close()

	if err := tmpl.Execute(out, g.cfg); err != nil {
		return fmt.Errorf("execute template %s: %w", tmplPath, err)
	}

	fmt.Printf("  ✔ %s\n", outRel)
	return nil
}

// resolveOutputPath maps a template path to its output path.
// Path is relative to the templates root (already fs.Sub'd), so no "templates/" prefix.
func (g *Generator) resolveOutputPath(tmplPath string) string {
	rel := tmplPath

	// Strip the "service/" prefix — service files go directly into the root
	rel = strings.TrimPrefix(rel, "service/")

	// Special cases
	switch {
	case strings.HasPrefix(rel, "ci/github-actions"):
		rel = ".github/workflows/ci.yml"
	case strings.HasPrefix(rel, "ci/gitlab-ci"):
		rel = ".gitlab-ci.yml"
	case strings.HasPrefix(rel, "docker/"):
		// keep as-is: docker/Dockerfile, docker/docker-compose.yml, docker/configs/...
	}

	// Remove .tmpl extension
	rel = strings.TrimSuffix(rel, ".tmpl")

	// Apply architectural folder mapping
	rel = strings.Replace(rel, "internal/service", g.cfg.Arch.Service, 1)
	rel = strings.Replace(rel, "internal/repository", g.cfg.Arch.Repository, 1)
	rel = strings.Replace(rel, "internal/transport", g.cfg.Arch.Transport, 1)
	rel = strings.Replace(rel, "internal/domain", g.cfg.Arch.Domain, 1)

	return rel
}

// templateFuncs returns custom template functions.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"title": strings.Title, //nolint:staticcheck
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"replace": func(s, old, new string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"isPostgres": func(db config.DBType) bool { return db == config.DBPostgres },
		"isMongo":    func(db config.DBType) bool { return db == config.DBMongo },
		"isNoDB":     func(db config.DBType) bool { return db == config.DBNone },
		"hasDB":      func(db config.DBType) bool { return db != config.DBNone },
		"isKafka":    func(b config.BrokerType) bool { return b == config.BrokerKafka },
		"isRabbitMQ": func(b config.BrokerType) bool { return b == config.BrokerRabbitMQ },
		"isNATS":     func(b config.BrokerType) bool { return b == config.BrokerNATS },
		"hasBroker":  func(b config.BrokerType) bool { return b != config.BrokerNone && b != "" },
		"isHTTP":     func(t config.TransportType) bool { return t == config.TransportHTTP || t == config.TransportBoth },
		"isGRPC":     func(t config.TransportType) bool { return t == config.TransportGRPC || t == config.TransportBoth },
	}
}
