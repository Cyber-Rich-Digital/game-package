CREATE Table
    Banks (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(255) NOT NULL,
        code VARCHAR(255) NOT NULL,
        icon_url VARCHAR(255) NOT NULL,
        type_flag VARCHAR(8) NOT NULL DEFAULT '00000000',
        created_at DATETIME NOT NULL DEFAULT NOW()
    );

ALTER TABLE `banks`
	ADD UNIQUE INDEX `uni_code` (`code`);

INSERT INTO `banks` (`name`, `code`, `icon_url`, `type_flag`) VALUES
	('ธนาคารกสิกรไทย', 'KBANK', '', '00001111'),
	('ธนาคารไทยพาณิชย์', 'SCB', '', '00001111'),
	('ธนาคารกรุงศรีอยุธยา', 'BAY', '', '00000011'),
	('ธนาคารกรุงเทพ', 'BBL', '', '00000011');

CREATE Table
    Bank_account_types (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(255) NOT NULL,
        limit_flag VARCHAR(8) NOT NULL DEFAULT '00000000',
        created_at DATETIME NOT NULL DEFAULT NOW()
    );

CREATE Table
    Bank_accounts (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        bank_id BIGINT NOT NULL,
        account_type_id BIGINT NOT NULL,
        account_name VARCHAR(255) NOT NULL,
        account_number VARCHAR(255) NOT NULL,
        transfer_priority VARCHAR(255) NOT NULL,
        account_status VARCHAR(255) NOT NULL,
        device_uid VARCHAR(255) NOT NULL,
        pin_code VARCHAR(255) NOT NULL,
        conection_status VARCHAR(255) NOT NULL,
        auto_credit_flag VARCHAR(255) NOT NULL,
        auto_withdraw_flag VARCHAR(255) NOT NULL,
        auto_withdraw_max_amount VARCHAR(255) NOT NULL,
        auto_transfer_max_amount VARCHAR(255) NOT NULL,
        qr_wallet_status VARCHAR(255) NOT NULL,
        created_at DATETIME NOT NULL DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

CREATE Table
    Bank_statements (
        id INT PRIMARY KEY AUTO_INCREMENT,
        description VARCHAR(255) NOT NULL,
        transfer_type VARCHAR(255) NOT NULL,
        amount DECIMAL(14,2) NOT NULL,
        transfer_at DATETIME NOT NULL,
        from_account_id BIGINT NOT NULL,
        from_bank_id BIGINT NOT NULL,
        from_account_name VARCHAR(255) NOT NULL,
        from_account_number VARCHAR(255) NOT NULL,
        to_account_id BIGINT NOT NULL,
        to_bank_id BIGINT NOT NULL,
        to_account_name VARCHAR(255) NOT NULL,
        to_account_number VARCHAR(255) NOT NULL,
        create_by_username VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );
