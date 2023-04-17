package service

import (
	"bytes"
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	UpdateBankAccount(id int64, data model.BankAccountUpdateBody) error
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

	GetExternalBankAccounts(data model.BankAccountListRequest) (*model.SuccessWithPagination, error)
	GetExternalBankAccountBalance(query model.ExternalBankAccountStatusRequest) (*model.ExternalBankAccountBalance, error)
	GetExternalBankAccountStatus(query model.ExternalBankAccountStatusRequest) (*model.ExternalBankAccountStatus, error)
	CreateExternalBankAccount(data model.ExternalBankAccountCreateBody) error
	UpdateExternalBankAccount(data model.ExternalBankAccountCreateBody) error
	EnableExternalBankAccount(query model.ExternalBankAccountEnableRequest) (*model.ExternalBankAccountStatus, error)
	DeleteExternalBankAccount(query model.ExternalBankAccountStatusRequest) error

	GetExternalBankAccountsLogs(data model.BankAccountListRequest) (*model.SuccessWithPagination, error)
	GetExternalBankStatements(data model.BankAccountListRequest) (*model.SuccessWithPagination, error)

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

	s.UpdateBankAccountBotStatus(data.Id)

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

func (s *accountingService) GetBankAccounts(data model.BankAccountListRequest) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&data.Page, &data.Limit); err != nil {
		return nil, badRequest(err.Error())
	}
	accounting, err := s.repo.GetBankAccounts(data)
	if err != nil {
		return nil, internalServerError(err.Error())
	}
	return accounting, nil
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

	var account model.BankAccountCreateBody
	account.BankId = bank.Id
	account.AccountTypeId = accountType.Id
	account.AccountName = data.AccountName
	account.AccountNumber = acNo
	account.AccountPriority = data.AccountPriority
	account.AutoCreditFlag = data.AutoCreditFlag
	account.AutoWithdrawFlag = data.AutoWithdrawFlag
	account.AutoTransferMaxAmount = data.AutoTransferMaxAmount
	account.AutoWithdrawMaxAmount = data.AutoWithdrawMaxAmount
	account.DeviceUid = data.DeviceUid
	account.PinCode = data.PinCode
	account.QrWalletStatus = data.QrWalletStatus
	account.AccountStatus = data.AccountStatus
	account.AccountBalance = 0
	account.ConnectionStatus = "disconnected"

	if err := s.repo.CreateBankAccount(account); err != nil {
		return internalServerError(err.Error())
	}

	// FASTBANK
	// fmt.Println(acNo)
	if acNo == "5014327339" {
		// AccountNo        string `json:"accountNo"`
		// BankCode         string `json:"bankCode"`
		// DeviceId         string `json:"deviceId"`
		// Password         string `json:"password"`
		// Pin              string `json:"pin"`
		// Username         string `json:"username"`
		// WebhookNotifyUrl string `json:"webhookNotifyUrl"`
		// WebhookUrl       string `json:"webhookUrl"`
		var body model.ExternalBankAccountCreateBody
		body.AccountNo = acNo
		body.BankCode = bank.Code
		body.DeviceId = data.DeviceUid
		// body.Password = data.Password
		body.Pin = data.PinCode
		// body.Username = data.Username
		body.WebhookNotifyUrl = os.Getenv("ACCOUNTING_LOCAL_WEBHOOK_ENDPOINT") + "/api//accounting/webhooks/noti"
		body.WebhookUrl = os.Getenv("ACCOUNTING_LOCAL_WEBHOOK_ENDPOINT") + "/api/accounting/webhooks/action"

		if err := s.CreateExternalBankAccount(body); err != nil {
			// todo: delete created account
			return internalServerError(err.Error())
		}
		// fmt.Println(body)
	}

	return nil
}

