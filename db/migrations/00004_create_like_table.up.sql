create table image_gallery.like(
    user_id INT NOT NULL,
    picture_id INT NOT NULL,
    PRIMARY KEY(user_id, picture_id)
);

CREATE INDEX like_user_pictures_idx ON image_gallery.like (picture_id);