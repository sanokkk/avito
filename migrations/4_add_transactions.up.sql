CREATE TABLE histories (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    from_user_id UUID NOT NULL,
    to_user_id UUID NOT NULL,
    amount int NOT NULL
);