export function isUniqueConstraintErr(err: unknown): boolean {
    if (!err) return false;
    const msg = err instanceof Error ? err.message : String(err);
    return msg.includes("UNIQUE constraint failed") || msg.includes("SQLITE_CONSTRAINT");
}

export function toTitleCase(str: string): string {
    return str
        .toLowerCase()
        .split(" ")
        .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
        .join(" ");
}
