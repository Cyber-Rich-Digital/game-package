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
	GetBankStatementById(id int64) (*model.BankStatement, error)
	GetBankStatements(req model.BankStatementListRequest) (*model.SuccessWithPagination, error)
	CreateBankStatement(data model.BankStatementCreateBody) error
	UpdateBankStatement(id int64, data model.BankStatementUpdateBody) error
	DeleteBankStatement(id int64) error

	GetBankTransactionById(id int64) (*model.BankTransaction, error)
	GetBankTransactions(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error)
	CreateBankTransaction(data model.BankTransactionCreateBody) error
	CreateBonusTransaction(data model.BonusTransactionCreateBody) error
	UpdateBankTransaction(id int64, data model.BankTransactionUpdateBody) error
	DeleteBankTransaction(id int64) error

	GetPendingDepositTransactions(req model.PendingDepositTransactionListRequest) (*model.SuccessWithPagination, error)
	GetPendingWithdrawTransactions(req model.PendingWithdrawTransactionListRequest) (*model.SuccessWithPagination, error)
	GetFinishedTransactions(req model.FinishedTransactionListRequest) (*model.SuccessWithPagination, error)
	RemoveFinishedTransaction(id int64, data model.BankTransactionRemoveBody) error
}

func (r repo) GetBankStatementById(id int64) (*model.BankStatement, error) {
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

func (r repo) GetBankStatements(req model.BankStatementListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankStatementResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_statements as statements")
	count = count.Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = statements.account_id")
	count = count.Select("statements.id")
	if req.AccountId != "" {
		count = count.Where("statements.account_id = ?", req.AccountId)
	}
	if req.Amount != "" {
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
		if req.AccountId != "" {
			query = query.Where("statements.account_id = ?", req.AccountId)
		}
		if req.Amount != "" {
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

func (r repo) CreateBankStatement(data model.BankStatementCreateBody) error {
	if err := r.db.Table("Bank_statements").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) UpdateBankStatement(id int64, data model.BankStatementUpdateBody) error {
	if err := r.db.Table("Bank_statements").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DeleteBankStatement(id int64) error {
	if err := r.db.Table("Bank_statements").Where("id = ?", id).Delete(&model.BankStatement{}).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) GetBankTransactionById(id int64) (*model.BankTransaction, error) {
	var record model.BankTransaction
	selectedFields := "transactions.id, transactions.member_code, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
	selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.before_amount, transactions.after_amount, transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.is_auto_credit"
	selectedFields += ", transactions.created_at, transactions.updated_at"
	if err := r.db.Table("Bank_transactions as transactions").
		Select(selectedFields).
		Where("transactions.id = ?", id).
		Where("transactions.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) GetBankTransactions(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankTransactionResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_transactions as transactions")
	count = count.Select("transactions.id")
	if req.MemberCode != "" {
		count = count.Where("transactions.member_code = ?", req.MemberCode)
	}
	if req.UserId != "" {
		count = count.Where("transactions.user_id = ?", req.UserId)
	}
	if req.FromTransferDate != "" {
		count = count.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where("transactions.from_account_name LIKE ?", search_like)
		count = count.Or("transactions.from_account_number LIKE ?", search_like)
		count = count.Or("transactions.to_account_name LIKE ?", search_like)
		count = count.Or("transactions.to_account_number LIKE ?", search_like)
	}

	if err = count.
		Where("transactions.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "transactions.id, transactions.member_code, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.before_amount, transactions.after_amount, transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)
		if req.MemberCode != "" {
			query = query.Where("transactions.member_code = ?", req.MemberCode)
		}
		if req.UserId != "" {
			query = query.Where("transactions.user_id = ?", req.UserId)
		}
		if req.FromTransferDate != "" {
			query = query.Where("transactions.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("transactions.transfer_at <= ?", req.ToTransferDate)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where("transactions.from_account_name LIKE ?", search_like)
			query = query.Or("transactions.from_account_number LIKE ?", search_like)
			query = query.Or("transactions.to_account_name LIKE ?", search_like)
			query = query.Or("transactions.to_account_number LIKE ?", search_like)
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

func (r repo) CreateBankTransaction(data model.BankTransactionCreateBody) error {
	if err := r.db.Table("Bank_transactions").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) CreateBonusTransaction(data model.BonusTransactionCreateBody) error {
	if err := r.db.Table("Bank_transactions").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) UpdateBankTransaction(id int64, data model.BankTransactionUpdateBody) error {
	if err := r.db.Table("Bank_transactions").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DeleteBankTransaction(id int64) error {
	if err := r.db.Table("Bank_transactions").Where("id = ?", id).Delete(&model.BankTransaction{}).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) GetPendingDepositTransactions(req model.PendingDepositTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankTransactionResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_transactions as transactions")
	count = count.Select("transactions.id")

	if req.FromTransferDate != "" {
		count = count.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where("transactions.from_account_name LIKE ?", search_like)
		count = count.Or("transactions.from_account_number LIKE ?", search_like)
		count = count.Or("transactions.to_account_name LIKE ?", search_like)
		count = count.Or("transactions.to_account_number LIKE ?", search_like)
	}

	if err = count.
		Where("transactions.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "transactions.id, transactions.member_code, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.before_amount, transactions.after_amount, transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)

		if req.FromTransferDate != "" {
			query = query.Where("transactions.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("transactions.transfer_at <= ?", req.ToTransferDate)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where("transactions.from_account_name LIKE ?", search_like)
			query = query.Or("transactions.from_account_number LIKE ?", search_like)
			query = query.Or("transactions.to_account_name LIKE ?", search_like)
			query = query.Or("transactions.to_account_number LIKE ?", search_like)
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

func (r repo) GetPendingWithdrawTransactions(req model.PendingWithdrawTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankTransactionResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_transactions as transactions")
	count = count.Select("transactions.id")

	if req.FromTransferDate != "" {
		count = count.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where("transactions.from_account_name LIKE ?", search_like)
		count = count.Or("transactions.from_account_number LIKE ?", search_like)
		count = count.Or("transactions.to_account_name LIKE ?", search_like)
		count = count.Or("transactions.to_account_number LIKE ?", search_like)
	}

	if err = count.
		Where("transactions.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "transactions.id, transactions.member_code, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.before_amount, transactions.after_amount, transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)

		if req.FromTransferDate != "" {
			query = query.Where("transactions.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("transactions.transfer_at <= ?", req.ToTransferDate)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where("transactions.from_account_name LIKE ?", search_like)
			query = query.Or("transactions.from_account_number LIKE ?", search_like)
			query = query.Or("transactions.to_account_name LIKE ?", search_like)
			query = query.Or("transactions.to_account_number LIKE ?", search_like)
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

func (r repo) GetFinishedTransactions(req model.FinishedTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankTransactionResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_transactions as transactions")
	count = count.Select("transactions.id")

	if req.FromTransferDate != "" {
		count = count.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where("transactions.from_account_name LIKE ?", search_like)
		count = count.Or("transactions.from_account_number LIKE ?", search_like)
		count = count.Or("transactions.to_account_name LIKE ?", search_like)
		count = count.Or("transactions.to_account_number LIKE ?", search_like)
	}

	if err = count.
		Where("transactions.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "transactions.id, transactions.member_code, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.before_amount, transactions.after_amount, transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)

		if req.FromTransferDate != "" {
			query = query.Where("transactions.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("transactions.transfer_at <= ?", req.ToTransferDate)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where("transactions.from_account_name LIKE ?", search_like)
			query = query.Or("transactions.from_account_number LIKE ?", search_like)
			query = query.Or("transactions.to_account_name LIKE ?", search_like)
			query = query.Or("transactions.to_account_number LIKE ?", search_like)
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

func (r repo) RemoveFinishedTransaction(id int64, data model.BankTransactionRemoveBody) error {
	if err := r.db.Table("Bank_transactions").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}
