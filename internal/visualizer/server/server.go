package server

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/DanielFasel/sporecaster/internal/spore"
)

//go:embed assets
var assets embed.FS

func Handler(s *spore.Spore) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/spore", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s)
	})

	static, _ := fs.Sub(assets, "assets")
	mux.Handle("/", http.FileServer(http.FS(static)))

	return mux
}
