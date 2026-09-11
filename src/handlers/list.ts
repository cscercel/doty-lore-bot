import { ChatInputCommandInteraction, EmbedBuilder } from "discord.js";
import { listCards, listCardsByType } from "../db/queries/cards";
import { toTitleCase } from "./utils";

export async function handleList(interaction: ChatInputCommandInteraction) {
    const typeFilter = interaction.options.getString("type") ?? "";

    let cards;
    let title = "All Lore Cards";

    try {
        if (typeFilter !== "") {
            cards = await listCardsByType(typeFilter);
            title = `${toTitleCase(typeFilter)} Cards`;
        } else {
            cards = await listCards();
        }
    } catch (err) {
        console.error("list cards error:", err);
        await interaction.reply({
            content: "Something went wrong listing cards.",
            ephemeral: true,
        });
        return;
    }

    if (cards.length === 0) {
        await interaction.reply({
            content: "No cards found.",
            ephemeral: true,
        });
        return;
    }

    const description = cards
        .map((c) => `**${c.name}** - _${c.type}_`)
        .join("\n");

    const embed = new EmbedBuilder()
        .setTitle(title)
        .setDescription(description)
        .setColor(0x8b5cf6);

    try {
        const dmChannel = await interaction.user.createDM();
        await dmChannel.send({ embeds: [embed] });
    } catch (err) {
        console.error("dm send error:", err);
        await interaction.reply({
            content: "Couldn't send you a DM — check your privacy settings allow DMs from server members.",
            ephemeral: true,
        });
        return;
    }

    await interaction.reply({
        content: "Sent the card list to your DMs.",
        ephemeral: true,
    });

    console.log(`cards listed via dm: type=${typeFilter || "none"} count=${cards.length} by=${interaction.user.username}`,
    );
}
