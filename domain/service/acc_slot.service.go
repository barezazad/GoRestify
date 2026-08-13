package service

import (
	"fmt"

	"GoRestify/domain/acc/acc_model"
	"GoRestify/domain/acc/acc_repo"
	"GoRestify/domain/acc/acc_term"
	"GoRestify/internal/core"

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

// GetAll of slots — not Redis-cached (ledger grows without bound)
func (s *AccSlotServ) GetAll(params param.Param) (slots []acc_model.Slot, err error) {
	if slots, err = s.Repo.GetAll(); err != nil {
		pkg_log.CheckError(err, "error in slots list")
	}
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

// Create a slot — used by DoTransaction only.
func (s *AccSlotServ) Create(tx tx.Tx, slot acc_model.Slot) (createdSlot acc_model.Slot, err error) {

	if err = validator.ValidateModel(slot, acc_term.Slot, validator.Create); err != nil {
		err = pkg_err.TickValidate(err, "E1134338", pkg_err.ValidationFailed, slot)
		return
	}

	if createdSlot, err = s.Repo.Create(tx, slot); err != nil {
		pkg_err.Log(err, "E1128262", "error in creating slot", slot)
		return
	}

	return
}

// Save is disabled — slots are immutable ledger lines.
func (s *AccSlotServ) Save(tx tx.Tx, slot acc_model.Slot) (updatedSlot, slotBefore acc_model.Slot, err error) {
	err = pkg_err.New(acc_term.LedgerRecordsAreImmutable, "E1201030").
		Message(acc_term.LedgerRecordsAreImmutable).Custom(pkg_err.ValidationFailedErr).Build()
	return
}

// Delete is disabled — slots are immutable ledger lines.
func (s *AccSlotServ) Delete(tx tx.Tx, id uint) (slot acc_model.Slot, err error) {
	err = pkg_err.New(acc_term.LedgerRecordsAreImmutable, "E1201031").
		Message(acc_term.LedgerRecordsAreImmutable).Custom(pkg_err.ValidationFailedErr).Build()
	return
}
