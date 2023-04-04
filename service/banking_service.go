package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"fmt"
)

type BankingService interface {
	GetBankStatementById(req model.BankStatementGetRequest) (*model.BankStatement, error)
	GetBankStatements(req model.BankStatementListRequest) (*model.SuccessWithPagination, error)
	CreateBankStatement(data model.BankStatementCreateBody) error
	DeleteBankStatement(id int64) error

	GetBankTransactionById(req model.BankTransactionGetRequest) (*model.BankTransaction, error)
	GetBankTransactions(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error)
	CreateBankTransaction(data model.BankTransactionCreateBody) error
	DeleteBankTransaction(id int64) error
}

var bankStatementferNotFound = "Statement not found"
var bankTransactionferNotFound = "Transaction not found"

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

func (s *bankingService) GetBankStatementById(req model.BankStatementGetRequest) (*model.BankStatement, error) {

	banking, err := s.repoBanking.GetBankStatementById(req.Id)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(bankStatementferNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return banking, nil
}

func (s *bankingService) GetBankStatements(req model.BankStatementListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	banking, err := s.repoBanking.GetBankStatements(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return banking, nil
}

func (s *bankingService) CreateBankStatement(data model.BankStatementCreateBody) error {

	toAccount, err := s.repoAccounting.GetBankAccountById(data.AccountId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid Bank Account")
	}

	var body model.BankStatementCreateBody
	body.AccountId = toAccount.Id
	body.Amount = data.Amount
	body.TransferAt = data.TransferAt
	body.Status = "pending"

	if err := s.repoBanking.CreateBankStatement(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) DeleteBankStatement(id int64) error {

	_, err := s.repoBanking.GetBankStatementById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if err := s.repoBanking.DeleteBankStatement(id); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) GetBankTransactionById(req model.BankTransactionGetRequest) (*model.BankTransaction, error) {

	banking, err := s.repoBanking.GetBankTransactionById(req.Id)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(bankTransactionferNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return banking, nil
}

func (s *bankingService) GetBankTransactions(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	banking, err := s.repoBanking.GetBankTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return banking, nil
}

func (s *bankingService) CreateBankTransaction(data model.BankTransactionCreateBody) error {

	fromAccount, err := s.repoAccounting.GetBankAccountById(data.FromAccountId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid Bank Account")
	}

	var body model.BankTransactionCreateBody
	body.FromAccountId = fromAccount.Id
	// body.Amount = data.Amount
	// body.TransferAt = data.TransferAt
	body.Status = "pending"

	if err := s.repoBanking.CreateBankTransaction(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) DeleteBankTransaction(id int64) error {

	_, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if err := s.repoBanking.DeleteBankTransaction(id); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}
