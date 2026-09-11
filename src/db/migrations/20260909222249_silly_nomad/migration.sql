-- Current sql file was generated after introspecting the database
-- If you want to run this migration please uncomment this code before executing migrations
/*
CREATE TABLE `goose_db_version` (
	`id` integer AUTOINCREMENT,
	`version_id` integer NOT NULL,
	`is_applied` integer NOT NULL,
	`tstamp` TIMESTAMP DEFAULT datetime('now'),
	CONSTRAINT `goose_db_version_pk` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE TABLE `lore_cards` (
	`id` integer AUTOINCREMENT,
	`name` text NOT NULL,
	`type` text NOT NULL,
	`summary` text NOT NULL,
	`body` text,
	`image_url` text,
	`tags` text,
	`created_at` text DEFAULT datetime('now') NOT NULL,
	`updated_at` text DEFAULT datetime('now') NOT NULL,
	CONSTRAINT `lore_cards_pk` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE UNIQUE INDEX `idx_lore_cards_name_unique` ON `lore_cards` (`name`);--> statement-breakpoint
CREATE INDEX `idx_lore_cards_type` ON `lore_cards` (`type`);--> statement-breakpoint
CREATE INDEX `idx_lore_cards_name` ON `lore_cards` (`name`);
*/