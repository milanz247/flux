package framework

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"flux/config"
)

// Inertia implements the server half of the Inertia protocol: the first
// request renders the root HTML document with the page object embedded;
// every navigation after that (marked with the X-Inertia header) receives
// the page object as JSON and the Vue client swaps components in place.
type Inertia struct {
	cfg      *config.Config
	template *template.Template
	version  string
}

// page is the Inertia page object exchanged with the Vue client.
type page struct {
	Component string         `json:"component"`
	Props     map[string]any `json:"props"`
	URL       string         `json:"url"`
	Version   string         `json:"version"`
}

const rootTemplatePath = "resources/views/app.html"

// NewInertia loads the root template and computes the asset version (a hash
// of the Vite manifest, so deploys with new assets trigger a full reload).
func NewInertia(cfg *config.Config) (*Inertia, error) {
	tmpl, err := template.ParseFiles(rootTemplatePath)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", rootTemplatePath, err)
	}

	return &Inertia{
		cfg:      cfg,
		template: tmpl,
		version:  assetVersion(cfg),
	}, nil
}

// Render responds with either the Inertia JSON page object (client-side
// navigation) or the full HTML document (first load / hard refresh).
func (i *Inertia) Render(c *gin.Context, component string, props map[string]any) {
	p := page{
		Component: component,
		Props:     props,
		URL:       c.Request.URL.RequestURI(),
		Version:   i.version,
	}

	if c.GetHeader("X-Inertia") == "true" {
		// Stale asset version on a GET → tell the client to do a full visit.
		if c.Request.Method == http.MethodGet &&
			c.GetHeader("X-Inertia-Version") != i.version {
			c.Header("X-Inertia-Location", p.URL)
			c.Status(http.StatusConflict)
			return
		}
		c.Header("X-Inertia", "true")
		c.Header("Vary", "X-Inertia")
		c.JSON(http.StatusOK, p)
		return
	}

	pageJSON, err := json.Marshal(p)
	if err != nil {
		c.String(http.StatusInternalServerError, "flux: failed to encode page object")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)
	_ = i.template.Execute(c.Writer, map[string]any{
		"AppName": i.cfg.App.Name,
		"Page":    string(pageJSON),
		"Vite":    i.viteTags(),
	})
}

// viteTags returns the script/link tags for the frontend bundle: the Vite
// dev server in local development, the built manifest entries in production.
func (i *Inertia) viteTags() template.HTML {
	if i.cfg.App.Env == "local" {
		dev := i.cfg.Vite.DevServer
		return template.HTML(fmt.Sprintf(
			`<script type="module" src="%s/@vite/client"></script>`+"\n"+
				`<script type="module" src="%s/resources/js/app.ts"></script>`,
			dev, dev,
		))
	}

	manifest, err := readViteManifest(i.cfg)
	if err != nil {
		return template.HTML("<!-- flux: vite manifest missing; run `flux build` -->")
	}

	entry, ok := manifest["resources/js/app.ts"]
	if !ok {
		return template.HTML("<!-- flux: entry resources/js/app.ts not in manifest -->")
	}

	tags := ""
	for _, css := range entry.CSS {
		tags += fmt.Sprintf(`<link rel="stylesheet" href="/%s/%s">`+"\n", i.cfg.Vite.BuildDir, css)
	}
	tags += fmt.Sprintf(`<script type="module" src="/%s/%s"></script>`, i.cfg.Vite.BuildDir, entry.File)
	return template.HTML(tags)
}

type manifestEntry struct {
	File string   `json:"file"`
	CSS  []string `json:"css"`
}

func manifestPath(cfg *config.Config) string {
	return filepath.Join("public", cfg.Vite.BuildDir, ".vite", "manifest.json")
}

func readViteManifest(cfg *config.Config) (map[string]manifestEntry, error) {
	raw, err := os.ReadFile(manifestPath(cfg))
	if err != nil {
		return nil, err
	}
	var manifest map[string]manifestEntry
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return nil, err
	}
	return manifest, nil
}

// assetVersion hashes the Vite manifest so the Inertia client can detect
// stale assets after a deploy. In development the version is constant.
func assetVersion(cfg *config.Config) string {
	raw, err := os.ReadFile(manifestPath(cfg))
	if err != nil {
		return "dev"
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:8])
}

// structToMap converts a DTO struct (or map) into the props map Inertia
// sends to Vue, honouring the DTO's json tags.
func structToMap(v any) map[string]any {
	result := map[string]any{}
	if v == nil {
		return result
	}
	if m, ok := v.(map[string]any); ok {
		for k, val := range m {
			result[k] = val
		}
		return result
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return result
	}
	_ = json.Unmarshal(raw, &result)
	return result
}
