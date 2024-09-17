package acc_api

import (
	"GoRestify/domain/acc"
	"GoRestify/domain/acc/acc_model"
	"GoRestify/domain/acc/acc_term"
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

// SlotAPI for injecting slot service
type SlotAPI struct {
	Service service.AccSlotServ
	Engine  *core.Engine
}

// ProvideSlotAPI .
func ProvideSlotAPI(c service.AccSlotServ) SlotAPI {
	return SlotAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch an slot by its id
func (a *SlotAPI) FindByID(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.SlotTable)
	var err error
	var slot acc_model.Slot
	var id uint

	if id, err = resp.GetID(c.Param("slotID"), "E1135221", acc_term.Slot); err != nil {
		return
	}

	if slot, err = a.Service.FindByID(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(acc.ViewSlot)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VInfo, acc_term.Slot).
		JSON(slot)
}

// GetAll list of slots
func (a *SlotAPI) GetAll(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.SlotTable)
	var slots []acc_model.Slot
	var err error

	if slots, err = a.Service.GetAll(params); err != nil {
		err = pkg_err.Take(err, "E1127268").Message(pkg_err.SomethingWentWrong).Build()
		resp.Error(err).JSON()
		return
	}

	resp.Record(acc.ListSlot)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, acc_term.Slots).
		JSON(slots)
}

// List of slots
func (a *SlotAPI) List(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.SlotTable)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = a.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(acc.ListSlot)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, acc_term.Slots).
		JSON(data)
}

// Create slot
func (a *SlotAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.SlotTable)
	var slot, createdSlot acc_model.Slot
	var err error

	if err = resp.Bind(&slot, "E1111986", acc_term.Slot); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"slot"), "rollback recover create slot")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1199756").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if createdSlot, err = a.Service.Create(params.Tx, slot); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(acc.CreateSlot, slot, createdSlot)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VCreatedSuccessfully, acc_term.Slot).
		JSON(createdSlot)
}

// Update slot
func (a *SlotAPI) Update(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.SlotTable)
	var err error
	var slot, slotBefore, slotUpdated acc_model.Slot

	if err = resp.Bind(&slot, "E1115577", acc_term.Slot); err != nil {
		return
	}

	if slot.ID, err = resp.GetID(c.Param("slotID"), "E1124539", acc_term.Slot); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"slot"), "rollback recover create slot")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1129224").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if slotUpdated, slotBefore, err = a.Service.Save(params.Tx, slot); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(acc.UpdateSlot, slotBefore, slot)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VUpdatedSuccessfully, acc_term.Slot).
		JSON(slotUpdated)
}

// Delete slot
func (a *SlotAPI) Delete(c *gin.Context) {
	resp, params := response.NewParam(c, acc_model.SlotTable)
	var err error
	var slot acc_model.Slot
	var id uint

	if id, err = resp.GetID(c.Param("slotID"), "E1171272", acc_term.Slot); err != nil {
		return
	}

	// start tx
	params.Tx.DB = a.Engine.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			pkg_log.LogError(fmt.Errorf("panic happened in tx mode for %v",
				"slot"), "rollback recover create slot")
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1127820").
				Message(pkg_err.SomethingWentWrong).Build()
			// rollback tx
			params.Tx.DB.Rollback()
			return
		}
	}()

	if slot, err = a.Service.Delete(params.Tx, id); err != nil {
		resp.Error(err).JSON()
		// rollback tx
		params.Tx.DB.Rollback()
		return
	}

	// commit tx
	params.Tx.DB.Commit()

	resp.Record(acc.DeleteSlot, slot)
	resp.Status(http.StatusOK).
		Message(pkg_terms.VDeletedSuccessfully, acc_term.Slot).
		JSON()
}
