package ui

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Aro-M/go-micro-gen/internal/config"
	"github.com/Aro-M/go-micro-gen/internal/generator"
)

//go:embed assets/*
var content embed.FS

type GenerateRequest struct {
	Name       string `json:"name"`
	Module     string `json:"module"`
	Database   string `json:"db"`
	Broker     string `json:"broker"`
	Cloud      string `json:"cloud"`
	Serverless bool   `json:"serverless"`
	GraphQL    bool   `json:"graphql"`
	JWT        bool   `json:"jwt"`
	Seeding    bool   `json:"seeding"`
}

func StartServer(port string) error {
	// Strip the outer "assets" directory from the embed FS mappings
	assets, err := fs.Sub(content, "assets")
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	// Mount static file server
	mux.Handle("/", http.FileServer(http.FS(assets)))

	// Mount logic endpoint
	mux.HandleFunc("/api/generate", handleGenerate)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	return server.ListenAndServe()
}

func handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusBadRequest)
		return
	}

	// Prepare config bindings mapped backwards to CLI structures
	cfg := &config.ServiceConfig{
		ServiceName:       req.Name,
		ModulePath:        req.Module,
		Database:          config.DBType(req.Database),
		Broker:            config.BrokerType(req.Broker),
		Transport:         "http", // Stardard HTTP adapter default
		Architecture:      "standard",
		Cloud:             config.CloudProvider(req.Cloud),
		IncludeServerless: req.Serverless,
		IncludeGraphQL:    req.GraphQL,
		IncludeJWT:        req.JWT,
		IncludeSeeding:    req.Seeding,
		IncludeDocker:     true, // Web UI builds Dockerfiles by default
		IncludeK8s:        false,
		IncludeHelm:       false,
		CI:                "none",
	}

	cwd, err := os.Getwd()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "could not determine cwd: %v"}`, err), http.StatusInternalServerError)
		return
	}
	cfg.OutputDir = filepath.Join(cwd, req.Name)

	g := generator.New(cfg)
	if err := g.Generate(); err != nil {
		// Log to term and send API response
		fmt.Printf("❌ UI generation failed: %v\n", err)
		http.Error(w, fmt.Sprintf(`{"error": "Generation failed: %v"}`, err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("✅ Project '%s' successfully generated via UI into %s\n", req.Name, cfg.OutputDir)

	w.Header().Set("Content-Type", "application/json")
	metadata := map[string]string{
		"status":  "success",
		"message": "Microservice successfully placed into " + cfg.OutputDir,
	}
	json.NewEncoder(w).Encode(metadata)
}
