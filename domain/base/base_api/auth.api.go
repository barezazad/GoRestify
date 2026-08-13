package base_api

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_term"
	"GoRestify/domain/service"
	"GoRestify/internal/core"
	"GoRestify/pkg/middleware"
	"GoRestify/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthAPI for injecting auth service
type AuthAPI struct {
	Service service.BaseAuthServ
	Engine  *core.Engine
}

// ProvideAuthAPI .
func ProvideAuthAPI(c service.BaseAuthServ) AuthAPI {
	return AuthAPI{Service: c, Engine: c.Engine}
}

// Login auth
func (a *AuthAPI) Login(c *gin.Context) {
	var auth base_model.Auth
	resp, params := response.NewParam(c, base_model.UserTable)
	var err error

	if err = resp.Bind(&auth, "E1170734", base_term.User); err != nil {
		middleware.RecordLoginFailure(c)
		resp.Error(err).JSON()
		return
	}

	user, err := a.Service.Login(params.Tx, auth)
	if err != nil {
		middleware.RecordLoginFailure(c)
		resp.Error(err).JSON()
		return
	}

	middleware.ResetLoginLimit(c)

	resp.Status(http.StatusOK).
		Message(base_term.UserLoggedInSuccessfully).
		JSON(user)
}

// RefreshToken exchanges a valid refresh token for a new access/refresh pair.
func (a *AuthAPI) RefreshToken(c *gin.Context) {
	var req base_model.RefreshTokenRequest
	resp, _ := response.NewParam(c, base_model.UserTable)

	if err := resp.Bind(&req, "E1170735", base_term.User); err != nil {
		return
	}

	account, err := a.Service.RefreshToken(req.RefreshToken)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Status(http.StatusOK).
		Message(base_term.UserLoggedInSuccessfully).
		JSON(account)
}
