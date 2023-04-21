package service

import (
	"bytes"
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type AccountingService interface {
	CheckCurrentAdminId(input any) (*int64, error)
	CheckCurrentUsername(input any) (*string, error)
	CheckConfirmationPassword(data model.ConfirmRequest) (*bool, error)

	GetBanks(data model.BankListRequest) (*model.SuccessWithPagination, error)
	GetAccountTypes(data model.AccountTypeListRequest) (*model.SuccessWithPagination, error)

	GetBankAccountById(data model.BankAccountParam) (*model.BankAccount, error)
	GetBankAccounts(data model.BankAccountListRequest) (*model.SuccessWithPagination, error)
	CreateBankAccount(data model.BankAccountCreateBody) error
	UpdateBankAccount(id int64, data model.BankAccountUpdateRequest) error
	DeleteBankAccount(id int64) error

	GetTransactionById(data model.BankAccountTransactionParam) (*model.BankAccountTransaction, error)
	GetTransactions(data model.BankAccountTransactionListRequest) (*model.SuccessWithPagination, error)
	CreateTransaction(data model.BankAccountTransactionBody) error
	UpdateTransaction(id int64, data model.BankAccountTransactionBody) error
	DeleteTransaction(id int64) error

	GetTransferById(data model.BankAccountTransferParam) (*model.BankAccountTransfer, error)
	GetTransfers(data model.BankAccountTransferListRequest) (*model.SuccessWithPagination, error)
	CreateTransfer(data model.BankAccountTransferBody) error
	ConfirmTransfer(id int64, actorId int64) error
	DeleteTransfer(id int64) error

	GetCustomerAccountsInfo(data model.CustomerAccountInfoRequest) (*model.CustomerAccountInfo, error)
	GetExternalAccounts() (*model.SuccessWithPagination, error)
	GetExternalAccountBalance(query model.ExternalAccountStatusRequest) (*model.ExternalAccountBalance, error)
	GetExternalAccountStatus(query model.ExternalAccountStatusRequest) (*model.ExternalAccountStatus, error)
	CreateExternalAccount(data model.ExternalAccountCreateBody) (*model.ExternalAccountCreateResponse, error)
	UpdateExternalAccount(data model.ExternalAccountUpdateBody) (*model.ExternalAccountCreateResponse, error)
	EnableExternalAccount(query model.ExternalAccountEnableRequest) (*model.ExternalAccountStatus, error)
	DeleteExternalAccount(query model.ExternalAccountStatusRequest) error
	TransferExternalAccount(data model.ExternalAccountTransferRequest) error
	CreateBankStatementFromWebhook(data model.WebhookStatement) error
	CreateBotaccountConfig(data model.BotAccountConfigCreateBody) error

	GetExternalAccountLogs(data model.BankAccountListRequest) (*model.SuccessWithPagination, error)
	GetExternalAccountStatements(data model.BankAccountListRequest) (*model.SuccessWithPagination, error)

	CreateWebhookLog(logType string, jsonRequest string) error
}

type accountingService struct {
	repo repository.AccountingRepository
}

var invalidConfirmation = "Invalid confirmation password"
var invalidCurrentAdminId = "Invalid current user id"

var recordNotFound = "record not found"
var bankNotFound = "Bank not found"
var bankAccountNotFound = "Account not found"
var transactionNotFound = "Transsaction not found"
var transferNotFound = "Transfer not found"

func NewAccountingService(
	repo repository.AccountingRepository,
) AccountingService {
	return &accountingService{repo}
}

func (s *accountingService) CheckCurrentAdminId(input any) (*int64, error) {

	// input := c.MustGet("adminId")
	if input == nil {
		return nil, badRequest(invalidCurrentAdminId)
	}
	var adminId = int64(input.(float64))
	if adminId <= 0 {
		return nil, badRequest(invalidCurrentAdminId)
	}
	return &adminId, nil
}

func (s *accountingService) CheckCurrentUsername(input any) (*string, error) {

	// input := c.MustGet("username")
	if input == nil {
		return nil, badRequest(invalidCurrentAdminId)
	}
	var username = input.(string)
	// if username == "" {
	// 	return nil, badRequest(invalidCurrentAdminId)
	// }
	return &username, nil
}

func (s *accountingService) CheckConfirmationPassword(data model.ConfirmRequest) (*bool, error) {

	user, err := s.repo.GetAdminById(data.UserId)
	if err != nil {
		fmt.Println(data)
		return nil, notFound(invalidConfirmation)
	}
	if user == nil {
		return nil, badRequest(invalidConfirmation)
	}
	if err := helper.CompareAdminPassword(data.Password, user.Password); err != nil {
		return nil, badRequest(invalidConfirmation)
	}
	token := true
	return &token, nil
}

func (s *accountingService) GetBanks(params model.BankListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&params.Page, &params.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	records, err := s.repo.GetBanks(params)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return records, nil
}

func (s *accountingService) GetAccountTypes(params model.AccountTypeListRequest) (*model.SuccessWithPagination, error) {

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

	err := s.UpdateBankAccountBotStatusById(data.Id)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	record, err := s.repo.GetBankAccountById(data.Id)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, notFound(bankAccountNotFound)
		}
		if err.Error() == "Account not found" {
			return nil, notFound(bankAccountNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return record, nil
}

func (s *accountingService) GetBankAccounts(data model.BankAccountListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&data.Page, &data.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	s.UpdateAllBankAccountBotStatus()

	list, err := s.repo.GetBankAccounts(data)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return list, nil
}

func (s *accountingService) CreateBankAccount(data model.BankAccountCreateBody) error {

	bank, err := s.repo.GetBankById(data.BankId)
	if err != nil {
		fmt.Println(err)
		if err.Error() == recordNotFound {
			return notFound(bankNotFound)
		}
		return badRequest("Invalid Bank")
	}

	accountType, err := s.repo.GetAccounTypeById(data.AccountTypeId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid Account Type")
	}

	acNo := helper.StripAllButNumbers(data.AccountNumber)
	exist, err := s.repo.HasBankAccount(acNo)
	if err != nil {
		fmt.Println(err)
		return internalServerError(err.Error())
	}
	if exist {
		return badRequest("Account already exist")
	}

	var createBody model.BankAccountCreateBody
	createBody.BankId = bank.Id
	createBody.AccountTypeId = accountType.Id
	createBody.AccountName = data.AccountName
	createBody.AccountNumber = acNo
	createBody.AccountPriority = data.AccountPriority
	createBody.AutoCreditFlag = data.AutoCreditFlag
	createBody.AutoWithdrawFlag = data.AutoWithdrawFlag
	createBody.AutoTransferMaxAmount = data.AutoTransferMaxAmount
	createBody.AutoWithdrawMaxAmount = data.AutoWithdrawMaxAmount
	createBody.DeviceUid = data.DeviceUid
	// อัพเดทหลังจากเรียกบอท createBody.PinCode = data.PinCode
	createBody.QrWalletStatus = data.QrWalletStatus
	createBody.AccountStatus = data.AccountStatus
	createBody.AccountBalance = 0
	createBody.ConnectionStatus = "disconnected"
	if err := s.repo.CreateBankAccount(createBody); err != nil {
		return internalServerError(err.Error())
	}

	allowCreateExternalAccount := false
	config, _ := s.GetExternalAccountConfig("allow_create_external_account")
	if config != nil {
		if config.ConfigVal == "list" {
			accountConfig, errConfig := s.HasExternalAccountConfig("allow_external_account_number", acNo)
			if errConfig != nil {
				return nil
			}
			if accountConfig.ConfigVal == acNo {
				allowCreateExternalAccount = true
			}
		} else if config.ConfigVal == "all" {
			allowCreateExternalAccount = true
		}
	}

	if allowCreateExternalAccount && data.DeviceUid != "" && data.PinCode != "" && !s.HasExternalAccount(acNo) {
		if _, err := s.HasExternalAccountConfig("allow_external_account_number", acNo); err != nil {
			return nil
		}
		// FASTBANK
		var createExternalBody model.ExternalAccountCreateBody
		createExternalBody.AccountNo = acNo
		createExternalBody.BankCode = bank.Code
		createExternalBody.DeviceId = data.DeviceUid
		// ไม่ได้ใช้ createExternalBody.Password = data.Password
		createExternalBody.Pin = data.PinCode
		// ไม่ได้ใช้ createExternalBody.Username = data.Username
		// ไม่ได้ใช้ createExternalBody.WebhookNotifyUrl = os.Getenv("ACCOUNTING_LOCAL_WEBHOOK_ENDPOINT") + "/api/accounting/webhooks/noti"
		createExternalBody.WebhookUrl = os.Getenv("ACCOUNTING_LOCAL_WEBHOOK_ENDPOINT") + "/api/accounting/webhooks/action"
		if createResp, err := s.CreateExternalAccount(createExternalBody); err != nil {
			s.CreateWebhookLog("CreateBankAccount.CreateExternalAccount, ERROR", helper.StructJson(struct {
				data model.BankAccountCreateBody
				err  error
			}{data, err}))
			return internalServerError(err.Error())
		} else {
			// Update EncryptionPin
			account, err := s.repo.GetBankAccountByAccountNumber(acNo)
			if err != nil {
				s.CreateWebhookLog("CreateBankAccount.GetBankAccountByAccountNumber, ERROR", helper.StructJson(struct {
					data model.BankAccountCreateBody
					err  error
				}{data, err}))
				return internalServerError(err.Error())
			}
			var updateBody model.BankAccountUpdateBody
			updateBody.PinCode = &createResp.Pin
			updateBody.ExternalId = &createResp.Id
			if err := s.repo.UpdateBankAccount(account.Id, updateBody); err != nil {
				s.CreateWebhookLog("CreateBankAccount.UpdateBankAccount, ERROR", helper.StructJson(struct {
					data model.BankAccountUpdateBody
					err  error
				}{updateBody, err}))
				return internalServerError(err.Error())
			}
		}
	}

	return nil
}

func (s *accountingService) UpdateBankAccount(id int64, req model.BankAccountUpdateRequest) error {

	account, err := s.repo.GetBankAccountById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	var updateBody model.BankAccountUpdateBody
	var updateExBody model.ExternalAccountUpdateBody
	onExternalChange := false

	// Validate
	if req.BankId != nil && account.BankId != *req.BankId {
		bank, err := s.repo.GetBankById(*req.BankId)
		if err != nil {
			fmt.Println(err)
			if err.Error() == recordNotFound {
				return notFound(bankNotFound)
			}
			return badRequest("Invalid Bank")
		}
		updateBody.BankId = &bank.Id
		// onExternalChange = true
		updateExBody.BankCode = bank.Code
	}
	if req.AccountTypeId != nil && account.AccountTypeId != *req.AccountTypeId {
		accountType, err := s.repo.GetAccounTypeById(*req.AccountTypeId)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Account Type")
		}
		updateBody.AccountTypeId = &accountType.Id
	}
	if req.AccountName != nil && account.AccountName != *req.AccountName {
		updateBody.AccountName = req.AccountName
	}
	if req.AccountNumber != nil && account.AccountNumber != *req.AccountNumber {
		acNo := helper.StripAllButNumbers(*req.AccountNumber)
		if acNo != "" {
			check, err := s.repo.HasBankAccount(acNo)
			if err != nil {
				return internalServerError(err.Error())
			}
			if !check {
				fmt.Println(acNo)
				return notFound("Account already exist")
			}
			updateBody.AccountNumber = &acNo
			// onExternalChange = true
			account.AccountNumber = acNo
		} else {
			updateBody.AccountNumber = &account.AccountNumber
		}
	}
	if req.DeviceUid != nil && account.DeviceUid != *req.DeviceUid {
		updateBody.DeviceUid = req.DeviceUid
		// onExternalChange = true
		updateExBody.DeviceId = &account.DeviceUid
	}
	if req.PinCode != nil {
		// updateBody.PinCode = req.PinCode
		onExternalChange = true
		updateExBody.Pin = req.PinCode
	}
	if req.AutoCreditFlag != nil && account.AutoCreditFlag != *req.AutoCreditFlag {
		updateBody.AutoCreditFlag = req.AutoCreditFlag
	}
	if req.AutoWithdrawFlag != nil && account.AutoWithdrawFlag != *req.AutoWithdrawFlag {
		updateBody.AutoWithdrawFlag = req.AutoWithdrawFlag
	}
	if req.AutoWithdrawMaxAmount != nil && account.AutoWithdrawMaxAmount != *req.AutoWithdrawMaxAmount {
		updateBody.AutoWithdrawMaxAmount = req.AutoWithdrawMaxAmount
	}
	if req.AutoTransferMaxAmount != nil && account.AutoTransferMaxAmount != *req.AutoTransferMaxAmount {
		updateBody.AutoTransferMaxAmount = req.AutoTransferMaxAmount
	}
	if req.AccountPriority != nil && account.AccountPriority != *req.AccountPriority {
		updateBody.AccountPriority = req.AccountPriority
	}
	if req.QrWalletStatus != nil && account.QrWalletStatus != *req.QrWalletStatus {
		updateBody.QrWalletStatus = req.QrWalletStatus
	}
	if req.AccountStatus != nil && account.AccountStatus != *req.AccountStatus {
		updateBody.AccountStatus = req.AccountStatus
	}

	if onExternalChange {
		if updateExBody.DeviceId == nil {
			updateExBody.DeviceId = &account.DeviceUid
		}
		// if updateExBody.Pin == nil {
		// 	updateExBody.Pin = &account.PinCode
		// }
		// Create if not exist
		if !s.HasExternalAccount(account.AccountNumber) {
			var createExternalBody model.ExternalAccountCreateBody
			createExternalBody.AccountNo = account.AccountNumber
			createExternalBody.BankCode = account.BankCode
			createExternalBody.DeviceId = *updateExBody.DeviceId
			// ไม่ได้ใช้ createExternalBody.Password = data.Password
			createExternalBody.Pin = *updateExBody.Pin
			// ไม่ได้ใช้ createExternalBody.Username = data.Username
			// ไม่ได้ใช้ createExternalBody.WebhookNotifyUrl = os.Getenv("ACCOUNTING_LOCAL_WEBHOOK_ENDPOINT") + "/api/accounting/webhooks/noti"
			createExternalBody.WebhookUrl = os.Getenv("ACCOUNTING_LOCAL_WEBHOOK_ENDPOINT") + "/api/accounting/webhooks/action"
			if createResp, err := s.CreateExternalAccount(createExternalBody); err != nil {
				s.CreateWebhookLog("UpdateBankAccount.CreateExternalAccount, ERROR", helper.StructJson(struct {
					req model.BankAccountUpdateRequest
					err error
				}{req, err}))
				return internalServerError(err.Error())
			} else {
				// Update EncryptionPin
				updateBody.PinCode = &createResp.Pin
				updateBody.ExternalId = &createResp.Id
			}
		} else {
			updateExBody.AccountNo = account.AccountNumber
			// ไม่ได้ใช้ updateExBody.WebhookNotifyUrl = os.Getenv("ACCOUNTING_LOCAL_WEBHOOK_ENDPOINT") + "/api/accounting/webhooks/noti"
			updateExBody.WebhookUrl = os.Getenv("ACCOUNTING_LOCAL_WEBHOOK_ENDPOINT") + "/api/accounting/webhooks/action"
			if externalCreateResp, err := s.UpdateExternalAccount(updateExBody); err != nil {
				s.CreateWebhookLog("UpdateBankAccount, ERROR", helper.StructJson(struct {
					id  int64
					req model.BankAccountUpdateRequest
					err error
				}{id, req, err}))
				return internalServerError(err.Error())
			} else {
				// Update EncryptionPin
				updateBody.PinCode = &externalCreateResp.Pin
				updateBody.ExternalId = &externalCreateResp.Id
			}
		}
	}

	if err := s.repo.UpdateBankAccount(account.Id, updateBody); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) UpdateBankAccountBotStatusById(id int64) error {

	account, err := s.repo.GetBankAccountById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	now := time.Now()
	if account.LastConnUpdateAt != nil {
		// fmt.Println(now.Sub(*account.LastConnUpdateAt).Seconds())
		if now.Sub(*account.LastConnUpdateAt).Seconds() < 30 {
			return nil
		}
	}

	status_active := "active"
	status_disconnected := "disconnected"
	var data model.BankAccountUpdateBody
	data.LastConnUpdateAt = &now
	data.ConnectionStatus = &status_disconnected
	// data.AccountBalance = 0

	// FASTBANK
	if account.AccountNumber == "5014327339" {
		var query model.ExternalAccountStatusRequest
		query.AccountNumber = account.AccountNumber
		statusResp, err := s.GetExternalAccountStatus(query)
		if err != nil {
			return internalServerError(err.Error())
		}
		if statusResp.Status == "online" {
			data.ConnectionStatus = &status_active
		} else {
			fmt.Println("statusResp", statusResp)
			data.ConnectionStatus = &status_disconnected
		}

		balaceResp, err := s.GetExternalAccountBalance(query)
		if err != nil {
			return internalServerError(err.Error())
		}

		if balaceResp.AccountNo == account.AccountNumber {
			balance, _ := strconv.ParseFloat(strings.TrimSpace(balaceResp.AccountBalance), 64)
			data.AccountBalance = &balance
		} else {
			fmt.Println("ERROR, balaceResp: ", balaceResp)
			return internalServerError(err.Error())
		}
	}

	if err := s.repo.UpdateBankAccount(id, data); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}

func (s *accountingService) UpdateAllBankAccountBotStatus() error {

	var query model.BankAccountListRequest
	query.Limit = 100
	query.Page = 0
	accounts, err := s.repo.GetBotBankAccounts(query)
	if err != nil {
		return internalServerError(err.Error())
	}
	now := time.Now()
	error_delay := time.Now().Add(time.Minute * 5)
	status_active := "active"
	status_disconnected := "disconnected"
	for _, account := range accounts.List.([]model.BankAccountResponse) {

		if account.LastConnUpdateAt != nil {
			if now.Sub(*account.LastConnUpdateAt).Seconds() < 30 {
				continue
			}
		}
		var data model.BankAccountUpdateBody
		data.LastConnUpdateAt = &now
		data.ConnectionStatus = &status_disconnected
		// data.AccountBalance = 0

		// FASTBANK
		var query model.ExternalAccountStatusRequest
		query.AccountNumber = account.AccountNumber
		statusResp, err := s.GetExternalAccountStatus(query)
		if err != nil {
			data.LastConnUpdateAt = &error_delay
			// fmt.Println("ERROR", err.Error())
		} else {
			if statusResp.Status == "online" {
				data.ConnectionStatus = &status_active
			} else {
				fmt.Println("statusResp", statusResp)
				data.ConnectionStatus = &status_disconnected
			}
		}

		balaceResp, err := s.GetExternalAccountBalance(query)
		if err != nil {
			data.LastConnUpdateAt = &error_delay
			// fmt.Println("ERROR", err.Error())
		} else {
			if balaceResp.AccountNo == account.AccountNumber {
				balance, _ := strconv.ParseFloat(strings.TrimSpace(balaceResp.AccountBalance), 64)
				data.AccountBalance = &balance
			} else {
				data.LastConnUpdateAt = &error_delay
				// fmt.Println("ERROR, balaceResp: ", balaceResp)
			}
		}

		if err := s.repo.UpdateBankAccount(account.Id, data); err != nil {
			fmt.Println("ERROR, UPDATE ", err.Error())
		}
	}

	return nil
}

func (s *accountingService) DeleteBankAccount(id int64) error {

	account, err := s.repo.GetBankAccountById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if account.ExternalId != "" && s.HasExternalAccount(account.AccountNumber) {
		var query model.ExternalAccountStatusRequest
		query.AccountNumber = account.AccountNumber
		if err := s.DeleteExternalAccount(query); err != nil {
			return internalServerError(err.Error())
		}
	}

	var updateBody model.BankAccountDeleteBody
	updateBody.AccountNumber = fmt.Sprintf("%s_del%d", account.AccountNumber, account.Id)
	updateBody.DeletedAt = time.Now()
	if err := s.repo.DeleteBankAccount(id, updateBody); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) GetTransactionById(data model.BankAccountTransactionParam) (*model.BankAccountTransaction, error) {

	record, err := s.repo.GetTransactionById(data.Id)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(transactionNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return record, nil
}

func (s *accountingService) GetTransactions(data model.BankAccountTransactionListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&data.Page, &data.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	list, err := s.repo.GetTransactions(data)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return list, nil
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
	transaction.CreatedByUsername = data.CreatedByUsername

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

	record, err := s.repo.GetTransferById(data.Id)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(transferNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return record, nil
}

func (s *accountingService) GetTransfers(data model.BankAccountTransferListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&data.Page, &data.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	list, err := s.repo.GetTransfers(data)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return list, nil
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
	body.FromBankId = fromAccount.BankId
	body.FromAccountName = fromAccount.AccountName
	body.FromAccountNumber = fromAccount.AccountNumber
	body.ToAccountId = toAccount.Id
	body.ToBankId = toAccount.BankId
	body.ToAccountName = toAccount.AccountName
	body.ToAccountNumber = toAccount.AccountNumber
	body.Amount = data.Amount
	body.TransferAt = data.TransferAt
	body.CreatedByUsername = data.CreatedByUsername
	body.Status = "pending"

	if err := s.repo.CreateTransfer(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) ConfirmTransfer(id int64, actorId int64) error {

	transfer, err := s.repo.GetTransferById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if transfer.Status == "pending" {
		var body model.BankAccountTransferConfirmBody
		body.Status = "confirmed"
		body.ConfirmedAt = time.Now()
		body.ConfirmedByUserId = actorId
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

func (s *accountingService) HasExternalAccount(accountNumber string) bool {

	data, err := s.GetExternalAccounts()
	if err != nil {
		return true
	}

	for _, account := range data.List.([]model.ExternalAccount) {
		if account.AccountNo == accountNumber {
			return true
		}
	}
	return false
}

func (s *accountingService) GetExternalAccountConfig(key string) (*model.BotAccountConfig, error) {

	var query model.BotAccountConfigListRequest
	query.SearchKey = &key

	data, err := s.repo.GetBotaccountConfigs(query)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	for _, record := range data.List.([]model.BotAccountConfig) {
		if record.ConfigKey == key {
			return &record, nil
		}
	}
	return nil, notFound("Config not found")
}

func (s *accountingService) HasExternalAccountConfig(key string, value string) (*model.BotAccountConfig, error) {

	var query model.BotAccountConfigListRequest
	query.SearchKey = &key
	query.SearchValue = &value

	data, err := s.repo.GetBotaccountConfigs(query)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	for _, record := range data.List.([]model.BotAccountConfig) {
		if record.ConfigKey == key {
			return &record, nil
		}
	}
	return nil, notFound("Config not found")
}

func (s *accountingService) GetCustomerAccountsInfo(body model.CustomerAccountInfoRequest) (*model.CustomerAccountInfo, error) {

	botAccount, err := s.repo.GetActiveExternalAccount()
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	body.AccountFrom = botAccount.AccountNumber
	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		return nil, internalServerError("Error from JSON")
	}
	fmt.Println(string(b))

	client := &http.Client{}
	// curl -X POST "https://api.fastbankapi.com/api/v2/statement/verifyTransfer" -H "accept: */*" -H "apiKey: aa.bb" -H "Content-Type: application/json" -d "{ \"accountFrom\": \"cccc\", \"accountTo\": \"dddd\", \"bankCode\": \"bay\"}"
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/statement/verifyTransfer", bytes.NewBuffer(data))
	req.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Println(response)
		return nil, internalServerError("Error from external API")
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result model.CustomerAccountInfoReponse
	errJson := json.Unmarshal(responseData, &result)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}
	return &result.Data, nil
}

func (s *accountingService) GetExternalAccounts() (*model.SuccessWithPagination, error) {

	client := &http.Client{}
	req, _ := http.NewRequest("GET", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/site/bankAccount", nil)
	req.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	response, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Println(response)
		return nil, internalServerError("Error from external API")
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var list []model.ExternalAccount
	errJson := json.Unmarshal(responseData, &list)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = list
	result.Total = int64(len(list))
	return &result, nil
}

func (s *accountingService) GetExternalAccountBalance(query model.ExternalAccountStatusRequest) (*model.ExternalAccountBalance, error) {

	client := &http.Client{}
	// curl -X GET "https://api.fastbankapi.com/api/v2/statement/balance?accountNo=4281243019" -H "accept: */*" -H "apiKey: 559a37455f1b3f1ece5e7e452b75bed8.805cc14b876f857784acf00d78eedcb8"
	req, _ := http.NewRequest("GET", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/statement/balance?accountNo="+query.AccountNumber, nil)
	req.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	response, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Println(response)
		return nil, internalServerError("Error from external API")
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result model.ExternalAccountBalance
	errJson := json.Unmarshal(responseData, &result)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}
	if result.AccountNo != query.AccountNumber {
		s.CreateWebhookLog("GetExternalAccountBalance, ERROR:", string(responseData))
		return nil, notFound("Bank account not found")
	}
	return &result, nil
}

func (s *accountingService) GetExternalAccountStatus(query model.ExternalAccountStatusRequest) (*model.ExternalAccountStatus, error) {

	client := &http.Client{}
	req, _ := http.NewRequest("GET", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/site/bank-status?accountNo="+query.AccountNumber, nil)
	req.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	response, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != 200 {
		s.CreateWebhookLog("GetExternalAccountStatus, ERROR", helper.StructJson(struct {
			query        model.ExternalAccountStatusRequest
			responseJson string
		}{query, string(responseData)}))
		return nil, notFound("External account not found")
	}

	var result model.ExternalAccountStatus
	errJson := json.Unmarshal(responseData, &result)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}
	return &result, nil
}

func (s *accountingService) CreateExternalAccount(body model.ExternalAccountCreateBody) (*model.ExternalAccountCreateResponse, error) {

	client := &http.Client{}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/site/bankAccount", bytes.NewBuffer(data))
	req.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != 200 {
		s.CreateWebhookLog("CreateExternalAccount, ERROR:", string(responseData))
		return nil, internalServerError("Error from external API")
	}

	var result model.ExternalAccountCreateResponse
	errJson := json.Unmarshal(responseData, &result)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}
	jsonResult, err := json.Marshal(result)
	if err == nil {
		s.CreateWebhookLog("CreateExternalAccount, SUCCESS", string(jsonResult))
	}
	return &result, nil
}

func (s *accountingService) UpdateExternalAccount(body model.ExternalAccountUpdateBody) (*model.ExternalAccountCreateResponse, error) {

	client := &http.Client{}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/site/bankAccount", bytes.NewBuffer(data))
	req.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != 200 {
		s.CreateWebhookLog("UpdateExternalAccount, ERROR:", string(responseData))
		return nil, internalServerError("Error from external API")
	}
	var result model.ExternalAccountCreateResponse
	errJson := json.Unmarshal(responseData, &result)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}
	jsonResult, err := json.Marshal(result)
	if err == nil {
		s.CreateWebhookLog("UpdateExternalAccount, SUCCESS", string(jsonResult))
	}
	return &result, nil
}

func (s *accountingService) DeleteExternalAccount(query model.ExternalAccountStatusRequest) error {

	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/site/bankAccount/"+query.AccountNumber, nil)
	req.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	response, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	if response.StatusCode != 200 {
		fmt.Println(response)
		return internalServerError("Error from external API")
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	s.CreateWebhookLog("DeleteExternalAccount, responseData:", string(responseData))

	return nil
}

func (s *accountingService) EnableExternalAccount(body model.ExternalAccountEnableRequest) (*model.ExternalAccountStatus, error) {

	client := &http.Client{}
	// curl -X POST "https://api.fastbankapi.com/api/v2/site/enable-bank" -H "accept: */*" -H "apiKey: 123" -H "Content-Type: application/json" -d "{ \"accountNo\": \"string\", \"enable\": true}"
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/site/enable-bank", bytes.NewBuffer(data))
	req.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	if response.StatusCode != 200 {
		fmt.Println(response)
		return nil, internalServerError("Error from external API")
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("EnableExternalAccount:", string(responseData))
	// {"success":true,"enable":true,"status":"online"}
	// {"success":true,"enable":false,"status":"offline"}
	var result model.ExternalAccountStatus
	errJson := json.Unmarshal(responseData, &result)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}
	return &result, nil
}

func (s *accountingService) GetExternalAccountLogs(query model.BankAccountListRequest) (*model.SuccessWithPagination, error) {

	fmt.Println("BankAccountListRequest", query)

	client := &http.Client{}
	// curl -X GET "https://api.fastbankapi.com/api/v2/site/bankAccount/logs?accountNo=aaaaaaaaaaaaaa&page=0&size=10" -H "accept: */*" -H "apiKey: xxxxxxxxxx.yyyyyyyyyyy"
	queryString := fmt.Sprintf("&page=%d&size=%d", query.Page, query.Limit)
	fullPath := os.Getenv("ACCOUNTING_API_ENDPOINT") + "/api/v2/site/bankAccount/logs?accountNo=" + query.AccountNumber + queryString
	req, _ := http.NewRequest("GET", fullPath, nil)
	req.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	response, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Println(response)
		return nil, internalServerError("Error from external API")
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var externalList model.ExternalListWithPagination
	errJson := json.Unmarshal(responseData, &externalList)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}
	// fmt.Println("response", string(responseData))

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = externalList.Content
	result.Total = externalList.TotalElements
	return &result, nil
}

func (s *accountingService) GetExternalAccountStatements(query model.BankAccountListRequest) (*model.SuccessWithPagination, error) {

	fmt.Println("BankAccountListRequest", query)

	client := &http.Client{}
	// https://api.fastbankapi.com/api/v2/statement?accountNo=5014327339&page=0&size=10&txnCode=all
	// curl -X GET "https://api.fastbankapi.com/api/v2/statement?accountNo=4281243019&page=0&size=10&txnCode=all" -H "accept: */*" -H "apiKey: 559a37455f1b3f1ece5e7e452b75bed8.805cc14b876f857784acf00d78eedcb8"
	queryString := fmt.Sprintf("&page=%d&size=%d&txnCode=all", query.Page, query.Limit)
	fullPath := os.Getenv("ACCOUNTING_API_ENDPOINT") + "/api/v2/statement?accountNo=" + query.AccountNumber + queryString
	req, _ := http.NewRequest("GET", fullPath, nil)
	req.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	response, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	if response.StatusCode != 200 {
		fmt.Println(response)
		return nil, internalServerError("Error from external API")
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var externalList model.ExternalListWithPagination
	errJson := json.Unmarshal(responseData, &externalList)
	if errJson != nil {
		return nil, internalServerError("Error from JSON response")
	}
	// fmt.Println("response", string(responseData))

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = externalList.Content
	result.Total = externalList.TotalElements
	return &result, nil
}

func (s *accountingService) TransferExternalAccount(req model.ExternalAccountTransferRequest) error {

	var body model.ExternalAccountTransferBody
	systemAccount, err := s.repo.GetBankAccountById(req.SystemAccountId)
	if err != nil {
		if err.Error() == recordNotFound {
			return notFound(bankAccountNotFound)
		}
		return internalServerError(err.Error())
	}
	body.AccountForm = systemAccount.AccountNumber
	body.AccountTo = req.AccountNumber
	body.Amount = req.Amount
	body.BankCode = req.BankCode
	body.Pin = systemAccount.PinCode

	client := &http.Client{}
	// curl -X POST "https://api.fastbankapi.com/api/v2/statement/transfer" -H "accept: */*" -H "apiKey: xxxxxxxxxx.yyyyyyyyyyy"
	//-H "Content-Type: application/json" -d "{ \"accountFrom\": \"aaaaaaaaaaaaaaaa\", \"accountTo\": \"bbbbbbbbbbbbbb\", \"amount\": \"8\", \"bankCode\": \"bay\", \"pin\": \"ccccc\"}"
	data, _ := json.Marshal(body)
	reqHttp, _ := http.NewRequest("POST", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/statement/transfer", bytes.NewBuffer(data))
	reqHttp.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	reqHttp.Header.Set("Content-Type", "application/json")
	response, err := client.Do(reqHttp)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("response", string(responseData))

	if response.StatusCode != 200 {
		var errorModel model.ExternalAccountError
		errJson := json.Unmarshal(responseData, &errorModel)
		if errJson != nil {
			return internalServerError("Error from JSON response")
		}
		fmt.Println("errorModel", errorModel)
		if errorModel.Error != "" {
			return internalServerError(errorModel.Error)
		}
		return internalServerError("Error from external API")
	}
	return nil
}

func (s *accountingService) CreateBankStatementFromWebhook(data model.WebhookStatement) error {

	systemAccount, err := s.repo.GetBankAccountByExternalId(data.BankAccountId)
	if err != nil {
		fmt.Println(err)
		return badRequest("Invalid Bank Account")
	}

	_, errOldStatement := s.repo.GetWebhookStatementByExternalId(data.Id)
	if errOldStatement != nil && errOldStatement.Error() == recordNotFound {
		var bodyCreateState model.BankStatementCreateBody
		bodyCreateState.AccountId = systemAccount.Id
		bodyCreateState.ExternalId = data.Id
		if data.TxnCode == "X1" || data.TxnCode == "CR" {
			bodyCreateState.StatementType = "transfer_in"
			bodyCreateState.Amount = data.Amount
		} else if data.TxnCode == "X2" || data.TxnCode == "DR" {
			bodyCreateState.StatementType = "transfer_out"
			bodyCreateState.Amount = data.Amount * -1
		} else {
			return badRequest("Invalid TxnCode")
		}

		bankId, _ := s.GetBankIdFromWebhook(data)
		bodyCreateState.FromBankId = bankId

		accountNumber, _ := s.GetAccountNoFromWebhook(data)
		bodyCreateState.FromAccountNumber = accountNumber

		bodyCreateState.Detail = data.TxnDescription + " " + data.Info
		bodyCreateState.TransferAt = data.DateTime
		bodyCreateState.Status = "pending"

		insertId, err := s.repo.CreateWebhookStatement(bodyCreateState)
		if err != nil {
			return internalServerError(err.Error())
		}

		// Auto Match if == 1
		var reqPosibleList model.MemberPossibleListRequest
		statement, err := s.repo.GetBankStatementById(*insertId)
		if err != nil {
			return nil
		}
		reqPosibleList.UnknownStatementId = statement.Id
		reqPosibleList.UserBankId = &statement.FromBankId
		reqPosibleList.UserAccountNumber = &statement.FromAccountNumber

		records, err := s.repo.GetPossibleStatementOwners(reqPosibleList)
		if err != nil {
			return nil
		}
		if records.Total == 1 {
			for _, record := range records.List.([]model.Member) {
				// Auto create transaction
				if bodyCreateState.StatementType == "transfer_in" {
					var createDepositBody model.BankTransactionCreateBody
					createDepositBody.MemberCode = record.MemberCode
					createDepositBody.TransferType = "deposit"
					createDepositBody.CreditAmount = bodyCreateState.Amount
					createDepositBody.TransferAt = bodyCreateState.TransferAt
					createDepositBody.IsAutoCredit = true
					// promotionId  bonusAmount
					createDepositBody.ToAccountId = &systemAccount.Id
					transId, err := s.CreateBankTransaction(createDepositBody)
					if err != nil {
						return internalServerError(err.Error())
					}
					var body model.BankConfirmDepositRequest
					body.TransferAt = &bodyCreateState.TransferAt
					body.BonusAmount = 0
					if err := s.ConfirmDepositTransaction(*transId, body); err != nil {
						return internalServerError(err.Error())
					}
				} else if bodyCreateState.StatementType == "transfer_out" {
					var createWithdrawBody model.BankTransactionCreateBody
					createWithdrawBody.MemberCode = record.MemberCode
					createWithdrawBody.TransferType = "withdraw"
					createWithdrawBody.CreditAmount = bodyCreateState.Amount
					createWithdrawBody.TransferAt = bodyCreateState.TransferAt
					createWithdrawBody.FromAccountId = &systemAccount.Id
					transId, err := s.CreateBankTransaction(createWithdrawBody)
					if err != nil {
						return internalServerError(err.Error())
					}
					var body model.BankConfirmWithdrawRequest
					body.CreditAmount = bodyCreateState.Amount
					body.BankChargeAmount = 0
					if err := s.ConfirmWithdrawTransaction(*transId, body); err != nil {
						return internalServerError(err.Error())
					}
				}
			}
			var body model.BankStatementUpdateBody
			body.Status = "confirmed"
		}
	}

	return nil
}

func (s *accountingService) CreateBankTransaction(data model.BankTransactionCreateBody) (*int64, error) {

	var body model.BankTransactionCreateBody

	if data.TransferType == "deposit" {
		member, err := s.repo.GetUserByMemberCode(data.MemberCode)
		if err != nil {
			fmt.Println(err)
			return nil, badRequest("Invalid Member code")
		}
		bank, err := s.repo.GetBankByCode(member.Bankname)
		if err != nil {
			fmt.Println(err)
			return nil, badRequest("Invalid User Bank")
		}
		body.MemberCode = member.MemberCode
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
			return nil, badRequest("Input Bank Account")
		}
		toAccount, err := s.repo.GetDepositAccountById(*data.ToAccountId)
		if err != nil {
			fmt.Println(err)
			return nil, badRequest("Invalid Bank Account")
		}
		body.ToAccountId = &toAccount.Id
		body.ToBankId = &toAccount.BankId
		body.ToAccountName = &toAccount.AccountName
		body.ToAccountNumber = &toAccount.AccountNumber

		// todo: createBonus + refDeposit
		body.PromotionId = data.PromotionId

	} else if data.TransferType == "withdraw" {
		member, err := s.repo.GetUserByMemberCode(data.MemberCode)
		if err != nil {
			fmt.Println(err)
			return nil, badRequest("Invalid Member code")
		}
		bank, err := s.repo.GetBankByCode(member.Bankname)
		if err != nil {
			fmt.Println(err)
			return nil, badRequest("Invalid User Bank")
		}
		body.MemberCode = member.MemberCode
		body.UserId = member.Id
		body.CreditAmount = data.CreditAmount
		body.TransferType = data.TransferType

		fromAccount, err := s.repo.GetWithdrawAccountById(*data.FromAccountId)
		if err != nil {
			fmt.Println(err)
			return nil, badRequest("Invalid Bank Account")
		}
		body.FromAccountId = &fromAccount.Id
		body.FromBankId = &fromAccount.BankId
		body.FromAccountName = &fromAccount.AccountName
		body.FromAccountNumber = &fromAccount.AccountNumber

		body.ToBankId = &bank.Id
		body.ToAccountName = &member.Fullname
		body.ToAccountNumber = &member.BankAccount
	} else {
		return nil, badRequest("Invalid Transfer Type")
	}

	body.TransferAt = data.TransferAt
	body.CreatedByUserId = data.CreatedByUserId
	body.CreatedByUsername = data.CreatedByUsername
	body.Status = "pending"

	insertId, err := s.repo.CreateBankTransaction(body)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return insertId, nil
}

func (s *accountingService) ConfirmDepositTransaction(id int64, req model.BankConfirmDepositRequest) error {

	record, err := s.repo.GetBankTransactionById(id)
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

	var createData model.BankTransactionCreateConfirmBody
	createData.TransactionId = record.Id
	createData.UserId = record.UserId
	createData.TransferType = record.TransferType
	createData.FromAccountId = record.FromAccountId
	createData.ToAccountId = record.ToAccountId
	createData.JsonBefore = string(jsonBefore)
	if req.TransferAt == nil {
		createData.TransferAt = record.TransferAt
	} else {
		TransferAt := req.TransferAt
		createData.TransferAt = *TransferAt
		updateData.TransferAt = *TransferAt
	}
	createData.SlipUrl = req.SlipUrl
	createData.BonusAmount = req.BonusAmount
	createData.ConfirmedAt = req.ConfirmedAt
	createData.ConfirmedByUserId = req.ConfirmedByUserId
	createData.ConfirmedByUsername = req.ConfirmedByUsername
	if err := s.repo.CreateConfirmTransaction(createData); err != nil {
		return internalServerError(err.Error())
	}
	if err := s.repo.ConfirmPendingTransaction(id, updateData); err != nil {
		return internalServerError(err.Error())
	}
	if err := s.IncreaseMemberCredit(record.UserId, record.CreditAmount); err != nil {
		return internalServerError(err.Error())
	}
	// todo: Bonus
	// commit
	return nil
}

func (s *accountingService) ConfirmWithdrawTransaction(id int64, req model.BankConfirmWithdrawRequest) error {

	record, err := s.repo.GetBankTransactionById(id)
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

	var createData model.BankTransactionCreateConfirmBody
	createData.TransactionId = record.Id
	createData.UserId = record.UserId
	createData.TransferType = record.TransferType
	createData.FromAccountId = record.FromAccountId
	createData.ToAccountId = record.ToAccountId
	createData.JsonBefore = string(jsonBefore)
	createData.TransferAt = record.TransferAt
	createData.CreditAmount = req.CreditAmount
	createData.BankChargeAmount = req.BankChargeAmount
	createData.ConfirmedAt = req.ConfirmedAt
	createData.ConfirmedByUserId = req.ConfirmedByUserId
	createData.ConfirmedByUsername = req.ConfirmedByUsername
	if err := s.repo.CreateConfirmTransaction(createData); err != nil {
		return internalServerError(err.Error())
	}
	if err := s.repo.ConfirmPendingTransaction(id, updateData); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) IncreaseMemberCredit(userId int64, creditAmount float32) error {

	// record, err := s.repoBanking.GetBankTransactionById(id)
	// if err != nil {
	// 	return internalServerError(err.Error())
	// }

	if err := s.repo.IncreaseMemberCredit(userId, creditAmount); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}

func (s *accountingService) GetBankIdFromWebhook(data model.WebhookStatement) (int64, error) {

	// todo :

	return 4, nil
}

func (s *accountingService) GetAccountNoFromWebhook(data model.WebhookStatement) (string, error) {

	// todo :

	return "3002", nil
}

func (s *accountingService) CreateWebhookLog(logType string, jsonRequest string) error {

	var body model.WebhookLogCreateBody
	body.JsonRequest = jsonRequest
	body.JsonPayload = "{}"
	body.LogType = logType
	body.Status = "success"

	if err := s.repo.CreateWebhookLog(body); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}

func (s *accountingService) CreateBotaccountConfig(data model.BotAccountConfigCreateBody) error {

	if err := s.repo.CreateBotaccountConfig(data); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}
