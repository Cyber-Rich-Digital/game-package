
CREATE Table 
    Botaccount_logs (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        external_id BIGINT NOT NULL,
        client_name VARCHAR(255) NOT NULL,
        log_type VARCHAR(255) NOT NULL,
        message TEXT NOT NULL,
        external_create_date VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

ALTER TABLE `Botaccount_logs`
    ADD INDEX `idx_external_id` (`external_id`);

CREATE Table 
    Botaccount_statements (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        external_id BIGINT NOT NULL,
        bank_account_id BIGINT NOT NULL,
        bank_code VARCHAR(255) NOT NULL,
        amount DECIMAL(14,2) NOT NULL,
        date_time DATETIME NOT NULL,
        raw_date_time DATETIME NOT NULL,
        info VARCHAR(255) NOT NULL,
        channel_code VARCHAR(255) NOT NULL,
        channel_description VARCHAR(255) NOT NULL,
        txn_code VARCHAR(255) NOT NULL,
        txn_description VARCHAR(255) NOT NULL,
        checksum VARCHAR(255) NOT NULL,
        is_read BOOLEAN NOT NULL,
        external_create_date VARCHAR(255) NOT NULL,
        extermal_update_date VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

ALTER TABLE `Botaccount_statements`
    ADD INDEX `idx_external_id` (`external_id`),
    ADD INDEX `idx_bank_account_id` (`bank_account_id`);