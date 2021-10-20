package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"gorm.io/gorm"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	CoreModel
	Hash      []byte    `json:"-"`
	Plaintext string    `json:"token" gorm:"-"`
	Scope     string    `json:"-"`
	Expiry    time.Time `json:"expiry"`
	UserID    int64     `json:"user_id"`
	User      User      `json:"user"`
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	// example plain token: Y3QMGX3PJ3WLRL2YRTQGQ6KRHU
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)

	// one way hash with no salt, user will send plain token...
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

type TokenModel struct {
	DB *gorm.DB
}

func (m TokenModel) Insert(token *Token) error {
	return m.DB.Create(&token).Error
}

func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}
