package repositories

import (
	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
)

func GetEquipmentByID(id string) (*models.Equipment, error) {
	var eq models.Equipment

	result := database.DB.Preload("Business").First(&eq, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &eq, nil
}

func CreateEquipment(equipment *models.Equipment) error {
	return database.DB.Create(equipment).Error
}
