import {
    ChatInputCommandInteraction,
    FileUploadBuilder,
    LabelBuilder,
    ModalBuilder,
    ModalSubmitInteraction,
    StringSelectMenuBuilder,
    TextInputBuilder,
    TextInputStyle,
} from "discord.js";
import { cardTypes } from "../commands";
import { NewCard } from "../db/schema";
import { createCard } from "../db/queries/cards";


export function isUniqueConstraintErr(err: unknown): boolean {
    if (!err) return false;
    const msg = err instanceof Error ? err.message : String(err);
    return msg.includes("UNIQUE constraint failed") || msg.includes("SQLITE_CONSTRAINT");
}


export async function handleCreateStart(interaction: ChatInputCommandInteraction) {
    const name = interaction.options.getString("name", true);

    const modal = new ModalBuilder()
        .setCustomId(`lore_create_modal:${name}`)
        .setTitle("New Lore Card");

    const typeSelect = new StringSelectMenuBuilder()
        .setCustomId("type")
        .setPlaceholder("Select card type...")
        .setRequired(true)
        .addOptions(...cardTypes);
    const typeLabel = new LabelBuilder()
        .setLabel("Type")
        .setStringSelectMenuComponent(typeSelect);

    const summaryInput = new TextInputBuilder()
        .setCustomId("summary")
        .setStyle(TextInputStyle.Short)
        .setRequired(true)
        .setMaxLength(200);
    const summaryLabel = new LabelBuilder()
        .setLabel("Summary")
        .setTextInputComponent(summaryInput);

    const bodyInput = new TextInputBuilder()
        .setCustomId("body")
        .setStyle(TextInputStyle.Paragraph)
        .setRequired(true)
        .setMaxLength(4000);
    const bodyLabel = new LabelBuilder()
        .setLabel("Body")
        .setTextInputComponent(bodyInput);

    const tagsInput = new TextInputBuilder()
        .setCustomId("tags")
        .setStyle(TextInputStyle.Short)
        .setRequired(false)
        .setMaxLength(200);
    const tagsLabel = new LabelBuilder()
        .setLabel("Tags (comma-separated)")
        .setTextInputComponent(tagsInput);

    const imageUpload = new FileUploadBuilder()
        .setCustomId("image")
        .setMinValues(0)
        .setMaxValues(1)
        .setRequired(false);
    const imageLabel = new LabelBuilder()
        .setLabel("Image")
        .setDescription("Optional: Upload an image for card")
        .setFileUploadComponent(imageUpload);

    modal.addLabelComponents(
        typeLabel,
        summaryLabel,
        bodyLabel,
        tagsLabel,
        imageLabel
    );

    await interaction.showModal(modal);
}


export async function handleCreateSubmit(interaction: ModalSubmitInteraction) {
    const parts = interaction.customId.split(":");
    if (parts.length !== 2) {
        await interaction.reply({
            content: "Something went wrong = please try creating the card again.",
            ephemeral: true,
        });
        return;
    }
    const [, name] = parts;

    const cardType = interaction.fields.getStringSelectValues("type")[0];
    const summary = interaction.fields.getTextInputValue("summary");

    const rawBody = interaction.fields.getTextInputValue("body");
    const body = rawBody !== "" ? rawBody : null;

    const rawTags = interaction.fields.getTextInputValue("tags");
    const tags = rawTags !== ""
        ? rawTags.split(",").map((t) => t.trim()).filter((t) => t !== "")
        : [];
    const tagsJSON = tags.length > 0 ? JSON.stringify(tags) : null;

    const uploadedFiles = interaction.fields.getUploadedFiles("image");
    const files = uploadedFiles ? [...uploadedFiles.values()] : [];

    // Check if file is an image (will be added to FileUploadBuilder in a future release)
    const invalid = files.some(f => !f.contentType?.startsWith("image/"));
    if (invalid) {
        return interaction.reply({
            content: "Please only upload image files.",
            ephemeral: true,
        });
    }

    const imageUrl = files.length > 0 ? files[0].url : null;

    const params: NewCard = {
        name: name,
        type: cardType,
        summary: summary,
        body: body,
        tags: tagsJSON,
        imageUrl: imageUrl,
    };

    let card;
    try {
        card = await createCard(params);
    } catch (err) {
        if (isUniqueConstraintErr(err)) {
            await interaction.reply({
                content: `A card named "${name}" already exists.`,
                ephemeral: true,
            });
            return;
        }
        console.error("create card error:", err);
        await interaction.reply({
            content: "Something went wrong creating the card.",
            ephemeral: true,
        });
        return;
    }

    console.log(`card created: id=${card.id} name=${card.name} type=${card.type} by=${interaction.user.username}`);

    await interaction.reply({
        content: `**${card.name}** created.`,
    });
}
