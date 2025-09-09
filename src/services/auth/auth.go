package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/plamen-v/tic-tac-toe/src/config"
	"github.com/plamen-v/tic-tac-toe/src/models"
	"github.com/plamen-v/tic-tac-toe/src/services/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	AUTHORIZATION_HEADER        string = "Authorization"
	AUTHORIZATION_HEADER_PREFIX string = "Bearer "
)

type AuthenticationService interface {
	CreateToken(player *models.Player) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	AuthenticatePlayer(string, string) (*models.Player, error)
}

func NewAuthenticationService(config *config.AppConfiguration, playerRepo repository.PlayerRepository) AuthenticationService {
	return &authenticationService{
		config:     config,
		playerRepo: playerRepo,
	}
}

type authenticationService struct {
	config     *config.AppConfiguration
	playerRepo repository.PlayerRepository
}

func (s *authenticationService) CreateToken(player *models.Player) (string, error) {
	claims := jwt.MapClaims{
		"sub": s.config.AppName, //TODO!
		"iss": fmt.Sprintf("%d", player.ID),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour).Unix(), //TODO!
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
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.config.Secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return jwtToken, nil
}

func (s *authenticationService) AuthenticatePlayer(login string, password string) (*models.Player, error) {
	player, err := s.playerRepo.GetByLogin(login)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return player, nil
}
