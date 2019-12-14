ALTER TABLE image_gallery.picture
    ADD CONSTRAINT pictures_users_fk FOREIGN KEY (user_id) REFERENCES image_gallery.user (user_id) ON DELETE CASCADE,
    ADD CONSTRAINT images_pictures_fk FOREIGN KEY (image_id) REFERENCES image_gallery.image (image_id);

ALTER TABLE image_gallery.user
    ADD CONSTRAINT images_users_fk FOREIGN KEY (avatar_id) REFERENCES image_gallery.image (image_id) ON DELETE SET NULL;

ALTER TABLE image_gallery.like
    ADD CONSTRAINT likes_users_fk FOREIGN KEY (user_id) REFERENCES image_gallery.user (user_id) ON DELETE CASCADE,
    ADD CONSTRAINT likes_pictures_fk FOREIGN KEY (picture_id) REFERENCES image_gallery.picture (picture_id) ON DELETE CASCADE;

ALTER TABLE image_gallery.comment
    ADD CONSTRAINT comments_users_fk FOREIGN KEY (user_id) REFERENCES image_gallery.user (user_id) ON DELETE CASCADE,
    ADD CONSTRAINT comments_pictures_fk FOREIGN KEY (picture_id) REFERENCES image_gallery.picture (picture_id) ON DELETE CASCADE;

CREATE TABLE image_gallery.gallery_has_picture (
    gallery_id INT NOT NULL,
    picture_id INT NOT NULL,
    CONSTRAINT gallery_contains_pictures_fk FOREIGN KEY (gallery_id) REFERENCES image_gallery.gallery (gallery_id),
    CONSTRAINT picture_belongs_to_gallery_fk FOREIGN KEY (picture_id) REFERENCES image_gallery.picture (picture_id)
)



