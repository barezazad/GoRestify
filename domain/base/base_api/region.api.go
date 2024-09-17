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

// RegionAPI for injecting region service
type RegionAPI struct {
	Service service.BaseRegionServ
	Engine  *core.Engine
}

// ProvideRegionAPI .
func ProvideRegionAPI(c service.BaseRegionServ) RegionAPI {
	return RegionAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch an region by its id
func (a *RegionAPI) FindByID(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.RegionTable)
	var err error
	var region base_model.Region
	var id uint

	if id, err = resp.GetID(c.Param("regionID"), "E1136276", base_term.Region); err != nil {
		return
	}

	if region, err = a.Service.FindByID(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ViewRegion)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VInfo, base_term.Region).
		JSON(region)
}

// GetAll list of regions
func (a *RegionAPI) GetAll(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.RegionTable)
	var regions []base_model.Region
	var err error

	if regions, err = a.Service.GetAll(params); err != nil {
		err = pkg_err.Take(err, "E1198115").Message(pkg_err.SomethingWentWrong).Build()
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListRegion)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, base_term.Regions).
		JSON(regions)
}

// List of regions
func (a *RegionAPI) List(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.RegionTable)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = a.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListRegion)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, base_term.Regions).
		JSON(data)
}

// Create region
func (a *RegionAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.RegionTable)
	var region, createdRegion base_model.Region
	var err error

	if err = resp.Bind(&region, "E1126350", base_term.Region); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"region"), "rollback recover create region")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1157983").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if createdRegion, err = a.Service.Create(params.Tx, region); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(base.CreateRegion, region, createdRegion)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VCreatedSuccessfully, base_term.Region).
		JSON(createdRegion)
}

// Update region
func (a *RegionAPI) Update(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.RegionTable)
	var err error
	var region, regionBefore, regionUpdated base_model.Region

	if err = resp.Bind(&region, "E1166586", base_term.Region); err != nil {
		return
	}

	if region.ID, err = resp.GetID(c.Param("regionID"), "E1117252", base_term.Region); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"region"), "rollback recover create region")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1191752").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if regionUpdated, regionBefore, err = a.Service.Save(params.Tx, region); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(base.UpdateRegion, regionBefore, region)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VUpdatedSuccessfully, base_term.Region).
		JSON(regionUpdated)
}

// Delete region
func (a *RegionAPI) Delete(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.RegionTable)
	var err error
	var region base_model.Region
	var id uint

	if id, err = resp.GetID(c.Param("regionID"), "E1126058", base_term.Region); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"region"), "rollback recover create region")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1138261").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if region, err = a.Service.Delete(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(base.DeleteRegion, region)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VDeletedSuccessfully, base_term.Region).
		JSON()
}
