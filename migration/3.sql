CREATE Table
    Bank_statements (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        account_id BIGINT NOT NULL,
        amount DECIMAL(14,2) NOT NULL,
        detail VARCHAR(255) NOT NULL,
        statement_type VARCHAR(255) NOT NULL,
        transfer_at DATETIME NOT NULL,
        status VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

CREATE Table
    Bank_transactions (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        member_code VARCHAR(255) NOT NULL,
        user_id BIGINT NOT NULL,
        transfer_type VARCHAR(255) NOT NULL,
        promotion_id BIGINT NULL,
        from_account_id BIGINT NULL,
        from_bank_id BIGINT NULL,
        from_account_name VARCHAR(255) NULL,
        from_account_number VARCHAR(255) NULL,
        to_account_id BIGINT NULL,
        to_bank_id BIGINT NULL,
        to_account_name VARCHAR(255) NULL,
        to_account_number VARCHAR(255) NULL,
        credit_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        paid_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        over_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        deposit_channel VARCHAR(255) NULL,
        bonus_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        bonus_reason VARCHAR(255) NULL,
        before_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        after_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        bank_charge_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        transfer_at DATETIME NOT NULL,
        created_by_user_id BIGINT NOT NULL,
        created_by_username VARCHAR(255) NOT NULL,
        cancel_remark VARCHAR(255) NULL,
        canceled_at DATETIME NULL,
        canceled_by_user_id BIGINT NULL,
        canceled_by_username VARCHAR(255) NULL,
        confirmed_at DATETIME NULL,
        confirmed_by_user_id BIGINT NULL,
        confirmed_by_username VARCHAR(255) NULL,
        removed_at DATETIME NULL,
        removed_by_user_id BIGINT NULL,
        removed_by_username VARCHAR(255) NULL,
        status VARCHAR(255) NOT NULL,
        status_detail VARCHAR(255) NULL,
        is_auto_credit TINYINT NOT NULL DEFAULT 0,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

CREATE Table
    Bank_confirm_transactions (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        transaction_id BIGINT NOT NULL,
        user_id BIGINT NOT NULL,
        transfer_type VARCHAR(255) NOT NULL,
        from_account_id BIGINT NULL,
        to_account_id BIGINT NULL,
        json_before TEXT NULL,
        transfer_at DATETIME NOT NULL,
        slip_url VARCHAR(255) NULL,
        bonus_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        confirmed_at DATETIME NULL,
        confirmed_by_user_id BIGINT NULL,
        confirmed_by_username VARCHAR(255) NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

ALTER TABLE `bank_confirm_transactions`
	ADD UNIQUE INDEX `uni_transaction_id` (`transaction_id`);