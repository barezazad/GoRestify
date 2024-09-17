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

// BaseRegionServ for injecting  base_repo
type BaseRegionServ struct {
	Repo   base_repo.RegionRepo
	Engine *core.Engine
}

// ProvideBaseRegionService for region is used in wire
func ProvideBaseRegionService(regionRepo base_repo.RegionRepo) BaseRegionServ {
	return BaseRegionServ{
		Repo:   regionRepo,
		Engine: regionRepo.Engine,
	}
}

// FindByID for getting region by its id
func (s *BaseRegionServ) FindByID(tx tx.Tx, id uint) (region base_model.Region, err error) {

	key := fmt.Sprintf("%v-%v", base_term.Region, id)
	if ok := s.Engine.RedisCacheAPI.GetCache(tx, key, &region); ok {
		return
	}

	if region, err = s.Repo.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1189966", "can't fetch the region", id)
		return
	}

	err = s.Engine.RedisCacheAPI.Set(key, region)

	return
}

// GetAll of regions, it supports pagination and search and return count
func (s *BaseRegionServ) GetAll(params param.Param) (regions []base_model.Region, err error) {

	if ok := s.Engine.RedisCacheAPI.GetCache(params.Tx, base_term.Regions, &regions); ok {
		return
	}

	params.Limit = 100000
	if regions, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in regions list")
		return
	}

	err = s.Engine.RedisCacheAPI.Set(base_term.Regions, regions)

	return
}

// List of regions, it supports pagination and search and return count
func (s *BaseRegionServ) List(params param.Param) (regions []base_model.Region,
	count int64, err error) {

	if regions, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in regions list")
		return
	}

	if count, err = s.Repo.Count(params); err != nil {
		pkg_log.CheckError(err, "error in regions count")
	}

	return
}

// Create a region
func (s *BaseRegionServ) Create(tx tx.Tx, region base_model.Region) (createdRegion base_model.Region, err error) {

	if err = validator.ValidateModel(region, base_term.Region, validator.Create); err != nil {
		err = pkg_err.TickValidate(err, "E1118046", pkg_err.ValidationFailed, region)
		return
	}

	if createdRegion, err = s.Repo.Create(tx, region); err != nil {
		pkg_err.Log(err, "E1121705", "error in creating region", region)
		return
	}

	s.Engine.RedisCacheAPI.Delete(base_term.Regions)

	return
}

// Save a region, if it is exists update it, if not create it
func (s *BaseRegionServ) Save(tx tx.Tx, region base_model.Region) (updatedRegion, regionBefore base_model.Region, err error) {

	if err = validator.ValidateModel(region, base_term.Region, validator.Update); err != nil {
		err = pkg_err.TickValidate(err, "E1197976", pkg_err.ValidationFailed, region)
		return
	}

	if regionBefore, err = s.FindByID(tx, region.ID); err != nil {
		pkg_err.Log(err, "E1111888", "can't fetch region by id for saving it", region.ID)
		return
	}

	if updatedRegion, err = s.Repo.Save(tx, region); err != nil {
		pkg_err.Log(err, "E1139209", "region not saved")
		return
	}

	key := fmt.Sprintf("%v-%v", base_term.Region, updatedRegion.ID)
	if err = s.Engine.RedisCacheAPI.Delete(key); err != nil {
		return
	}

	s.Engine.RedisCacheAPI.Delete(base_term.Regions)

	return
}

// Delete region, it is soft delete
func (s *BaseRegionServ) Delete(tx tx.Tx, id uint) (region base_model.Region, err error) {

	if region, err = s.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1114367", "region not found for deleting")
		return
	}

	if err = s.Repo.Delete(tx, region); err != nil {
		pkg_err.Log(err, "E1162176", "region not deleted")
		return
	}

	key := fmt.Sprintf("%v-%v", base_term.Region, region.ID)
	s.Engine.RedisCacheAPI.Delete(key)
	s.Engine.RedisCacheAPI.Delete(base_term.Regions)

	return
}
