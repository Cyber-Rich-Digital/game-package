package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"fmt"
	"time"
)

type AccountingService interface {
	GetBanks(data model.BankListRequest) (*model.Pagination, error)

	GetAccountTypes(data model.AccountTypeListRequest) (*model.Pagination, error)

	GetBankAccountById(data model.BankAccountParam) (*model.BankAccount, error)
	GetBankAccounts(data model.BankAccountListRequest) (*model.Pagination, error)
	CreateBankAccount(data model.BankAccountBody) error
	UpdateBankAccount(id int64, data model.BankAccountBody) error
	DeleteBankAccount(id int64) error

	GetTransactionById(data model.BankAccountTransactionParam) (*model.BankAccountTransaction, error)
	GetTransactions(data model.BankAccountTransactionListRequest) (*model.Pagination, error)
	CreateTransaction(data model.BankAccountTransactionBody) error
	UpdateTransaction(id int64, data model.BankAccountTransactionBody) error
	DeleteTransaction(id int64) error

	GetTransferById(data model.BankAccountTransferParam) (*model.BankAccountTransfer, error)
	GetTransfers(data model.BankAccountTransferListRequest) (*model.Pagination, error)
	CreateTransfer(data model.BankAccountTransferBody) error
	ConfirmTransfer(id int64, data model.BankAccountTransferConfirmBody) error
	DeleteTransfer(id int64) error
}

type accountingService struct {
	repo repository.AccountingRepository
}

var bankNotFound = "Bank not found"
var bankAccountNotFound = "Account not found"

func NewAccountingService(
	repo repository.AccountingRepository,
) AccountingService {
	return &accountingService{repo}
}

func (s *accountingService) GetBanks(params model.BankListRequest) (*model.Pagination, error) {

	if err := helper.Pagination(&params.Page, &params.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	records, err := s.repo.GetBanks(params)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *accountingService) GetAccountTypes(params model.AccountTypeListRequest) (*model.Pagination, error) {

	if err := helper.Pagination(&params.Page, &params.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	records, err := s.repo.GetAccountTypes(params)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *accountingService) GetBankAccountById(data model.BankAccountParam) (*model.BankAccount, error) {

	accounting, err := s.repo.GetBankAccountById(data.Id)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, notFound(bankAccountNotFound)
		}
		if err.Error() == "Account not found" {
			return nil, notFound(bankAccountNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return accounting, nil
}

func (s *accountingService) GetBankAccounts(data model.BankAccountListRequest) (*model.Pagination, error) {

	if err := helper.Pagination(&data.Page, &data.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	accounting, err := s.repo.GetBankAccounts(data)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return accounting, nil
}

func (s *accountingService) CreateBankAccount(data model.BankAccountBody) error {

	bank, err := s.repo.GetBankById(data.BankId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid Bank")
	}

	accountType, err := s.repo.GetAccounTypeById(data.AccountTypeId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid Account Type")
	}

	exist, err := s.repo.HasBankAccount(data.AccountNumber)
	if err != nil {
		fmt.Println(err)
		return internalServerError(err.Error())
	}
	if exist {
		return badRequest("Account already exist")
	}

	var account model.BankAccountBody
	account.BankId = bank.Id
	account.AccountTypeId = accountType.Id
	account.AccountName = data.AccountName
	account.AccountNumber = data.AccountNumber

	if err := s.repo.CreateBankAccount(account); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) UpdateBankAccount(id int64, data model.BankAccountBody) error {

	account, err := s.repo.GetBankAccountById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	// Validate
	if account.BankId != data.BankId {
		bank, err := s.repo.GetBankById(data.BankId)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Bank")
		}
		data.BankId = bank.Id
	}
	if account.AccountTypeId != data.AccountTypeId {
		accountType, err := s.repo.GetAccounTypeById(data.AccountTypeId)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Account Type")
		}
		data.AccountTypeId = accountType.Id
	}
	if account.AccountNumber != data.AccountNumber {
		check, err := s.repo.HasBankAccount(data.AccountNumber)
		if err != nil {
			return internalServerError(err.Error())
		}
		if !check {
			return notFound("Account already exist")
		}
	}

	if err := s.repo.UpdateBankAccount(id, data); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) DeleteBankAccount(id int64) error {

	_, err := s.repo.GetBankAccountById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if err := s.repo.DeleteBankAccount(id); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) GetTransactionById(data model.BankAccountTransactionParam) (*model.BankAccountTransaction, error) {

	accounting, err := s.repo.GetTransactionById(data.Id)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, notFound(bankAccountNotFound)
		}
		if err.Error() == "Account not found" {
			return nil, notFound(bankAccountNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return accounting, nil
}

func (s *accountingService) GetTransactions(data model.BankAccountTransactionListRequest) (*model.Pagination, error) {

	if err := helper.Pagination(&data.Page, &data.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	accounting, err := s.repo.GetTransactions(data)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return accounting, nil
}

func (s *accountingService) CreateTransaction(data model.BankAccountTransactionBody) error {

	account, err := s.repo.GetBankAccountById(data.AccountId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid Bank Account")
	}

	var transaction model.BankAccountTransactionBody
	transaction.AccountId = account.Id
	transaction.Description = data.Description
	transaction.TransferType = data.TransferType
	transaction.Amount = data.Amount
	transaction.TransferAt = data.TransferAt

	if err := s.repo.CreateTransaction(transaction); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) UpdateTransaction(id int64, data model.BankAccountTransactionBody) error {

	_, err := s.repo.GetTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	// no Update
	return notFound("Function not found")
}

func (s *accountingService) DeleteTransaction(id int64) error {

	_, err := s.repo.GetTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if err := s.repo.DeleteTransaction(id); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) GetTransferById(data model.BankAccountTransferParam) (*model.BankAccountTransfer, error) {

	accounting, err := s.repo.GetTransferById(data.Id)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, notFound(bankAccountNotFound)
		}
		if err.Error() == "Account not found" {
			return nil, notFound(bankAccountNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return accounting, nil
}

func (s *accountingService) GetTransfers(data model.BankAccountTransferListRequest) (*model.Pagination, error) {

	if err := helper.Pagination(&data.Page, &data.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	accounting, err := s.repo.GetTransfers(data)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return accounting, nil
}

func (s *accountingService) CreateTransfer(data model.BankAccountTransferBody) error {

	fromAccount, err := s.repo.GetBankAccountById(data.FromAccountId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid source Bank Account")
	}

	toAccount, err := s.repo.GetBankAccountById(data.ToAccountId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid destination Bank Account")
	}

	var body model.BankAccountTransferBody
	body.FromAccountId = fromAccount.Id
	body.ToAccountId = toAccount.Id
	body.Amount = data.Amount
	body.TransferAt = data.TransferAt
	body.Status = "pending"

	if err := s.repo.CreateTransfer(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) ConfirmTransfer(id int64, data model.BankAccountTransferConfirmBody) error {

	transfer, err := s.repo.GetTransferById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if transfer.Status == "pending" {
		var body model.BankAccountTransferConfirmBody
		body.Status = "confirmed"
		body.ConfirmedAt = time.Now()
		body.ConfirmedByUsername = data.ConfirmedByUsername
		if err := s.repo.ConfirmTransfer(id, body); err != nil {
			return internalServerError(err.Error())
		}
	} else {
		return badRequest("Transfer not in pending status")
	}
	return nil
}

func (s *accountingService) DeleteTransfer(id int64) error {

	_, err := s.repo.GetTransferById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if err := s.repo.DeleteTransfer(id); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}
