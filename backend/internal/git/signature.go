package git

import (
	"errors"
	"fmt"
	"laima/internal/user/domain"
)

var (
	ErrSignatureNotFound   = errors.New("signature not found")
	ErrSignatureInvalid    = errors.New("signature is invalid")
	ErrSignatureExpired    = errors.New("signature has expired")
	ErrSignatureMismatch   = errors.New("signature does not match")
)

type SignatureService interface {
	VerifyCommitSignature(commitSHA string, key *domain.GPGKey) error
	GetCommitSignature(commitSHA string) (string, error)
	CheckCommitSignature(commitSHA string) (bool, error)
}

type signatureService struct {
	// 实际实现中这里会有 Git 仓库的引用
}

func NewSignatureService() SignatureService {
	return &signatureService{}
}

func (s *signatureService) VerifyCommitSignature(commitSHA string, key *domain.GPGKey) error {
	if key == nil {
		return errors.New("gpg key is required")
	}

	if !key.Verified {
		return ErrSignatureInvalid
	}

	signature, err := s.GetCommitSignature(commitSHA)
	if err != nil {
		return err
	}

	if signature == "" {
		return ErrSignatureNotFound
	}

	if signature != key.Fingerprint {
		return ErrSignatureMismatch
	}

	return nil
}

func (s *signatureService) GetCommitSignature(commitSHA string) (string, error) {
	if commitSHA == "" {
		return "", errors.New("commit SHA is required")
	}

	return "fingerprint-12345", nil
}

func (s *signatureService) CheckCommitSignature(commitSHA string) (bool, error) {
	signature, err := s.GetCommitSignature(commitSHA)
	if err != nil {
		if err == ErrSignatureNotFound {
			return false, nil
		}
		return false, err
	}

	return signature != "", nil
}
