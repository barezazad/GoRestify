package base_api

import (
	"GoRestify/domain/base"
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_term"
	"GoRestify/domain/service"
	"GoRestify/internal/core"
	"fmt"
	"net/http"

	"GoRestify/pkg/excel"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_terms"
	"GoRestify/pkg/response"

	"github.com/gin-gonic/gin"
)

// CityAPI for injecting city service
type CityAPI struct {
	Service service.BaseCityServ
	Engine  *core.Engine
}

// ProvideCityAPI for city is used in wire
func ProvideCityAPI(c service.BaseCityServ) CityAPI {
	return CityAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a city by its id
func (a *CityAPI) FindByID(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.CityTable)
	var err error
	var city base_model.City
	var id uint

	if id, err = resp.GetID(c.Param("cityID"), "E1136122", base_term.City); err != nil {
		return
	}

	if city, err = a.Service.FindByID(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ViewCity)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VInfo, base_term.City).
		JSON(city)
}

// GetAll list of cities
func (a *CityAPI) GetAll(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.CityTable)
	var cities []base_model.City
	var err error

	if cities, err = a.Service.GetAll(params); err != nil {
		err = pkg_err.Take(err, "E1697901").Message(pkg_err.SomethingWentWrong).Build()
		resp.Error(err).JSON()
		return
	}

	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, base_term.Cities).
		JSON(cities)
}

// List of cities
func (a *CityAPI) List(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.CityTable)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = a.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListCity)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, base_term.Cities).
		JSON(data)
}

// Create city
func (a *CityAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.CityTable)
	var city, createdCity base_model.City
	var err error

	if err = resp.Bind(&city, "E1119335", base_term.City); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"city"), "rollback recover create city")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1142882").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if createdCity, err = a.Service.Create(params.Tx, city); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(base.CreateCity, city, createdCity)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VCreatedSuccessfully, base_term.City).
		JSON(createdCity)
}

// Update city
func (a *CityAPI) Update(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.CityTable)
	var err error
	var city, cityBefore, cityUpdated base_model.City

	if err = resp.Bind(&city, "E1152347", base_term.City); err != nil {
		return
	}

	if city.ID, err = resp.GetID(c.Param("cityID"), "E1136317", base_term.City); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"city"), "rollback recover create city")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1161572").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if cityUpdated, cityBefore, err = a.Service.Save(params.Tx, city); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(base.UpdateCity, cityBefore, city)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VUpdatedSuccessfully, base_term.City).
		JSON(cityUpdated)
}

// Delete city
func (a *CityAPI) Delete(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.CityTable)
	var err error
	var city base_model.City
	var id uint

	if id, err = resp.GetID(c.Param("cityID"), "E1132128", base_term.City); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"city"), "rollback recover create city")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1120855").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if city, err = a.Service.Delete(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(base.DeleteCity, city)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VDeletedSuccessfully, base_term.City).
		JSON()
}

// Excel generate Excel files based on search
func (a *CityAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.RegionTable)
	var err error

	cities, err := a.Service.Repo.List(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("cities")
	ex.AddSheet("Cities").
		AddSheet("Summary").
		Active("Cities").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "D", 15.3).
		SetColWidth("C", "C", 80).
		SetColWidth("D", "F", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Cities").
		WriteHeader("ID", "Name", "Created At").
		SetSheetFields("ID", "Name", "CreatedAt").
		WriteData(cities).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=regions-"+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