func (s *accountingService) UpdateBankAccount(id int64, data model.BankAccountUpdateBody) error {

	account, err := s.repo.GetBankAccountById(id)
	if err != nil {
		return internalServerError(err.Error())
	}
	var updateExBody model.ExternalBankAccountCreateBody

	// Validate
	if data.BankId != 0 && account.BankId != data.BankId {
		bank, err := s.repo.GetBankById(data.BankId)
		if err != nil {
			fmt.Println(err)
			if err.Error() == recordNotFound {
				return notFound(bankNotFound)
			}
			return badRequest("Invalid Bank")
		}
		data.BankId = bank.Id
		updateExBody.BankCode = bank.Code
	}
	if account.AccountTypeId != data.AccountTypeId {
		accountType, err := s.repo.GetAccounTypeById(data.AccountTypeId)
		if err != nil {
			fmt.Println(err)
			return badRequest("Invalid Account Type")
		}
		data.AccountTypeId = accountType.Id
	}
	acNo := helper.StripAllButNumbers(data.AccountNumber)
	if acNo != "" && account.AccountNumber != acNo {
		check, err := s.repo.HasBankAccount(acNo)
		if err != nil {
			return internalServerError(err.Error())
		}
		if !check {
			fmt.Println(acNo)
			return notFound("Account already exist")
		}
	}

	if err := s.repo.UpdateBankAccount(id, data); err != nil {
		return internalServerError(err.Error())
	}

	// FASTBANK
	if acNo == "5014327339" {
		// AccountNo        string `json:"accountNo"`
		// BankCode         string `json:"bankCode"`
		// DeviceId         string `json:"deviceId"`
		// Password         string `json:"password"`
		// Pin              string `json:"pin"`
		// Username         string `json:"username"`
		// WebhookNotifyUrl string `json:"webhookNotifyUrl"`
		// WebhookUrl       string `json:"webhookUrl"`
		// updateExBody.AccountNo = acNo
		updateExBody.DeviceId = data.DeviceUid
		// body.Password = data.Password
		updateExBody.Pin = data.PinCode
		// body.Username = data.Username
		updateExBody.WebhookNotifyUrl = "http://143.198.211.247:3001/api/accounting/bankaccounts2/list"
		updateExBody.WebhookUrl = "http://143.198.211.247:3001/api/accounting/bankaccounts2/list"
		if err := s.UpdateExternalBankAccount(updateExBody); err != nil {
			// todo: delete created account
			return internalServerError(err.Error())
		}
	}

	return nil
}

