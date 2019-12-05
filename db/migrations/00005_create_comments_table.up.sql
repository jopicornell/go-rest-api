CREATE TABLE image_gallery.comments (
    comment_id INT NOT NULL,
    user_id INT NOT NULL,
    picture_id INT NOT NULL,
    created TIMESTAMP,
    comment TEXT
);