CREATE TABLE items (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    cost int NOT NULL
);

CREATE TABLE user_items (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    item_id UUID NOT NULL,
    title TEXT NOT NULL,
    quantity int NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items (id) ON DELETE NO ACTION,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE NO ACTION
);