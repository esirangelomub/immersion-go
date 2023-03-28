CREATE TABLE folders (
    id  SERIAL,
    parent_id INT,
    name VARCHAR(60) NOT NULL,
    created_at TIMESTAMP default current_timestamp,
    modified_at TIMESTAMP NOT NULL,
    deleted BOOL NOT NULL DEFAULT false,
    PRIMARY KEY(id),
    CONSTRAINT fk_parent
        foreign key(parent_id)
            references folders(id)
);

