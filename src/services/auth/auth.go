package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/plamen-v/tic-tac-toe-models/models"
	"github.com/plamen-v/tic-tac-toe/src/config"
	"github.com/plamen-v/tic-tac-toe/src/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	AUTHORIZATION_HEADER        string = "Authorization"
	AUTHORIZATION_HEADER_PREFIX string = "Bearer "
	PLAYER_ID_CLAIM_KEY         string = "player_id"
)

type AuthenticationService interface {
	ValidateToken(token string) (*jwt.Token, error)
	Authenticate(context.Context, string, string) (*models.Player, string, error)
}

type ExtendedClaims struct {
	PlayerID uuid.NullUUID `json:"player_id"`
	jwt.RegisteredClaims
}

func (c ExtendedClaims) Validate() error {
	if !c.PlayerID.Valid {
		return models.NewAuthorizationError("player_id claim is invalid")
	}

	return nil
}

func NewAuthenticationService(config *config.AppConfiguration, db *sql.DB) AuthenticationService {
	return &authenticationServiceImpl{
		config: config,
		db:     db,
	}
}

type authenticationServiceImpl struct {
	config *config.AppConfiguration
	db     *sql.DB
}

func (s *authenticationServiceImpl) playerRepositoryFactory(q repository.Querier) repository.PlayerRepository {
	return repository.NewPlayerRepository(q)
}

func (s *authenticationServiceImpl) ValidateToken(tokenString string) (*jwt.Token, error) {
	if jwtToken, err := jwt.ParseWithClaims(tokenString, &ExtendedClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Secret), nil
	}); err != nil {
		return nil, models.NewAuthorizationError(err.Error())
	} else {
		return jwtToken, nil
	}
}

func (s *authenticationServiceImpl) Authenticate(ctx context.Context, login string, password string) (*models.Player, string, error) {
	player, err := s.playerRepositoryFactory(s.db).GetByLogin(ctx, login)
	if err != nil {
		return nil, "", models.NewAuthorizationError(err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(password))
	if err != nil {
		return nil, "", models.NewAuthorizationErrorf("invalid password")
	}

	token, err := s.createToken(player)
	if err != nil {
		return nil, "", models.NewGenericError(err.Error())
	}
	return player, token, nil
}

func (s *authenticationServiceImpl) createToken(player *models.Player) (string, error) {
	claims := ExtendedClaims{
		PlayerID: uuid.NullUUID{UUID: player.ID, Valid: true},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   s.config.AppName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3600 * time.Second)),
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
