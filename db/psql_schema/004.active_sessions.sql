CREATE VIEW active_sessions AS SELECT * FROM sessions WHERE NOW() < expires_at;
