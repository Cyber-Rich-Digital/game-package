package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"encoding/json"
	"fmt"
)

type BankingService interface {
	GetBankStatementById(req model.BankStatementGetRequest) (*model.BankStatement, error)
	GetBankStatements(req model.BankStatementListRequest) (*model.SuccessWithPagination, error)
	GetBankStatementSummary(req model.BankStatementListRequest) (*model.BankStatementSummary, error)
	CreateBankStatement(data model.BankStatementCreateBody) error
	MatchStatementOwner(id int64, req model.BankStatementMatchRequest) error
	IgnoreStatementOwner(id int64, req model.BankStatementMatchRequest) error
	DeleteBankStatement(id int64) error

	GetBankTransactionById(req model.BankTransactionGetRequest) (*model.BankTransaction, error)
	GetBankTransactions(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error)
	CreateBankTransaction(data model.BankTransactionCreateBody) error
	CreateBonusTransaction(data model.BonusTransactionCreateBody) error
	DeleteBankTransaction(id int64) error

	GetPendingDepositTransactions(req model.PendingDepositTransactionListRequest) (*model.SuccessWithPagination, error)
	GetPendingWithdrawTransactions(req model.PendingWithdrawTransactionListRequest) (*model.SuccessWithPagination, error)
	ConfirmDepositTransaction(id int64, data model.BankConfirmDepositRequest) error
	ConfirmWithdrawTransaction(id int64, data model.BankConfirmWithdrawRequest) error
	CancelPendingTransaction(id int64, data model.BankTransactionCancelBody) error
	GetFinishedTransactions(req model.FinishedTransactionListRequest) (*model.SuccessWithPagination, error)
	RemoveFinishedTransaction(id int64, data model.BankTransactionRemoveBody) error
	GetRemovedTransactions(req model.RemovedTransactionListRequest) (*model.SuccessWithPagination, error)

	GetMemberByCode(code string) (*model.Member, error)
	GetMembers(req model.MemberListRequest) (*model.SuccessWithPagination, error)
	GetPossibleStatementOwners(req model.MemberPossibleListRequest) (*model.SuccessWithPagination, error)
	GetMemberTransactions(req model.MemberTransactionListRequest) (*model.SuccessWithPagination, error)
	GetMemberTransactionSummary(req model.MemberTransactionListRequest) (*model.MemberTransactionSummary, error)
}

var memberNotFound = "Member not found"
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

	record, err := s.repoBanking.GetBankStatementById(req.Id)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(bankStatementferNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return record, nil
}

