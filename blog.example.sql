CREATE TABLE posts (
	id TEXT PRIMARY KEY NOT NULL,
	title TEXT NOT NULL,
	created TEXT NOT NULL,
	edited TEXT NOT NULL,
	author TEXT NOT NULL,
	description TEXT NOT NULL,
	commentable INTEGER NOT NULL,
	visible INTEGER NOT NULL,
	pinToTop INTEGER NOT NULL
);

CREATE TABLE tags (
	post_id TEXT NOT NULL,
	tag TEXT NOT NULL,
	PRIMARY KEY (post_id, tag)
);
