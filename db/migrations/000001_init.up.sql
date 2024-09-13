CREATE TABLE accounts (
    id bigint NOT NULL UNIQUE,
    chat_id bigserial PRIMARY KEY,
    username varchar(255) NOT NULL,
    balance bigint NOT NULL,
    currency varchar(3) NOT NULL
);

CREATE TABLE messages (
    chat_id bigint NOT NULL,
    message_id bigserial PRIMARY KEY,
    text varchar(255) NOT NULL,
);

CREATE TABLE entries (
    id bigserial PRIMARY KEY,
    chat_id bigint NOT NULL,
    account_id bigint NOT NULL,
    amount bigint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT(now())
);

CREATE TABLE transfers (
    id bigserial PRIMARY KEY,
    chat_id bigint NOT NULL,
    from_chat_id bigint NOT NULL,
    to_chat_id bigint NOT NULL,
    amount bigint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT(now())
);

ALTER TABLE messages ADD FOREIGN KEY (chat_id) REFERENCES accounts (chat_id);
ALTER TABLE entries ADD FOREIGN KEY (chat_id) REFERENCES accounts (chat_id);
ALTER TABLE entries ADD FOREIGN KEY (account_id) REFERENCES accounts (id);
ALTER TABLE transfers ADD FOREIGN KEY (chat_id) REFERENCES accounts (chat_id);
ALTER TABLE transfers ADD FOREIGN KEY (from_chat_id) REFERENCES accounts (chat_id);
ALTER TABLE transfers ADD FOREIGN KEY (to_chat_id) REFERENCES accounts (chat_id);

CREATE INDEX ON entries (chat_id);
CREATE INDEX ON transfers (chat_id);
CREATE INDEX ON transfers (from_chat_id, to_chat_id);

COMMENT ON COLUMN entries.amount IS 'can be negative or positive';
COMMENT ON COLUMN transfers.amount IS 'must be positive';
