import { eq, like, sql } from "drizzle-orm";
import { db } from "../index.js";
import { loreCards, NewCard, Card } from "../schema.js";


export async function createCard(card: NewCard) {
    const rows = await db
        .insert(loreCards)
        .values(card)
        .onConflictDoNothing()
        .returning();

    if (rows.length === 0) {
        throw new Error("Failed to create card");
    }

    return rows[0];
}

export async function getCardByName(name: string) {
    const rows = await db
        .select()
        .from(loreCards)
        .where(eq(loreCards.name, name))

    return rows.length > 0 ? rows[0] : null;
}

export async function listCards() {
    const rows = await db
        .select()
        .from(loreCards)
        .orderBy(loreCards.name)

    return rows;
}

export async function listCardsByType(typ: string) {
    const rows = await db
        .select()
        .from(loreCards)
        .where(eq(loreCards.type, typ))
        .orderBy(loreCards.name)

    return rows;
}

export async function updateCard(
    id: Card['id'],
    data: Partial<Omit<Card, 'id' | 'updatedAt'>>
) {
    const [result] = await db
        .update(loreCards)
        .set({
            ...data,
            updatedAt: sql`datetime('now')`,
        })
        .where(eq(loreCards.id, id))
        .returning();

    if (!result) {
        throw new Error(`Card with ID ${id} not found.`);
    }

    return result;
}

export async function searchCardsByName(searchTerm: string) {
    const rows = await db
        .select()
        .from(loreCards)
        .where(like(loreCards.name, `${searchTerm}`))
        .orderBy(loreCards.name)
        .limit(25);

    return rows;
}

export async function deleteCard(id: number) {
    await db.delete(loreCards).where(eq(loreCards.id, id));
}

export async function deleteCardByName(name: string) {
    await db.delete(loreCards).where(eq(loreCards.name, name));
}
