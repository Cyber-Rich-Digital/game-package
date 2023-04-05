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
	CreateBonusTransaction(data model.BonusTransactionCreateBody) error
	DeleteBankTransaction(id int64) error

	GetPendingDepositTransactions(req model.PendingDepositTransactionListRequest) (*model.SuccessWithPagination, error)
	GetPendingWithdrawTransactions(req model.PendingWithdrawTransactionListRequest) (*model.SuccessWithPagination, error)
	GetFinishedTransactions(req model.FinishedTransactionListRequest) (*model.SuccessWithPagination, error)
	RemoveFinishedTransaction(id int64, data model.BankTransactionRemoveBody) error
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

	var body model.BankTransactionCreateBody

	if data.TransferType == "deposit" {
		member, err := s.repoAccounting.GetUserByMemberCode(data.MemberCode)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Member code")
		}
		body.MemberCode = member.MemberCode
		body.UserId = member.Id
		body.CreditAmount = data.CreditAmount
		body.TransferType = data.TransferType
		body.DepositChannel = data.DepositChannel
		body.OverAmount = data.OverAmount
		body.IsAutoCredit = data.IsAutoCredit

		body.FromAccountId = 0
		body.FromBankId = member.BankId
		body.FromAccountName = member.AccountName
		body.FromAccountNumber = member.AccountNumber
		toAccount, err := s.repoAccounting.GetBankAccountById(data.ToAccountId)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Bank Account")
		}
		body.ToAccountId = toAccount.Id
		body.ToBankId = toAccount.BankId
		body.ToAccountName = toAccount.AccountName
		body.ToAccountNumber = toAccount.AccountNumber
		// body.PromotionId = data.PromotionId

	} else if data.TransferType == "withdraw" {
		member, err := s.repoAccounting.GetUserByMemberCode(data.MemberCode)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Member code")
		}
		body.MemberCode = member.MemberCode
		body.UserId = member.Id
		body.CreditAmount = data.CreditAmount
		body.TransferType = data.TransferType

		fromAccount, err := s.repoAccounting.GetBankAccountById(data.FromAccountId)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Bank Account")
		}
		body.FromAccountId = fromAccount.Id
		body.FromBankId = fromAccount.BankId
		body.FromAccountName = fromAccount.AccountName
		body.FromAccountNumber = fromAccount.AccountNumber

		body.ToAccountId = 0
		body.ToBankId = member.BankId
		body.ToAccountName = member.AccountName
		body.ToAccountNumber = member.AccountNumber

	} else if data.TransferType == "getcreditback" {
		member, err := s.repoAccounting.GetUserByMemberCode(data.MemberCode)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Member code")
		}
		body.MemberCode = member.MemberCode
		body.UserId = member.Id
		body.CreditAmount = data.CreditAmount
		body.TransferType = data.TransferType

		// body.ToAccountId = 0
		// body.ToBankId = member.BankId
		// body.ToAccountName = member.AccountName
		// body.ToAccountNumber = member.AccountNumber

	} else {
		return badRequest("Invalid Transfer Type")
	}

	body.TransferAt = data.TransferAt
	body.CreatedByUserId = data.CreatedByUserId
	body.CreatedByUsername = data.CreatedByUsername
	body.Status = "pending"

	if err := s.repoBanking.CreateBankTransaction(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) CreateBonusTransaction(data model.BonusTransactionCreateBody) error {

	member, err := s.repoAccounting.GetUserByMemberCode(data.MemberCode)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid Member code")
	}

	var body model.BonusTransactionCreateBody
	body.MemberCode = member.MemberCode
	body.UserId = member.Id
	body.TransferType = "bonus"
	body.ToAccountId = 0
	body.ToBankId = member.BankId
	body.ToAccountName = member.AccountName
	body.ToAccountNumber = member.AccountNumber
	// body.BeforeAmount = data.BeforeAmount
	// body.AfterAmount = data.AfterAmount
	body.BonusAmount = data.BonusAmount
	body.BonusReason = data.BonusReason
	body.TransferAt = data.TransferAt
	body.CreatedByUserId = data.CreatedByUserId
	body.CreatedByUsername = data.CreatedByUsername
	body.Status = "pending"

	if err := s.repoBanking.CreateBonusTransaction(body); err != nil {
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

func (s *bankingService) GetPendingDepositTransactions(req model.PendingDepositTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	banking, err := s.repoBanking.GetPendingDepositTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return banking, nil
}

func (s *bankingService) GetPendingWithdrawTransactions(req model.PendingWithdrawTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	banking, err := s.repoBanking.GetPendingWithdrawTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return banking, nil
}

func (s *bankingService) GetFinishedTransactions(req model.FinishedTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	banking, err := s.repoBanking.GetFinishedTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return banking, nil
}
func (s *bankingService) RemoveFinishedTransaction(id int64, data model.BankTransactionRemoveBody) error {

	record, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if record.Status != "finished" {
		return badRequest("Transaction is not finished")
	}
	if record.Status == "removed" {
		return badRequest("Transaction is already removed")
	}

	if err := s.repoBanking.RemoveFinishedTransaction(id, data); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}
