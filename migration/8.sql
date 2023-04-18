CREATE TABLE
    `Line_notify` (
        `id` bigint PRIMARY KEY AUTO_INCREMENT,
        `start_credit` DECIMAL(14, 2) NULL DEFAULT 0.00,
        `token` varchar(255) DEFAULT NULL,
        `notify_id` bigint (20),
        status ENUM ('ACTIVE', 'DEACTIVE') DEFAULT 'ACTIVE',
        `created_at` DATETIME DEFAULT current_timestamp(),
        `updated_at` DATETIME DEFAULT NULL
    );

ALTER TABLE `Line_notify` ADD INDEX `idx_notify_id` (`notify_id`);

CREATE TABLE
    `Type_notify` (
        `id` bigint PRIMARY KEY AUTO_INCREMENT,
        `name` varchar(255) DEFAULT NULL,
        `created_at` DATETIME DEFAULT current_timestamp(),
        `updated_at` DATETIME DEFAULT NULL
    );

INSERT INTO
    `Type_notify` (`name`)
VALUES
    ('แจ้ง สมัครสมาชิก'),
    ('แจ้งฝาก ก่อนปรับเครดิต'),
    ('แจ้งฝาก หลังปรับเครดิต'),
    ('แจ้งถอน ก่อนปรับเครดิต'),
    ('แจ้งถอน รอโอนเงิน'),
    ('แจ้งถอน หลังปรับเครดิต')