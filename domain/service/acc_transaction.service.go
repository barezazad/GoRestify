package service

import (
	"fmt"
	"strings"

	"GoRestify/domain/acc/acc_model"
	"GoRestify/domain/acc/acc_repo"
	"GoRestify/domain/acc/acc_term"
	"GoRestify/internal/core"

	"GoRestify/pkg/decimal"
	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/tx"
	"GoRestify/pkg/validator"

	"github.com/google/uuid"
)

// AccTransactionServ for injecting  acc_repo
type AccTransactionServ struct {
	Repo   acc_repo.TransactionRepo
	Engine *core.Engine
}

// ProvideAccTransactionService for transaction is used in wire
func ProvideAccTransactionService(transactionRepo acc_repo.TransactionRepo) AccTransactionServ {
	return AccTransactionServ{
		Repo:   transactionRepo,
		Engine: transactionRepo.Engine,
	}
}

// FindByID for getting transaction by its id (includes slots)
func (s *AccTransactionServ) FindByID(tx tx.Tx, id uint) (transaction acc_model.Transaction, err error) {

	key := fmt.Sprintf("%v-%v", acc_term.Transaction, id)
	if ok := s.Engine.RedisCacheAPI.GetCache(tx, key, &transaction); ok {
		return
	}

	if transaction, err = s.Repo.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1155388", "can't fetch the transaction", id)
		return
	}

	if transaction.Slots, err = AccSlotService.Repo.ListByTransactionID(tx, transaction.ID); err != nil {
		pkg_err.Log(err, "E1155389", "can't fetch transaction slots", id)
		return
	}

	err = s.Engine.RedisCacheAPI.Set(key, transaction)

	return
}

// GetAll of transactions — not Redis-cached (ledger grows without bound)
func (s *AccTransactionServ) GetAll(params param.Param) (transactions []acc_model.Transaction, err error) {
	if transactions, err = s.Repo.GetAll(); err != nil {
		pkg_log.CheckError(err, "error in transactions list")
	}
	return
}

// List of transactions, it supports pagination and search and return count
func (s *AccTransactionServ) List(params param.Param) (transactions []acc_model.Transaction,
	count int64, err error) {

	if transactions, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in transactions list")
		return
	}

	if count, err = s.Repo.Count(params); err != nil {
		pkg_log.CheckError(err, "error in transactions count")
	}

	return
}

