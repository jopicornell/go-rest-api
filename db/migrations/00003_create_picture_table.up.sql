CREATE TABLE image_gallery.picture (
    picture_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    image_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    created TIMESTAMP DEFAULT NOW(),
    num_likes INT NOT NULL DEFAULT 0,
    num_comments INT NOT NULL DEFAULT 0
)