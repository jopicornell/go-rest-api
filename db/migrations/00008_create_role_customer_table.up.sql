CREATE TABLE image_gallery.user_has_roles (
    user_id INT NOT NULL,
    role VARCHAR(255) NOT NULL,
    CONSTRAINT user_has_roles_fk FOREIGN KEY (user_id) REFERENCES image_gallery.user (user_id)
)