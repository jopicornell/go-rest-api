CREATE TABLE image_gallery.comment (
    comment_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    picture_id INT NOT NULL,
    created TIMESTAMP DEFAULT NOW(),
    comment TEXT NOT NULL
);