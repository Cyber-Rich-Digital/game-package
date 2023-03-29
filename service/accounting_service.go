package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"fmt"
)

type AccountingService interface {
	GetBanks(data model.BankListRequest) (*model.Pagination, error)

	GetAccountTypes(data model.AccountTypeListRequest) (*model.Pagination, error)

	GetBankAccountById(data model.BankAccountParam) (*model.BankAccount, error)
	GetBankAccounts(data model.BankAccountListRequest) (*model.Pagination, error)
	CreateBankAccount(data model.BankAccountBody) error
	UpdateBankAccount(id int64, data model.BankAccountBody) error
	DeleteBankAccount(id int64) error
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

	check, err := s.repo.HasBankAccount(data.AccountNumber)
	if err != nil {
		return internalServerError(err.Error())
	}

	if check {
		return badRequest("Account already exist")
	}

	if err := s.repo.UpdateBankAccount(id, data); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}

func (s *accountingService) DeleteBankAccount(id int64) error {

	if err := s.repo.DeleteBankAccount(id); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}
