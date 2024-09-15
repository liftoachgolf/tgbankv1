CREATE TABLE accounts(
    id bigserial PRIMARY KEY,
    chat_id bigint NOT NULL UNIQUE,
    username varchar(255) NOT NULL,
    balance bigint NOT NULL,
    currency varchar(3) NOT NULL
);

CREATE TABLE messages (
    chat_id bigint NOT NULL,
    message_id bigint NOT NULL,
    text varchar(255) NOT NULL,
    FOREIGN KEY (chat_id) REFERENCES accounts (chat_id)
);

CREATE TABLE entries (
    id bigserial PRIMARY KEY,
    chat_id bigint NOT NULL,
    amount bigint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    FOREIGN KEY (chat_id) REFERENCES accounts (chat_id),
    CHECK (amount != 0)
);

CREATE TABLE transfers (
    id bigserial PRIMARY KEY,
    from_chat_id bigint NOT NULL,
    to_chat_id bigint NOT NULL,
    amount bigint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now()),
    FOREIGN KEY (from_chat_id) REFERENCES accounts (chat_id),
    FOREIGN KEY (to_chat_id) REFERENCES accounts (chat_id),
    CHECK (amount > 0)
);

CREATE INDEX ON entries (chat_id);
CREATE INDEX ON transfers (from_chat_id, to_chat_id);

COMMENT ON COLUMN entries.amount IS 'can be negative or positive';
COMMENT ON COLUMN transfers.amount IS 'must be positive';