func (s *accountingService) UpdateBankAccountBotStatus(id int64) error {

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

	var data model.BankAccountUpdateBody
	data.LastConnUpdateAt = &now
	data.ConnectionStatus = "disconnected"
	// data.AccountBalance = 0

	// FASTBANK
	if account.AccountNumber == "5014327339" {
		var query model.ExternalBankAccountStatusRequest
		query.AccountNumber = account.AccountNumber
		statusResp, err := s.GetExternalBankAccountStatus(query)
		if err != nil {
			return internalServerError(err.Error())
		}
		if statusResp.Status == "online" {
			data.ConnectionStatus = "active"
		} else {
			fmt.Println("statusResp", statusResp)
			data.ConnectionStatus = "disconnected"
		}

		balaceResp, err := s.GetExternalBankAccountBalance(query)
		if err != nil {
			return internalServerError(err.Error())
		}

		if balaceResp.AccountNo == account.AccountNumber {
			balance, _ := strconv.ParseFloat(strings.TrimSpace(balaceResp.AccountBalance), 64)
			data.AccountBalance = balance
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

func (s *accountingService) UpdateAllBankAccountBotStatus(id int64) error {

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

	var data model.BankAccountUpdateBody
	data.LastConnUpdateAt = &now
	data.ConnectionStatus = "disconnected"
	// data.AccountBalance = 0

	// FASTBANK
	if account.AccountNumber == "5014327339" {
		var query model.ExternalBankAccountStatusRequest
		query.AccountNumber = account.AccountNumber
		statusResp, err := s.GetExternalBankAccountStatus(query)
		if err != nil {
			return internalServerError(err.Error())
		}
		if statusResp.Status == "online" {
			data.ConnectionStatus = "active"
		} else {
			fmt.Println("statusResp", statusResp)
			data.ConnectionStatus = "disconnected"
		}

		balaceResp, err := s.GetExternalBankAccountBalance(query)
		if err != nil {
			return internalServerError(err.Error())
		}

		if balaceResp.AccountNo == account.AccountNumber {
			balance, _ := strconv.ParseFloat(strings.TrimSpace(balaceResp.AccountBalance), 64)
			data.AccountBalance = balance
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
		if err.Error() == recordNotFound {
			return nil, notFound(transactionNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return accounting, nil
}

func (s *accountingService) GetTransactions(data model.BankAccountTransactionListRequest) (*model.SuccessWithPagination, error) {

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

	accounting, err := s.repo.GetTransferById(data.Id)
	if err != nil {
		if err.Error() == recordNotFound {
			return nil, notFound(transferNotFound)
		}
		return nil, internalServerError(err.Error())
	}
	return accounting, nil
}

func (s *accountingService) GetTransfers(data model.BankAccountTransferListRequest) (*model.SuccessWithPagination, error) {

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

func (s *accountingService) GetExternalBankAccounts(data model.BankAccountListRequest) (*model.SuccessWithPagination, error) {

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
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var list []model.ExternalBankAccount
	json.Unmarshal(responseData, &list)

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = list
	result.Total = int64(len(list))
	return &result, nil
}

func (s *accountingService) GetExternalBankAccountBalance(query model.ExternalBankAccountStatusRequest) (*model.ExternalBankAccountBalance, error) {

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
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result model.ExternalBankAccountBalance
	json.Unmarshal(responseData, &result)
	if result.AccountNo != query.AccountNumber {
		fmt.Println("response", string(responseData))
		return nil, notFound("Bank account not found")
	}
	return &result, nil
}

func (s *accountingService) GetExternalBankAccountStatus(query model.ExternalBankAccountStatusRequest) (*model.ExternalBankAccountStatus, error) {

	client := &http.Client{}
	req, _ := http.NewRequest("GET", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/site/bank-status?accountNo="+query.AccountNumber, nil)
	req.Header.Set("apiKey", os.Getenv("ACCOUNTING_API_KEY"))
	response, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	if response.StatusCode != 200 {
		fmt.Println(response)
		return nil, notFound("Bank account not found")
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result model.ExternalBankAccountStatus
	json.Unmarshal(responseData, &result)
	return &result, nil
}

func (s *accountingService) CreateExternalBankAccount(body model.ExternalBankAccountCreateBody) error {

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
	if response.StatusCode != 200 {
		fmt.Println(response)
		return internalServerError("Error from external API")
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result model.ExternalBankAccountCreateResponse
	json.Unmarshal(responseData, &result)
	fmt.Println("response", result)

	return nil
}

func (s *accountingService) UpdateExternalBankAccount(body model.ExternalBankAccountCreateBody) error {

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
	if response.StatusCode != 200 {
		fmt.Println(response)
		return internalServerError("Error from external API")
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result model.ExternalBankAccountCreateResponse
	json.Unmarshal(responseData, &result)
	fmt.Println("response", result)

	return nil
}

func (s *accountingService) DeleteExternalBankAccount(query model.ExternalBankAccountStatusRequest) error {

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
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("response", string(responseData))

	return nil
}

func (s *accountingService) EnableExternalBankAccount(body model.ExternalBankAccountEnableRequest) (*model.ExternalBankAccountStatus, error) {

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
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("EnableExternalBankAccount:", string(responseData))
	// {"success":true,"enable":true,"status":"online"}
	// {"success":true,"enable":false,"status":"offline"}
	var result model.ExternalBankAccountStatus
	json.Unmarshal(responseData, &result)
	return &result, nil
}

func (s *accountingService) GetExternalBankAccountsLogs(query model.BankAccountListRequest) (*model.SuccessWithPagination, error) {

	client := &http.Client{}
	// curl -X GET "https://api.fastbankapi.com/api/v2/site/bankAccount/logs?accountNo=123&page=0&size=10" -H "accept: */*" -H "apiKey: 123"
	queryString := fmt.Sprintf("&page=%d&size=%d", query.Page, query.Limit)
	req, _ := http.NewRequest("GET", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/site/bankAccount/logs?accountNo="+query.AccountNumber+queryString, nil)
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
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var list []model.ExternalBankAccount
	// json.Unmarshal(responseData, &list)
	fmt.Println("response", string(responseData))

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = list
	result.Total = int64(len(list))
	return &result, nil
}

func (s *accountingService) GetExternalBankStatements(query model.BankAccountListRequest) (*model.SuccessWithPagination, error) {

	client := &http.Client{}
	// curl -X GET "https://api.fastbankapi.com/api/v2/statement?accountNo=4281243019&page=0&size=10&txnCode=all" -H "accept: */*" -H "apiKey: 559a37455f1b3f1ece5e7e452b75bed8.805cc14b876f857784acf00d78eedcb8"
	queryString := fmt.Sprintf("&page=%d&size=%d&txnCode=all", query.Page, query.Limit)
	req, _ := http.NewRequest("GET", os.Getenv("ACCOUNTING_API_ENDPOINT")+"/api/v2/statement?accountNo="+query.Search+queryString, nil)
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
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var list []model.ExternalBankAccount
	// json.Unmarshal(responseData, &list)
	fmt.Println("response", string(responseData))

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = list
	result.Total = int64(len(list))
	return &result, nil
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
