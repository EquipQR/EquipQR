package repositories

import (
	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
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
