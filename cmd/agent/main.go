package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	embedui "github.com/pvc-explorer-operator/pvc-explorer-agent"
	"github.com/pvc-explorer-operator/pvc-explorer-agent/agent"
	"github.com/pvc-explorer-operator/pvc-explorer-agent/agent/pvcwatch"
)

var buildVersion = ""

func version() string {
	if buildVersion != "" {
		return buildVersion
	}
	exe, err := os.Executable()
	if err != nil {
		return "unknown"
	}
	info, err := os.Stat(exe)
	if err != nil {
		return "unknown"
	}
	return fmt.Sprintf("%d", info.ModTime().Unix())
}

type modeState struct {
	mu      sync.RWMutex
	forceRW bool
	watcher *pvcwatch.Watcher
}

func (s *modeState) isReadonly(ctx context.Context) bool {
	s.mu.RLock()
	forced := s.forceRW
	s.mu.RUnlock()
	if forced {
		return false
	}
	if s.watcher == nil {
		return false
	}
	inUse, err := s.watcher.PVCInUse(ctx)
	if err != nil {
		log.Printf("pvcwatch: %v — defaulting to readonly", err)
		return true
	}
	return inUse
}

func (s *modeState) setForceRW(v bool) {
	s.mu.Lock()
	s.forceRW = v
	s.mu.Unlock()
}

func main() {
	var root, pvcName, clusterName, uiOverlay string
	flag.StringVar(&root, "root", "/data", "Root directory to serve")
	flag.StringVar(&pvcName, "pvc", "", "PVC claim name to watch for conflicts")
	flag.StringVar(&clusterName, "cluster", "", "Cluster name shown in UI")
	flag.StringVar(&	uiOverlay, "ui-overlay", "/config", "Directory of UI files to serve instead of embedded defaults")
	flag.Parse()

	namespace := os.Getenv("POD_NAMESPACE")
	podName := os.Getenv("POD_NAME")
	authToken := os.Getenv("AUTH_TOKEN")

	mode := &modeState{}
	if pvcName != "" {
		w, err := pvcwatch.New(pvcName)
		if err != nil {
			log.Printf("pvcwatch init failed (%v) — running without PVC conflict detection", err)
		} else {
			mode.watcher = w
			if podName == "" {
				podName = w.PodName()
			}
		}
	}

	mux := http.NewServeMux()
	agent.RegisterRoutes(mux, root, func(r *http.Request) bool {
		return mode.isReadonly(r.Context())
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ro := mode.isReadonly(r.Context())
		mode.mu.RLock()
		forceRW := mode.forceRW
		mode.mu.RUnlock()
		json.NewEncoder(w).Encode(map[string]interface{}{
			"readonly":  ro,
			"forceRW":   forceRW,
			"pvcWatch":  mode.watcher != nil,
			"cluster":   clusterName,
			"namespace": namespace,
			"pvc":       pvcName,
			"pod":       podName,
			"version":   version(),
		})
	})

	mux.HandleFunc("/api/mode", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			ForceRW bool `json:"forceRW"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		mode.setForceRW(req.ForceRW)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"forceRW": req.ForceRW})
	})

	mux.Handle("/", embedui.Handler(uiOverlay))

	// Wrap the entire mux with optional token authentication.
	// When AUTH_TOKEN is set, every request must present a valid
	// Authorization: Bearer <token> header. Standalone (no token)
	// mode remains fully open.
	handler := agent.WithAuth(mux, authToken)

	addr := ":8081"
	if authToken != "" {
		fmt.Printf("PVC Exporter Agent listening on %s with token authentication\n", addr)
	} else {
		fmt.Printf("PVC Exporter Agent listening on %s (no authentication)\n", addr)
	}
	fmt.Printf("  root: %s, pvc: %s, cluster: %s\n", root, pvcName, clusterName)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
