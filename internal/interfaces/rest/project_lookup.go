package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/SpaceUp-UPC/Setrya-SpaceManagement-Service/internal/application"
)

type SpaceLookupHandler struct {
	App application.Service
}

func (h SpaceLookupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/v1/projects/")
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}

	space, err := h.App.Get(r.Context(), id)
	if err == application.ErrNotFound {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "could not retrieve space", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(space)
}
