package repository

import (
	"cybergame-api/model"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func NewAccountingRepository(db *gorm.DB) AccountingRepository {
	return &repo{db}
}

type AccountingRepository interface {
	GetAdminById(id int64) (*model.Admin, error)

	GetBanks(req model.BankListRequest) (*model.SuccessWithPagination, error)
	GetBankById(id int64) (*model.Bank, error)
	GetBankByCode(code string) (*model.Bank, error)

	GetAccountTypes(req model.AccountTypeListRequest) (*model.SuccessWithPagination, error)
	GetAccounTypeById(id int64) (*model.AccountType, error)

	GetUserByMemberCode(memberCode string) (*model.TempUser, error)

	HasBankAccount(accountNumber string) (bool, error)
	GetBankAccountById(id int64) (*model.BankAccount, error)
	GetBankAccounts(data model.BankAccountListRequest) (*model.SuccessWithPagination, error)
	CreateBankAccount(data model.BankAccountCreateBody) error
	UpdateBankAccount(id int64, data model.BankAccountUpdateBody) error
	DeleteBankAccount(id int64) error

	GetTransactionById(id int64) (*model.BankAccountTransaction, error)
	GetTransactions(data model.BankAccountTransactionListRequest) (*model.SuccessWithPagination, error)
	CreateTransaction(data model.BankAccountTransactionBody) error
	UpdateTransaction(id int64, data model.BankAccountTransactionBody) error
	DeleteTransaction(id int64) error

	GetTransferById(id int64) (*model.BankAccountTransfer, error)
	GetTransfers(data model.BankAccountTransferListRequest) (*model.SuccessWithPagination, error)
	CreateTransfer(data model.BankAccountTransferBody) error
	ConfirmTransfer(id int64, data model.BankAccountTransferConfirmBody) error
	DeleteTransfer(id int64) error
}

func (r repo) GetAdminById(id int64) (*model.Admin, error) {
	var admin model.Admin

	if err := r.db.Table("Admins").
		Select("id, username, phone, password, email, role").
		Where("id = ?", id).
		First(&admin).
		Error; err != nil {
		return nil, err
	}
	// "record not found"
	// if admin.Id == 0 {
	// 	return nil, errors.New("Admin not found")
	// }
	return &admin, nil
}

func (r repo) GetBanks(req model.BankListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Banks")
	count = count.Select("id")
	if req.Search != "" {
		count = count.Where("code = ?", req.Search)
	}
	if err = count.
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if total > 0 {
		// SELECT //
		query := r.db.Table("Banks")
		query = query.Select("id, name, code, icon_url, icon_url, type_flag")
		if req.Search != "" {
			query = query.Where("code = ?", req.Search)
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}
		if err = query.
			Limit(req.Limit).
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	if list == nil {
		list = []model.BankResponse{}
	}
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetBankById(id int64) (*model.Bank, error) {

	var result *model.Bank
	if err := r.db.Table("Banks").
		Select("id, name, code, icon_url, type_flag").
		Where("id = ?", id).
		First(&result).
		Error; err != nil {
		return nil, err
	}

	if result.Id == 0 {
		return nil, errors.New("Bank not found")
	}
	return result, nil
}

func (r repo) GetBankByCode(code string) (*model.Bank, error) {

	var result *model.Bank
	if err := r.db.Table("Banks").
		Select("id, name, code, icon_url, type_flag").
		Where("code = ?", code).
		First(&result).
		Error; err != nil {
		return nil, err
	}

	if result.Id == 0 {
		return nil, errors.New("Bank not found")
	}
	return result, nil
}

func (r repo) GetAccountTypes(req model.AccountTypeListRequest) (*model.SuccessWithPagination, error) {

	var list []model.AccountTypeResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_account_types")
	count = count.Select("id")
	if req.Search != "" {
		count = count.Where("name = ?", req.Search)
	}
	if err = count.
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if total > 0 {
		// SELECT //
		query := r.db.Table("Bank_account_types")
		query = query.Select("id, name, limit_flag")
		if req.Search != "" {
			query = query.Where("name = ?", req.Search)
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}
		if err = query.
			Limit(req.Limit).
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	if list == nil {
		list = []model.AccountTypeResponse{}
	}
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetAccounTypeById(id int64) (*model.AccountType, error) {

	var result *model.AccountType
	if err := r.db.Table("Bank_account_types").
		Select("id, name, limit_flag").
		Where("id = ?", id).
		First(&result).
		Error; err != nil {
		return nil, err
	}

	if result.Id == 0 {
		return nil, errors.New("Account type not found")
	}
	return result, nil
}

func (r repo) GetUserByMemberCode(memberCode string) (*model.TempUser, error) {

	var user model.TempUser
	user.Id = 88
	user.MemberCode = memberCode
	user.BankId = 2
	user.AccountName = "MOCK TEMP ACCOUNT NAME"
	user.AccountNumber = "000-111-222-333"

	return &user, nil
}

func (r repo) HasBankAccount(accountNumber string) (bool, error) {
	var count int64
	if err := r.db.Table("Bank_accounts").
		Select("id").
		Where("account_number = ?", accountNumber).
		Where("deleted_at IS NULL").
		Limit(1).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r repo) GetBankAccountById(id int64) (*model.BankAccount, error) {

	var accounting model.BankAccount
	selectedFields := "accounts.id, accounts.bank_id, accounts.account_type_id, accounts.account_name, accounts.account_number, accounts.account_balance, accounts.account_priority, accounts.account_status, accounts.created_at, accounts.updated_at"
	selectedFields += ",banks.name as bank_name, banks.code, banks.icon_url as bank_icon_url, banks.type_flag"
	selectedFields += ",account_types.name as account_type_name, account_types.limit_flag"
	if err := r.db.Table("Bank_accounts as accounts").
		Select(selectedFields).
		Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id").
		Joins("LEFT JOIN Bank_account_types AS account_types ON account_types.id = accounts.account_type_id").
		Where("accounts.id = ?", id).
		Where("accounts.deleted_at IS NULL").
		First(&accounting).
		Error; err != nil {
		return nil, err
	}

	if accounting.Id == 0 {
		return nil, errors.New("Account not found")
	}
	return &accounting, nil
}

func (r repo) GetBankAccounts(req model.BankAccountListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankAccountResponse
	var total int64
	var err error

	// SELECT //
	query := r.db.Table("Bank_accounts AS accounts")
	selectedFields := "accounts.id, accounts.bank_id, accounts.account_type_id, accounts.account_name, accounts.account_number, accounts.account_balance, accounts.account_priority, accounts.account_status, accounts.created_at, accounts.updated_at"
	selectedFields += ",banks.name as bank_name, banks.code, banks.icon_url as bank_icon_url, banks.type_flag"
	selectedFields += ",account_types.name as account_type_name, account_types.limit_flag"
	query = query.Select(selectedFields)
	query = query.Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id")
	query = query.Joins("LEFT JOIN Bank_account_types AS account_types ON account_types.id = accounts.account_type_id")
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		query = query.Where("accounts.account_name LIKE ?", search_like).Or("accounts.account_number LIKE ?", search_like)
	}

	// Sort by ANY //
	req.SortCol = strings.TrimSpace(req.SortCol)
	if req.SortCol != "" {
		if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
			req.SortAsc = "DESC"
		} else {
			req.SortAsc = "ASC"
		}
		query = query.Order(req.SortCol + " " + req.SortAsc)
	}

	if err = query.
		Where("accounts.deleted_at IS NULL").
		Limit(req.Limit).
		Offset(req.Page * req.Limit).
		Scan(&list).
		Error; err != nil {
		return nil, err
	}

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_accounts")
	count = count.Select("id")
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where("account_name LIKE ?", search_like).Or("account_number LIKE ?", search_like)
	}
	if err = count.
		Where("deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) CreateBankAccount(data model.BankAccountCreateBody) error {
	if err := r.db.Table("Bank_accounts").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) UpdateBankAccount(id int64, data model.BankAccountUpdateBody) error {
	if err := r.db.Table("Bank_accounts").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DeleteBankAccount(id int64) error {
	if err := r.db.Table("Bank_accounts").Where("id = ?", id).Delete(&model.BankAccount{}).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) GetTransactionById(id int64) (*model.BankAccountTransaction, error) {
	var record model.BankAccountTransaction
	selectedFields := "transactions.id, transactions.account_id, transactions.description, transactions.transfer_type, transactions.amount, transactions.transfer_at, transactions.created_by_username, transactions.created_at, transactions.updated_at"
	selectedFields += ",accounts.bank_id, accounts.account_type_id, accounts.account_name, accounts.account_number, accounts.account_balance, accounts.account_priority, accounts.account_status, accounts.created_at, accounts.updated_at"
	selectedFields += ",banks.name as bank_name, banks.code, banks.icon_url as bank_icon_url, banks.type_flag"
	if err := r.db.Table("Bank_account_transactions as transactions").
		Select(selectedFields).
		Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = transactions.account_id").
		Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id").
		Where("transactions.id = ?", id).
		Where("transactions.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) GetTransactions(req model.BankAccountTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankAccountTransactionResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_account_transactions as transactions")
	count = count.Select("transactions.id")
	count = count.Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = transactions.account_id")
	count = count.Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id")
	if req.AccountId != 0 {
		count = count.Where("transactions.account_id = ?", req.AccountId)
	}
	if req.FromCreatedDate != "" {
		count = count.Where("transactions.created_at >= ?", req.FromCreatedDate)
	}
	if req.ToCreatedDate != "" {
		count = count.Where("transactions.created_at <= ?", req.ToCreatedDate)
	}
	if req.TransferType != "" {
		count = count.Where("transactions.transfer_type = ?", req.TransferType)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where("transactions.description LIKE ?", search_like)
		count = count.Or("accounts.account_name LIKE ?", search_like)
		count = count.Or("accounts.account_number LIKE ?", search_like)
	}
	if err = count.
		Where("transactions.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	if total > 0 {
		// SELECT //
		selectedFields := "transactions.id, transactions.account_id, transactions.description, transactions.transfer_type, transactions.amount, transactions.transfer_at, transactions.created_by_username, transactions.created_at, transactions.updated_at"
		selectedFields += ",accounts.bank_id, accounts.account_type_id, accounts.account_name, accounts.account_number, accounts.account_balance, accounts.account_priority, accounts.account_status, accounts.created_at, accounts.updated_at"
		selectedFields += ",banks.name as bank_name, banks.code, banks.icon_url as bank_icon_url, banks.type_flag"
		query := r.db.Table("Bank_account_transactions as transactions")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = transactions.account_id")
		query = query.Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id")
		if req.AccountId != 0 {
			query = query.Where("transactions.account_id = ?", req.AccountId)
		}
		if req.FromCreatedDate != "" {
			query = query.Where("transactions.created_at >= ?", req.FromCreatedDate)
		}
		if req.ToCreatedDate != "" {
			query = query.Where("transactions.created_at <= ?", req.ToCreatedDate)
		}
		if req.TransferType != "" {
			query = query.Where("transactions.transfer_type = ?", req.TransferType)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where("transactions.description LIKE ?", search_like)
			query = query.Or("accounts.account_name LIKE ?", search_like)
			query = query.Or("accounts.account_number LIKE ?", search_like)
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}

		if err = query.
			Where("transactions.deleted_at IS NULL").
			Limit(req.Limit).
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) CreateTransaction(data model.BankAccountTransactionBody) error {
	if err := r.db.Table("Bank_account_transactions").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) UpdateTransaction(id int64, data model.BankAccountTransactionBody) error {
	if err := r.db.Table("Bank_account_transactions").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DeleteTransaction(id int64) error {
	if err := r.db.Table("Bank_account_transactions").Where("id = ?", id).Delete(&model.BankAccountTransaction{}).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) GetTransferById(id int64) (*model.BankAccountTransfer, error) {
	var record model.BankAccountTransfer
	selectedFields := "transfers.id, transfers.from_account_id, transfers.from_bank_id, transfers.from_account_name, transfers.from_account_number"
	selectedFields += ",transfers.to_account_id, transfers.to_bank_id, transfers.to_account_name, transfers.to_account_number"
	selectedFields += ",transfers.amount, transfers.transfer_at, transfers.created_by_username, transfers.status, transfers.confirmed_at, transfers.confirmed_by_user_id, transfers.created_at, transfers.updated_at"
	selectedFields += ",from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
	selectedFields += ",to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
	if err := r.db.Table("Bank_account_transfers as transfers").
		Select(selectedFields).
		Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = transfers.from_bank_id").
		Joins("LEFT JOIN Banks AS to_banks ON to_banks.id = transfers.to_bank_id").
		Where("transfers.id = ?", id).
		Where("transfers.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) GetTransfers(req model.BankAccountTransferListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankAccountTransferResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_account_transfers as transfers")
	count = count.Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = transfers.from_bank_id")
	count = count.Joins("LEFT JOIN Banks AS to_banks ON to_banks.id = transfers.to_bank_id")
	count = count.Select("transfers.id")
	if req.AccountId != 0 {
		count = count.Where("transfers.from_account_id = ?", req.AccountId)
	}
	if req.FromCreatedDate != "" {
		count = count.Where("transfers.created_at >= ?", req.FromCreatedDate)
	}
	if req.ToCreatedDate != "" {
		count = count.Where("transfers.created_at <= ?", req.ToCreatedDate)
	}
	if req.ToAccountId != 0 {
		count = count.Where("transfers.to_account_id = ?", req.ToAccountId)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where("transfers.description LIKE ?", search_like)
		count = count.Or("transfers.from_account_name LIKE ?", search_like)
		count = count.Or("transfers.from_account_number LIKE ?", search_like)
		count = count.Or("transfers.to_account_name LIKE ?", search_like)
		count = count.Or("transfers.to_account_number LIKE ?", search_like)
	}

	if err = count.
		Where("transfers.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "transfers.id, transfers.from_account_id, transfers.from_bank_id, transfers.from_account_name, transfers.from_account_number"
		selectedFields += ",transfers.to_account_id, transfers.to_bank_id, transfers.to_account_name, transfers.to_account_number"
		selectedFields += ",transfers.amount, transfers.transfer_at, transfers.created_by_username, transfers.status, transfers.confirmed_at, transfers.confirmed_by_user_id, transfers.created_at, transfers.updated_at"
		selectedFields += ",from_banks.name as from_bank_name, to_banks.name as to_bank_name"
		query := r.db.Table("Bank_account_transfers as transfers")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = transfers.from_bank_id")
		query = query.Joins("LEFT JOIN Banks AS to_banks ON to_banks.id = transfers.to_bank_id")
		if req.AccountId != 0 {
			query = query.Where("transfers.from_account_id = ?", req.AccountId)
		}
		if req.FromCreatedDate != "" {
			query = query.Where("transfers.created_at >= ?", req.FromCreatedDate)
		}
		if req.ToCreatedDate != "" {
			query = query.Where("transfers.created_at <= ?", req.ToCreatedDate)
		}
		if req.ToAccountId != 0 {
			query = query.Where("transfers.to_account_id = ?", req.ToAccountId)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where("transfers.description LIKE ?", search_like)
			query = query.Or("transfers.from_account_name LIKE ?", search_like)
			query = query.Or("transfers.from_account_number LIKE ?", search_like)
			query = query.Or("transfers.to_account_name LIKE ?", search_like)
			query = query.Or("transfers.to_account_number LIKE ?", search_like)
		}

		// Sort by ANY //
		req.SortCol = strings.TrimSpace(req.SortCol)
		if req.SortCol != "" {
			if strings.ToLower(strings.TrimSpace(req.SortAsc)) == "desc" {
				req.SortAsc = "DESC"
			} else {
				req.SortAsc = "ASC"
			}
			query = query.Order(req.SortCol + " " + req.SortAsc)
		}

		if err = query.
			Where("transfers.deleted_at IS NULL").
			Limit(req.Limit).
			Offset(req.Page * req.Limit).
			Scan(&list).
			Error; err != nil {
			return nil, err
		}
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.SuccessWithPagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) CreateTransfer(data model.BankAccountTransferBody) error {
	if err := r.db.Table("Bank_account_transfers").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) ConfirmTransfer(id int64, data model.BankAccountTransferConfirmBody) error {
	if err := r.db.Table("Bank_account_transfers").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DeleteTransfer(id int64) error {
	if err := r.db.Table("Bank_account_transfers").Where("id = ?", id).Delete(&model.BankAccountTransfer{}).Error; err != nil {
		return err
	}
	return nil
}
