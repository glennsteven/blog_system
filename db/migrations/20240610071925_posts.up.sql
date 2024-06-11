CREATE TABLE IF NOT EXISTS posts(
    id SERIAL PRIMARY KEY ,
    title TEXT NULL ,
    content TEXT NOT NULL ,
    status smallint NOT NULL DEFAULT 1,
    drafting INT,
    publisher INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_drafting FOREIGN KEY(drafting) REFERENCES users(id),
    CONSTRAINT fk_publisher FOREIGN KEY(publisher) REFERENCES users(id)
);

CREATE INDEX idx_status ON posts USING BTREE (status);
CREATE INDEX idx_created_at ON posts USING BRIN (created_at);
