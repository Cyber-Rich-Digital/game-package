package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
)

type AccountingService interface {
	GetBankAccountAndTags(data model.BankAccountParam, userId int) (*model.BankAccount, error)
	GetBankAccount(data model.BankAccountParam, userId int) (*model.BankAccount, error)
	GetBankAccounts(data model.BankAccountQuery) (*model.Pagination, error)
	GetBankAccountTotals(body model.BankAccountDate) (*[]model.BankAccountListResponse, error)
	CreateBankAccount(data model.BankAccountBody, userId int) error
	UpdateBankAccount(id int64, data model.BankAccountBody) error
	DeleteBankAccount(id int) error
}

type accountingService struct {
	repo repository.AccountingRepository
}

var bankAccountNotFound = "Accounting not found"

func NewAccountingService(
	repo repository.AccountingRepository,
) AccountingService {
	return &accountingService{repo}
}

func (s *accountingService) GetBankAccountAndTags(data model.BankAccountParam, userId int) (*model.BankAccount, error) {

	accounting, err := s.repo.GetBankAccountAndTags(data.Id)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound(bankAccountNotFound)
		}

		if err.Error() == "Account not found" {
			return nil, notFound(bankAccountNotFound)
		}

		return nil, internalServerError(err.Error())
	}

	// if accounting.UserId != userId {
	// 	return nil, badRequest("You don't have permission to access this accounting")
	// }

	return accounting, nil
}

func (s *accountingService) GetBankAccount(data model.BankAccountParam, userId int) (*model.BankAccount, error) {

	accounting, err := s.repo.GetBankAccount(data.Id)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound(bankAccountNotFound)
		}

		if err.Error() == "Account not found" {
			return nil, notFound(bankAccountNotFound)
		}

		return nil, internalServerError(err.Error())
	}

	// if accounting.UserId != userId {
	// 	return nil, badRequest("You don't have permission to access this accounting")
	// }

	return accounting, nil
}

func (s *accountingService) GetBankAccountTotals(body model.BankAccountDate) (*[]model.BankAccountListResponse, error) {

	accounting, err := s.repo.GetBankAccountTotals(body)
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

func (s *accountingService) GetBankAccounts(data model.BankAccountQuery) (*model.Pagination, error) {

	if err := helper.Pagination(&data.Page, &data.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	accounting, err := s.repo.GetBankAccounts(data)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	return accounting, nil
}

func (s *accountingService) CreateBankAccount(data model.BankAccountBody, userId int) error {

	exist, err := s.repo.CheckBankAccount(data.DomainName)
	if err != nil {
		return internalServerError(err.Error())
	}

	if exist {
		return badRequest("Account already exist")
	}

	key := helper.GenKey(50)

	var accounting model.BankAccount
	accounting.Title = data.Title
	accounting.DomainName = data.DomainName
	accounting.UserId = userId
	accounting.ApiKey = key

	if err := s.repo.CreateBankAccount(accounting); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}

func (s *accountingService) UpdateBankAccount(id int64, data model.BankAccountBody) error {

	// if check, err := s.repo.GetBankAccountsByIdAndUserIds(id); err != nil {
	// 	return internalServerError(err.Error())
	// } else if !check {
	// 	return notFound("Account not found")
	// }
	check, err := s.repo.CheckBankAccount(data.DomainName)
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

func (s *accountingService) DeleteBankAccount(id int) error {

	if err := s.repo.DeleteBankAccount(id); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}
