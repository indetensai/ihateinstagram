CREATE TYPE visibility AS ENUM ('public', 'followers', 'private');
CREATE TABLE IF NOT EXISTS app_user(
		user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(64) NOT NULL
);
CREATE TABLE IF NOT EXISTS session(
		session_id VARCHAR(64) PRIMARY KEY UNIQUE,
		user_id UUID REFERENCES app_user(user_id) ON DELETE CASCADE,
		created_at TIMESTAMP NOT NULL DEFAULT now()	
);
CREATE TABLE IF NOT EXISTS post(
		post_id UUID PRIMARY KEY UNIQUE DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES app_user(user_id) ON DELETE CASCADE,
		description TEXT,
		vision visibility DEFAULT 'private',
		created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS following(
		follower_id UUID PRIMARY KEY NOT NULL REFERENCES app_user(user_id) ON DELETE CASCADE,
		user_id UUID REFERENCES app_user(user_id) ON DELETE CASCADE,
		following_since TIMESTAMP NOT NULL DEFAULT now(),
		UNIQUE (user_id,follower_id),
		CHECK (user_id!=follower_id)
);
CREATE TABLE IF NOT EXISTS liking(
		user_id UUID NOT NULL REFERENCES app_user(user_id) ON DELETE CASCADE,
		post_id UUID NOT NULL REFERENCES post(post_id) ON DELETE CASCADE,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		PRIMARY KEY (user_id,post_id),
		UNIQUE (user_id,post_id)
);