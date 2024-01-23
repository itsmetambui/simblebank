package token

import (
	"fmt"
	"time"

	"golang.org/x/crypto/chacha20poly1305"

	"github.com/o1egl/paseto"
)

type PastoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PastoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (m *PastoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}
	token, err := m.paseto.Encrypt(m.symmetricKey, payload, nil)
	return token, payload, err
}

func (m *PastoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := m.paseto.Decrypt(token, m.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if err := payload.Valid(); err != nil {
		return nil, err
	}
	return payload, nil
}
