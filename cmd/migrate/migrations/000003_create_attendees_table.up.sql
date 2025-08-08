CREATE TABLE IF NOT EXISTS attendees (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    event_id INTEGER NOT NULL,
    Foreign Key (user_id) REFERENCES users (id) ON DELETE CASCADE,
    Foreign Key (event_id) REFERENCES events (id) ON DELETE CASCADE
);