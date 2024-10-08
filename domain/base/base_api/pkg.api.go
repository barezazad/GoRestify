package base_api

import (
	"GoRestify/domain/acc/acc_term"
	"GoRestify/domain/base"
	"GoRestify/domain/base/base_term"
	"GoRestify/internal/core"
	"fmt"
	"net/http"

	"GoRestify/pkg/activity"
	"GoRestify/pkg/models"
	"GoRestify/pkg/pkg_terms"
	"GoRestify/pkg/response"
	"GoRestify/pkg/setting"

	"github.com/gin-gonic/gin"
)

// PkgAPI for injecting setting service
type PkgAPI struct {
	Engine *core.Engine
}

// ProvidePkgAPI .
func ProvidePkgAPI(engine *core.Engine) PkgAPI {
	return PkgAPI{Engine: engine}
}

// ------------------------------------------ Settings APIs ------------------------------------------

// SettingList of settings
func (a *PkgAPI) SettingList(c *gin.Context) {
	resp, params := response.NewParam(c, models.SettingTable)

	data := make(map[string]interface{})
	var err error

	if data["list"], err = setting.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	if data["count"], err = setting.Count(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ViewSetting)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, pkg_terms.Settings).
		JSON(data)
}

// SettingUpdate setting
func (a *PkgAPI) SettingUpdate(c *gin.Context) {
	resp, _ := response.NewParam(c, models.SettingTable)
	var settingPayload, settingUpdated models.Setting
	var err error

	if err = resp.Bind(&settingPayload, "E1168843", pkg_terms.Settings); err != nil {
		return
	}

	if settingPayload.ID, err = resp.GetID(c.Param("settingID"), "E1129420", pkg_terms.Settings); err != nil {
		return
	}

	if settingUpdated, err = setting.Save(settingPayload); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.UpdateSetting, settingPayload, settingUpdated)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VUpdatedSuccessfully, pkg_terms.Settings).
		JSON(settingUpdated)
}

// ------------------------------------------ Activities APIs ------------------------------------------

// ActivitiesList of activities
func (a *PkgAPI) ActivitiesList(c *gin.Context) {
	resp, params := response.NewParam(c, models.ActivityTable)

	data := make(map[string]interface{})
	var err error

	if data["list"], err = activity.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	if data["count"], err = activity.Count(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, pkg_terms.Activities).
		JSON(data)
}

// ------------------------------------------ Redis Cache APIs ------------------------------------------

// RedisResetCacheByKey reset redis cache api
func (a *PkgAPI) RedisResetCacheByKey(c *gin.Context) {
	resp, _ := response.NewParam(c, models.SettingTable)

	key := c.Param("key")
	var keyPattern string

	switch key {

	case base_term.Region:
		keyPattern = fmt.Sprintf("%v-*", base_term.Region)
		a.Engine.RedisCacheAPI.Delete(base_term.Regions)

	case base_term.City:
		keyPattern = fmt.Sprintf("%v-*", base_term.City)
		a.Engine.RedisCacheAPI.Delete(base_term.Cities)

	case base_term.Role:
		keyPattern = fmt.Sprintf("%v-*", base_term.Role)
		a.Engine.RedisCacheAPI.Delete(base_term.Roles)

	case base_term.User:
		keyPattern = fmt.Sprintf("%v-*", base_term.User)
		a.Engine.RedisCacheAPI.Delete(base_term.Users)

	case base_term.Account:
		keyPattern = fmt.Sprintf("%v-*", base_term.Account)
		a.Engine.RedisCacheAPI.Delete(base_term.Accounts)

	case acc_term.Currency:
		keyPattern = fmt.Sprintf("%v-*", acc_term.Currency)
		a.Engine.RedisCacheAPI.Delete(acc_term.Currencies)

	case acc_term.Slot:
		keyPattern = fmt.Sprintf("%v-*", acc_term.Slot)
		a.Engine.RedisCacheAPI.Delete(acc_term.Slots)

	case acc_term.Transaction:
		keyPattern = fmt.Sprintf("%v-*", acc_term.Transaction)
		a.Engine.RedisCacheAPI.Delete(acc_term.Transactions)

	case acc_term.AccountCredit:
		keyPattern = fmt.Sprintf("%v-*", acc_term.AccountCredit)
		a.Engine.RedisCacheAPI.Delete(acc_term.AccountCredits)

	}

	if keyPattern != "" {
		a.Engine.RedisCacheAPI.ResetCacheByKeyPatten(keyPattern)
	}

	resp.Record(base.ClearCache, key)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VInCacheHasBeenCleared, key).
		JSON(key)
}

// RedisClearCacheToUser reset redis cache api
func (a *PkgAPI) RedisClearCacheToUser(c *gin.Context) {
	resp, _ := response.NewParam(c, models.SettingTable)

	userID := c.Param("userID")

	// its just example
	a.Engine.RedisCacheAPI.Delete(fmt.Sprintf("%v-%v", base_term.User, userID))
	a.Engine.RedisCacheAPI.Delete(fmt.Sprintf("%v-%v", base_term.Account, userID))

	resp.Record(base.ClearCacheUser, userID)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VInCacheHasBeenCleared, userID).
		JSON(userID)
}
