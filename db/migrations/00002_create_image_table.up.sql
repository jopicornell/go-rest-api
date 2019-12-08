CREATE TABLE IF NOT EXISTS image_gallery.image (
    image_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    thumb_url VARCHAR(255) NOT NULL,
    low_res_url VARCHAR(255) NOT NULL,
    high_res_url VARCHAR(255) NOT NULL
);