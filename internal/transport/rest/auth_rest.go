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
// @Produce string
// @Success 200 {string} string "Successfully signed up"
// @Failure 400 {string}
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var user domain.User
	if err := c.BindJSON(&user); err != nil {
		log.WithField("handler", "signUp").Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.UserService.SignUp(context.TODO(), user); err != nil {
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
// @Success 200 {json} json "{"token": token_string}"
// @Failure 400 {string}
// @Failure 500 {json} json "{"error": error_string}"
// @Router /auth/sign-in [get]
func (h *Handler) signIn(c *gin.Context) {
	var user domain.UserSignIn
	if err := c.BindJSON(&user); err != nil {
		log.WithField("handler", "signIn").Error(err)
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.UserService.SignIn(context.TODO(), user)
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

	c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
