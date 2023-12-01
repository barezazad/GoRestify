package service

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_repo"
	"GoRestify/domain/base/base_term"
	"GoRestify/internal/core"
	"fmt"

	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_password"
	"GoRestify/pkg/tx"

	"GoRestify/pkg/validator"
)

// BaseUserServ for injecting  base_repo
type BaseUserServ struct {
	Repo   base_repo.UserRepo
	Engine *core.Engine
}

// ProvideBaseUserService for user is used in wire
func ProvideBaseUserService(userRepo base_repo.UserRepo) BaseUserServ {
	return BaseUserServ{
		Repo:   userRepo,
		Engine: userRepo.Engine,
	}
}

// FindByID for getting user by its id
func (s *BaseUserServ) FindByID(tx tx.Tx, id uint) (user base_model.User, err error) {

	key := fmt.Sprintf("%v-%v", base_term.User, id)
	if ok := s.Engine.RedisCacheAPI.GetCache(tx, key, &user); ok {
		return
	}

	if user, err = s.Repo.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1167635", "can't fetch the user", id)
		return
	}

	err = s.Engine.RedisCacheAPI.Set(key, user)

	return
}

// FindByUsername for getting user by its username
func (s *BaseUserServ) FindByUsername(tx tx.Tx, username string) (user base_model.User, err error) {

	if user, err = s.Repo.FindByUsername(tx, username); err != nil {
		pkg_err.Log(err, "E1167315", "can't fetch the user", username)
		return
	}

	return
}

// GetAll of users, it supports pagination and search and return count
func (s *BaseUserServ) GetAll(params param.Param) (users []base_model.User, err error) {

	if ok := s.Engine.RedisCacheAPI.GetCache(params.Tx, base_term.Users, &users); ok {
		return
	}

	params.Pagination.Limit = 1000000
	if users, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in users list")
		return
	}

	for i := range users {
		users[i].Password = ""
	}

	err = s.Engine.RedisCacheAPI.Set(base_term.Users, users)

	return
}

// List of users, it supports pagination and search and return count
func (s *BaseUserServ) List(params param.Param) (users []base_model.User,
	count int64, err error) {

	if users, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in users list")
		return
	}

	for i := range users {
		users[i].Password = ""
	}

	if count, err = s.Repo.Count(params); err != nil {
		pkg_log.CheckError(err, "error in users count")
	}

	return
}

// Create a user
func (s *BaseUserServ) Create(tx tx.Tx, user base_model.User) (createdUser base_model.User, err error) {

	if err = validator.ValidateModel(user, base_term.User, validator.Create); err != nil {
		err = pkg_err.TickValidate(err, "E1196830", pkg_err.ValidationFailed, user)
		return
	}

	if user.Password, err = pkg_password.Hash(user.Password, s.Engine.Envs[core.PasswordSalt]); err != nil {
		err = pkg_err.Log(err, "E1181138", "error in hashing password", user)
		return
	}

	if createdUser, err = s.Repo.Create(tx, user); err != nil {
		pkg_err.Log(err, "E1198014", "error in creating user", user)
		return
	}

	createdUser.Password = ""

	s.Engine.RedisCacheAPI.Delete(base_term.Users)

	return
}

// Save a user, if it is exists update it, if not create it
func (s *BaseUserServ) Save(tx tx.Tx, user base_model.User) (updatedUser, userBefore base_model.User, err error) {

	if err = validator.ValidateModel(user, base_term.User, validator.Update); err != nil {
		err = pkg_err.TickValidate(err, "E1191128", pkg_err.ValidationFailed, user)
		return
	}

	if userBefore, err = s.FindByID(tx, user.ID); err != nil {
		pkg_err.Log(err, "E1114861", "can't fetch user by id for saving it", user.ID)
		return
	}

	if user.Password != "" {
		if user.Password, err = pkg_password.Hash(user.Password, s.Engine.Envs[core.PasswordSalt]); err != nil {
			err = pkg_err.Log(err, "E1115211", "error in hashing password", user)
			return
		}
	} else {
		user.Password = userBefore.Password
	}

	if updatedUser, err = s.Repo.Save(tx, user); err != nil {
		pkg_err.Log(err, "E1141555", "user not saved")
		return
	}

	updatedUser.Password = ""

	key := fmt.Sprintf("%v-%v", base_term.User, updatedUser.ID)
	if err = s.Engine.RedisCacheAPI.Delete(key); err != nil {
		return
	}

	s.Engine.RedisCacheAPI.Delete(base_term.Users)

	return
}

// Delete user, it is soft delete
func (s *BaseUserServ) Delete(tx tx.Tx, id uint) (user base_model.User, err error) {

	if user, err = s.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1123745", "user not found for deleting")
		return
	}

	if err = s.Repo.Delete(tx, user); err != nil {
		pkg_err.Log(err, "E1188941", "user not deleted")
		return
	}

	key := fmt.Sprintf("%v-%v", base_term.User, user.ID)
	s.Engine.RedisCacheAPI.Delete(key)
	s.Engine.RedisCacheAPI.Delete(base_term.Users)

	return
}
