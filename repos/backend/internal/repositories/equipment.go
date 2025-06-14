package repositories

import (
	"errors"

	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
)

func GetEquipmentByID(id string) (*models.Equipment, error) {
	var eq models.Equipment

	result := database.DB.First(&eq, "id = ?", id)
	if errors.Is(result.Error, database.DB.Error) || result.Error != nil {
		return nil, result.Error
	}

	return &eq, nil
}
