package repositories

import (
	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/google/uuid"
)

func GetAllPendingJoinsForBusiness(businessID uuid.UUID) ([]models.PendingJoinRequest, error) {
	var requests []models.PendingJoinRequest
	err := database.DB.Where("business_id = ?", businessID).Find(&requests).Error
	return requests, err
}

func ApprovePendingJoin(requestID uuid.UUID) error {
	req, err := GetPendingJoinRequestByID(requestID)
	if err != nil {
		return err
	}

	if err := AddUserToBusiness(req.UserID.String(), req.BusinessID.String(), false); err != nil {
		return err
	}

	if err := DeletePendingJoinRequestByID(requestID); err != nil {
		return err
	}

	return nil
}

func GetPendingJoinRequestByID(id uuid.UUID) (*models.PendingJoinRequest, error) {
	var req models.PendingJoinRequest
	err := database.DB.First(&req, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func DeletePendingJoinRequestByID(id uuid.UUID) error {
	return database.DB.
		Where("id = ?", id).
		Delete(&models.PendingJoinRequest{}).Error
}
