CREATE TABLE IF NOT EXISTS test_outputs (
	tid 	INTEGER 	NOT NULL UNIQUE REFERENCES tests(id),
	data 	BLOB 		NOT NULL
);