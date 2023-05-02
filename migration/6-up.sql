
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

ALTER TABLE `Bank_accounts`
	ADD COLUMN `external_id` BIGINT NULL AFTER `pin_code`;

ALTER TABLE `Bank_accounts`
    ADD INDEX `idx_external_id` (`external_id`);

ALTER TABLE `Bank_statements`
	ADD COLUMN `external_id` BIGINT(19) NOT NULL AFTER `account_id`,
	ADD COLUMN `from_bank_id` BIGINT(19) NULL AFTER `detail`,
	ADD COLUMN `from_account_number` VARCHAR(255) NULL AFTER `from_bank_id`;

ALTER TABLE `Bank_statements`
    ADD UNIQUE `uni_external_id` (`external_id`);

CREATE TABLE 
    `Botaccount_config` (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        config_key VARCHAR(255) NOT NULL,
        config_val VARCHAR(255) NOT NULL
    );

ALTER TABLE `Botaccount_config`
    ADD INDEX `idx_config_key` (`config_key`);


INSERT INTO `Botaccount_config` (`config_key`, `config_val`) VALUES
	('allow_create_external_account', '_all'),
	('allow_create_external_account', '_list'),
	('allow_external_account_number', 'set to list and set account number');

ALTER TABLE `Bank_account_types`
	ADD COLUMN `allow_deposit` TINYINT NOT NULL DEFAULT 0 AFTER `limit_flag`,
	ADD COLUMN `allow_withdraw` TINYINT NOT NULL DEFAULT 0 AFTER `allow_deposit`;

UPDATE `Bank_account_types` SET `allow_deposit`=1, `allow_withdraw`=0 WHERE `limit_flag`='00001000';
UPDATE `Bank_account_types` SET `allow_deposit`=0, `allow_withdraw`=1 WHERE `limit_flag`='00000100';
UPDATE `Bank_account_types` SET `allow_deposit`=1, `allow_withdraw`=1 WHERE `limit_flag`='00001100';

ALTER TABLE `Bank_confirm_transactions`
	ADD COLUMN `credit_amount` DECIMAL(14,2) NULL DEFAULT NULL AFTER `slip_url`,
	CHANGE COLUMN `bonus_amount` `bonus_amount` DECIMAL(14,2) NOT NULL DEFAULT 0 AFTER `credit_amount`;

ALTER TABLE `Bank_confirm_transactions`
	ADD COLUMN `bank_charge_amount` DECIMAL(14,2) NOT NULL DEFAULT '0.00' AFTER `bonus_amount`;

CREATE Table 
    Bank_confirm_statements (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        statement_id BIGINT NOT NULL,
        action_type VARCHAR(255) NOT NULL,
        user_id BIGINT NULL,
        account_id BIGINT NOT NULL,
        json_before TEXT NOT NULL,
        confirmed_at DATETIME NULL,
        confirmed_by_user_id BIGINT NULL,
        confirmed_by_username VARCHAR(255) NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

ALTER TABLE `Bank_confirm_statements`
	ADD UNIQUE INDEX `uni_statement_id` (`statement_id`),
    ADD INDEX `idx_account_id` (`account_id`);