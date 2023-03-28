CREATE Table
    Tags (
        id INT(11) PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(255) NULL,
        website_id INT(11) NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE Table
    Websites (
        id INT(11) PRIMARY KEY AUTO_INCREMENT,
        title VARCHAR(255) NULL,
        domain_name VARCHAR(255) NULL,
        api_key VARCHAR(255) NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE Table
    Messages (
        id INT(11) PRIMARY KEY AUTO_INCREMENT,
        message VARCHAR(255) NULL,
        tag_id INT(11) NULL,
        website_id INT(11) NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE Table
    Users (
        id INT PRIMARY KEY AUTO_INCREMENT,
        username VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL,
        role ENUM('ADMIN', 'USER', 'SUPER-ADMIN') NOT NULL DEFAULT 'USER',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP NULL ON UPDATE CURRENT_TIMESTAMP
    );

CREATE Table
    Devices (
        id INT(11) PRIMARY KEY AUTO_INCREMENT,
        os VARCHAR(255) NULL,
        version VARCHAR(255) NULL,
        fcm_token VARCHAR(255) NULL,
        hardware_id VARCHAR(255) NULL,
        website_id INT(11) NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );