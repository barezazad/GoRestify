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

// RoleAPI for injecting role service
type RoleAPI struct {
	Service service.BaseRoleServ
	Engine  *core.Engine
}

// ProvideRoleAPI .
func ProvideRoleAPI(c service.BaseRoleServ) RoleAPI {
	return RoleAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch an role by its id
func (a *RoleAPI) FindByID(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.RoleTable)
	var err error
	var role base_model.Role
	var id uint

	if id, err = resp.GetID(c.Param("roleID"), "E1172071", base_term.Role); err != nil {
		return
	}

	if role, err = a.Service.FindByID(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ViewRole)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VInfo, base_term.Role).
		JSON(role)
}

// GetResources list of roles
func (a *RoleAPI) GetResources(c *gin.Context) {
	resp, _ := response.NewParam(c, base_model.RoleTable)

	resp.Record(base.ListRole)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, base_term.Roles).
		JSON(service.AllResources())
}

// GetAll list of roles
func (a *RoleAPI) GetAll(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.RoleTable)
	var roles []base_model.Role
	var err error

	if roles, err = a.Service.GetAll(params); err != nil {
		err = pkg_err.Take(err, "E1117072").Message(pkg_err.SomethingWentWrong).Build()
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListRole)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, base_term.Roles).
		JSON(roles)
}

// List of roles
func (a *RoleAPI) List(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.RoleTable)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = a.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListRole)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, base_term.Roles).
		JSON(data)
}

// Create role
func (a *RoleAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.RoleTable)
	var role, createdRole base_model.Role
	var err error

	if err = resp.Bind(&role, "E1131721", base_term.Role); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"role"), "rollback recover create role")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1199923").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if createdRole, err = a.Service.Create(params.Tx, role); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(base.CreateRole, role, createdRole)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VCreatedSuccessfully, base_term.Role).
		JSON(createdRole)
}

// Update role
func (a *RoleAPI) Update(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.RoleTable)
	var err error
	var role, roleBefore, roleUpdated base_model.Role

	if err = resp.Bind(&role, "E1177796", base_term.Role); err != nil {
		return
	}

	if role.ID, err = resp.GetID(c.Param("roleID"), "E1162114", base_term.Role); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"role"), "rollback recover create role")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1117850").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if roleUpdated, roleBefore, err = a.Service.Save(params.Tx, role); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(base.UpdateRole, roleBefore, role)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VUpdatedSuccessfully, base_term.Role).
		JSON(roleUpdated)
}

// Delete role
func (a *RoleAPI) Delete(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.RoleTable)
	var err error
	var role base_model.Role
	var id uint

	if id, err = resp.GetID(c.Param("roleID"), "E1166796", base_term.Role); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"role"), "rollback recover create role")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1190534").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if role, err = a.Service.Delete(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(base.DeleteRole, role)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VDeletedSuccessfully, base_term.Role).
		JSON()
}
