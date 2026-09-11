import { describe, expect, it } from "vitest";
import { isUniqueConstraintErr, toTitleCase } from "./utils";

describe("toTitleCase", () => {
    it("capitalizes a single lowercase word", () => {
        expect(toTitleCase("npc")).toBe("Npc");
    });

    it("capitalizes each word in a multi-word string", () => {
        expect(toTitleCase("location event")).toBe("Location Event");
    });

    it("normalizes mixed-case input", () => {
        expect(toTitleCase("nPc LoCaTion")).toBe("Npc Location");
    });

    it("handles already title-cased input", () => {
        expect(toTitleCase("Event")).toBe("Event");
    });

    it("returns an empty string for empty input", () => {
        expect(toTitleCase("")).toBe("");
    });

    it("collapses multiple spaces into empty-word segments without crashing", () => {
        expect(toTitleCase("npc  location")).toBe("Npc  Location");
    });
});

describe("isUniqueConstraintErr", () => {
    it("returns false for null/undefined", () => {
        expect(isUniqueConstraintErr(null)).toBe(false);
        expect(isUniqueConstraintErr(undefined)).toBe(false);
    });

    it("returns false for an unrelated Error", () => {
        expect(isUniqueConstraintErr(new Error("connection refused"))).toBe(false);
    });

    it("returns true for an Error containing 'UNIQUE constraint failed'", () => {
        expect(
            isUniqueConstraintErr(new Error("UNIQUE constraint failed: cards.name")),
        ).toBe(true);
    });

    it("returns true for an Error containing 'SQLITE_CONSTRAINT'", () => {
        expect(isUniqueConstraintErr(new Error("SQLITE_CONSTRAINT: ..."))).toBe(true);
    });

    it("returns true for a plain string containing the marker text", () => {
        expect(isUniqueConstraintErr("UNIQUE constraint failed: cards.name")).toBe(true);
    });

    it("returns false for a plain string without the marker text", () => {
        expect(isUniqueConstraintErr("some other error")).toBe(false);
    });

    it("returns false for non-error, non-string values", () => {
        expect(isUniqueConstraintErr(42)).toBe(false);
        expect(isUniqueConstraintErr({})).toBe(false);
    });
});
