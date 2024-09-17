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
	"GoRestify/pkg/tx"
	"GoRestify/pkg/validator"
)

// BaseCityServ for injecting auth base_repo
type BaseCityServ struct {
	Repo   base_repo.CityRepo
	Engine *core.Engine
}

// ProvideBaseCityService for city is used in wire
func ProvideBaseCityService(cityRepo base_repo.CityRepo) BaseCityServ {
	return BaseCityServ{
		Repo:   cityRepo,
		Engine: cityRepo.Engine,
	}
}

// FindByID for getting city by its id
func (s *BaseCityServ) FindByID(tx tx.Tx, id uint) (city base_model.City, err error) {

	key := fmt.Sprintf("%v-%v", base_term.City, id)
	if ok := s.Engine.RedisCacheAPI.GetCache(tx, key, &city); ok {
		return
	}

	if city, err = s.Repo.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1117192", "can't fetch the city", id)
		return
	}

	err = s.Engine.RedisCacheAPI.Set(key, city)

	return
}

// GetAll of cities, it supports pagination and search and return count
func (s *BaseCityServ) GetAll(params param.Param) (cities []base_model.City, err error) {
	if ok := s.Engine.RedisCacheAPI.GetCache(params.Tx, base_term.Cities, &cities); ok {
		return
	}

	params.Limit = 100000
	if cities, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in cities list")
		return
	}

	err = s.Engine.RedisCacheAPI.Set(base_term.Cities, cities)

	return
}

// List of cities, it supports pagination and search and return count
func (s *BaseCityServ) List(params param.Param) (cities []base_model.City,
	count int64, err error) {

	if cities, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in cities list")
		return
	}

	if count, err = s.Repo.Count(params); err != nil {
		pkg_log.CheckError(err, "error in cities count")
	}

	return
}

// Create a city
func (s *BaseCityServ) Create(tx tx.Tx, city base_model.City) (createdCity base_model.City, err error) {
	if err = validator.ValidateModel(city, base_term.City, validator.Create); err != nil {
		err = pkg_err.TickValidate(err, "E1159448", pkg_err.ValidationFailed, city)
		return
	}

	if createdCity, err = s.Repo.Create(tx, city); err != nil {
		pkg_err.Log(err, "E1175150", "city not saved")
		return
	}

	key := fmt.Sprintf("%v-%v", base_term.City, createdCity.ID)
	if err = s.Engine.RedisCacheAPI.Set(key, createdCity); err != nil {
		return
	}

	s.Engine.RedisCacheAPI.Delete(base_term.Cities)

	return
}

// Save a city, if it is exists update it, if not create it
func (s *BaseCityServ) Save(tx tx.Tx, city base_model.City) (savedCity, cityBefore base_model.City, err error) {
	if err = validator.ValidateModel(city, base_term.City, validator.Update); err != nil {
		err = pkg_err.TickValidate(err, "E1147437", pkg_err.ValidationFailed, city)
		return
	}

	if cityBefore, err = s.FindByID(tx, city.ID); err != nil {
		pkg_err.Log(err, "E1178144", "can't fetch city by id for saving it", city.ID)
		return
	}

	if savedCity, err = s.Repo.Save(tx, city); err != nil {
		pkg_err.Log(err, "E1128646", "city not saved")
		return
	}

	key := fmt.Sprintf("%v-%v", base_term.City, savedCity.ID)
	if err = s.Engine.RedisCacheAPI.Set(key, savedCity); err != nil {
		return
	}

	s.Engine.RedisCacheAPI.Delete(base_term.Cities)

	return
}

// Delete city
func (s *BaseCityServ) Delete(tx tx.Tx, id uint) (city base_model.City, err error) {
	if city, err = s.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1185319", "city not found for deleting")
		return
	}

	if err = s.Repo.Delete(tx, city); err != nil {
		pkg_err.Log(err, "E1142739", "city not deleted")
		return
	}

	key := fmt.Sprintf("%v-%v", base_term.City, city.ID)
	s.Engine.RedisCacheAPI.Delete(key)
	s.Engine.RedisCacheAPI.Delete(base_term.Cities)

	return
}
