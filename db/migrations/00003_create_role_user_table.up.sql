CREATE TABLE `role_user` (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    role VARCHAR(255) NOT NULL,

    CONSTRAINT `fk_role_user`
        FOREIGN KEY (user_id) REFERENCES users (id)
            ON DELETE RESTRICT
) ENGINE=INNODB