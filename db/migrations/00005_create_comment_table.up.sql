CREATE TABLE image_gallery.comment (
    comment_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    picture_id INT NOT NULL,
    created TIMESTAMP DEFAULT NOW(),
    comment TEXT NOT NULL
);

CREATE INDEX comment_user_picture_idx ON image_gallery.comment (picture_id);