func (s *bankingService) GetBankStatements(req model.BankStatementListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetBankStatements(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetBankStatementSummary(req model.BankStatementListRequest) (*model.BankStatementSummary, error) {

	records, err := s.repoBanking.GetBankStatementSummary(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) CreateBankStatement(data model.BankStatementCreateBody) error {

	toAccount, err := s.repoAccounting.GetBankAccountById(data.AccountId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid Bank Account")
	}
	var body model.BankStatementCreateBody
	body.AccountId = toAccount.Id
	if data.StatementType == "transfer_in" {
		body.Amount = data.Amount
	} else if data.StatementType == "transfer_out" {
		body.Amount = data.Amount * -1
	} else {
		return badRequest("Invalid Transfer Type")
	}
	body.Detail = data.Detail
	body.StatementType = data.StatementType
	body.TransferAt = data.TransferAt
	body.Status = "pending"

	if err := s.repoBanking.CreateBankStatement(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) MatchStatementOwner(id int64, req model.BankStatementMatchRequest) error {

	statement, err := s.repoBanking.GetBankStatementById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if statement.Status != "pending" {
		return badRequest("Statement is not pending")
	}
	member, err := s.repoBanking.GetMemberById(req.UserId)
	if err != nil {
		return badRequest("Invalid Member")
	}
	jsonBefore, _ := json.Marshal(statement)

	// TransAction
	var createBody model.CreateBankStatementActionBody
	createBody.StatementId = statement.Id
	createBody.UserId = req.UserId
	createBody.ActionType = "confirmed"
	createBody.AccountId = statement.AccountId
	createBody.JsonBefore = string(jsonBefore)
	createBody.ConfirmedAt = req.ConfirmedAt
	createBody.ConfirmedByUserId = req.ConfirmedByUserId
	createBody.ConfirmedByUsername = req.ConfirmedByUsername
	if err := s.repoBanking.CreateStatementAction(createBody); err == nil {
		var body model.BankStatementUpdateBody
		body.Status = "confirmed"
		// todo:
		member.Id = 0

		if err := s.repoBanking.MatchStatementOwner(id, body); err != nil {
			return internalServerError(err.Error())
		}
	} else {
		return internalServerError(err.Error())
	}

	return nil
}

func (s *bankingService) IgnoreStatementOwner(id int64, req model.BankStatementMatchRequest) error {

	statement, err := s.repoBanking.GetBankStatementById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if statement.Status != "pending" {
		return badRequest("Statement is not pending")
	}

	var body model.BankStatementUpdateBody
	body.Status = "ignored"
	jsonBefore, _ := json.Marshal(statement)

	// TransAction
	var createBody model.CreateBankStatementActionBody
	createBody.StatementId = statement.Id
	createBody.ActionType = "ignored"
	createBody.AccountId = statement.AccountId
	createBody.JsonBefore = string(jsonBefore)
	createBody.ConfirmedAt = req.ConfirmedAt
	createBody.ConfirmedByUserId = req.ConfirmedByUserId
	createBody.ConfirmedByUsername = req.ConfirmedByUsername
	if err := s.repoBanking.CreateStatementAction(createBody); err == nil {
		if err := s.repoBanking.IgnoreStatementOwner(id, body); err != nil {
			return internalServerError(err.Error())
		}
	} else {
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

	record, err := s.repoBanking.GetBankTransactionById(req.Id)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(bankTransactionferNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return record, nil
}

func (s *bankingService) GetBankTransactions(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetBankTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) CreateBankTransaction(data model.BankTransactionCreateBody) error {

	var body model.BankTransactionCreateBody

	if data.TransferType == "deposit" {
		member, err := s.repoAccounting.GetUserByMemberCode(data.MemberCode)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Member code")
		}
		bank, err := s.repoAccounting.GetBankByCode(member.Bankname)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid User Bank")
		}
		body.MemberCode = *member.MemberCode
		body.UserId = member.Id
		body.CreditAmount = data.CreditAmount
		body.TransferType = data.TransferType
		body.DepositChannel = data.DepositChannel
		body.OverAmount = data.OverAmount
		body.IsAutoCredit = data.IsAutoCredit

		body.FromBankId = &bank.Id
		body.FromAccountName = &member.Fullname
		body.FromAccountNumber = &member.BankAccount
		if data.ToAccountId == nil {
			return badRequest("Input Bank Account")
		}
		toAccount, err := s.repoAccounting.GetDepositAccountById(*data.ToAccountId)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Bank Account")
		}
		body.ToAccountId = &toAccount.Id
		body.ToBankId = &toAccount.BankId
		body.ToAccountName = &toAccount.AccountName
		body.ToAccountNumber = &toAccount.AccountNumber

		// todo: createBonus + refDeposit
		body.PromotionId = data.PromotionId

	} else if data.TransferType == "withdraw" {
		member, err := s.repoAccounting.GetUserByMemberCode(data.MemberCode)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Member code")
		}
		bank, err := s.repoAccounting.GetBankByCode(member.Bankname)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid User Bank")
		}
		body.MemberCode = *member.MemberCode
		body.UserId = member.Id
		body.CreditAmount = data.CreditAmount
		body.TransferType = data.TransferType

		fromAccount, err := s.repoAccounting.GetWithdrawAccountById(*data.FromAccountId)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Bank Account")
		}
		body.FromAccountId = &fromAccount.Id
		body.FromBankId = &fromAccount.BankId
		body.FromAccountName = &fromAccount.AccountName
		body.FromAccountNumber = &fromAccount.AccountNumber

		body.ToBankId = &bank.Id
		body.ToAccountName = &member.Fullname
		body.ToAccountNumber = &member.BankAccount
	} else if data.TransferType == "getcreditback" {
		// ดึงยอดสลายไปเลย
		member, err := s.repoAccounting.GetUserByMemberCode(data.MemberCode)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Member code")
		}
		bank, err := s.repoAccounting.GetBankByCode(member.Bankname)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid User Bank")
		}
		body.MemberCode = *member.MemberCode
		body.UserId = member.Id
		body.CreditAmount = data.CreditAmount
		body.TransferType = data.TransferType

		body.FromBankId = &bank.Id
		body.FromAccountName = &member.Fullname
		body.FromAccountNumber = &member.BankAccount
	} else {
		return badRequest("Invalid Transfer Type")
	}

	body.TransferAt = data.TransferAt
	body.CreatedByUserId = data.CreatedByUserId
	body.CreatedByUsername = data.CreatedByUsername
	body.Status = "pending"

	if _, err := s.repoBanking.CreateBankTransaction(body); err != nil {
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
	bank, err := s.repoAccounting.GetBankByCode(member.Bankname)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid User Bank")
	}

	var body model.BonusTransactionCreateBody
	body.MemberCode = *member.MemberCode
	body.UserId = member.Id
	body.TransferType = "bonus"
	body.ToAccountId = 0
	body.ToBankId = bank.Id
	body.ToAccountName = member.Fullname
	body.ToAccountNumber = member.BankAccount
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
	records, err := s.repoBanking.GetPendingDepositTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetPendingWithdrawTransactions(req model.PendingWithdrawTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetPendingWithdrawTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) CancelPendingTransaction(id int64, data model.BankTransactionCancelBody) error {

	record, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if record.Status != "pending" {
		return badRequest("Transaction is not pending")
	}

	if err := s.repoBanking.CancelPendingTransaction(id, data); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) ConfirmDepositTransaction(id int64, req model.BankConfirmDepositRequest) error {

	record, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if record.Status != "pending" {
		return badRequest("Transaction is not pending")
	}
	if record.TransferType != "deposit" {
		return badRequest("Transaction is not deposit")
	}
	jsonBefore, _ := json.Marshal(record)

	var updateData model.BankTransactionConfirmBody
	updateData.Status = "finished"
	updateData.ConfirmedAt = req.ConfirmedAt
	updateData.ConfirmedByUserId = req.ConfirmedByUserId
	updateData.ConfirmedByUsername = req.ConfirmedByUsername
	updateData.BonusAmount = req.BonusAmount

	var createBody model.CreateBankTransactionActionBody
	createBody.TransactionId = record.Id
	createBody.UserId = record.UserId
	createBody.TransferType = record.TransferType
	createBody.FromAccountId = record.FromAccountId
	createBody.ToAccountId = record.ToAccountId
	createBody.JsonBefore = string(jsonBefore)
	if req.TransferAt == nil {
		createBody.TransferAt = record.TransferAt
	} else {
		TransferAt := req.TransferAt
		createBody.TransferAt = *TransferAt
		updateData.TransferAt = *TransferAt
	}
	createBody.SlipUrl = req.SlipUrl
	createBody.BonusAmount = req.BonusAmount
	createBody.ConfirmedAt = req.ConfirmedAt
	createBody.ConfirmedByUserId = req.ConfirmedByUserId
	createBody.ConfirmedByUsername = req.ConfirmedByUsername
	if err := s.repoBanking.CreateTransactionAction(createBody); err != nil {
		return internalServerError(err.Error())
	}
	if err := s.repoBanking.ConfirmPendingTransaction(id, updateData); err != nil {
		return internalServerError(err.Error())
	}
	if err := s.IncreaseMemberCredit(record.UserId, record.CreditAmount); err != nil {
		return internalServerError(err.Error())
	}
	// todo: Bonus
	// commit
	return nil
}

