PRAGMA foreign_keys = ON;
---
CREATE TABLE IF NOT EXISTS user_metadata (
	artist_pastel_id TEXT PRIMARY KEY UNIQUE,
	real_name TEXT,
	facebook_link TEXT,
	twitter_link TEXT,
	native_currency TEXT,
	location TEXT,
	primary_language TEXT,
	categories TEXT,
	biography TEXT,
	timestamp INTEGER NOT NULL,
	signature TEXT NOT NULL,
	previous_block_hash TEXT NOT NULL,
	user_data_hash TEXT NOT NULL,
	avatar_image BLOB,
	avatar_filename TEXT,
	cover_photo_image BLOB,
	cover_photo_filename TEXT
);
---
CREATE TABLE IF NOT EXISTS art_metadata (
	art_id TEXT PRIMARY KEY UNIQUE,
	artist_pastel_id TEXT NOT NULL,
	copies INTEGER NOT NULL,
	created_timestamp INTEGER NOT NULL,
	FOREIGN KEY(artist_pastel_id) REFERENCES user_metadata(artist_pastel_id)
);
---
CREATE INDEX IF NOT EXISTS idx_art_artist_pastel_id ON art_metadata(artist_pastel_id);
---
CREATE TABLE IF NOT EXISTS art_instance_metadata (
	instance_id TEXT PRIMARY KEY UNIQUE,
	art_id TEXT NOT NULL,
	owner_pastel_id TEXT NOT NULL,
	price FLOAT NOT NULL,
	FOREIGN KEY(art_id) REFERENCES art_metadata(art_id)
	FOREIGN KEY(owner_pastel_id) REFERENCES user_metadata(artist_pastel_id)
);
---
CREATE INDEX IF NOT EXISTS idx_instance_art_id ON art_instance_metadata(art_id);
---
CREATE INDEX IF NOT EXISTS idx_instance_owner_id ON art_instance_metadata(owner_pastel_id);
---
CREATE TRIGGER IF NOT EXISTS num_of_instance_check
BEFORE INSERT ON art_instance_metadata
WHEN (SELECT COUNT(*) FROM art_instance_metadata WHERE art_id = NEW.art_id) >= (SELECT copies FROM art_metadata WHERE art_id = NEW.art_id)
BEGIN
	SELECT RAISE(FAIL, "maximum number of copies reached");
END;
---
CREATE TRIGGER IF NOT EXISTS owner_id_auto_fill
AFTER INSERT ON art_instance_metadata
WHEN (SELECT COUNT(*) FROM user_metadata WHERE artist_pastel_id = NEW.owner_pastel_id) = 0
BEGIN
	UPDATE art_instance_metadata
	SET owner_pastel_id = (SELECT artist_pastel_id FROM art_metadata WHERE art_id = NEW.art_id)
	WHERE instance_id = NEW.instance_id;
END;
---
CREATE TABLE IF NOT EXISTS transactions (
	transaction_id TEXT PRIMARY KEY UNIQUE,
	instance_id TEXT NOT NULL,
    timestamp INTEGER NOT NULL,
    seller_pastel_id TEXT NOT NULL,
    buyer_pastel_id TEXT NOT NULL,
	price FLOAT NOT NULL,
	FOREIGN KEY(instance_id) REFERENCES art_instance_metadata(instance_id)
	FOREIGN KEY(seller_pastel_id) REFERENCES user_metadata(artist_pastel_id)
	FOREIGN KEY(buyer_pastel_id) REFERENCES user_metadata(artist_pastel_id)
);
---
CREATE TRIGGER IF NOT EXISTS check_sale_by_owner
	BEFORE INSERT ON transactions
WHEN (
	(SELECT owner_pastel_id FROM art_instance_metadata WHERE instance_id = NEW.instance_id) <> NEW.seller_pastel_id
)
BEGIN
	SELECT RAISE(FAIL, "NFT instance is not sold by the owner");
END;
---
CREATE TRIGGER IF NOT EXISTS update_instance_table 
	AFTER INSERT ON transactions
BEGIN
	UPDATE art_instance_metadata
	SET
		owner_pastel_id = NEW.buyer_pastel_id,
		price = NEW.price
	WHERE instance_id = NEW.instance_id;
END;
---
CREATE TABLE IF NOT EXISTS follow_relations (
	follower_pastel_id TEXT NOT NULL,
	followee_pastel_id TEXT NOT NULL,
	PRIMARY KEY (follower_pastel_id, followee_pastel_id),
	FOREIGN KEY(follower_pastel_id) REFERENCES user_metadata(artist_pastel_id),
	FOREIGN KEY(followee_pastel_id) REFERENCES user_metadata(artist_pastel_id)
);
---
CREATE TRIGGER IF NOT EXISTS check_self_follow
	BEFORE INSERT ON follow_relations
WHEN (
	NEW.follower_pastel_id = NEW.followee_pastel_id
)
BEGIN
	SELECT RAISE(FAIL, "user cannot follow him/her self");
END;
---
CREATE TABLE IF NOT EXISTS art_like (
	art_id TEXT PRIMARY KEY,
	pastel_id TEXT NOT NULL,
	FOREIGN KEY(art_id) REFERENCES art_metadata(art_id),
	FOREIGN KEY(pastel_id) REFERENCES user_metadata(artist_pastel_id)
);
