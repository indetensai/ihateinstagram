CREATE TABLE image(
		image_id UUID PRIMARY KEY UNIQUE DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES app_user(user_id) ON DELETE CASCADE,
		post_id UUID REFERENCES post(post_id) ON DELETE SET NULL,
		content BYTEA,
        thumbnail BYTEA
);