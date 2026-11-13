package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/SpaceUp-UPC/Setrya-SpaceManagement-Service/internal/application"
	"github.com/SpaceUp-UPC/Setrya-SpaceManagement-Service/internal/domain"
)

type Handler struct { App application.Service }

type spaceRequest struct {
	OwnerID string `json:"ownerId"`
	Name string `json:"name"`
	Description string `json:"description"`
	Location domain.Location `json:"location"`
	Dimensions domain.Dimensions `json:"dimensions"`
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/v1/spaces"), "/")
	if path == "" { h.create(w, r); return }
	parts := strings.Split(path, "/")
	if len(parts) == 2 && parts[1] == "devices" {
		h.devices(w, r, parts[0], "")
		return
	}
	if len(parts) == 4 && parts[1] == "devices" && parts[3] == "deactivate" {
		h.devices(w, r, parts[0], parts[2])
		return
	}
	if len(parts) == 2 && parts[1] == "reviews" {
		h.reviews(w, r, parts[0])
		return
	}
	if len(parts) == 2 && parts[1] == "availability" { h.setAvailability(w, r, parts[0]); return }
	if len(parts) == 2 && parts[1] == "pause" { h.pause(w, r, parts[0]); return }
	if len(parts) != 1 || parts[0] == "" { http.NotFound(w, r); return }
	switch r.Method {
	case http.MethodGet: h.get(w, r, parts[0])
	case http.MethodPut: h.update(w, r, parts[0])
	case http.MethodDelete: h.delete(w, r, parts[0])
	default: http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h Handler) create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	var request spaceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil { http.Error(w, "invalid JSON", http.StatusBadRequest); return }
	id, err := newID()
	if err != nil { http.Error(w, "could not create space id", http.StatusInternalServerError); return }
	space, err := domain.NewSpace(id, request.OwnerID, request.Name, request.Description, request.Location, request.Dimensions)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
	if err := h.App.Create(r.Context(), space); err != nil { http.Error(w, "could not save space", http.StatusInternalServerError); return }
	w.Header().Set("Location", "/api/v1/spaces/"+space.ID)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(space)
}

func (h Handler) get(w http.ResponseWriter, r *http.Request, id string) {
	space, err := h.App.Get(r.Context(), id)
	if err == application.ErrNotFound { http.Error(w, err.Error(), http.StatusNotFound); return }
	if err != nil { http.Error(w, "could not retrieve space", http.StatusInternalServerError); return }
	_ = json.NewEncoder(w).Encode(space)
}

func (h Handler) update(w http.ResponseWriter, r *http.Request, id string) {
	var request spaceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil { http.Error(w, "invalid JSON", http.StatusBadRequest); return }
	if err := h.App.Update(r.Context(), id, request.Name, request.Description, request.Location, request.Dimensions); err != nil {
		if err == application.ErrNotFound { http.Error(w, err.Error(), http.StatusNotFound) } else { http.Error(w, err.Error(), http.StatusBadRequest) }
		return
	}
	h.get(w, r, id)
}

func (h Handler) delete(w http.ResponseWriter, r *http.Request, spaceID string) {
	if err := h.App.Delete(r.Context(), spaceID); err != nil {
		if err == application.ErrNotFound { http.NotFound(w, r) } else { http.Error(w, "could not deactivate space", http.StatusInternalServerError) }
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h Handler) setAvailability(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodPatch { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	var request struct { Active *bool `json:"active"` }
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.Active == nil { http.Error(w, "active is required and must be boolean", http.StatusBadRequest); return }
	if err := h.App.SetStatus(r.Context(), id, *request.Active); err != nil {
		if err == application.ErrNotFound { http.Error(w, err.Error(), http.StatusNotFound) } else { http.Error(w, "could not update availability", http.StatusInternalServerError) }
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h Handler) pause(w http.ResponseWriter, r *http.Request, spaceID string) {
	if r.Method != http.MethodPatch { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	if err := h.App.SetStatus(r.Context(), spaceID, false); err != nil {
		if err == application.ErrNotFound { http.NotFound(w, r) } else { http.Error(w, "could not pause space", http.StatusInternalServerError) }
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h Handler) devices(w http.ResponseWriter, r *http.Request, spaceID, deviceID string) {
	if deviceID == "" && r.Method == http.MethodPost {
		var request struct {
			DeviceID   string `json:"deviceId"`
			DeviceType string `json:"deviceType"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		id, err := newID()
		if err != nil { http.Error(w, "could not create link id", http.StatusInternalServerError); return }
		device, err := domain.NewLinkedIoTDevice(id, spaceID, request.DeviceID, request.DeviceType)
		if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
		if err := h.App.LinkDevice(r.Context(), device); err != nil { http.Error(w, "could not link device", http.StatusConflict); return }
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(device)
		return
	}
	if deviceID == "" && r.Method == http.MethodGet {
		devices, err := h.App.Devices(r.Context(), spaceID)
		if err == application.ErrNotFound { http.Error(w, err.Error(), http.StatusNotFound); return }
		if err != nil { http.Error(w, "could not retrieve devices", http.StatusInternalServerError); return }
		_ = json.NewEncoder(w).Encode(devices)
		return
	}
	if deviceID != "" && r.Method == http.MethodPatch {
		if err := h.App.DeactivateDevice(r.Context(), spaceID, deviceID); err != nil {
			if err == application.ErrDeviceNotFound || err == application.ErrNotFound { http.NotFound(w, r) } else { http.Error(w, "could not deactivate device link", http.StatusInternalServerError) }
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func (h Handler) reviews(w http.ResponseWriter, r *http.Request, spaceID string) {
	if r.Method == http.MethodGet {
		reviews, err := h.App.Reviews(r.Context(), spaceID)
		if err == application.ErrNotFound { http.NotFound(w, r); return }
		if err != nil { http.Error(w, "could not retrieve reviews", http.StatusInternalServerError); return }
		_ = json.NewEncoder(w).Encode(reviews)
		return
	}
	if r.Method == http.MethodPost {
		var request struct {
			AuthorID string `json:"authorId"`
			Rating   int    `json:"rating"`
			Comment  string `json:"comment"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil { http.Error(w, "invalid JSON", http.StatusBadRequest); return }
		id, err := newID()
		if err != nil { http.Error(w, "could not create review id", http.StatusInternalServerError); return }
		review, err := domain.NewReview(id, spaceID, request.AuthorID, request.Rating, request.Comment)
		if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
		if err := h.App.AddReview(r.Context(), review); err != nil {
			if err == application.ErrNotFound { http.NotFound(w, r) } else { http.Error(w, err.Error(), http.StatusConflict) }
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(review)
		return
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}



