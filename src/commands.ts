import { config } from "./config";
import { REST, Routes, SlashCommandBuilder } from "discord.js";


export const cardTypes = [
    { name: "Location", label: "Location", value: "location" },
    { name: "NPC", label: "NPC", value: "npc" },
    { name: "Rule", label: "Rule", value: "rule" },
    { name: "Religion", label: "Religion", value: "religion" },
    { name: "History", label: "History", value: "history" },
    { name: "Event", label: "Event", value: "event" },
    { name: "Organization", label: "Organization", value: "organization" },
    { name: "Character", label: "Character", value: "character" },
    { name: "Artifact", label: "Artifact", value: "artifact" },
    { name: "Creature", label: "Creature", value: "creature" },
    { name: "Spell", label: "Spell", value: "spell" },
]


const data = new SlashCommandBuilder()
    .setName("lore")
    .setDescription("Manage D&D lore cards")
    .addSubcommand((subcommand) =>
        subcommand
            .setName("create")
            .setDescription("Create a new lore card")
            .addStringOption((option) =>
                option
                    .setName("name")
                    .setDescription("Card name")
                    .setRequired(true)
            )
            .addStringOption((option) =>
                option
                    .setName("type")
                    .setDescription("Card type")
                    .setRequired(true)
                    .addChoices(...cardTypes)
            )
    )
    .addSubcommand((subcommand) =>
        subcommand
            .setName("view")
            .setDescription("View a lore card")
            .addStringOption((option) =>
                option
                    .setName("name")
                    .setDescription("Card name")
                    .setRequired(true)
                    .setAutocomplete(true)
            )
    )
    .addSubcommand((subcommand) =>
        subcommand
            .setName("list")
            .setDescription("List lore cards")
            .addStringOption((option) =>
                option
                    .setName("type")
                    .setDescription("Filter by card type")
                    .setRequired(false)
                    .addChoices(...cardTypes)
            )
    )
    .addSubcommand((subcommand) =>
        subcommand
            .setName("edit")
            .setDescription("Edit an existing lore card")
            .addStringOption((option) =>
                option
                    .setName("name")
                    .setDescription("Card to edit")
                    .setRequired(true)
                    .setAutocomplete(true)
            )
    )
    .addSubcommand((subcommand) =>
        subcommand
            .setName("delete")
            .setDescription("Delete a lore card")
            .addStringOption((option) =>
                option
                    .setName("name")
                    .setDescription("Card to delete")
                    .setRequired(true)
                    .setAutocomplete(true)
            )
    );

// Registration logic
const rest = new REST({ version: "10" }).setToken(config.discord_token || '');

export async function registerCommands(clientId: string, guildId?: string) {
    try {
        const route = guildId
            ? Routes.applicationGuildCommands(clientId, guildId)
            : Routes.applicationCommands(clientId);

        await rest.put(route, { body: [data.toJSON()] });
    } catch (error) {
        console.error("Failed to register commands:", error);
    }
}
