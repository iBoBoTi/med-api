package server

import (
	"context"
	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/services/jwt"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"

	"github.com/decagonhq/meddle-api/errors"
	"github.com/decagonhq/meddle-api/models"
	"github.com/decagonhq/meddle-api/server/response"
	"github.com/gin-gonic/gin"
)

func (s *Server) HandleSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		userResponse, err := s.AuthService.SignupUser(&user)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "user created successfully", http.StatusCreated, userResponse, nil)
	}
}

func (s *Server) handleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest models.LoginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			response.JSON(c, "", errors.ErrBadRequest.Status, nil, err)
			return
		}
		userResponse, err := s.AuthService.LoginUser(&loginRequest)
		if err != nil {
			response.JSON(c, "", err.Status, nil, err)
			return
		}
		response.JSON(c, "login successful", http.StatusOK, userResponse, nil)
	}
}

func GetValuesFromContext(c *gin.Context) (string, *models.User, *errors.Error) {
	var tokenI, userI interface{}
	var tokenExists, userExists bool

	if tokenI, tokenExists = c.Get("access_token"); !tokenExists {
		return "", nil, errors.New("forbidden", http.StatusForbidden)
	}
	if userI, userExists = c.Get("user"); !userExists {
		return "", nil, errors.New("forbidden", http.StatusForbidden)
	}

	token, ok := tokenI.(string)
	if !ok {
		return "", nil, errors.New("internal server error", http.StatusInternalServerError)
	}
	user, ok := userI.(*models.User)
	if !ok {
		return "", nil, errors.New("internal server error", http.StatusInternalServerError)
	}
	return token, user, nil
}

func (s *Server) handleLogout() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, user, err := GetValuesFromContext(c)
		if err != nil {
			response.JSON(c, "", err.Status, nil, err)
			return
		}
		claims, errr := jwt.ValidateAndGetClaims(token, s.Config.JWTSecret)
		if errr != nil {
			response.JSON(c, "", http.StatusUnauthorized, nil, errr)
			return
		}
		convertClaims, _ := claims["exp"].(int64) //jwt pkg to validate
		if convertClaims < time.Now().Unix() {
			accBlacklist := &models.BlackList{
				Email: user.Email,
				Token: token,
			}
			if err := s.AuthRepository.AddToBlackList(accBlacklist); err != nil {
				log.Printf("can't add access token to blacklist: %v\n", err)
				response.JSON(c, "logout failed", http.StatusInternalServerError, nil, errors.New("can't add access token to blacklist", http.StatusInternalServerError))
				return
			}
		}
		response.JSON(c, "logout successful", http.StatusOK, nil, nil)
	}
}

func (s *Server) handleFBLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := config.GetFacebookOAuthConfig(s.Config.FacebookClientID, s.Config.FacebookClientSecret, s.Config.FacebookRedirectURL)
		state, err := jwt.GenerateToken("", s.Config.JWTSecret)
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, err)
			return
		}
		url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

func (s *Server) fbCallbackHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var state = c.Query("state")
		var code = c.Query("code")

		_, err := jwt.ValidateToken(state, s.Config.JWTSecret)
		if err != nil {
			respondAndAbort(c, "", http.StatusUnauthorized, nil, errors.New("invalid login", http.StatusUnauthorized))
			return
		}

		var OAuth2Config = config.GetFacebookOAuthConfig(s.Config.FacebookClientID, s.Config.FacebookClientSecret, s.Config.FacebookRedirectURL)

		token, err := OAuth2Config.Exchange(context.Background(), code)
		if err != nil || token == nil {
			respondAndAbort(c, "", http.StatusUnauthorized, nil, errors.New("invalid token", http.StatusUnauthorized))
			return
		}

		authToken, errr := s.AuthService.FacebookSignInUser(token.AccessToken)
		if errr != nil {
			respondAndAbort(c, "", http.StatusUnauthorized, nil, errors.New("invalid authToken", http.StatusUnauthorized))
			return
		}

		response.JSON(c, "facebook sign in successful", http.StatusOK, authToken, nil)
	}
}

func (s *Server) handleGetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.JSON(c, "successful", http.StatusOK, nil, nil)
	}
}

func (s *Server) handleUpdateUserDetails() gin.HandlerFunc {
	return func(c *gin.Context) {

		response.JSON(c, "successful", http.StatusOK, nil, nil)
	}
}

func (s *Server) handleShowProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

		response.JSON(c, "successful", http.StatusOK, nil, nil)
	}
}

func (s *Server) HandleVerifyEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		paramToken := c.Param("token")
		err := s.AuthService.VerifyEmail(paramToken)
		if err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})

		//response.JSON(c, "user successfully verified", http.StatusOK, nil, nil)
	}
}
