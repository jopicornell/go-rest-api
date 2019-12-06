CREATE TABLE image_gallery.pictures (
    picture_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    image_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    created TIMESTAMP DEFAULT NOW(),
    num_likes INT NOT NULL,
    num_comments INT
)