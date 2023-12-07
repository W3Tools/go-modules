package gmjwt

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JwtClient struct {
	SignedKey  string
	Expiration int64
}

type JwtClaims struct {
	Seed      interface{} `json:"seed"`
	Timestamp int64       `json:"timestamp"`
	jwt.StandardClaims
}

func InitJwtClient(signedKey string, expiration int64) (client *JwtClient) {
	if strings.EqualFold(signedKey, "") {
		signedKey = uuid.New().String()
	}
	cli := &JwtClient{
		SignedKey:  signedKey,
		Expiration: expiration,
	}
	return cli
}

func (c *JwtClient) NewJwtClaims(seed interface{}) *JwtClaims {
	claims := &JwtClaims{
		Seed:      seed,
		Timestamp: time.Now().Unix(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(c.Expiration) * time.Minute).Unix(),
		},
	}
	return claims
}

func (c *JwtClient) NewJwtToken(seed interface{}) (string, error) {
	claims := c.NewJwtClaims(seed)

	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(c.SignedKey))
}

func (c *JwtClient) ParseJwtToken(token string) (*JwtClaims, error) {
	if strings.EqualFold(token, "") {
		return nil, fmt.Errorf("jwt token not found")
	}

	claims := &JwtClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.SignedKey), nil
	})

	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired {
			return claims, fmt.Errorf("expired jwt token")
		}

		return nil, err
	}
	return claims, nil
}

func (c *JwtClaims) ParseJwtTokenSeed(v any) error {
	data, err := json.Marshal(c.Seed)
	if err != nil {
		return fmt.Errorf("invalid seed of jwt token")
	}

	return json.Unmarshal(data, v)
}
