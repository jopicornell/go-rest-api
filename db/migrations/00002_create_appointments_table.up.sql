CREATE TABLE IF NOT EXISTS appointments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    start_date DATETIME NOT NULL,
    duration integer not null,
    end_date DATETIME NOT NULL,
    status varchar(255) NOT NULL,
    user_id integer not null,
    resource_id integer not null,
    CONSTRAINT `fk_appointments_user`
        FOREIGN KEY (user_id) REFERENCES users (id)
            ON DELETE RESTRICT,
    CONSTRAINT `fk_appointments_resource`
        FOREIGN KEY (resource_id) REFERENCES users (id)
            ON DELETE RESTRICT
)  ENGINE=INNODB;