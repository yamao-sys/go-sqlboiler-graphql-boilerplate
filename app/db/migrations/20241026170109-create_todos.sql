
-- +migrate Up
CREATE TABLE IF NOT EXISTS todos(
	id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
	user_id INT NOT NULL,
	title VARCHAR(255) NOT NULL,
	content TEXT,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL,
	index index_user_id (user_id),
	CONSTRAINT fk_todos_users FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +migrate Down
DROP TABLE IF EXISTS todos;
