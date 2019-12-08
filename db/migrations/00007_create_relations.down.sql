ALTER TABLE image_gallery.picture
    DROP CONSTRAINT pictures_users_fk,
    DROP CONSTRAINT images_pictures_fk;

ALTER TABLE image_gallery.customer
    DROP CONSTRAINT images_users_fk;

ALTER TABLE image_gallery.like
    DROP CONSTRAINT likes_users_fk,
    DROP CONSTRAINT likes_pictures_fk;

ALTER TABLE image_gallery.comment
    DROP CONSTRAINT comments_users_fk,
    DROP CONSTRAINT comments_pictures_fk;

DROP TABLE image_gallery.gallery_has_picture;



