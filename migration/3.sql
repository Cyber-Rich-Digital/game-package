CREATE Table
    Bank_statements (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        account_id BIGINT NOT NULL,
        amount DECIMAL(14,2) NOT NULL,
        transfer_at DATETIME NOT NULL,
        status VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

CREATE Table
    Bank_transactions (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        transfer_type VARCHAR(255) NOT NULL,
        promotion_id BIGINT NOT NULL,
        from_account_id BIGINT NOT NULL,
        from_bank_id BIGINT NOT NULL,
        from_account_name VARCHAR(255) NOT NULL,
        from_account_number VARCHAR(255) NOT NULL,
        to_account_id BIGINT NOT NULL,
        to_bank_id BIGINT NOT NULL,
        to_account_name VARCHAR(255) NOT NULL,
        to_account_number VARCHAR(255) NOT NULL,
        credit_amount DECIMAL(14,2) NOT NULL,
        paid_amount DECIMAL(14,2) NOT NULL,
        over_amount DECIMAL(14,2) NOT NULL,
        before_amount DECIMAL(14,2) NOT NULL,
        after_amount DECIMAL(14,2) NOT NULL,
        transfer_at DATETIME NOT NULL,
        created_by_user_id BIGINT NOT NULL,
        created_by_username VARCHAR(255) NOT NULL,
        removed_at DATETIME NOT NULL,
        removed_by_user_id BIGINT NOT NULL,
        removed_by_username VARCHAR(255) NOT NULL,
        status VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

