CREATE TABLE IF NOT EXISTS users (
                                     id INT AUTO_INCREMENT PRIMARY KEY,
                                     email VARCHAR(255) NOT NULL,
                                     name VARCHAR(255) NOT NULL,
                                     password BINARY(60) NOT NULL,
                                     created_at TIMESTAMP DEFAULT NOW(),
                                     updated_at TIMESTAMP DEFAULT NOW(),
                                     deleted_at TIMESTAMP,
                                     active BOOL NOT NULL
)  ENGINE=INNODB;