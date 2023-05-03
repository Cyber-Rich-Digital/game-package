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
	GetBankStatementSummary(req model.BankStatementListRequest) (*model.BankStatementSummary, error)
	CreateBankStatement(data model.BankStatementCreateBody) error
	UpdateBankStatement(id int64, data model.BankStatementUpdateBody) error
	MatchStatementOwner(id int64, data model.BankStatementUpdateBody) error
	IgnoreStatementOwner(id int64, data model.BankStatementUpdateBody) error
	DeleteBankStatement(id int64) error

	GetBankTransactionById(id int64) (*model.BankTransaction, error)
	GetBankTransactions(req model.BankTransactionListRequest) (*model.SuccessWithPagination, error)
	CreateBankTransaction(data model.BankTransactionCreateBody) (*int64, error)
	CreateBonusTransaction(data model.BonusTransactionCreateBody) error
	UpdateBankTransaction(id int64, data model.BankTransactionUpdateBody) error
	DeleteBankTransaction(id int64) error

	GetPendingDepositTransactions(req model.PendingDepositTransactionListRequest) (*model.SuccessWithPagination, error)
	GetPendingWithdrawTransactions(req model.PendingWithdrawTransactionListRequest) (*model.SuccessWithPagination, error)
	CreateTransactionAction(data model.CreateBankTransactionActionBody) error
	CreateStatementAction(data model.CreateBankStatementActionBody) error
	ConfirmPendingDepositTransaction(id int64, data model.BankDepositTransactionConfirmBody) error
	ConfirmPendingWithdrawTransaction(id int64, data model.BankWithdrawTransactionConfirmBody) error
	CancelPendingTransaction(id int64, data model.BankTransactionCancelBody) error
	GetFinishedTransactions(req model.FinishedTransactionListRequest) (*model.SuccessWithPagination, error)
	RemoveFinishedTransaction(id int64, data model.BankTransactionRemoveBody) error
	GetRemovedTransactions(req model.RemovedTransactionListRequest) (*model.SuccessWithPagination, error)

	GetMemberById(id int64) (*model.Member, error)
	GetMemberByCode(code string) (*model.Member, error)
	GetMembers(req model.MemberListRequest) (*model.SuccessWithPagination, error)
	GetPossibleStatementOwners(req model.MemberPossibleListRequest) (*model.SuccessWithPagination, error)
	GetMemberTransactions(req model.MemberTransactionListRequest) (*model.SuccessWithPagination, error)
	GetMemberTransactionSummary(req model.MemberTransactionListRequest) (*model.MemberTransactionSummary, error)
	IncreaseMemberCredit(id int64, amount float32) error
	DecreaseMemberCredit(id int64, amount float32) error
}

