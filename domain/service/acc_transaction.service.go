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

// FindByID for getting transaction by its id
func (s *AccTransactionServ) FindByID(tx tx.Tx, id uint) (transaction acc_model.Transaction, err error) {

	key := fmt.Sprintf("%v-%v", acc_term.Transaction, id)
	if ok := s.Engine.RedisCacheAPI.GetCache(tx, key, &transaction); ok {
		return
	}

	if transaction, err = s.Repo.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1155388", "can't fetch the transaction", id)
		return
	}

	err = s.Engine.RedisCacheAPI.Set(key, transaction)

	return
}

// GetAll of transactions, it supports pagination and search and return count
func (s *AccTransactionServ) GetAll(params param.Param) (transactions []acc_model.Transaction, err error) {

	if ok := s.Engine.RedisCacheAPI.GetCache(params.Tx, acc_term.Transactions, &transactions); ok {
		return
	}

	params.Limit = 100000
	if transactions, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in transactions list")
		return
	}

	err = s.Engine.RedisCacheAPI.Set(acc_term.Transactions, transactions)

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

// DoTransaction a transaction
func (s *AccTransactionServ) DoTransaction(tx tx.Tx, transaction acc_model.Transaction) (createdTransaction acc_model.Transaction, err error) {

	if err = validator.ValidateModel(transaction, acc_term.Transaction, validator.Create); err != nil {
		err = pkg_err.TickValidate(err, "E1163777", pkg_err.ValidationFailed, transaction)
		return
	}

	// calculate fee and net amount
	transaction.NetAmount = transaction.Amount.Num().Sub(transaction.Fee.Num()).Build()

	// set transaction hash
	transaction.Hash = uuid.New().String()

	if createdTransaction, err = s.Repo.Create(tx, transaction); err != nil {
		pkg_err.Log(err, "E1126470", "error in creating transaction", transaction)
		return
	}

	for _, v := range transaction.Slots {

		// get account Credit
		var accountCredit acc_model.AccountCredit
		tx.IsLock = true
		if accountCredit, err = AccAccountCreditService.FindByAccountIDAndCurrency(tx, v.AccountID, v.CurrencyID); err != nil {
			err = pkg_err.Log(err, "E1118793", "sender account not found")
			return
		}
		tx.IsLock = false

		// check if does have enough balance
		if !v.Credit.Num().IsZero() {
			// check balance with the amount
			if accountCredit.Balance.Num().LessThan(v.Credit.Num()) {
				err = pkg_err.New(acc_term.YouHaveNotEnoughCredit, "E1099045", "you have not enough balance", accountCredit.AccountID).
					Message(acc_term.YouHaveNotEnoughCredit).Custom(pkg_err.BadRequestErr).Build()
				return
			}
		}

		// save latest credit
		remainBalance := accountCredit.Balance.Num().Add(v.Debit.Num()).Sub(v.Credit.Num()).Build()
		accountCredit.Balance = remainBalance

		if _, _, err = AccAccountCreditService.Save(tx, accountCredit); err != nil {
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1167734", "error in updating credit account", accountCredit).
				Message(pkg_err.SomethingWentWrong).
				Custom(pkg_err.InternalServerErr).Build()
			return
		}

		// add slot
		slot := acc_model.Slot{
			TransactionID: createdTransaction.ID,
			AccountID:     v.AccountID,
			CurrencyID:    v.CurrencyID,
			Credit:        v.Credit,
			Debit:         v.Debit,
			Balance:       remainBalance,
		}

		if slot, err = AccSlotService.Create(tx, slot); err != nil {
			err = pkg_err.New(pkg_err.SomethingWentWrong, "E1159619", "error in creating sender slot", slot).
				Message(pkg_err.SomethingWentWrong).
				Custom(pkg_err.InternalServerErr).Build()
			return
		}

		createdTransaction.Slots = append(createdTransaction.Slots, slot)

	}

	// ignore slots
	transaction.Slots = nil

	return
}

// Create a transaction
func (s *AccTransactionServ) Create(tx tx.Tx, transaction acc_model.Transaction) (createdTransaction acc_model.Transaction, err error) {

	if err = validator.ValidateModel(transaction, acc_term.Transaction, validator.Create); err != nil {
		err = pkg_err.TickValidate(err, "E1152912", pkg_err.ValidationFailed, transaction)
		return
	}

	if createdTransaction, err = s.Repo.Create(tx, transaction); err != nil {
		pkg_err.Log(err, "E1171379", "error in creating transaction", transaction)
		return
	}

	s.Engine.RedisCacheAPI.Delete(acc_term.Transactions)

	return
}

// Save a transaction, if it is exists update it, if not create it
func (s *AccTransactionServ) Save(tx tx.Tx, transaction acc_model.Transaction) (updatedTransaction, transactionBefore acc_model.Transaction, err error) {

	if err = validator.ValidateModel(transaction, acc_term.Transaction, validator.Update); err != nil {
		err = pkg_err.TickValidate(err, "E1151288", pkg_err.ValidationFailed, transaction)
		return
	}

	if transactionBefore, err = s.FindByID(tx, transaction.ID); err != nil {
		pkg_err.Log(err, "E1168605", "can't fetch transaction by id for saving it", transaction.ID)
		return
	}

	if updatedTransaction, err = s.Repo.Save(tx, transaction); err != nil {
		pkg_err.Log(err, "E1186552", "transaction not saved")
		return
	}

	key := fmt.Sprintf("%v-%v", acc_term.Transaction, updatedTransaction.ID)
	if err = s.Engine.RedisCacheAPI.Delete(key); err != nil {
		return
	}

	s.Engine.RedisCacheAPI.Delete(acc_term.Transactions)

	return
}

// Delete transaction, it is soft delete
func (s *AccTransactionServ) Delete(tx tx.Tx, id uint) (transaction acc_model.Transaction, err error) {

	if transaction, err = s.FindByID(tx, id); err != nil {
		pkg_err.Log(err, "E1145937", "transaction not found for deleting")
		return
	}

	if err = s.Repo.Delete(tx, transaction); err != nil {
		pkg_err.Log(err, "E1160494", "transaction not deleted")
		return
	}

	key := fmt.Sprintf("%v-%v", acc_term.Transaction, transaction.ID)
	s.Engine.RedisCacheAPI.Delete(key)
	s.Engine.RedisCacheAPI.Delete(acc_term.Transactions)

	return
}
