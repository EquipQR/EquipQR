package repositories

import (
	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
)

func GetIssueByID(id string) (*models.Issue, error) {
	var issue models.Issue

	if err := database.DB.First(&issue, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &issue, nil
}

func CreateIssue(issue *models.Issue) error {
	return database.DB.Create(issue).Error
}
