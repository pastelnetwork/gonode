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
	green_nft BOOLEAN NOT NULL,
	rareness_score FLOAT NOT NULL,
	royalty_rate_percentage FLOAT NOT NULL,
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
	asking_price FLOAT,
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
	art_id TEXT NOT NULL,
	pastel_id TEXT NOT NULL,
	PRIMARY KEY (art_id, pastel_id),
	FOREIGN KEY(art_id) REFERENCES art_metadata(art_id),
	FOREIGN KEY(pastel_id) REFERENCES user_metadata(artist_pastel_id)
);
---
CREATE TABLE IF NOT EXISTS art_auction (
	auction_id INTEGER PRIMARY KEY AUTOINCREMENT,
	instance_id TEXT NOT NULL,
	lowest_price FLOAT NOT NULL,
	is_open BOOLEAN NOT NULL,
	start_time TIMESTAMP NOT NULL DEFAULT (strftime('%s', 'now')),
	end_time TIMESTAMP,
	first_price FLOAT,
	second_price FLOAT,
	FOREIGN KEY(instance_id) REFERENCES art_instance_metadata(instance_id)
);
---
CREATE INDEX IF NOT EXISTS idx_auction_instance_id ON art_auction(instance_id);
---
CREATE TABLE IF NOT EXISTS art_bidding (
	bid_id INTEGER PRIMARY KEY AUTOINCREMENT,
	auction_id INTEGER NOT NULL,
	pastel_id TEXT NOT NULL,
	bid_price FLOAT NOT NULL,
	FOREIGN KEY(auction_id) REFERENCES art_auction(auction_id),
	FOREIGN KEY(pastel_id) REFERENCES user_metadata(artist_pastel_id)
);
---
CREATE INDEX IF NOT EXISTS idx_bidding_auction_id ON art_bidding(auction_id);
---
CREATE TRIGGER IF NOT EXISTS check_auction_open
	BEFORE INSERT ON art_bidding
WHEN (
	(SELECT is_open FROM art_auction WHERE auction_id = NEW.auction_id) = 0
)
BEGIN
	SELECT RAISE(FAIL, "the auction is closed");
END;
---
CREATE TRIGGER IF NOT EXISTS check_bid_price
	BEFORE INSERT ON art_bidding
WHEN (
	(SELECT lowest_price FROM art_auction WHERE auction_id = NEW.auction_id) > NEW.bid_price
)
BEGIN
	SELECT RAISE(FAIL, "bid price is lower than lowest price");
END;
---
CREATE TRIGGER IF NOT EXISTS update_first_price
	AFTER INSERT ON art_bidding
BEGIN
	UPDATE art_auction
	SET
		first_price = (SELECT bid_price FROM art_bidding WHERE auction_id = NEW.auction_id ORDER BY bid_price DESC LIMIT 1)
	WHERE auction_id = NEW.auction_id;
END;
---
CREATE TRIGGER IF NOT EXISTS update_second_price
	AFTER INSERT ON art_bidding
WHEN (
	(SELECT COUNT(*) FROM art_bidding WHERE auction_id = NEW.auction_id) > 1
)
BEGIN
	UPDATE art_auction
	SET
		second_price = (SELECT bid_price FROM art_bidding WHERE auction_id = NEW.auction_id ORDER BY bid_price DESC LIMIT 1 OFFSET 1)
	WHERE auction_id = NEW.auction_id;
END;
---
CREATE TABLE IF NOT EXISTS sn_activities (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	query TEXT NOT NULL,
	activity_type TEXT NOT NULL,
	cnt INTEGER NOT NULL
);
---
CREATE UNIQUE INDEX IF NOT EXISTS idx_query_activity_type ON sn_activities(query, activity_type);
---
CREATE INDEX IF NOT EXISTS idx_acitivities_cnt ON sn_activities(cnt);
