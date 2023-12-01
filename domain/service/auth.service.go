package service

import (
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_term"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_password"
	"GoRestify/pkg/pkg_types"
	"GoRestify/pkg/tx"
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"GoRestify/internal/core"
	"GoRestify/pkg/validator"

	"github.com/golang-jwt/jwt/v5"
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
func (s *BaseAuthServ) Login(tx tx.Tx, auth base_model.Auth) (user base_model.User, err error) {

	if err = validator.ValidateModel(auth, "login", "login"); err != nil {
		err = pkg_err.TickValidate(err, "E1062875", "validation failed username and password is required", auth)
		return
	}

	jwtKey := s.Engine.Envs.ToByte(core.JWTSecretKey)

	if user, err = BaseUserService.FindByUsername(tx, auth.Username); err != nil {
		err = pkg_err.Log(err, "E1130152", "can't fetch the user by username")
		return
	}

	var role base_model.Role
	if role, err = BaseRoleService.FindByID(tx, user.RoleID); err != nil {
		err = pkg_err.Log(err, "E1140158", "can't fetch the role by id")
		return
	}

	if !pkg_password.Verify(auth.Password, user.Password, s.Engine.Envs[core.PasswordSalt]) {
		err = pkg_err.New(base_term.UsernameOrPasswordIsWrong, "E1169512").
			Message(base_term.UsernameOrPasswordIsWrong).Build()
		return
	}

	expirationTime := time.Now().Add(12 * time.Hour)
	claims := &pkg_types.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	if user.Token, err = token.SignedString(jwtKey); err != nil {
		err = pkg_err.Log(err, "E1147810", "error generating token")
		return
	}

	user.Resources = strings.Split(role.Resources, ",")
	user.Password = ""
	user.Role = role.Name

	key := fmt.Sprintf("%v-%v", base_term.Auth, user.ID)
	s.Engine.RedisCacheAPI.Delete(key)

	return
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