// DoTransaction posts a balanced money movement. Requires an open DB transaction (tx.DB).
// Optional Hash is an idempotency key; if already posted, returns the existing transaction.
func (s *AccTransactionServ) DoTransaction(txSession tx.Tx, transaction acc_model.Transaction) (createdTransaction acc_model.Transaction, err error) {

	if txSession.DB == nil {
		err = pkg_err.New(acc_term.TransactionDBTxRequired, "E1201001").
			Message(acc_term.TransactionDBTxRequired).Custom(pkg_err.BadRequestErr).Build()
		return
	}

	if err = validator.ValidateModel(transaction, acc_term.Transaction, validator.Create); err != nil {
		err = pkg_err.TickValidate(err, "E1163777", pkg_err.ValidationFailed, transaction)
		return
	}

	if transaction.Fee == nil {
		transaction.Fee = decimal.NewFromInt(0).Build()
	}

	if err = validateTransactionAmounts(transaction); err != nil {
		return
	}

	inputSlots := transaction.Slots
	if len(inputSlots) == 0 {
		err = pkg_err.New(acc_term.SlotsRequired, "E1201002").
			Message(acc_term.SlotsRequired).Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	if err = validateSlots(transaction, inputSlots); err != nil {
		return
	}

	// Idempotency: client-supplied hash returns existing post
	hash := strings.TrimSpace(transaction.Hash)
	if hash != "" {
		var existing acc_model.Transaction
		if existing, err = s.Repo.FindByHash(txSession, hash); err == nil {
			existing.Slots, _ = AccSlotService.Repo.ListByTransactionID(txSession, existing.ID)
			return existing, nil
		}
		err = nil
	} else {
		hash = uuid.New().String()
	}
	transaction.Hash = hash

	transaction.NetAmount = transaction.Amount.Num().Sub(transaction.Fee.Num()).Build()
	transaction.Slots = nil

	if createdTransaction, err = s.Repo.Create(txSession, transaction); err != nil {
		pkg_err.Log(err, "E1126470", "error in creating transaction", transaction)
		return
	}

	createdTransaction.Slots = nil
	for _, v := range inputSlots {
		var accountCredit acc_model.AccountCredit
		txSession.IsLock = true
		if accountCredit, err = AccAccountCreditService.Ensure(txSession, v.AccountID, v.CurrencyID); err != nil {
			err = pkg_err.Log(err, "E1118793", "account credit not found or created", v.AccountID, v.CurrencyID)
			return
		}
		txSession.IsLock = false

		if !v.Credit.Num().IsZero() {
			if accountCredit.Balance.Num().LessThan(v.Credit.Num()) {
				err = pkg_err.New(acc_term.YouHaveNotEnoughCredit, "E1099045", accountCredit.AccountID).
					Message(acc_term.YouHaveNotEnoughCredit).Custom(pkg_err.BadRequestErr).Build()
				return
			}
		}

		remainBalance := accountCredit.Balance.Num().Add(v.Debit.Num()).Sub(v.Credit.Num()).Build()
		accountCredit.Balance = remainBalance

		if _, _, err = AccAccountCreditService.Save(txSession, accountCredit); err != nil {
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1167734", "error in updating credit account", accountCredit).
				Message(pkg_err.SomethingWentWrong).
				Custom(pkg_err.InternalServerErr).Build()
			return
		}

		slot := acc_model.Slot{
			TransactionID: createdTransaction.ID,
			AccountID:     v.AccountID,
			CurrencyID:    v.CurrencyID,
			Credit:        v.Credit,
			Debit:         v.Debit,
			Balance:       remainBalance,
			Notes:         v.Notes,
		}

		if slot, err = AccSlotService.Create(txSession, slot); err != nil {
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1159619", "error in creating slot", slot).
				Message(pkg_err.SomethingWentWrong).
				Custom(pkg_err.InternalServerErr).Build()
			return
		}

		createdTransaction.Slots = append(createdTransaction.Slots, slot)
	}

	key := fmt.Sprintf("%v-%v", acc_term.Transaction, createdTransaction.ID)
	_ = s.Engine.RedisCacheAPI.Delete(key)
	_ = s.Engine.RedisCacheAPI.Delete(acc_term.Transactions)

	return
}

func validateTransactionAmounts(transaction acc_model.Transaction) error {
	if transaction.Amount == nil {
		return pkg_err.New(acc_term.InvalidTransactionAmounts, "E1201003").
			Message(acc_term.InvalidTransactionAmounts).Custom(pkg_err.ValidationFailedErr).Build()
	}
	amount := transaction.Amount.Num()
	fee := transaction.Fee.Num()
	if !amount.IsPositive() || fee.IsNegative() || fee.GreaterThan(amount) {
		return pkg_err.New(acc_term.InvalidTransactionAmounts, "E1201015").
			Message(acc_term.InvalidTransactionAmounts).Custom(pkg_err.ValidationFailedErr).Build()
	}
	return nil
}

func validateSlots(transaction acc_model.Transaction, slots []acc_model.Slot) error {
	sumDebit := decimal.NewFromInt(0)
	sumCredit := decimal.NewFromInt(0)

	for i := range slots {
		v := &slots[i]
		if v.AccountID == 0 {
			return pkg_err.New(acc_term.SlotAccountNotAllowed, "E1201004", i).
				Message(acc_term.SlotAccountNotAllowed).Custom(pkg_err.ValidationFailedErr).Build()
		}
		if v.AccountID != transaction.SenderID && v.AccountID != transaction.ReceiverID {
			return pkg_err.New(acc_term.SlotAccountNotAllowed, "E1201005", v.AccountID).
				Message(acc_term.SlotAccountNotAllowed).Custom(pkg_err.ValidationFailedErr).Build()
		}
		if v.CurrencyID != transaction.CurrencyID {
			return pkg_err.New(acc_term.SlotCurrencyMismatch, "E1201006", v.CurrencyID).
				Message(acc_term.SlotCurrencyMismatch).Custom(pkg_err.ValidationFailedErr).Build()
		}

		if v.Debit == nil {
			v.Debit = decimal.NewFromInt(0).Build()
		}
		if v.Credit == nil {
			v.Credit = decimal.NewFromInt(0).Build()
		}
		debit := v.Debit.Num()
		credit := v.Credit.Num()
		if debit.IsNegative() || credit.IsNegative() {
			return pkg_err.New(acc_term.InvalidSlotAmounts, "E1201007", i).
				Message(acc_term.InvalidSlotAmounts).Custom(pkg_err.ValidationFailedErr).Build()
		}
		if debit.IsZero() && credit.IsZero() {
			return pkg_err.New(acc_term.InvalidSlotAmounts, "E1201008", i).
				Message(acc_term.InvalidSlotAmounts).Custom(pkg_err.ValidationFailedErr).Build()
		}
		if !debit.IsZero() && !credit.IsZero() {
			return pkg_err.New(acc_term.InvalidSlotAmounts, "E1201009", i).
				Message(acc_term.InvalidSlotAmounts).Custom(pkg_err.ValidationFailedErr).Build()
		}

		sumDebit = sumDebit.Add(debit)
		sumCredit = sumCredit.Add(credit)
	}

	if !sumDebit.Equal(sumCredit) {
		return pkg_err.New(acc_term.SlotsMustBalance, "E1201010").
			Message(acc_term.SlotsMustBalance).Custom(pkg_err.ValidationFailedErr).Build()
	}
	if !sumDebit.Equal(transaction.Amount.Num()) {
		return pkg_err.New(acc_term.SlotsMustMatchAmount, "E1201011").
			Message(acc_term.SlotsMustMatchAmount).Custom(pkg_err.ValidationFailedErr).Build()
	}
	return nil
}

// Create is disabled for headers — use DoTransaction.
func (s *AccTransactionServ) Create(tx tx.Tx, transaction acc_model.Transaction) (createdTransaction acc_model.Transaction, err error) {
	err = pkg_err.New(acc_term.LedgerRecordsAreImmutable, "E1201012").
		Message(acc_term.LedgerRecordsAreImmutable).Custom(pkg_err.ValidationFailedErr).Build()
	return
}

// Save is disabled — ledger headers are immutable.
func (s *AccTransactionServ) Save(tx tx.Tx, transaction acc_model.Transaction) (updatedTransaction, transactionBefore acc_model.Transaction, err error) {
	err = pkg_err.New(acc_term.LedgerRecordsAreImmutable, "E1201013").
		Message(acc_term.LedgerRecordsAreImmutable).Custom(pkg_err.ValidationFailedErr).Build()
	return
}

// Delete is disabled — ledger headers are immutable (hard-delete remains in repo for repair tools only).
func (s *AccTransactionServ) Delete(tx tx.Tx, id uint) (transaction acc_model.Transaction, err error) {
	err = pkg_err.New(acc_term.LedgerRecordsAreImmutable, "E1201014").
		Message(acc_term.LedgerRecordsAreImmutable).Custom(pkg_err.ValidationFailedErr).Build()
	return
}