func (r repo) GetBankStatementById(id int64) (*model.BankStatement, error) {
	var record model.BankStatement
	selectedFields := "statements.id, statements.account_id, statements.detail, statements.statement_type, statements.transfer_at, statements.from_bank_id, statements.from_account_number, statements.amount, statements.status, statements.created_at, statements.updated_at"
	selectedFields += ",accounts.account_name, accounts.account_number, accounts.account_type_id, accounts.bank_id"
	selectedFields += ",banks.name as bank_name, banks.code as bank_code, banks.icon_url as bank_icon_url, banks.type_flag as bank_type_flag"
	selectedFields += ",from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url"
	if err := r.db.Table("Bank_statements as statements").
		Select(selectedFields).
		Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = statements.account_id").
		Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id").
		Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = statements.from_bank_id").
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
	count = count.Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id")
	count = count.Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = statements.from_bank_id")
	count = count.Select("statements.id")
	if req.AccountId != "" {
		count = count.Where("statements.account_id = ?", req.AccountId)
	}
	if req.FromTransferDate != "" {
		count = count.Where("statements.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("statements.transfer_at <= ?", req.ToTransferDate)
	}
	if req.StatementType != "" {
		count = count.Where("statements.statement_type = ?", req.StatementType)
	}
	if req.Status != "" {
		count = count.Where("statements.status = ?", req.Status)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("accounts.account_name LIKE ?", search_like).Or("accounts.account_number LIKE ?", search_like))
	}

	if err = count.
		Where("statements.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "statements.id, statements.account_id, statements.detail, statements.statement_type, statements.transfer_at, statements.from_bank_id, statements.from_account_number, statements.amount, statements.status, statements.created_at, statements.updated_at"
		selectedFields += ",accounts.account_name, accounts.account_number, accounts.account_type_id, accounts.bank_id"
		selectedFields += ",banks.name as bank_name, banks.code as bank_code, banks.icon_url as bank_icon_url, banks.type_flag as bank_type_flag"
		selectedFields += ",from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url"
		query := r.db.Table("Bank_statements as statements")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = statements.account_id")
		query = query.Joins("LEFT JOIN Banks AS banks ON banks.id = accounts.bank_id")
		query = query.Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = statements.from_bank_id")
		if req.AccountId != "" {
			query = query.Where("statements.account_id = ?", req.AccountId)
		}
		if req.FromTransferDate != "" {
			query = query.Where("statements.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("statements.transfer_at <= ?", req.ToTransferDate)
		}
		if req.StatementType != "" {
			query = query.Where("statements.statement_type = ?", req.StatementType)
		}
		if req.Status != "" {
			query = query.Where("statements.status = ?", req.Status)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where(r.db.Where("accounts.account_name LIKE ?", search_like).Or("accounts.account_number LIKE ?", search_like))
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

func (r repo) GetBankStatementSummary(req model.BankStatementListRequest) (*model.BankStatementSummary, error) {

	var result model.BankStatementSummary
	var totalPendingStatementCount int64
	var totalPendingDepositCount int64
	var totalPendingWithdrawCount int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_statements as statements")
	count = count.Joins("LEFT JOIN Bank_accounts AS accounts ON accounts.id = statements.account_id")
	count = count.Select("statements.id")
	count = count.Where("statements.status = ?", "pending")
	if req.AccountId != "" {
		count = count.Where("statements.account_id = ?", req.AccountId)
	}
	if req.StatementType != "" {
		count = count.Where("statements.statement_type = ?", req.StatementType)
	}
	if req.FromTransferDate != "" {
		count = count.Where("statements.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("statements.transfer_at <= ?", req.ToTransferDate)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("accounts.account_name LIKE ?", search_like).Or("accounts.account_number LIKE ?", search_like))
	}
	if err = count.
		Where("statements.deleted_at IS NULL").
		Count(&totalPendingStatementCount).
		Error; err != nil {
		return nil, err
	}

	// Count total records for pagination purposes (without limit and offset) //
	countDeposit := r.db.Table("Bank_transactions as transactions")
	countDeposit = countDeposit.Select("transactions.id")
	countDeposit = countDeposit.Where("transactions.status = ?", "pending")
	countDeposit = countDeposit.Where("transactions.transfer_type = ?", "deposit")
	if req.AccountId != "" {
		countDeposit = countDeposit.Where("transactions.to_account_id = ?", req.AccountId)
	}
	if req.FromTransferDate != "" {
		countDeposit = countDeposit.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		countDeposit = countDeposit.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if err = countDeposit.
		Where("transactions.deleted_at IS NULL").
		Count(&totalPendingDepositCount).
		Error; err != nil {
		return nil, err
	}

	// Count total records for pagination purposes (without limit and offset) //
	countWithdraw := r.db.Table("Bank_transactions as transactions")
	countWithdraw = countWithdraw.Select("transactions.id")
	countWithdraw = countWithdraw.Where("transactions.status = ?", "pending")
	countWithdraw = countWithdraw.Where("transactions.transfer_type = ?", "withdraw")
	if req.AccountId != "" {
		countWithdraw = countWithdraw.Where("transactions.from_account_id = ?", req.AccountId)
	}
	if req.FromTransferDate != "" {
		countWithdraw = countWithdraw.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		countWithdraw = countWithdraw.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if err = countWithdraw.
		Where("transactions.deleted_at IS NULL").
		Count(&totalPendingWithdrawCount).
		Error; err != nil {
		return nil, err
	}

	result.TotalPendingStatementCount = totalPendingStatementCount
	result.TotalPendingDepositCount = totalPendingDepositCount
	result.TotalPendingWithdrawCount = totalPendingWithdrawCount

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

func (r repo) MatchStatementOwner(id int64, data model.BankStatementUpdateBody) error {
	if err := r.db.Table("Bank_statements").Where("id = ?", id).Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) IgnoreStatementOwner(id int64, data model.BankStatementUpdateBody) error {
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
	selectedFields := "transactions.id, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
	selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.before_amount, transactions.after_amount, transactions.bank_charge_amount, transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.status_detail, transactions.is_auto_credit"
	selectedFields += ", transactions.created_at, transactions.updated_at"
	selectedFields += ", from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
	selectedFields += ", to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
	selectedFields += ", users.member_code as member_code, users.username as user_username, users.firstname as user_firstname, users.lastname as user_lastname, users.fullname as user_fullname, users.phone as user_phone"
	if err := r.db.Table("Bank_transactions as transactions").
		Select(selectedFields).
		Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = transactions.from_bank_id").
		Joins("LEFT JOIN Banks AS to_banks ON to_banks.id = transactions.to_bank_id").
		Joins("LEFT JOIN Users AS users ON users.id = transactions.user_id").
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
	count = count.Where("transactions.removed_at IS NULL")
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
	if req.TransferType != "" {
		count = count.Where("transactions.transfer_type = ?", req.TransferType)
	}
	if req.TransferStatus != "" {
		count = count.Where("transactions.status = ?", req.TransferStatus)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("transactions.from_account_name LIKE ?", search_like).Or("transactions.from_account_number LIKE ?", search_like).Or("transactions.to_account_name LIKE ?", search_like).Or("transactions.to_account_number LIKE ?", search_like))
	}

	if err = count.
		Where("transactions.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "transactions.id, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.before_amount, transactions.after_amount, transactions.bank_charge_amount, transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.status_detail, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		selectedFields += ", from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
		selectedFields += ", to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
		selectedFields += ", users.member_code as member_code, users.username as user_username, users.firstname as user_firstname, users.lastname as user_lastname, users.fullname as user_fullname, users.phone as user_phone"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks AS from_banks ON from_banks.id = transactions.from_bank_id")
		query = query.Joins("LEFT JOIN Banks AS to_banks ON to_banks.id = transactions.to_bank_id")
		query = query.Joins("LEFT JOIN Users AS users ON users.id = transactions.user_id")
		query = query.Where("transactions.removed_at IS NULL")
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
		if req.TransferType != "" {
			query = query.Where("transactions.transfer_type = ?", req.TransferType)
		}
		if req.TransferStatus != "" {
			query = query.Where("transactions.status = ?", req.TransferStatus)
		}

		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where(r.db.Where("transactions.from_account_name LIKE ?", search_like).Or("transactions.from_account_number LIKE ?", search_like).Or("transactions.to_account_name LIKE ?", search_like).Or("transactions.to_account_number LIKE ?", search_like))
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

func (r repo) CreateBankTransaction(data model.BankTransactionCreateBody) (*int64, error) {
	if err := r.db.Table("Bank_transactions").Create(&data).Error; err != nil {
		return nil, err
	}
	return &data.Id, nil
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
	count = count.Where("transactions.transfer_type = ?", "deposit")
	count = count.Where("transactions.status = ?", "pending")
	count = count.Where("transactions.removed_at IS NULL")
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
		selectedFields := "transactions.id, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.before_amount, transactions.after_amount, transactions.bank_charge_amount, transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.status_detail, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		selectedFields += ", from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
		selectedFields += ", to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
		selectedFields += ", users.member_code as member_code, users.username as user_username, users.firstname as user_firstname, users.lastname as user_lastname, users.fullname as user_fullname, users.phone as user_phone"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks as from_banks ON from_banks.id = transactions.from_bank_id")
		query = query.Joins("LEFT JOIN Banks as to_banks ON to_banks.id = transactions.to_bank_id")
		query = query.Joins("LEFT JOIN Users as users ON users.id = transactions.user_id")
		query = query.Where("transactions.transfer_type = ?", "deposit")
		query = query.Where("transactions.status = ?", "pending")
		query = query.Where("transactions.removed_at IS NULL")
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
	count = count.Where("transactions.transfer_type = ?", "withdraw")
	count = count.Where("transactions.status = ?", "pending")
	count = count.Where("transactions.removed_at IS NULL")
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
		selectedFields := "transactions.id, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.before_amount, transactions.after_amount, transactions.bank_charge_amount, transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.status_detail, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		selectedFields += ", from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
		selectedFields += ", to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
		selectedFields += ", users.member_code as member_code, users.username as user_username, users.firstname as user_firstname, users.lastname as user_lastname, users.fullname as user_fullname, users.phone as user_phone"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks as from_banks ON from_banks.id = transactions.from_bank_id")
		query = query.Joins("LEFT JOIN Banks as to_banks ON to_banks.id = transactions.to_bank_id")
		query = query.Joins("LEFT JOIN Users as users ON users.id = transactions.user_id")
		query = query.Where("transactions.transfer_type = ?", "withdraw")
		query = query.Where("transactions.status = ?", "pending")
		query = query.Where("transactions.removed_at IS NULL")
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

func (r repo) CancelPendingTransaction(id int64, data model.BankTransactionCancelBody) error {
	if err := r.db.Table("Bank_transactions").Where("id = ?", id).Where("status = ?", "pending").Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) CreateTransactionAction(data model.CreateBankTransactionActionBody) error {
	if err := r.db.Table("Bank_confirm_transactions").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) CreateStatementAction(data model.CreateBankStatementActionBody) error {
	if err := r.db.Table("Bank_confirm_statements").Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) ConfirmPendingDepositTransaction(id int64, body model.BankDepositTransactionConfirmBody) error {
	data := map[string]interface{}{
		"transfer_at":           body.TransferAt,
		"bonus_amount":          body.BonusAmount,
		"status":                body.Status,
		"confirmed_at":          body.ConfirmedAt,
		"confirmed_by_user_id":  body.ConfirmedByUserId,
		"confirmed_by_username": body.ConfirmedByUsername,
	}
	if err := r.db.Table("Bank_transactions").Where("id = ?", id).Where("status = ?", "pending").Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) ConfirmPendingWithdrawTransaction(id int64, body model.BankWithdrawTransactionConfirmBody) error {
	data := map[string]interface{}{
		"transfer_at":           body.TransferAt,
		"credit_amount":         body.CreditAmount,
		"bank_charge_amount":    body.BankChargeAmount,
		"status":                body.Status,
		"confirmed_at":          body.ConfirmedAt,
		"confirmed_by_user_id":  body.ConfirmedByUserId,
		"confirmed_by_username": body.ConfirmedByUsername,
	}
	if err := r.db.Table("Bank_transactions").Where("id = ?", id).Where("status = ?", "pending").Updates(&data).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) GetFinishedTransactions(req model.FinishedTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankTransactionResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_transactions as transactions")
	count = count.Select("transactions.id")
	count = count.Where("transactions.status = ?", "finished")
	count = count.Where("transactions.removed_at IS NULL")
	if req.AccountId != "" {
		count = count.Where(r.db.Where("transactions.from_account_id = ?", req.AccountId).Or("transactions.to_account_id = ?", req.AccountId))
	}
	if req.FromTransferDate != "" {
		count = count.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if req.TransferType != "" {
		count = count.Where("transactions.transfer_type = ?", req.TransferType)
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
		selectedFields := "transactions.id, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.before_amount, transactions.after_amount, transactions.bank_charge_amount, transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.status_detail, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		selectedFields += ", from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
		selectedFields += ", to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
		selectedFields += ", users.member_code as member_code, users.username as user_username, users.firstname as user_firstname, users.lastname as user_lastname, users.fullname as user_fullname, users.phone as user_phone"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks as from_banks ON from_banks.id = transactions.from_bank_id")
		query = query.Joins("LEFT JOIN Banks as to_banks ON to_banks.id = transactions.to_bank_id")
		query = query.Joins("LEFT JOIN Users as users ON users.id = transactions.user_id")
		query = query.Where("transactions.status = ?", "finished")
		query = query.Where("transactions.removed_at IS NULL")
		if req.AccountId != "" {
			query = query.Where(r.db.Where("transactions.from_account_id = ?", req.AccountId).Or("transactions.to_account_id = ?", req.AccountId))
		}
		if req.FromTransferDate != "" {
			query = query.Where("transactions.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("transactions.transfer_at <= ?", req.ToTransferDate)
		}
		if req.TransferType != "" {
			query = query.Where("transactions.transfer_type = ?", req.TransferType)
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

func (r repo) GetRemovedTransactions(req model.RemovedTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.BankTransactionResponse
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_transactions as transactions")
	count = count.Select("transactions.id")
	count = count.Where("transactions.removed_at IS NOT NULL")
	if req.FromTransferDate != "" {
		count = count.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if req.TransferType != "" {
		count = count.Where("transactions.transfer_type = ?", req.TransferType)
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
		selectedFields := "transactions.id, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.before_amount, transactions.after_amount, transactions.bank_charge_amount, transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.status_detail, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		selectedFields += ", from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
		selectedFields += ", to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
		selectedFields += ", users.member_code as member_code, users.username as user_username, users.firstname as user_firstname, users.lastname as user_lastname, users.fullname as user_fullname, users.phone as user_phone"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks as from_banks ON from_banks.id = transactions.from_bank_id")
		query = query.Joins("LEFT JOIN Banks as to_banks ON to_banks.id = transactions.to_bank_id")
		query = query.Joins("LEFT JOIN Users as users ON users.id = transactions.user_id")
		query = query.Where("transactions.removed_at IS NOT NULL")
		if req.FromTransferDate != "" {
			query = query.Where("transactions.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("transactions.transfer_at <= ?", req.ToTransferDate)
		}
		if req.TransferType != "" {
			query = query.Where("transactions.transfer_type = ?", req.TransferType)
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

func (r repo) GetMemberById(id int64) (*model.Member, error) {
	var record model.Member

	selectedFields := "users.id, users.member_code, users.username, users.phone, users.firstname, users.lastname, users.fullname, users.credit, users.bankname, users.bank_account, users.promotion, users.status, users.channel, users.true_wallet, users.note, users.turnover_limit, users.created_at"
	if err := r.db.Table("Users as users").
		Select(selectedFields).
		Where("users.id = ?", id).
		Where("users.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) GetMemberByCode(memberCode string) (*model.Member, error) {
	var record model.Member

	selectedFields := "users.id, users.member_code, users.username, users.phone, users.firstname, users.lastname, users.fullname, users.credit, users.bankname, users.bank_account, users.promotion, users.status, users.channel, users.true_wallet, users.note, users.turnover_limit, users.created_at"
	if err := r.db.Table("Users as users").
		Select(selectedFields).
		Where("users.member_code = ?", memberCode).
		Where("users.deleted_at IS NULL").
		First(&record).
		Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r repo) GetMembers(req model.MemberListRequest) (*model.SuccessWithPagination, error) {

	var list []model.Member
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Users as users")
	count = count.Select("users.id")
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("users.username LIKE ?", search_like).Or("users.phone LIKE ?", search_like).Or("users.fullname LIKE ?", search_like).Or("users.bankname LIKE ?", search_like).Or("users.bank_account LIKE ?", search_like))
	}

	if err = count.
		Where("users.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "users.id, users.member_code, users.username, users.phone, users.firstname, users.lastname, users.fullname, users.credit, users.bankname, users.bank_account, users.promotion, users.status, users.channel, users.true_wallet, users.note, users.turnover_limit, users.created_at"
		query := r.db.Table("Users as users")
		query = query.Select(selectedFields)
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where(r.db.Where("users.username LIKE ?", search_like).Or("users.phone LIKE ?", search_like).Or("users.fullname LIKE ?", search_like).Or("users.bankname LIKE ?", search_like).Or("users.bank_account LIKE ?", search_like))
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
			Where("users.deleted_at IS NULL").
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

func (r repo) GetPossibleStatementOwners(req model.MemberPossibleListRequest) (*model.SuccessWithPagination, error) {

	var list []model.Member
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Users as users")
	count = count.Select("users.id")
	if req.UserBankCode != nil {
		count = count.Where("users.bank_code = ?", *req.UserBankCode)
	}
	if req.UserAccountNumber != nil && *req.UserAccountNumber != "" {
		search_like := fmt.Sprintf("%%%s%%", *req.UserAccountNumber)
		count = count.Where("users.bank_account LIKE ?", search_like)
	}

	if err = count.
		Where("users.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "users.id, users.member_code, users.username, users.phone, users.firstname, users.lastname, users.fullname, users.credit, users.bankname, users.bank_account, users.promotion, users.status, users.channel, users.true_wallet, users.note, users.turnover_limit, users.created_at"
		query := r.db.Table("Users as users")
		query = query.Select(selectedFields)
		if req.UserBankCode != nil {
			query = query.Where("users.bank_code = ?", req.UserBankCode)
		}
		if req.UserAccountNumber != nil && *req.UserAccountNumber != "" {
			search_like := fmt.Sprintf("%%%s%%", *req.UserAccountNumber)
			query = query.Where("users.bank_account LIKE ?", search_like)
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
			Where("users.deleted_at IS NULL").
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

func (r repo) GetMemberTransactions(req model.MemberTransactionListRequest) (*model.SuccessWithPagination, error) {

	var list []model.MemberTransaction
	var total int64
	var err error

	// Count total records for pagination purposes (without limit and offset) //
	count := r.db.Table("Bank_transactions as transactions")
	count = count.Select("transactions.id")
	count = count.Where("transactions.removed_at IS NULL")
	if req.UserId != "" {
		count = count.Where("transactions.user_id = ?", req.UserId)
	}
	if req.FromTransferDate != "" {
		count = count.Where("transactions.transfer_at >= ?", req.FromTransferDate)
	}
	if req.ToTransferDate != "" {
		count = count.Where("transactions.transfer_at <= ?", req.ToTransferDate)
	}
	if req.TransferType != "" {
		count = count.Where("transactions.transfer_type = ?", req.TransferType)
	}
	if req.Search != "" {
		search_like := fmt.Sprintf("%%%s%%", req.Search)
		count = count.Where(r.db.Where("transactions.member_code LIKE ?", search_like).Or("transactions.from_account_name LIKE ?", search_like).Or("transactions.from_account_number LIKE ?", search_like).Or("transactions.to_account_name LIKE ?", search_like).Or("transactions.to_account_number LIKE ?", search_like))
	}

	if err = count.
		Where("transactions.deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}
	if total > 0 {
		// SELECT //
		selectedFields := "transactions.id, transactions.user_id, transactions.transfer_type, transactions.promotion_id, transactions.from_account_id, transactions.from_bank_id, transactions.from_account_name, transactions.from_account_number, transactions.to_account_id, transactions.to_bank_id, transactions.to_account_name, transactions.to_account_number"
		selectedFields += ", transactions.credit_amount, transactions.paid_amount, transactions.over_amount, transactions.before_amount, transactions.after_amount, transactions.bank_charge_amount, transactions.transfer_at, transactions.created_by_user_id, transactions.created_by_username, transactions.removed_at, transactions.removed_by_user_id, transactions.removed_by_username, transactions.status, transactions.status_detail, transactions.is_auto_credit"
		selectedFields += ", transactions.created_at, transactions.updated_at"
		selectedFields += ", from_banks.name as from_bank_name, from_banks.code as from_bank_code, from_banks.icon_url as from_bank_icon_url, from_banks.type_flag as from_bank_type_flag"
		selectedFields += ", to_banks.name as to_bank_name, to_banks.code as to_bank_code, to_banks.icon_url as to_bank_icon_url, to_banks.type_flag as to_bank_type_flag"
		selectedFields += ", users.member_code as member_code, users.username as user_username, users.firstname as user_firstname, users.lastname as user_lastname, users.fullname as user_fullname, users.phone as user_phone"
		query := r.db.Table("Bank_transactions as transactions")
		query = query.Select(selectedFields)
		query = query.Joins("LEFT JOIN Banks as from_banks ON from_banks.id = transactions.from_bank_id")
		query = query.Joins("LEFT JOIN Banks as to_banks ON to_banks.id = transactions.to_bank_id")
		query = query.Joins("LEFT JOIN Users as users ON users.id = transactions.user_id")
		query = query.Where("transactions.removed_at IS NULL")
		if req.UserId != "" {
			query = query.Where("transactions.user_id = ?", req.UserId)
		}
		if req.FromTransferDate != "" {
			query = query.Where("transactions.transfer_at >= ?", req.FromTransferDate)
		}
		if req.ToTransferDate != "" {
			query = query.Where("transactions.transfer_at <= ?", req.ToTransferDate)
		}
		if req.TransferType != "" {
			query = query.Where("transactions.transfer_type = ?", req.TransferType)
		}
		if req.Search != "" {
			search_like := fmt.Sprintf("%%%s%%", req.Search)
			query = query.Where(r.db.Where("transactions.member_code LIKE ?", search_like).Or("transactions.from_account_name LIKE ?", search_like).Or("transactions.from_account_number LIKE ?", search_like).Or("transactions.to_account_name LIKE ?", search_like).Or("transactions.to_account_number LIKE ?", search_like))
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

func (r repo) GetMemberTransactionSummary(req model.MemberTransactionListRequest) (*model.MemberTransactionSummary, error) {

	var result model.MemberTransactionSummary
	var err error

	// SELECT //
	selectedFields := "SUM(case when transfer_type = 'deposit' then credit_amount else 0 end) as total_deposit_amount, SUM(case when transfer_type = 'withdraw' then credit_amount else 0 end) as total_withdraw_amount, SUM(bonus_amount) as total_bom_amount"
	query := r.db.Table("Bank_transactions as transactions")
	query = query.Select(selectedFields)
	query = query.Joins("LEFT JOIN Banks as from_banks ON from_banks.id = transactions.from_bank_id")
	query = query.Joins("LEFT JOIN Banks as to_banks ON to_banks.id = transactions.to_bank_id")
	query = query.Where("transactions.removed_at IS NULL")
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
		query = query.Where(r.db.Where("transactions.member_code LIKE ?", search_like).Or("transactions.from_account_name LIKE ?", search_like).Or("transactions.from_account_number LIKE ?", search_like).Or("transactions.to_account_name LIKE ?", search_like).Or("transactions.to_account_number LIKE ?", search_like))
	}

	if err = query.
		Where("transactions.deleted_at IS NULL").
		Scan(&result).
		Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r repo) IncreaseMemberCredit(userId int64, amount float32) error {

	if err := r.db.Table("Users").Where("id = ?", userId).UpdateColumn("credit", gorm.Expr("credit + ?", amount)).Error; err != nil {
		return err
	}
	return nil
}

func (r repo) DecreaseMemberCredit(userId int64, amount float32) error {

	if err := r.db.Table("Users").Where("id = ?", userId).UpdateColumn("credit", gorm.Expr("credit - ?", amount)).Error; err != nil {
		return err
	}
	return nil
}
