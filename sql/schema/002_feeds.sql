-- +goose Up
CREATE TABLE feeds(
	id UUID PRIMARY KEY,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	name TEXT,
	url TEXT,
	user_id UUID,
	UNIQUE (user_id, url),
	FOREIGN KEY(user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
