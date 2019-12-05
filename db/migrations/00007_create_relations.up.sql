ALTER TABLE image_gallery.pictures
    ADD CONSTRAINT pictures_users_fk FOREIGN KEY (user_id) REFERENCES image_gallery.users (user_id),
    ADD CONSTRAINT images_pictures_fk FOREIGN KEY (image_id) REFERENCES image_gallery.images (image_id);

ALTER TABLE image_gallery.users
    ADD CONSTRAINT images_users_fk FOREIGN KEY (avatar_id) REFERENCES image_gallery.images (image_id);

ALTER TABLE image_gallery.likes
    ADD CONSTRAINT likes_users_fk FOREIGN KEY (user_id) REFERENCES image_gallery.users (user_id),
    ADD CONSTRAINT likes_pictures_fk FOREIGN KEY (picture_id) REFERENCES image_gallery.pictures (picture_id);

ALTER TABLE image_gallery.comments
    ADD CONSTRAINT comments_users_fk FOREIGN KEY (user_id) REFERENCES image_gallery.users (user_id),
    ADD CONSTRAINT comments_pictures_fk FOREIGN KEY (picture_id) REFERENCES image_gallery.pictures (picture_id);

CREATE TABLE image_gallery.gallery_has_picture (
    gallery_id INT NOT NULL,
    picture_id INT NOT NULL,
    CONSTRAINT gallery_contains_pictures_fk FOREIGN KEY (gallery_id) REFERENCES image_gallery.galleries (gallery_id),
    CONSTRAINT picture_belongs_to_gallery_fk FOREIGN KEY (picture_id) REFERENCES image_gallery.pictures (picture_id)
)



