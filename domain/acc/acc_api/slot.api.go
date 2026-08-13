package acc_api

import (
	"net/http"

	"GoRestify/domain/acc"
	"GoRestify/domain/acc/acc_model"
	"GoRestify/domain/acc/acc_term"
	"GoRestify/domain/service"
	"GoRestify/internal/core"

	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_terms"
	"GoRestify/pkg/response"

	"github.com/gin-gonic/gin"
)

// SlotAPI for injecting slot service (read-only HTTP surface)
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
