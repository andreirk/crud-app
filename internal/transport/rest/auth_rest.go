package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackietana/crud-app/internal/domain"
	log "github.com/sirupsen/logrus"
)

// @Summary Sign Up
// @Description sign up method
// @Tags auth
// @Accept json
// @Produce plain
// @Success 200 {string} string "Successfully signed up"
// @Failure 400
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var user domain.User
	if err := c.BindJSON(&user); err != nil {
		log.WithField("handler", "signUp").Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userService.SignUp(context.TODO(), user); err != nil {
		log.WithField("handler", "signUp").Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	c.String(http.StatusOK, "Successfully signed up")
}

// @Summary Sign In
// @Description sign in method
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {string} string
// @Failure 400
// @Failure 500 {string} string
// @Router /auth/sign-in [get]
func (h *Handler) signIn(c *gin.Context) {
	var user domain.UserSignIn
	if err := c.BindJSON(&user); err != nil {
		log.WithField("handler", "signIn").Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.userService.SignIn(context.TODO(), user)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			log.WithField("handler", "signIn").Error("user not found")
			c.JSON(http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		log.WithField("handler", "signIn").Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	c.SetCookie("refresh-token", refreshToken, 180, "/auth", "localhost", false, true)
	c.JSON(http.StatusOK, map[string]string{
		"token": accessToken,
	})
}

// TODO: add swagger
func (h *Handler) refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		log.WithField("handler", "refresh").Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	log.WithField("handler", "refresh").Infof("%s", cookie)

	accessToken, refreshToken, err := h.userService.RefreshTokens(c.Request.Context(), cookie)
	if err != nil {
		log.WithField("handler", "refresh").Error(err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.SetCookie("refresh-token", refreshToken, 180, "/auth", "localhost", false, true)
	c.JSON(http.StatusOK, map[string]string{
		"token": accessToken,
	})
}
