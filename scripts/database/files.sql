CREATE TABLE files (
    id SERIAL,
    folder_id INT,
    owner_id INT,
    name VARCHAR(200) NOT NULL,
    type VARCHAR(50) NOT NULL,
    path VARCHAR(250) NOT NULL,
    created_at TIMESTAMP default current_timestamp,
    modified_at TIMESTAMP NOT NULL,
    deleted BOOL NOT NULL DEFAULT false,
    PRIMARY KEY(id),
    CONSTRAINT fk_folders
        foreign key(folder_id)
            references folders(id),
    CONSTRAINT fk_owner
        foreign key(owner_id)
            references users(id)

)