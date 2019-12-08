CREATE TABLE image_gallery.customer_has_roles (
    user_id INT NOT NULL,
    role VARCHAR(255) NOT NULL,
    CONSTRAINT customer_has_roles_fk FOREIGN KEY (user_id) REFERENCES image_gallery.customer (user_id)
)