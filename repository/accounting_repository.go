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
	GetBanks(req model.BankListRequest) (*model.Pagination, error)
	GetBankByCode(code string) (*model.BankAccount, error)

	HasBankAccount(accountNumber string) (bool, error)
	GetBankAccountById(id int64) (*model.BankAccount, error)
	GetBankAccounts(data model.BankAccountListRequest) (*model.Pagination, error)
	CreateBankAccount(data model.BankAccount) error
	UpdateBankAccount(id int64, data model.BankAccountBody) error
	DeleteBankAccount(id int64) error
}

func (r repo) GetBanks(req model.BankListRequest) (*model.Pagination, error) {

	var list []model.BankResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("banks")
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
		query := r.db.Table("banks")
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
	var result model.Pagination
	if list == nil {
		list = []model.BankResponse{}
	}
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) GetBankByCode(code string) (*model.BankAccount, error) {

	var result *model.BankAccount

	if err := r.db.Table("banks").
		Select("id, name, code, icon_url, icon_url, type_flag").
		Where("code = ?", code).
		First(&result).
		Error; err != nil {
		return nil, err
	}

	if result.Id == 0 {
		return nil, errors.New("Account not found")
	}

	return result, nil
}

func (r repo) HasBankAccount(accountNumber string) (bool, error) {
	var count int64
	if err := r.db.Table("bank_accounts").
		Select("id").
		Where("account_number = ?", accountNumber).
		Limit(1).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r repo) GetBankAccountById(id int64) (*model.BankAccount, error) {

	var accounting model.BankAccount
	if err := r.db.Table("bank_accounts").
		Select("id, bank_id, account_type_id, account_name, account_number, transfer_priority, account_status, created_at, updated_at").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&accounting).
		Error; err != nil {
		return nil, err
	}

	if accounting.Id == 0 {
		return nil, errors.New("Account not found")
	}

	return &accounting, nil
}

func (r repo) GetBankAccounts(req model.BankAccountListRequest) (*model.Pagination, error) {

	var list []model.BankAccountResponse
	var total int64
	var err error

	// SELECT //
	query := r.db.Table("bank_accounts")
	query = query.Select("id, bank_id, account_type_id, account_name, account_number, transfer_priority, account_status, created_at, updated_at")

	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		query = query.Where("account_name LIKE ?", search_like).Or("account_number LIKE ?", search_like)
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

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("bank_accounts")
	count = count.Select("id")
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where("account_name LIKE ?", search_like).Or("account_number LIKE ?", search_like)
	}
	if err = count.
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	// End count total records for pagination purposes (without limit and offset) //
	var result model.Pagination
	result.List = list
	result.Total = total
	return &result, nil
}

func (r repo) CreateBankAccount(data model.BankAccount) error {
	if err := r.db.Table("bank_accounts").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) UpdateBankAccount(id int64, data model.BankAccountBody) error {
	if err := r.db.Table("bank_accounts").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DeleteBankAccount(id int64) error {
	if err := r.db.Table("bank_accounts").Where("id = ?", id).Delete(&model.BankAccount{}).Error; err != nil {
		return err
	}
	return nil
}
