CREATE TABLE auth.role_user (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    role VARCHAR(255) NOT NULL,

    CONSTRAINT fk_role_user
        FOREIGN KEY (user_id) REFERENCES auth.users (id)
            ON DELETE RESTRICT
);