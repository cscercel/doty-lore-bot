import { sqliteTable, uniqueIndex, index, integer, text, customType } from "drizzle-orm/sqlite-core"
import { sql } from "drizzle-orm"

export const gooseDbVersion = sqliteTable("goose_db_version", {
	id: integer().primaryKey({ autoIncrement: true }),
	versionId: integer("version_id").notNull(),
	isApplied: integer("is_applied").notNull(),
	tstamp: customType({ dataType: () => 'TIMESTAMP' })().default("datetime('now')"),
});

export const loreCards = sqliteTable("lore_cards", {
	id: integer().primaryKey({ autoIncrement: true }),
	name: text().notNull(),
	type: text().notNull(),
	summary: text().notNull(),
	body: text(),
	imageUrl: text("image_url"),
	tags: text(),
	createdAt: text("created_at").default(sql`datetime('now')`).notNull(),
	updatedAt: text("updated_at").default(sql`datetime('now')`).notNull(),
},
(table) => [uniqueIndex("idx_lore_cards_name_unique").on(table.name),
index("idx_lore_cards_type").on(table.type),
index("idx_lore_cards_name").on(table.name),
]);

