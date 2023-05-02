DROP TABLE IF EXISTS `Botaccount_logs`;
DROP TABLE IF EXISTS `Botaccount_statements`;

ALTER TABLE `Bank_accounts`
    DROP COLUMN `external_id`,
    DROP INDEX `idx_external_id`;

ALTER TABLE `Bank_statements`
    DROP COLUMN `external_id`,
    DROP COLUMN `from_bank_id`,
    DROP COLUMN `from_account_number`,
    DROP INDEX `uni_external_id`;

DROP TABLE IF EXISTS `Botaccount_config`;

ALTER TABLE `Bank_account_types`
    DROP COLUMN `allow_deposit`,
    DROP COLUMN `allow_withdraw`;

ALTER TABLE `Bank_confirm_transactions`
	DROP COLUMN `credit_amount`,
	DROP COLUMN `bank_charge_amount`;