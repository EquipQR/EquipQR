package repositories

import (
	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/google/uuid"
)

func ApprovePendingJoin(userID uuid.UUID, businessID uuid.UUID) error {
	_, err := GetPendingJoinRequest(userID, businessID)
	if err != nil {
		return err
	}

	if err := AddUserToBusiness(userID.String(), businessID.String(), false); err != nil {
		return err
	}

	if err := DeletePendingJoinRequest(userID, businessID); err != nil {
		return err
	}

	return nil
}

func GetPendingJoinRequest(userID uuid.UUID, businessID uuid.UUID) (*models.PendingJoinRequest, error) {
	var req models.PendingJoinRequest
	err := database.DB.First(&req, "user_id = ? AND business_id = ?", userID, businessID).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func DeletePendingJoinRequest(userID uuid.UUID, businessID uuid.UUID) error {
	return database.DB.
		Where("user_id = ? AND business_id = ?", userID, businessID).
		Delete(&models.PendingJoinRequest{}).Error
}
