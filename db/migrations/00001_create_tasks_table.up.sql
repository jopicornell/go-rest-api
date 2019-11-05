CREATE TABLE IF NOT EXISTS tasks (
                                     task_id INT AUTO_INCREMENT PRIMARY KEY,
                                     title VARCHAR(255) NOT NULL,
                                     date DATETIME,
                                     completed BOOL NOT NULL,
                                     description TEXT
)  ENGINE=INNODB;