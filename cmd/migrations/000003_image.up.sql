CREATE TABLE images(
		image_id UUID PRIMARY KEY UNIQUE DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
		post_id UUID REFERENCES posts(post_id) ON DELETE SET NULL,
		content BYTEA,
        thumbnail BYTEA
);