package domain

import (
	"errors"
	"strings"
	"time"
)

type LinkedIoTDevice struct {
	ID         string    `json:"id"`
	SpaceID    string    `json:"spaceId"`
	DeviceID   string    `json:"deviceId"`
	DeviceType string    `json:"deviceType"`
	Status     string    `json:"status"`
	LinkedAt   time.Time `json:"linkedAt"`
}

func NewLinkedIoTDevice(id, spaceID, deviceID, deviceType string) (LinkedIoTDevice, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(spaceID) == "" || strings.TrimSpace(deviceID) == "" {
		return LinkedIoTDevice{}, errors.New("id, spaceId y deviceId son requeridos")
	}
	if strings.TrimSpace(deviceType) == "" {
		return LinkedIoTDevice{}, errors.New("deviceType es requerido")
	}
	return LinkedIoTDevice{
		ID: id, SpaceID: spaceID, DeviceID: deviceID,
		DeviceType: deviceType, Status: "active", LinkedAt: time.Now().UTC(),
	}, nil
}

func (d *LinkedIoTDevice) Deactivate() { d.Status = "inactive" }

type Review struct {
	ID        string    `json:"id"`
	SpaceID   string    `json:"spaceId"`
	AuthorID  string    `json:"authorId"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewReview(id, spaceID, authorID string, rating int, comment string) (Review, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(spaceID) == "" || strings.TrimSpace(authorID) == "" {
		return Review{}, errors.New("id, spaceId y authorId son requeridos")
	}
	if rating < 1 || rating > 5 {
		return Review{}, errors.New("rating debe estar entre 1 y 5")
	}
	return Review{ID: id, SpaceID: spaceID, AuthorID: authorID, Rating: rating, Comment: strings.TrimSpace(comment), CreatedAt: time.Now().UTC()}, nil
}

func (d LinkedIoTDevice) IsLinkedTo(spaceID string) bool {
	return d.SpaceID == spaceID && d.Status == "active"
}

func (d *LinkedIoTDevice) Activate() { d.Status = "active" }

func (r Review) BelongsToSpace(spaceID string) bool {
	return r.SpaceID == spaceID
}

func (r *Review) UpdateContent(rating int, comment string) error {
	if rating < 1 || rating > 5 {
		return errors.New("rating debe estar entre 1 y 5")
	}
	r.Rating = rating
	r.Comment = strings.TrimSpace(comment)
	return nil
}
