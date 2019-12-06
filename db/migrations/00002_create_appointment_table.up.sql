CREATE TABLE IF NOT EXISTS image_gallery.images (
    image_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    thumb_url varchar(255) not null,
    low_res_url varchar(255) NOT NULL,
    high_res_url varchar(255) NOT NULL
);