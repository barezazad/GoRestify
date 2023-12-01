package base_api

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_term"
	"GoRestify/domain/service"
	"GoRestify/internal/core"
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
		return
	}

	user, err := a.Service.Login(params.Tx, auth)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Status(http.StatusOK).
		Message(base_term.UserLoggedInSuccessfully).
		JSON(user)
}
