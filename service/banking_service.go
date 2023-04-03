package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"fmt"
)

type BankingService interface {
	GetStatementById(req model.BankStatementGetRequest) (*model.BankStatement, error)
	GetStatements(req model.BankStatementListRequest) (*model.SuccessWithPagination, error)
	CreateStatement(data model.BankStatementCreateBody) error
	DeleteStatement(id int64) error
}

type bankingService struct {
	repoBanking    repository.BankingRepository
	repoAccounting repository.AccountingRepository
}

func NewBankingService(
	repoBanking repository.BankingRepository,
	repoAccounting repository.AccountingRepository,
) BankingService {
	return &bankingService{repoBanking, repoAccounting}
}

func (s *bankingService) GetStatementById(req model.BankStatementGetRequest) (*model.BankStatement, error) {

	banking, err := s.repoBanking.GetStatementById(req.Id)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(transferNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return banking, nil
}

func (s *bankingService) GetStatements(req model.BankStatementListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	banking, err := s.repoBanking.GetStatements(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return banking, nil
}

func (s *bankingService) CreateStatement(data model.BankStatementCreateBody) error {

	toAccount, err := s.repoAccounting.GetBankAccountById(data.AccountId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid destination Bank Account")
	}

	var body model.BankStatementCreateBody
	body.AccountId = toAccount.Id
	body.Amount = data.Amount
	body.TransferAt = data.TransferAt
	body.Status = "pending"

	if err := s.repoBanking.CreateStatement(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) DeleteStatement(id int64) error {

	_, err := s.repoBanking.GetStatementById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if err := s.repoBanking.DeleteStatement(id); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}
