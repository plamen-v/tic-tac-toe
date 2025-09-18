package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe/src/config"
	"github.com/plamen-v/tic-tac-toe/src/services/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	AUTHORIZATION_HEADER        string = "Authorization"
	AUTHORIZATION_HEADER_PREFIX string = "Bearer "
	PLAYER_ID_CLAIM_KEY         string = "player_id"
)

type AuthenticationService interface {
	CreateToken(player *models.Player) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	Authenticate(string, string) (*models.Player, error)
}

type ExtendedClaims struct {
	PlayerID int64 `json:"player_id"`
	jwt.RegisteredClaims
}

func (c ExtendedClaims) Validate() error {
	if c.PlayerID <= 0 {
		return fmt.Errorf("user_id cannot be empty") //todo!
	}

	// You can add more custom checks here todo!
	return nil
}

func NewAuthenticationService(config *config.AppConfiguration, repo repository.Repository) AuthenticationService {
	return &authenticationService{
		config: config,
		repo:   repo,
	}
}

type authenticationService struct {
	config *config.AppConfiguration
	repo   repository.Repository
}

func (s *authenticationService) CreateToken(player *models.Player) (string, error) {
	claims := ExtendedClaims{
		PlayerID: player.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   s.config.AppName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Second)), //todo! config
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    player.Login,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := jwtToken.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authenticationService) ValidateToken(tokenString string) (*jwt.Token, error) {
	// Parse the token with the secret key
	jwtToken, err := jwt.ParseWithClaims(tokenString, &ExtendedClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return jwtToken, nil
}

func (s *authenticationService) Authenticate(login string, password string) (*models.Player, error) {
	player, err := s.repo.Players().GetByLogin(login, nil) //todo! not found error
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return player, nil
}
