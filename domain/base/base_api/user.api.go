package base_api

import (
	"GoRestify/domain/base"
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_term"
	"GoRestify/domain/service"
	"GoRestify/internal/core"
	"fmt"
	"net/http"

	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_terms"
	"GoRestify/pkg/response"

	"github.com/gin-gonic/gin"
)

// UserAPI for injecting user service
type UserAPI struct {
	Service service.BaseUserServ
	Engine  *core.Engine
}

// ProvideUserAPI .
func ProvideUserAPI(c service.BaseUserServ) UserAPI {
	return UserAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch an user by its id
func (a *UserAPI) FindByID(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.UserTable)
	var err error
	var user base_model.User
	var id uint

	if id, err = resp.GetID(c.Param("userID"), "E1161102", base_term.User); err != nil {
		return
	}

	if user, err = a.Service.FindByID(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ViewUser)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VInfo, base_term.User).
		JSON(user)
}

// GetAll list of users
func (a *UserAPI) GetAll(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.UserTable)
	var users []base_model.User
	var err error

	if users, err = a.Service.GetAll(params); err != nil {
		err = pkg_err.Take(err, "E1188617").Message(pkg_err.SomethingWentWrong).Build()
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListUser)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, base_term.Users).
		JSON(users)
}

// List of users
func (a *UserAPI) List(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.UserTable)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = a.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListUser)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, base_term.Users).
		JSON(data)
}

// Create user
func (a *UserAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.UserTable)
	var user, createdUser base_model.User
	var err error

	if err = resp.Bind(&user, "E1115226", base_term.User); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"user"), "rollback recover create user")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1190863").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if createdUser, err = a.Service.Create(params.Tx, user); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(base.CreateUser, user, createdUser)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VCreatedSuccessfully, base_term.User).
		JSON(createdUser)
}

// Update user
func (a *UserAPI) Update(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.UserTable)
	var err error
	var user, userBefore, userUpdated base_model.User

	if err = resp.Bind(&user, "E1150620", base_term.User); err != nil {
		return
	}

	if user.ID, err = resp.GetID(c.Param("userID"), "E1189417", base_term.User); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"user"), "rollback recover create user")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1119850").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if userUpdated, userBefore, err = a.Service.Save(params.Tx, user); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(base.UpdateUser, userBefore, user)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VUpdatedSuccessfully, base_term.User).
		JSON(userUpdated)
}

// Delete user
func (a *UserAPI) Delete(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.UserTable)
	var err error
	var user base_model.User
	var id uint

	if id, err = resp.GetID(c.Param("userID"), "E1135988", base_term.User); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"user"), "rollback recover create user")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1114769").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if user, err = a.Service.Delete(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(base.DeleteUser, user)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VDeletedSuccessfully, base_term.User).
		JSON()
}