func (s *bankingService) IncreaseMemberCredit(userId int64, creditAmount float32) error {

	// record, err := s.repoBanking.GetBankTransactionById(id)
	// if err != nil {
	// 	return internalServerError(err.Error())
	// }

	if err := s.repoBanking.IncreaseMemberCredit(userId, creditAmount); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}

func (s *bankingService) DecreaseMemberCredit(userId int64, creditAmount float32) error {

	// record, err := s.repoBanking.GetBankTransactionById(id)
	// if err != nil {
	// 	return internalServerError(err.Error())
	// }

	if err := s.repoBanking.DecreaseMemberCredit(userId, creditAmount); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}

func (s *bankingService) ConfirmWithdrawTransaction(id int64, req model.BankConfirmWithdrawRequest) error {

	record, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if record.Status != "pending" {
		return badRequest("Transaction is not pending")
	}
	if record.TransferType != "withdraw" {
		return badRequest("Transaction is not withdraw")
	}
	jsonBefore, _ := json.Marshal(record)

	var updateData model.BankTransactionConfirmBody
	updateData.Status = "finished"
	updateData.ConfirmedAt = req.ConfirmedAt
	updateData.ConfirmedByUserId = req.ConfirmedByUserId
	updateData.ConfirmedByUsername = req.ConfirmedByUsername
	updateData.BankChargeAmount = req.BankChargeAmount
	updateData.CreditAmount = req.CreditAmount

	var createBody model.CreateBankTransactionActionBody
	createBody.TransactionId = record.Id
	createBody.UserId = record.UserId
	createBody.TransferType = record.TransferType
	createBody.FromAccountId = record.FromAccountId
	createBody.ToAccountId = record.ToAccountId
	createBody.JsonBefore = string(jsonBefore)
	createBody.TransferAt = record.TransferAt
	createBody.CreditAmount = req.CreditAmount
	createBody.BankChargeAmount = req.BankChargeAmount
	createBody.ConfirmedAt = req.ConfirmedAt
	createBody.ConfirmedByUserId = req.ConfirmedByUserId
	createBody.ConfirmedByUsername = req.ConfirmedByUsername
	if err := s.repoBanking.CreateTransactionAction(createBody); err != nil {
		return internalServerError(err.Error())
	}
	if err := s.repoBanking.ConfirmPendingTransaction(id, updateData); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) GetFinishedTransactions(req model.FinishedTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetFinishedTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) RemoveFinishedTransaction(id int64, data model.BankTransactionRemoveBody) error {

	record, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if record.Status != "finished" {
		return badRequest("Transaction is not finished")
	}

	if err := s.repoBanking.RemoveFinishedTransaction(id, data); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *bankingService) GetRemovedTransactions(req model.RemovedTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetRemovedTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetMemberByCode(code string) (*model.Member, error) {

	if code == "" {
		return nil, badRequest("Code is required")
	}

	records, err := s.repoBanking.GetMemberByCode(code)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(memberNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetMembers(req model.MemberListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetMembers(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetPossibleStatementOwners(req model.MemberPossibleListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	statement, err := s.repoBanking.GetBankStatementById(req.UnknownStatementId)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(bankStatementferNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	req.UserBankCode = &statement.FromBankCode
	req.UserAccountNumber = &statement.FromAccountNumber

	records, err := s.repoBanking.GetPossibleStatementOwners(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetMemberTransactions(req model.MemberTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&req.Page, &req.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	records, err := s.repoBanking.GetMemberTransactions(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *bankingService) GetMemberTransactionSummary(req model.MemberTransactionListRequest) (*model.MemberTransactionSummary, error) {

	result, err := s.repoBanking.GetMemberTransactionSummary(req)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return result, nil
}

func (s *bankingService) MatchDepositTransaction(id int64, req model.BankConfirmDepositRequest) error {

	record, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if record.Status != "pending" {
		return badRequest("Transaction is not pending")
	}
	if record.TransferType != "deposit" {
		return badRequest("Transaction is not deposit")
	}
	// todo: Bonus
	// commit
	if err := s.ConfirmDepositTransaction(record.UserId, req); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}

func (s *bankingService) MatchWithdrawTransaction(id int64, req model.BankConfirmWithdrawRequest) error {

	record, err := s.repoBanking.GetBankTransactionById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	if record.Status != "pending" {
		return badRequest("Transaction is not pending")
	}
	if record.TransferType != "withdraw" {
		return badRequest("Transaction is not withdraw")
	}

	// todo: Match

	if err := s.ConfirmWithdrawTransaction(record.UserId, req); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}
