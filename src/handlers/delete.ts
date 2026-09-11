import { ChatInputCommandInteraction, MessageFlags } from "discord.js";
import { deleteCard, getCardByName } from "../db/queries/cards";


export async function handleDelete(interaction: ChatInputCommandInteraction) {
    const name = interaction.options.getString("name", true);

    let existing;
    try {
        existing = await getCardByName(name);
    } catch (err) {
        console.error("delete lookup error:", err);
        await interaction.reply({
            content: "Something went wrong looking card up.",
            flags: MessageFlags.Ephemeral,
        });
        return;
    }

    if (!existing) {
        await interaction.reply({
            content: `No card found named "${name}".`,
            flags: MessageFlags.Ephemeral,
        });
        return;
    }

    try {
        await deleteCard(existing.id);
    } catch (err) {
        console.error("delete card error:", err);
        await interaction.reply({
            content: "Something went wrong deleting card.",
            flags: MessageFlags.Ephemeral,
        });
        return;
    }

    await interaction.reply({
        content: `Deleted **${existing.name}**.`,
    });

    console.log(`card deleted name=${existing.name} by=${interaction.user.username}`);
}
