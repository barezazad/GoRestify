package service

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_term"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_jwt"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_password"
	"GoRestify/pkg/pkg_types"
	"GoRestify/pkg/tx"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"

	"GoRestify/internal/core"
	"GoRestify/pkg/validator"
)

// BaseAuthServ for injecting base_repo
type BaseAuthServ struct {
	Engine *core.Engine
}

// ProvideBaseAuthService for auth is used in wire
func ProvideBaseAuthService(engine *core.Engine) BaseAuthServ {
	return BaseAuthServ{
		Engine: engine,
	}
}

// Login User
func (s *BaseAuthServ) Login(tx tx.Tx, auth base_model.Auth) (account base_model.Account, err error) {

	if err = validator.ValidateModel(auth, "login", "login"); err != nil {
		err = pkg_err.TickValidate(err, "E1062875", "validation failed username and password is required", auth)
		return
	}

	account, findErr := BaseAccountService.FindByUsername(tx, auth.Username)
	if findErr != nil {
		pkg_log.CheckError(findErr, "login failed: account lookup")
		err = invalidLoginErr("E1162611")
		return
	}

	var user base_model.User
	if user, findErr = BaseUserService.FindByID(tx, account.ID); findErr != nil {
		pkg_log.CheckError(findErr, "login failed: user lookup")
		err = invalidLoginErr("E1153521")
		return
	}

	var role base_model.Role
	if role, findErr = BaseRoleService.FindByID(tx, user.RoleID); findErr != nil {
		pkg_log.CheckError(findErr, "login failed: role lookup")
		err = invalidLoginErr("E1140158")
		return
	}

	if !pkg_password.Verify(auth.Password, account.Password, s.Engine.Envs[core.PasswordSalt]) {
		err = invalidLoginErr("E1169512")
		return
	}

	pair, err := pkg_jwt.GenerateTokenPair(user.ID, account.Username)
	if err != nil {
		err = pkg_err.Log(err, "E1147810", "error generating token")
		return
	}

	if err = s.storeRefreshToken(pair.RefreshJTI, user.ID); err != nil {
		err = pkg_err.Log(err, "E1147811", "error storing refresh token")
		return
	}

	account.Token = pair.AccessToken
	account.RefreshToken = pair.RefreshToken
	account.ExpiresIn = pair.ExpiresIn
	account.RefreshExpiresIn = pair.RefreshExpiresIn
	account.Resources = strings.Split(role.Resources, ",")
	account.Password = ""
	account.Type = ""

	key := fmt.Sprintf("%v-%v", base_term.Auth, user.ID)
	s.Engine.RedisCacheAPI.Delete(key)

	return
}

func invalidLoginErr(code string) error {
	return pkg_err.New(base_term.UsernameOrPasswordIsWrong, code).
		Custom(pkg_err.UnauthorizedErr).
		Message(base_term.UsernameOrPasswordIsWrong).
		Build()
}

// RefreshToken validates a refresh token and issues a new token pair (rotation).
func (s *BaseAuthServ) RefreshToken(refreshToken string) (account base_model.Account, err error) {
	claims, err := pkg_jwt.ParseAndValidate(refreshToken)
	if err != nil || claims.TokenType != pkg_jwt.TokenTypeRefresh {
		err = pkg_err.New(pkg_err.TokenIsNotValid, "E1147820").
			Custom(pkg_err.UnauthorizedErr).
			Message(pkg_err.TokenIsNotValid).Build()
		return
	}

	redisKey := pkg_jwt.RefreshRedisKey(claims.ID)
	stored := s.Engine.RedisCacheAPI.Get(redisKey)
	if stored == "" || stored != fmt.Sprint(claims.UserID) {
		err = pkg_err.New(pkg_err.TokenIsNotValid, "E1147821").
			Custom(pkg_err.UnauthorizedErr).
			Message(pkg_err.TokenIsNotValid).Build()
		return
	}

	// rotate: invalidate old refresh token
	s.Engine.RedisCacheAPI.Delete(redisKey)

	if account, err = BaseAccountService.FindByID(tx.Tx{}, claims.UserID); err != nil {
		err = pkg_err.Log(err, "E1147822", "can't fetch account for refresh")
		return
	}

	var user base_model.User
	if user, err = BaseUserService.FindByID(tx.Tx{}, account.ID); err != nil {
		err = pkg_err.Log(err, "E1147823", "can't fetch user for refresh")
		return
	}

	var role base_model.Role
	if role, err = BaseRoleService.FindByID(tx.Tx{}, user.RoleID); err != nil {
		err = pkg_err.Log(err, "E1147824", "can't fetch role for refresh")
		return
	}

	pair, err := pkg_jwt.GenerateTokenPair(user.ID, account.Username)
	if err != nil {
		err = pkg_err.Log(err, "E1147825", "error generating refreshed tokens")
		return
	}

	if err = s.storeRefreshToken(pair.RefreshJTI, user.ID); err != nil {
		err = pkg_err.Log(err, "E1147826", "error storing refreshed token")
		return
	}

	account.Token = pair.AccessToken
	account.RefreshToken = pair.RefreshToken
	account.ExpiresIn = pair.ExpiresIn
	account.RefreshExpiresIn = pair.RefreshExpiresIn
	account.Resources = strings.Split(role.Resources, ",")
	account.Password = ""
	account.Type = ""
	return
}

func (s *BaseAuthServ) storeRefreshToken(jti string, userID uint) error {
	ttl := int(pkg_jwt.RefreshTTL().Seconds())
	return s.Engine.RedisCacheAPI.SetWithTTL(pkg_jwt.RefreshRedisKey(jti), fmt.Sprint(userID), ttl)
}

// CheckAccess is used in middleware to find if user has permission or not
func (s *BaseAuthServ) CheckAccess(userID uint, resource pkg_types.Resource) (isAllow bool) {
	var err error
	var resourceList base_model.ResourceList

	key := fmt.Sprintf("%v-%v", base_term.Auth, userID)
	if ok := s.Engine.RedisCacheAPI.GetCache(tx.Tx{}, key, &resourceList); ok {
		return slices.Contains(resourceList.ResourcesArray, resource.String())
	}

	if resourceList, err = BaseUserService.Repo.GetUserResources(userID); err != nil {
		pkg_log.CheckError(err, "error in finding the resources for user", userID)
		return
	}

	s.Engine.RedisCacheAPI.Set(key, &resourceList)

	return slices.Contains(resourceList.ResourcesArray, resource.String())
}
