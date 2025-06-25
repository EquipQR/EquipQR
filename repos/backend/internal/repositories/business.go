package repositories

import (
	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/google/uuid"
)

func GetBusinessByID(id string) (*models.Business, error) {
	var business models.Business

	if err := database.DB.First(&business, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &business, nil
}

func CreateBusiness(business *models.Business) error {
	return database.DB.Create(business).Error
}

func GetBusinessesPaginated(limit int, offset int) ([]models.Business, error) {
	var businesses []models.Business
	result := database.DB.Limit(limit).Offset(offset).Find(&businesses)
	return businesses, result.Error
}

func CountBusinessMembers(businessID string) (int64, error) {
	var count int64
	err := database.DB.
		Table("user_businesses").
		Where("business_id = ?", businessID).
		Count(&count).
		Error

	return count, err
}

func CreatePendingJoinRequest(userID uuid.UUID, businessID uuid.UUID) error {
	req := models.PendingJoinRequest{
		UserID:     userID,
		BusinessID: businessID,
	}
	return database.DB.Create(&req).Error
}

func AddUserToBusiness(userID string, businessID string, isAdmin bool) error {
	entry := models.UserBusiness{
		UserID:     uuid.MustParse(userID),
		BusinessID: uuid.MustParse(businessID),
		IsAdmin:    isAdmin,
	}

	// Prevent duplicate entries
	var existing models.UserBusiness
	err := database.DB.
		Where("user_id = ? AND business_id = ?", entry.UserID, entry.BusinessID).
		First(&existing).Error

	if err == nil {
		// Already exists — optionally promote to admin
		if !existing.IsAdmin && isAdmin {
			existing.IsAdmin = true
			return database.DB.Save(&existing).Error
		}
		return nil
	}

	return database.DB.Create(&entry).Error
}
