ALTER TABLE `setting_web`
ADD COLUMN `auto_withdraw` VARCHAR(20) DEFAULT NULL AFTER `otp_register`;

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
        `client_id` varchar(255) DEFAULT NULL,
        `client_secret` varchar(255) DEFAULT NULL,
        `response_type` varchar(255) DEFAULT NULL,
        `redirect_uri` varchar(255) DEFAULT NULL,
        `scope` varchar(255) DEFAULT NULL,
        `state` varchar(255) DEFAULT NULL,
        `created_at` DATETIME DEFAULT current_timestamp(),
        `updated_at` DATETIME DEFAULT NULL
    );

INSERT INTO
    `Type_notify` (
        `name`,
        `client_id`,
        `client_Secret`,
        `response_type`,
        `redirect_uri`,
        `scope`,
        `state`
    )
VALUES
    (
        'แจ้ง สมัครสมาชิก',
        'mV1hDYHmq38qDq3gWe1v4c',
        'ecNcbiTrJW9qsoC0Mcd0ajlS9XdDUuqPsrV7R21vjYf',
        'code',
        'https://cyberrichdigital.com/',
        'notify',
        '1'
    ),
    (
        'แจ้งฝาก ก่อนปรับเครดิต',
        '261SlRtTayZXja0MCj1Qdt',
        'yupExXSc0rZczMhbdoUHsVQYoLXe6MqJrAPZ1w18gOO',
        'code',
        'https://cyberrichdigital.com/',
        'notify',
        '1'
    ),
    (
        'แจ้งฝาก หลังปรับเครดิต',
        'lq2YMfcpuIwUQPiSWL3ffL',
        'Gqur4SPg5xpqC3f77SiXZMc2utxGIg1ChYVKZoTiilu',
        'code',
        'https://cyberrichdigital.com/',
        'notify',
        '1'
    ),
    (
        'แจ้งถอน ก่อนปรับเครดิต',
        'RPgAqM8fb98Lq5lLsutpMJ',
        'Ik91e14JaKXvl3EPcgbm0uMpkgcR9ktUzNzyh2GjeBQ',
        'code',
        'https://cyberrichdigital.com/',
        'notify',
        '1'
    ),
    (
        'แจ้งถอน รอโอนเงิน',
        'i5U8AFnsH81k7qjTcZrD2V',
        '9GSXpAZVHbPxxT7zEbO9bdoxpDZg1XkAXwZZVr5sypi',
        'code',
        'https://cyberrichdigital.com/',
        'notify',
        '1'
    ),
    (
        'แจ้งถอน หลังปรับเครดิต',
        'i5U8AFnsH81k7qjTcZrD2V',
        '9GSXpAZVHbPxxT7zEbO9bdoxpDZg1XkAXwZZVr5sypi',
        'code',
        'https://cyberrichdigital.com/',
        'notify',
        '1'
    );

CREATE TABLE
    `User_linenotify` (
        `id` bigint PRIMARY KEY AUTO_INCREMENT,
        `user_id` bigint (20) DEFAULT NULL,
        `type_notify_id` varchar(255) DEFAULT NULL,
        `token` varchar(255) DEFAULT NULL,
        `created_at` DATETIME DEFAULT current_timestamp(),
        `updated_at` DATETIME DEFAULT NULL
    );

ALTER TABLE `User_linenotify` ADD INDEX `idx_user_id` (`user_id`);