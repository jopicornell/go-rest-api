CREATE TABLE IF NOT EXISTS time_management.appointments (
    id SERIAL PRIMARY KEY,
    start_date TIMESTAMP NOT NULL,
    duration integer not null,
    end_date TIMESTAMP NOT NULL,
    status varchar(255) NOT NULL,
    user_id integer not null,
    resource_id integer not null,
    CONSTRAINT fk_appointments_user
        FOREIGN KEY (user_id) REFERENCES auth.users (id)
            ON DELETE RESTRICT,
    CONSTRAINT fk_appointments_resource
        FOREIGN KEY (resource_id) REFERENCES auth.users (id)
            ON DELETE RESTRICT
);