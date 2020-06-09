package app

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"
)

var (
	ErrUserIsBlocked   = errors.New("user is blocked")
	ErrInvalidPassword = errors.New("invalid password")
)

func (s *apiserver) login(login, password string) (string, string, error) { // return access token and role
	expert, err := s.repository.GetExpertByEmail(strings.ToLower(login))
	if err == nil {
		if !checkPwds(password, expert.Password) {
			return "", "", ErrInvalidPassword
		}
		token, err := s.generateAccessToken(expert.ID, UserRoleExpert)
		if err != nil {
			return "", "", err
		}

		return token, UserRoleExpert, nil
	}

	admin, err := s.repository.GetAdminByLogin(login)
	if err != nil {
		return "", "", fmt.Errorf("failed to get user for auth: %v", err)
	}

	if admin.Status != model.AdminStatusActive {
		return "", "", fmt.Errorf("user is blocked")
	}

	if !checkPwds(password, admin.Password) {
		return "", "", fmt.Errorf("invalid password")
	}

	token, err := s.generateAccessToken(admin.ID, UserRoleAdmin)
	if err != nil {
		return "", "", err
	}
	return token, UserRoleAdmin, nil
}

func (s *apiserver) generateAccessToken(userID, role string) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role":   role,
		"userId": userID,
		"exp":    time.Now().Add(s.config.AccessTokenDuration).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWTTokenSignKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %v", err)
	}
	return accessTokenString, nil
}

func (s *apiserver) checkAccessToken(token string) (string, string, error) {
	var (
		err      error
		jwtToken *jwt.Token
	)

	if jwtToken, err = jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTTokenSignKey), nil
	}); err != nil {
		if validationError, ok := err.(*jwt.ValidationError); ok {
			if validationError.Errors == jwt.ValidationErrorExpired {
				return "", "", fmt.Errorf("token expired")
			}
			return "", "", fmt.Errorf("invalid token: %v", validationError)
		}
		return "", "", err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return "", "", fmt.Errorf("invalid access token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid claim 'role'")
	}

	userID, ok := claims["userId"].(string)
	if !ok {
		return "", "", fmt.Errorf("invalid claim 'userId'")
	}

	return userID, role, nil
}

func checkPwds(rawPassword, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword)) == nil
}

func (s *apiserver) authenticationInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(authCookieKey)
		if err != nil {
			c.Redirect(http.StatusFound, "/admin/login")
			return
		}

		userID, role, err := s.checkAccessToken(token)
		if err != nil {
			c.Redirect(http.StatusFound, "/admin/login")
			return
		}

		c.Set("userID", userID)
		c.Set("role", role)

		var status string

		switch role {
		case UserRoleAdmin:
			admin, err := s.repository.GetAdmin(userID)
			if err != nil {
				c.Redirect(http.StatusFound, "/admin")
				return
			}

			status = admin.Status
		case UserRoleExpert:
			expert, err := s.repository.GetExpert(userID)
			if err != nil {
				c.Redirect(http.StatusFound, "/admin")
				return
			}

			status = expert.Status
		default:
			c.Redirect(http.StatusFound, "/admin")
			return
		}

		if status == "" {
			status = model.ExpertStatusBlocked
		}

		c.Set("status", status)
	}
}

func (s *apiserver) authorizationInterceptor(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("role")
		if !ok {
			c.Redirect(http.StatusFound, "/admin")
			return
		}

		rolestr, ok := role.(string)
		if !ok {
			c.Redirect(http.StatusFound, "/admin")
			return
		}

		for _, allowed := range roles {
			if rolestr == allowed {
				return
			}
		}

		c.Redirect(http.StatusFound, "/admin")
		return
	}
}

func (s *apiserver) checkStatusInterceptor(statuses ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		iStatus, ok := c.Get("status")
		if !ok {
			c.Redirect(http.StatusFound, "/admin")
			return
		}

		status, ok := iStatus.(string)
		if !ok {
			c.Redirect(http.StatusFound, "/admin")
			return
		}

		for _, allowed := range statuses {
			if status == allowed {
				return
			}
		}

		c.Redirect(http.StatusFound, "/admin")
		return
	}
}
