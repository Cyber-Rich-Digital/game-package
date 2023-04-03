package repository

import (
	"cybergame-api/model"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func NewBankingRepository(db *gorm.DB) BankingRepository {
	return &repo{db}
}

type BankingRepository interface {
	GetStatementById(id int64) (*model.BankStatement, error)
	GetStatements(req model.BankStatementListRequest) (*model.SuccessWithPagination, error)
	CreateStatement(data model.BankStatementCreateBody) error
	UpdateStatement(id int64, data model.BankStatementUpdateBody) error
	DeleteStatement(id int64) error
}

func (r repo) GetStatementById(id int64) (*model.BankStatement, error) {
	var record model.BankStatement
	selectedFields := "statements.id, statements.account_id, statements.transfer_at, statements.amount, statements.status, statements.created_at, statements.updated_at"
	selectedFields += ",accounts.account_name, accounts.account_number, accounts.account_type_id, accounts.bank_id, banks.name as bank_name, banks.code as bank_code, banks.icon_url as bank_icon_url, banks.type_flag as bank_type_flag"
	if err := r.db.Table("Bank_statements as statements").
		Select(selectedFields).
		Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = statements.account_id").
		Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id").
		Where("statements.id = ?", id).
		Where("statements.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) GetStatements(req model.BankStatementListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankStatementResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_statements as statements")
	count = count.Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = statements.account_id")
	count = count.Select("statements.id")
	if req.AccountId != 0 {
		count = count.Where("statements.account_id = ?", req.AccountId)
	}
	if req.Amount >= 0 {
		count = count.Where("statements.amount = ?", req.Amount)
	}
	if req.FromTransferDate != "" {
		count = count.Where("statements.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("statements.transfer_at <= ?", req.ToTransferDate)
	}

	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where("accounts.account_name LIKE ?", search_like)
		count = count.Or("accounts.account_number LIKE ?", search_like)
	}

	if err = count.
		Where("statements.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "statements.id, statements.account_id, statements.transfer_at, statements.amount, statements.status, statements.created_at, statements.updated_at"
		selectedFields += ",accounts.account_name, accounts.account_number, accounts.account_type_id, accounts.bank_id, banks.name as bank_name, banks.code as bank_code, banks.icon_url as bank_icon_url, banks.type_flag as bank_type_flag"
		query := r.db.Table("Bank_statements as statements")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = statements.account_id")
		query = query.Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id")
		if req.AccountId != 0 {
			query = query.Where("statements.account_id = ?", req.AccountId)
		}
		if req.Amount >= 0 {
			query = query.Where("statements.amount = ?", req.Amount)
		}
		if req.FromTransferDate != "" {
			query = query.Where("statements.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("statements.transfer_at <= ?", req.ToTransferDate)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where("accounts.account_name LIKE ?", search_like)
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
			Where("statements.deleted_at IS NULL").
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

func (r repo) CreateStatement(data model.BankStatementCreateBody) error {
	if err := r.db.Table("Bank_statements").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) UpdateStatement(id int64, data model.BankStatementUpdateBody) error {
	if err := r.db.Table("Bank_statements").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DeleteStatement(id int64) error {
	if err := r.db.Table("Bank_statements").Where("id = ?", id).Delete(&model.BankStatement{}).Error; err != nil {
		return err
	}
	return nil
}
