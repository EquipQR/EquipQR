package repositories

import (
	"errors"
	"fmt"

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

func CreateIssueFromRequest(req utils.CreateIssueRequest, assigneeID uuid.UUID) (*models.Issue, error) {
	fmt.Printf("[CreateIssue] 🛠️ Received request: %+v\n", req)

	equipmentID, err := uuid.Parse(req.EquipmentID)
	if err != nil {
		fmt.Printf("[CreateIssue] ❌ Invalid equipment_id: %s\n", req.EquipmentID)
		return nil, errors.New("invalid equipment_id")
	}

	issue := models.Issue{
		Title:       req.Title,
		Description: req.Description,
		EquipmentID: equipmentID,
		AssigneeID:  assigneeID,
	}

	fmt.Printf("[CreateIssue] ✅ Creating issue in DB: %+v\n", issue)

	if err := database.DB.Create(&issue).Error; err != nil {
		fmt.Printf("[CreateIssue] ❌ Failed to create issue: %v\n", err)
		return nil, err
	}

	fmt.Printf("[CreateIssue] ✅ Successfully created issue with ID: %s\n", issue.ID)
	return &issue, nil
}
