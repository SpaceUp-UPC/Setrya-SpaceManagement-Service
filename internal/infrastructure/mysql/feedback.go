package mysql

import (
	"context"
	"database/sql"

	"github.com/SpaceUp-UPC/Setrya-SpaceManagement-Service/internal/application"
	"github.com/SpaceUp-UPC/Setrya-SpaceManagement-Service/internal/domain"
)

func (r Repository) LinkDevice(ctx context.Context, device domain.LinkedIoTDevice) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO linked_iot_devices (id, space_id, device_id, device_type, status, linked_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, device.ID, device.SpaceID, device.DeviceID, device.DeviceType, device.Status, device.LinkedAt)
	return err
}

func (r Repository) DevicesBySpace(ctx context.Context, spaceID string) ([]domain.LinkedIoTDevice, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, space_id, device_id, device_type, status, linked_at
		FROM linked_iot_devices WHERE space_id = ? ORDER BY linked_at DESC
	`, spaceID)
	if err != nil { return nil, err }
	defer rows.Close()

	devices := []domain.LinkedIoTDevice{}
	for rows.Next() {
		var device domain.LinkedIoTDevice
		if err := rows.Scan(&device.ID, &device.SpaceID, &device.DeviceID, &device.DeviceType, &device.Status, &device.LinkedAt); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, rows.Err()
}

func (r Repository) DeactivateDevice(ctx context.Context, spaceID, deviceID string) error {
	result, err := r.DB.ExecContext(ctx, `
		UPDATE linked_iot_devices SET status = 'inactive'
		WHERE space_id = ? AND device_id = ?
	`, spaceID, deviceID)
	if err != nil { return err }
	count, err := result.RowsAffected()
	if err != nil { return err }
	if count == 0 { return application.ErrDeviceNotFound }
	return nil
}

func (r Repository) SaveReview(ctx context.Context, review domain.Review) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO reviews (id, space_id, author_id, rating, comment, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, review.ID, review.SpaceID, review.AuthorID, review.Rating, review.Comment, review.CreatedAt)
	return err
}

func (r Repository) ReviewsBySpace(ctx context.Context, spaceID string) ([]domain.Review, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, space_id, author_id, rating, comment, created_at
		FROM reviews WHERE space_id = ? ORDER BY created_at DESC
	`, spaceID)
	if err != nil { return nil, err }
	defer rows.Close()

	reviews := []domain.Review{}
	for rows.Next() {
		var review domain.Review
		if err := rows.Scan(&review.ID, &review.SpaceID, &review.AuthorID, &review.Rating, &review.Comment, &review.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, rows.Err()
}

func (r Repository) FindReview(ctx context.Context, id string) (domain.Review, error) {
	var review domain.Review
	err := r.DB.QueryRowContext(ctx, `
		SELECT id, space_id, author_id, rating, comment, created_at
		FROM reviews WHERE id = ?
	`, id).Scan(&review.ID, &review.SpaceID, &review.AuthorID, &review.Rating, &review.Comment, &review.CreatedAt)
	if err == sql.ErrNoRows { return review, application.ErrNotFound }
	return review, err
}

