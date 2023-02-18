CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    email TEXT NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE posts (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    category TEXT NOT NULL,
    creation_time TIMESTAMP NOT NULL,
    author TEXT NOT NULL,
    likes INTEGER NOT NULL DEFAULT 0,
    dislikes INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE comments (
    id INTEGER PRIMARY KEY,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    likes INTEGER NOT NULL DEFAULT 0,
    dislikes INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

INSERT INTO users (email, username, password) VALUES
('user1@example.com', 'user1', '$2a$10$ZQkQV7cG.qfOCwV7vChyFO0iV92rOftOmyx9bKjhTpYhpf0soiH5W'),
('user2@example.com', 'user2', '$2a$10$rQqBvLRJsvxy5nmb.ZMZ1OMR5CFtbc5lxyLL.e1hBr.Rkw8E2Xqo.');

INSERT INTO posts (title, content, category, creation_time, author) VALUES
('First post', 'This is the content of the first post.', 'Uncategorized', '2023-02-19 15:30:00', 'user1'),
('Second post', 'This is the content of the second post.', 'News', '2023-02-20 16:45:00', 'user2'),
('Third post', 'This is the content of the third post.', 'Events', '2023-02-21 10:15:00', 'user1');

INSERT INTO comments (post_id, user_id, content, created_at, updated_at) VALUES
(1, 2, 'Great post!', '2023-02-20 14:00:00', '2023-02-20 14:00:00'),
(1, 1, 'I agree, this is a really good post.', '2023-02-21 10:00:00', '2023-02-21 10:00:00'),
(2, 1, 'Interesting article.', '2023-02-21 12:30:00', '2023-02-21 12:30:00'),
(3, 2, 'Thanks for sharing this event.', '2023-02-22 09:00:00', '2023-02-22 09:00:00');
