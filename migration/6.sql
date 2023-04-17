CREATE Table
    Webhook_logs (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        json_request TEXT NOT NULL,
        json_payload TEXT NOT NULL,
        log_type VARCHAR(255) NOT NULL,
        status VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );
