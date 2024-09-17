package service

import (
	"GoRestify/domain/acc/acc_model"
	"GoRestify/domain/acc/acc_repo"
	"GoRestify/domain/acc/acc_term"
	"GoRestify/internal/core"
	"fmt"

	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/tx"

	"GoRestify/pkg/validator"
)

// AccSlotServ for injecting  acc_repo
type AccSlotServ struct {
	Repo   acc_repo.SlotRepo
	Engine *core.Engine
}

// ProvideAccSlotService for slot is used in wire
func ProvideAccSlotService(slotRepo acc_repo.SlotRepo) AccSlotServ {
	return AccSlotServ{
		Repo:   slotRepo,
		Engine: slotRepo.Engine,
	}
}

// FindByID for getting slot by its id
func (s *AccSlotServ) FindByID(tx tx.Tx, id uint) (slot acc_model.Slot, err error) {

	key := fmt.Sprintf("%v-%v", acc_term.Slot, id)
	if ok := s.Engine.RedisCacheAPI.GetCache(tx, key, &slot); ok {
		return
	}

	if slot, err = s.Repo.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1112195", "can't fetch the slot", id)
		return
	}

	err = s.Engine.RedisCacheAPI.Set(key, slot)

	return
}

// GetAll of slots, it supports pagination and search and return count
func (s *AccSlotServ) GetAll(params param.Param) (slots []acc_model.Slot, err error) {

	if ok := s.Engine.RedisCacheAPI.GetCache(params.Tx, acc_term.Slots, &slots); ok {
		return
	}

	params.Limit = 100000
	if slots, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in slots list")
		return
	}

	err = s.Engine.RedisCacheAPI.Set(acc_term.Slots, slots)

	return
}

// List of slots, it supports pagination and search and return count
func (s *AccSlotServ) List(params param.Param) (slots []acc_model.Slot,
	count int64, err error) {

	if slots, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in slots list")
		return
	}

	if count, err = s.Repo.Count(params); err != nil {
		pkg_log.CheckError(err, "error in slots count")
	}

	return
}

// Create a slot
func (s *AccSlotServ) Create(tx tx.Tx, slot acc_model.Slot) (createdSlot acc_model.Slot, err error) {

	if err = validator.ValidateModel(slot, acc_term.Slot, validator.Create); err != nil {
		err = pkg_err.TickValidate(err, "E1134338", pkg_err.ValidationFailed, slot)
		return
	}

	if createdSlot, err = s.Repo.Create(tx, slot); err != nil {
		pkg_err.Log(err, "E1128262", "error in creating slot", slot)
		return
	}

	s.Engine.RedisCacheAPI.Delete(acc_term.Slots)

	return
}

// Save a slot, if it is exists update it, if not create it
func (s *AccSlotServ) Save(tx tx.Tx, slot acc_model.Slot) (updatedSlot, slotBefore acc_model.Slot, err error) {

	if err = validator.ValidateModel(slot, acc_term.Slot, validator.Update); err != nil {
		err = pkg_err.TickValidate(err, "E1130667", pkg_err.ValidationFailed, slot)
		return
	}

	if slotBefore, err = s.FindByID(tx, slot.ID); err != nil {
		pkg_err.Log(err, "E1126558", "can't fetch slot by id for saving it", slot.ID)
		return
	}

	if updatedSlot, err = s.Repo.Save(tx, slot); err != nil {
		pkg_err.Log(err, "E1132099", "slot not saved")
		return
	}

	key := fmt.Sprintf("%v-%v", acc_term.Slot, updatedSlot.ID)
	if err = s.Engine.RedisCacheAPI.Delete(key); err != nil {
		return
	}

	s.Engine.RedisCacheAPI.Delete(acc_term.Slots)

	return
}

// Delete slot, it is soft delete
func (s *AccSlotServ) Delete(tx tx.Tx, id uint) (slot acc_model.Slot, err error) {

	if slot, err = s.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1153236", "slot not found for deleting")
		return
	}

	if err = s.Repo.Delete(tx, slot); err != nil {
		pkg_err.Log(err, "E1192268", "slot not deleted")
		return
	}

	key := fmt.Sprintf("%v-%v", acc_term.Slot, slot.ID)
	s.Engine.RedisCacheAPI.Delete(key)
	s.Engine.RedisCacheAPI.Delete(acc_term.Slots)

	return
}
