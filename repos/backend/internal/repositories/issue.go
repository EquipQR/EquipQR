package repositories

import (
	"errors"

	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/google/uuid"
)

func GetIssueByID(id string) (*models.Issue, error) {
	var issue models.Issue
	if err := database.DB.First(&issue, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &issue, nil
}

func CreateIssueFromRequest(req utils.CreateIssueRequest) (*models.Issue, error) {
	equipmentID, err := uuid.Parse(req.EquipmentID)
	if err != nil {
		return nil, errors.New("invalid equipment_id")
	}

	assigneeID, err := uuid.Parse(req.AssigneeID)
	if err != nil {
		return nil, errors.New("invalid assignee_id")
	}

	issue := models.Issue{
		Title:       req.Title,
		Description: req.Description,
		EquipmentID: equipmentID,
		AssigneeID:  assigneeID,
	}

	if err := database.DB.Create(&issue).Error; err != nil {
		return nil, err
	}

	return &issue, nil
}
