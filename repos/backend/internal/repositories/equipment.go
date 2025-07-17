package repositories

import (
	"encoding/json"
	"errors"

	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func GetEquipmentByID(id string) (*models.Equipment, error) {
	var eq models.Equipment
	result := database.DB.Preload("Business").First(&eq, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &eq, nil
}

func GetIssuesByEquipmentID(id string) ([]models.Issue, error) {
	var issues []models.Issue
	result := database.DB.Where("equipment_id = ?", id).Find(&issues)
	return issues, result.Error
}

func CreateEquipmentEntry(reqID string, status string, eqType string, location string, moreFields map[string]any) (*models.Equipment, error) {
	businessID, err := uuid.Parse(reqID)
	if err != nil {
		return nil, errors.New("invalid business_id")
	}

	_, err = GetBusinessByID(reqID)
	if err != nil {
		return nil, errors.New("business not found")
	}

	moreFieldsJSON, err := json.Marshal(moreFields)
	if err != nil {
		return nil, errors.New("invalid more_fields format")
	}

	equipment := models.Equipment{
		BusinessID: businessID,
		Status:     status,
		Type:       eqType,
		Location:   location,
		MoreFields: datatypes.JSON(moreFieldsJSON),
	}

	if err := database.DB.Create(&equipment).Error; err != nil {
		return nil, err
	}

	return GetEquipmentByID(equipment.ID.String())
}
