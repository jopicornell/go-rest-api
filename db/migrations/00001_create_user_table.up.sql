CREATE TABLE IF NOT EXISTS image_gallery.customer (
                                     user_id SERIAL PRIMARY KEY,
                                     username VARCHAR(255) NOT NULL,
                                     full_name VARCHAR(255) NOT NULL,
                                     password VARCHAR(255) NOT NULL,
                                     num_pictures INT NOT NULL DEFAULT 0,
                                     num_comments INT NOT NULL DEFAULT 0,
                                     num_likes INT NOT NULL DEFAULT 0,
                                     avatar_id INT
);