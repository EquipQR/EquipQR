package repositories

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/utils"
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

func DenyPendingJoin(requestID uuid.UUID) error {
	_, err := GetPendingJoinRequestByID(requestID)
	if err != nil {
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

func GenerateInviteLinkWithEmail(businessID uuid.UUID, email string, secret string, baseURL string, expiryMinutes int) (string, error) {
	inviteToken := uuid.New()
	expiry := time.Now().Add(time.Duration(expiryMinutes) * time.Minute).Unix()

	data := fmt.Sprintf("%s:%s:%s:%d", businessID.String(), inviteToken.String(), email, expiry)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	signature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	params := url.Values{}
	params.Set("business", businessID.String())
	params.Set("token", inviteToken.String())
	params.Set("email", email)
	params.Set("exp", fmt.Sprintf("%d", expiry))
	params.Set("sig", signature)

	link := fmt.Sprintf("%s/join?%s", baseURL, params.Encode())
	return link, nil
}

func ProcessInvite(params utils.InviteParams, userID string) error {
	config := utils.LoadConfigFromEnv()
	exp, err := strconv.ParseInt(params.Expiry, 10, 64)
	if err != nil || time.Now().Unix() > exp {
		return errors.New("invite expired or invalid expiry format")
	}

	data := params.BusinessID + ":" + params.Token + ":" + params.Email + ":" + params.Expiry
	mac := hmac.New(sha256.New, []byte(config.JWT_Secret)) // reusing JWT secret
	mac.Write([]byte(data))
	expectedSig := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expectedSig), []byte(params.Signature)) {
		return errors.New("invalid signature")
	}

	businessID, err := uuid.Parse(params.BusinessID)
	if err != nil {
		return errors.New("invalid business ID")
	}

	if err := AddUserToBusiness(userID, businessID.String(), false); err != nil {
		return errors.New("failed to add user to business")
	}

	return nil
}
