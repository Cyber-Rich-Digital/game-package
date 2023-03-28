package repository

import (
	"cybergame-api/model"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func NewAccountingRepository(db *gorm.DB) AccountingRepository {
	return &repo{db}
}

type AccountingRepository interface {
	CheckBankAccount(domainName string) (bool, error)
	GetBankAccountByApiKey(apiKey string) (*model.BankAccount, error)
	GetBankAccountByDomain(domainName string) (*model.BankAccount, error)
	GetBankAccount(id int) (*model.BankAccount, error)
	GetBankAccountAndTags(id int) (*model.BankAccount, error)
	GetBankAccountsByUserIds(id []int) (*[]model.BankAccountList, error)
	GetBankAccounts(data model.BankAccountQuery) (*model.Pagination, error)
	GetBankAccountTotals(data model.BankAccountDate) (*[]model.BankAccountListResponse, error)
	CreateBankAccount(data model.BankAccount) error
	UpdateBankAccount(id int64, data model.BankAccountBody) error
	DeleteBankAccount(id int) error
}

func (r repo) GetBankAccountByDomain(domainName string) (*model.BankAccount, error) {

	var result *model.BankAccount

	if err := r.db.Table("BankAccounts").
		Select("id, user_id").
		Where("domain_name = ?", domainName).
		Where("deleted_at IS NULL").
		First(&result).
		Error; err != nil {
		return nil, err
	}

	if result.Id == 0 {
		return nil, errors.New("Account not found")
	}

	return result, nil
}

func (r repo) GetBankAccountByApiKey(apiKey string) (*model.BankAccount, error) {

	var result *model.BankAccount

	if err := r.db.Table("BankAccounts").
		Select("id").
		Where("api_key = ?", apiKey).
		Where("deleted_at IS NULL").
		First(&result).
		Error; err != nil {
		return nil, err
	}

	if result.Id == 0 {
		return nil, errors.New("Account not found")
	}

	return result, nil
}

func (r repo) CheckBankAccount(domainName string) (bool, error) {

	var count int64

	if err := r.db.Table("BankAccounts").
		Select("id").
		Where("domain_name = ?", domainName).
		Where("deleted_at IS NULL").
		Limit(1).
		Count(&count).
		Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r repo) GetBankAccount(id int) (*model.BankAccount, error) {

	var accounting model.BankAccount
	if err := r.db.Table("BankAccounts").
		Select("id, title, domain_name, api_key, user_id, created_at").
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

func (r repo) GetBankAccountAndTags(id int) (*model.BankAccount, error) {

	var accounting model.BankAccount
	if err := r.db.Table("BankAccounts w").
		Select("w.id, w.title, w.domain_name, w.api_key, w.user_id, w.created_at").
		Joins("LEFT JOIN Tags AS t ON t.accounting_id = w.id AND t.deleted_at IS NULL").
		Where("w.id = ?", id).
		Where("w.deleted_at IS NULL").
		First(&accounting).
		Error; err != nil {
		return nil, err
	}

	if accounting.Id == 0 {
		return nil, errors.New("Account not found")
	}

	if err := r.db.Table("Tags").
		Select("id, name, accounting_id").
		Where("accounting_id = ?", id).
		Where("deleted_at IS NULL").
		Find(&accounting.Tags).
		Error; err != nil {
		return nil, err
	}

	return &accounting, nil
}

func (r repo) GetBankAccountsByUserIds(id []int) (*[]model.BankAccountList, error) {

	var accounting *[]model.BankAccountList
	if err := r.db.Table("BankAccounts").
		Select("id, title, domain_name, api_key, user_id, created_at").
		Where("id IN (?)", id).
		Where("deleted_at IS NULL").
		Find(&accounting).
		Error; err != nil {
		return nil, err
	}

	return accounting, nil
}

func (r repo) GetBankAccounts(data model.BankAccountQuery) (*model.Pagination, error) {

	var list []model.BankAccountResponse
	var total int64
	var err error

	selectFields := "w.id, w.title, w.domain_name, w.api_key, w.created_at, w.updated_at, COUNT(m.id) AS total"
	join := "LEFT OUTER JOIN Messages AS m ON w.id = m.accounting_id"
	group := "m.accounting_id, w.id, w.title, w.domain_name, w.api_key, w.created_at, w.updated_at"
	whereVal := fmt.Sprintf("%%%s%%", data.Search)

	// Get list of accountings //

	query := r.db.Table("BankAccounts w")
	query = query.Select(selectFields).
		Joins(join).
		Group(group)

	if data.Role == "USER" {
		query = query.Where("w.user_id = ?", data.UserId)
	}

	if data.Search != "" {
		query = query.Where("w.title LIKE ?", whereVal).
			Or("w.domain_name LIKE ?", whereVal)
	}

	// Sort by created_at //

	if data.Sort == 1 {
		query = query.Order("w.created_at ASC")
	} else {
		query = query.Order("w.created_at DESC")
	}

	if err = query.
		Where("w.deleted_at IS NULL").
		Limit(data.Limit).
		Offset(data.Page * data.Limit).
		Scan(&list).
		Error; err != nil {
		return nil, err
	}

	// Count total records for pagination purposes (without limit and offset)  //

	count := r.db.Table("BankAccounts w")
	count = count.Select(selectFields).
		Joins(join).
		Group(group)

	if data.Role == "USER" {
		count = count.Where("w.user_id = ?", data.UserId)
	}

	if data.Search != "" {
		count = count.Where("w.title LIKE ?", whereVal).
			Or("w.domain_name LIKE ?", whereVal)
	}

	if err = count.
		Where("w.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	// End count total records for pagination purposes (without limit and offset)  //

	var accounting model.Pagination
	accounting.List = list
	accounting.Total = total

	return &accounting, nil
}

func (r repo) GetBankAccountTotals(data model.BankAccountDate) (*[]model.BankAccountListResponse, error) {

	t := time.Unix(data.Date, 0)

	var list *[]model.BankAccountListResponse

	selectFields := "w.id, COUNT(m.id) AS total"
	join := "LEFT OUTER JOIN Messages AS m ON  w.id = m.accounting_id AND m.created_at >= ?"
	group := "m.accounting_id, w.id"

	if err := r.db.Table("BankAccounts w").
		Select(selectFields).
		Joins(join, t.Format("2006-01-02 15:04:05")).
		Group(group).
		Where("w.user_id = ?", data.UserId).
		Where("w.deleted_at IS NULL").
		Find(&list).
		Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r repo) CreateBankAccount(data model.BankAccount) error {
	if err := r.db.Table("BankAccounts").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) UpdateBankAccount(id int64, data model.BankAccountBody) error {
	if err := r.db.Table("BankAccounts").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DeleteBankAccount(id int) error {
	if err := r.db.Table("BankAccounts").Where("id = ?", id).Delete(&model.BankAccount{}).Error; err != nil {
		return err
	}
	return nil
}
