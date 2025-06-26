package repositories

import (
	"errors"

	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/lib/pq"
)

func GetUserByID(id string) (*models.User, error) {
	var user models.User

	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	if err := database.DB.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

func RegisterNewUser(req utils.CreateUserRequest) (*models.User, *models.Business, string, error) {
	hashedPassword, err := utils.GeneratePasswordHash(req.Password, utils.DefaultArgon2Config)
	if err != nil {
		return nil, nil, "", errors.New("failed to hash password")
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		IsActive: true,
	}

	if err := CreateUser(user); err != nil {
		return nil, nil, "", errors.New("failed to create user")
	}

	// Option A: Join existing business
	if req.BusinessID != "" {
		business, err := GetBusinessByID(req.BusinessID)
		if err != nil {
			return user, nil, "", errors.New("invalid business ID")
		}

		if err := CreatePendingJoinRequest(user.ID, business.ID); err != nil {
			return user, nil, "", errors.New("failed to request business approval")
		}

		// Return no business because it's pending
		return user, nil, "", nil
	}

	// Option B: Create new business
	if req.BusinessName == "" || req.BusinessType == "" {
		return nil, nil, "", errors.New("business name and type required for new business creation")
	}

	business := &models.Business{
		BusinessName:    req.BusinessName,
		BusinessEmail:   req.BusinessEmail,
		Phone:           req.Phone,
		CountryCode:     req.CountryCode,
		Type:            req.BusinessType,
		CompanySize:     req.CompanySize,
		Country:         req.Country,
		UserCanRegister: false,
		LoginMethods:    pq.StringArray{"password"},
	}

	if err := CreateBusiness(business); err != nil {
		return nil, nil, "", errors.New("failed to create business")
	}

	if err := AddUserToBusiness(user.ID.String(), business.ID.String(), true); err != nil {
		return nil, nil, "", errors.New("failed to assign user to business")
	}

	token, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		return nil, nil, "", errors.New("failed to generate token")
	}

	return user, business, token, nil
}
