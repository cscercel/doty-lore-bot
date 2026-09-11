import {
    AutocompleteInteraction,
    ChatInputCommandInteraction,
    EmbedBuilder,
    MessageFlags,
} from "discord.js";
import { searchCardsByName, getCardByName } from "../db/queries/cards";


export async function handleViewAutoComplete(interaction: AutocompleteInteraction) {
    const input = interaction.options.getFocused();
    const searchTerm = `${input}`;

    let cards;
    try {
        cards = await searchCardsByName(searchTerm);
    } catch (err) {
        console.error("autocomplete search error:", err);
        await interaction.respond([]);
        return;
    }

    const choices = cards.map((c) => ({ name: c.name, value: c.name }));

    await interaction.respond(choices);
}


export async function handleView(interaction: ChatInputCommandInteraction) {
    const name = interaction.options.getString("name", true);

    const card = await getCardByName(name);
    if (!card) {
        await interaction.reply({
            content: `No card found named "${name}".`,
            flags: MessageFlags.Ephemeral,
        });
        return;
    }

    const embed = new EmbedBuilder()
        .setTitle(card.name)
        .setDescription(card.summary)
        .setColor(0x8b5cf6)
        .addFields({ name: "Type", value: card.type, inline: true });

    if (card.body && card.body !== "") {
        embed.addFields({ name: "Details", value: card.body, inline: false });
    }

    if (card.imageUrl) {
        embed.setImage(card.imageUrl);
    }

    try {
        const dmChannel = await interaction.user.createDM();
        await dmChannel.send({ embeds: [embed] });
    } catch (err) {
        console.error("dm send error:", err);
        await interaction.reply({
            content: "Couldn't send you a DM — check your privacy settings allow DMs from server members.",
            flags: MessageFlags.Ephemeral,
        });
        return;
    }

    await interaction.reply({
        content: `Sent "${card.name}" to your DMs.`,
        flags: MessageFlags.Ephemeral,
    });

    console.log(`card viewed via dm: name=${card.name} by=${interaction.user.username}`);
}
