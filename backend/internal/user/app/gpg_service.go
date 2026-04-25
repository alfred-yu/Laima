package app

import (
	"errors"
	"time"

	"laima/internal/user/domain"
)

var (
	ErrGPGKeyNotFound      = errors.New("gpg key not found")
	ErrGPGKeyInvalid       = errors.New("gpg key is invalid")
	ErrGPGKeyAlreadyExists = errors.New("gpg key already exists")
)

type GPGService interface {
	AddGPGKey(userID int, req *domain.GPGKeyRequest) (*domain.GPGKey, error)
	VerifyGPGKey(userID int, keyID string) error
	ListGPGKeys(userID int) ([]*domain.GPGKey, error)
	DeleteGPGKey(userID int, keyID string) error
}

type gpgService struct {
	keys map[int][]*domain.GPGKey
}

func NewGPGService() GPGService {
	return &gpgService{
		keys: make(map[int][]*domain.GPGKey),
	}
}

func (s *gpgService) AddGPGKey(userID int, req *domain.GPGKeyRequest) (*domain.GPGKey, error) {
	if req.Key == "" {
		return nil, errors.New("gpg key is required")
	}

	if req.Name == "" {
		return nil, errors.New("key name is required")
	}

	for _, key := range s.keys[userID] {
		if key.Fingerprint == req.Fingerprint {
			return nil, ErrGPGKeyAlreadyExists
		}
	}

	key := &domain.GPGKey{
		ID:          len(s.keys[userID]) + 1,
		UserID:      userID,
		Name:        req.Name,
		Key:         req.Key,
		Fingerprint: req.Fingerprint,
		CreatedAt:   time.Now(),
		ExpiresAt:   req.ExpiresAt,
		Verified:    false,
	}

	s.keys[userID] = append(s.keys[userID], key)

	return key, nil
}

func (s *gpgService) VerifyGPGKey(userID int, keyID string) error {
	for _, key := range s.keys[userID] {
		if key.ID == keyID {
			key.Verified = true
			key.UpdatedAt = time.Now()
			return nil
		}
	}

	return ErrGPGKeyNotFound
}

func (s *gpgService) ListGPGKeys(userID int) ([]*domain.GPGKey, error) {
	return s.keys[userID], nil
}

func (s *gpgService) DeleteGPGKey(userID int, keyID string) error {
	keys := s.keys[userID]
	var newKeys []*domain.GPGKey

	for _, key := range keys {
		if key.ID != keyID {
			newKeys = append(newKeys, key)
		}
	}

	if len(newKeys) == len(keys) {
		return ErrGPGKeyNotFound
	}

	s.keys[userID] = newKeys
	return nil
}
