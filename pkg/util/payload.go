package util

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("время срока жизни токена истекло")
	ErrInvalidToken = errors.New("токен не действителен")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID int64 `json:"user_id"`
	Email     string    `json:"username"`
	Role string `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(userID int64, email string, role string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		UserID: userID,
		Email:     email,
		Role: role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}
