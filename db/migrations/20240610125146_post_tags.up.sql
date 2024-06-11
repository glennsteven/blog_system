CREATE TABLE IF NOT EXISTS post_tags(
    post_id INT,
    tag_id INT,

    CONSTRAINT fk_post_id FOREIGN KEY(post_id) REFERENCES posts(id),
    CONSTRAINT fk_tag_id FOREIGN KEY(tag_id) REFERENCES tags(id)
);