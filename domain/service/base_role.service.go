package service

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_repo"
	"GoRestify/domain/base/base_term"
	"GoRestify/internal/core"
	"fmt"
	"strings"

	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/tx"

	"GoRestify/pkg/validator"
)

// BaseRoleServ for injecting  base_repo
type BaseRoleServ struct {
	Repo   base_repo.RoleRepo
	Engine *core.Engine
}

// ProvideBaseRoleService for role is used in wire
func ProvideBaseRoleService(roleRepo base_repo.RoleRepo) BaseRoleServ {
	return BaseRoleServ{
		Repo:   roleRepo,
		Engine: roleRepo.Engine,
	}
}

// FindByID for getting role by its id
func (s *BaseRoleServ) FindByID(tx tx.Tx, id uint) (role base_model.Role, err error) {

	key := fmt.Sprintf("%v-%v", base_term.Role, id)
	if ok := s.Engine.RedisCacheAPI.GetCache(tx, key, &role); ok {
		return
	}

	if role, err = s.Repo.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1191380", "can't fetch the role", id)
		return
	}

	err = s.Engine.RedisCacheAPI.Set(key, role)

	return
}

// GetAll of roles, it supports pagination and search and return count
func (s *BaseRoleServ) GetAll(params param.Param) (roles []base_model.Role, err error) {

	if ok := s.Engine.RedisCacheAPI.GetCache(params.Tx, base_term.Roles, &roles); ok {
		return
	}

	params.Pagination.Limit = 1000000
	if roles, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in roles list")
		return
	}

	err = s.Engine.RedisCacheAPI.Set(base_term.Roles, roles)

	return
}

// List of roles, it supports pagination and search and return count
func (s *BaseRoleServ) List(params param.Param) (roles []base_model.Role,
	count int64, err error) {

	if roles, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in roles list")
		return
	}

	if count, err = s.Repo.Count(params); err != nil {
		pkg_log.CheckError(err, "error in roles count")
	}

	return
}

// Create a role
func (s *BaseRoleServ) Create(tx tx.Tx, role base_model.Role) (createdRole base_model.Role, err error) {

	if err = validator.ValidateModel(role, base_term.Role, validator.Create); err != nil {
		err = pkg_err.TickValidate(err, "E1140085", pkg_err.ValidationFailed, role)
		return
	}

	role.Resources = strings.ReplaceAll(role.Resources, " ", "")

	if createdRole, err = s.Repo.Create(tx, role); err != nil {
		pkg_err.Log(err, "E1195741", "error in creating role", role)
		return
	}

	s.Engine.RedisCacheAPI.Delete(base_term.Roles)

	return
}

// Save a role, if it is exists update it, if not create it
func (s *BaseRoleServ) Save(tx tx.Tx, role base_model.Role) (updatedRole, roleBefore base_model.Role, err error) {

	if err = validator.ValidateModel(role, base_term.Role, validator.Update); err != nil {
		err = pkg_err.TickValidate(err, "E1180267", pkg_err.ValidationFailed, role)
		return
	}

	if roleBefore, err = s.FindByID(tx, role.ID); err != nil {
		pkg_err.Log(err, "E1131161", "can't fetch role by id for saving it", role.ID)
		return
	}

	role.Resources = strings.ReplaceAll(role.Resources, " ", "")

	if updatedRole, err = s.Repo.Save(tx, role); err != nil {
		pkg_err.Log(err, "E1189853", "role not saved")
		return
	}

	key := fmt.Sprintf("%v-%v", base_term.Role, updatedRole.ID)
	if err = s.Engine.RedisCacheAPI.Delete(key); err != nil {
		return
	}

	s.Engine.RedisCacheAPI.Delete(base_term.Roles)

	return
}

// Delete role, it is soft delete
func (s *BaseRoleServ) Delete(tx tx.Tx, id uint) (role base_model.Role, err error) {

	if role, err = s.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1155143", "role not found for deleting")
		return
	}

	if err = s.Repo.Delete(tx, role); err != nil {
		pkg_err.Log(err, "E1168701", "role not deleted")
		return
	}

	key := fmt.Sprintf("%v-%v", base_term.Role, role.ID)
	s.Engine.RedisCacheAPI.Delete(key)
	s.Engine.RedisCacheAPI.Delete(base_term.Roles)

	return
}
