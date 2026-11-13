package application

import (
	"context"
	"errors"

	"github.com/SpaceUp-UPC/Setrya-SpaceManagement-Service/internal/domain"
)

type SpaceRepository interface {
	Save(context.Context, domain.Space) error
	Find(context.Context, string) (domain.Space, error)
	ByOwner(context.Context, string) ([]domain.Space, error)
	Delete(context.Context, string) error
}

type DeviceRepository interface {
	LinkDevice(context.Context, domain.LinkedIoTDevice) error
	DevicesBySpace(context.Context, string) ([]domain.LinkedIoTDevice, error)
	DeactivateDevice(context.Context, string, string) error
}

type ReviewRepository interface {
	SaveReview(context.Context, domain.Review) error
	FindReview(context.Context, string) (domain.Review, error)
	ReviewsBySpace(context.Context, string) ([]domain.Review, error)
}

var ErrNotFound = errors.New("space not found")
var ErrDeviceNotFound = errors.New("linked device not found")

type Service struct {
	Spaces  SpaceRepository
	DevicesRepo DeviceRepository
	ReviewsRepo ReviewRepository
}

func (s Service) Create(ctx context.Context, space domain.Space) error {
	return s.Spaces.Save(ctx, space)
}

func (s Service) Get(ctx context.Context, id string) (domain.Space, error) {
	return s.Spaces.Find(ctx, id)
}

func (s Service) List(ctx context.Context, ownerID string) ([]domain.Space, error) {
	return s.Spaces.ByOwner(ctx, ownerID)
}

func (s Service) Update(ctx context.Context, id, name, description string, location domain.Location, dimensions domain.Dimensions) error {
	space, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := space.Update(name, description, location, dimensions); err != nil {
		return err
	}
	return s.Spaces.Save(ctx, space)
}

func (s Service) SetStatus(ctx context.Context, id string, active bool) error {
	space, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if active {
		space.Publish()
	} else {
		space.Pause()
	}
	return s.Spaces.Save(ctx, space)
}

func (s Service) Delete(ctx context.Context, id string) error {
	return s.Spaces.Delete(ctx, id)
}

func (s Service) LinkDevice(ctx context.Context, device domain.LinkedIoTDevice) error {
	if _, err := s.Get(ctx, device.SpaceID); err != nil {
		return err
	}
	return s.DevicesRepo.LinkDevice(ctx, device)
}

func (s Service) Devices(ctx context.Context, spaceID string) ([]domain.LinkedIoTDevice, error) {
	if _, err := s.Get(ctx, spaceID); err != nil {
		return nil, err
	}
	return s.DevicesRepo.DevicesBySpace(ctx, spaceID)
}

func (s Service) DeactivateDevice(ctx context.Context, spaceID, deviceID string) error {
	if _, err := s.Get(ctx, spaceID); err != nil {
		return err
	}
	return s.DevicesRepo.DeactivateDevice(ctx, spaceID, deviceID)
}

func (s Service) AddReview(ctx context.Context, review domain.Review) error {
	space, err := s.Get(ctx, review.SpaceID)
	if err != nil {
		return err
	}
	if space.Status != domain.StatusActive {
		return errors.New("only active spaces can receive reviews")
	}
	return s.ReviewsRepo.SaveReview(ctx, review)
}

func (s Service) Reviews(ctx context.Context, spaceID string) ([]domain.Review, error) {
	if _, err := s.Get(ctx, spaceID); err != nil {
		return nil, err
	}
	return s.ReviewsRepo.ReviewsBySpace(ctx, spaceID)
}
