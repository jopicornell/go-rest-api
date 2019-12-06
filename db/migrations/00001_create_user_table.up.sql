CREATE TABLE IF NOT EXISTS image_gallery.users (
                                     user_id SERIAL PRIMARY KEY,
                                     username VARCHAR(255) NOT NULL,
                                     email VARCHAR(255) NOT NULL,
                                     full_name VARCHAR(255) NOT NULL,
                                     password BYTEA,
                                     num_pictures INT,
                                     num_comments INT,
                                     num_likes INT,
                                     avatar_id INT
);