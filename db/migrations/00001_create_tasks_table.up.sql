CREATE TABLE IF NOT EXISTS tasks (
                                     id INT AUTO_INCREMENT PRIMARY KEY,
                                     title VARCHAR(255) NOT NULL,
                                     start_date DATE,
                                     due_date DATE,
                                     status VARCHAR(50) NOT NULL,
                                     priority TINYINT NOT NULL,
                                     description TEXT,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)  ENGINE=INNODB;