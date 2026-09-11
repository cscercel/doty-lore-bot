import {
    ChatInputCommandInteraction,
    LabelBuilder,
    ModalBuilder,
    ModalSubmitInteraction,
    StringSelectMenuBuilder,
    TextInputBuilder,
    TextInputStyle,
    MessageFlags,
    ActionRowBuilder,
    FileUploadBuilder,
    StringSelectMenuInteraction,
} from "discord.js";
import { cardTypes } from "../commands";
import { getCardByName, updateCard } from "../db/queries/cards";
import { Card } from "../db/schema";
import { isUniqueConstraintErr } from "./utils";


export async function handleEditStart(interaction: ChatInputCommandInteraction) {
    const name = interaction.options.getString("name", true);

    const existing = await getCardByName(name);
    if (!existing) {
        await interaction.reply({
            content: `No card found named "${name}".`,
            flags: MessageFlags.Ephemeral,
        });
        return;
    }

    const modal = new ModalBuilder()
        .setCustomId(`lore_edit_modal:${existing.id}`)
        .setTitle(`Edit: ${existing.name}`);

    const nameInput = new TextInputBuilder()
        .setCustomId("name")
        .setStyle(TextInputStyle.Short)
        .setRequired(true)
        .setMaxLength(100)
        .setValue(existing.name)
    const nameLabel = new LabelBuilder()
        .setLabel("Name")
        .setTextInputComponent(nameInput);

    const typeSelect = new StringSelectMenuBuilder()
        .setCustomId("type")
        .setPlaceholder(existing.type)
        .setRequired(true)
        .addOptions(...cardTypes);
    const typeLabel = new LabelBuilder()
        .setLabel("Type")
        .setStringSelectMenuComponent(typeSelect);

    const summaryInput = new TextInputBuilder()
        .setCustomId("summary")
        .setStyle(TextInputStyle.Short)
        .setValue(existing.summary)
        .setRequired(true)
        .setMaxLength(200);
    const summaryLabel = new LabelBuilder()
        .setLabel("Summary")
        .setTextInputComponent(summaryInput);

    const bodyInput = new TextInputBuilder()
        .setCustomId("body")
        .setStyle(TextInputStyle.Paragraph)
        .setValue(existing.body || "")
        .setRequired(true)
        .setMaxLength(4000);
    const bodyLabel = new LabelBuilder()
        .setLabel("Body")
        .setTextInputComponent(bodyInput);

    const tagsInput = new TextInputBuilder()
        .setCustomId("tags")
        .setStyle(TextInputStyle.Short)
        .setValue(existing.tags || "")
        .setRequired(false)
        .setMaxLength(200);
    const tagsLabel = new LabelBuilder()
        .setLabel("Tags (comma-separated)")
        .setTextInputComponent(tagsInput);

    modal.addLabelComponents(
        nameLabel,
        typeLabel,
        summaryLabel,
        bodyLabel,
        tagsLabel,
    );

    await interaction.showModal(modal);
}

export async function handleEditSubmit(interaction: ModalSubmitInteraction) {
    const parts = interaction.customId.split(":");
    if (parts.length !== 2) {
        await interaction.reply({
            content: "Something went wrong - please try editing again.",
            flags: MessageFlags.Ephemeral,
        });
        return;
    }

    const id = Number(parts[1]);
    if (Number.isNaN(id)) {
        await interaction.reply({
            content: "Something went wrong - please try editing again.",
            flags: MessageFlags.Ephemeral,
        });
        return;
    }

    const newName = interaction.fields.getTextInputValue("name").trim();
    const type = interaction.fields.getStringSelectValues("type")[0];
    const summary = interaction.fields.getTextInputValue("summary");

    const rawBody = interaction.fields.getTextInputValue("body");
    const body = rawBody !== "" ? rawBody : null;

    const rawTags = interaction.fields.getTextInputValue("tags");
    const tags = rawTags !== ""
        ? rawTags.split(",").map((t) => t.trim()).filter((t) => t !== "")
        : [];
    const tagsJSON = tags.length > 0 ? JSON.stringify(tags) : null;

    const conflict = await getCardByName(newName);
    if (conflict && conflict.id !== id) {
        await interaction.reply({
            content: `A card named "${newName} already exists - pick a different name`,
            flags: MessageFlags.Ephemeral,
        });
        return;
    }

    const params: Partial<Omit<Card, "id" | "updatedAt">> = {
        name: newName,
        type: type,
        summary: summary,
        body: body,
        tags: tagsJSON,
    };

    let updated
    try {
        updated = await updateCard(id, params);
    } catch (err) {
        if (isUniqueConstraintErr(err)) {
            await interaction.reply({
                content: `A card named "${newName}" already exists.`,
                flags: MessageFlags.Ephemeral,
            });
            return;
        }
        console.error("update card error:", err);
        await interaction.reply({
            content: "Something went wrong updating card.",
            flags: MessageFlags.Ephemeral,
        });
        return;
    }

    console.log(`card edited: id=${updated.id} name=${updated.name} by=${interaction.user.username}`);

    await interaction.reply({
        content: `**${updated.name}** updated. Want to change the image too?`,
        flags: MessageFlags.Ephemeral,
        components: [
            new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(
                new StringSelectMenuBuilder()
                    .setCustomId(`lore_edit_image:${updated.id}`)
                    .setPlaceholder("Image options...")
                    .addOptions(
                        { label: "Keep existing image", value: "keep" },
                        { label: "Delete image", value: "delete" },
                        { label: "Change image", value: "change" },
                    )
            ),
        ],
    });
}


export async function handleEditImageOption(interaction: StringSelectMenuInteraction) {
    const parts = interaction.customId.split(":");
    if (parts.length !== 2) {
        await interaction.update({
            content: "Something went wrong - please try editing again.",
            components: [],
        });
        return;
    }

    const id = Number(parts[1]);
    if (Number.isNaN(id)) {
        await interaction.update({
            content: "Something went wrong - please try editing again.",
            components: [],
        });
        return;
    }

    const choice = interaction.values[0];

    if (choice === "keep") {
        await interaction.update({
            content: "Image unchanged.",
            components: [],
        });
        return;
    }

    if (choice === "delete") {
        try {
            await updateCard(id, { imageUrl: null });
        } catch (err) {
            console.error("delete image error:", err);
            await interaction.update({
                content: "Something went wrong removing the image.",
                components: [],
            });
            return;
        }
        await interaction.update({
            content: "Image removed.",
            components: [],
        });
        return;
    }

    // choice === "change" -> show a follow-up modal with the file upload
    const modal = new ModalBuilder()
        .setCustomId(`lore_edit_image_modal:${id}`)
        .setTitle("Change image");

    const imageUpload = new FileUploadBuilder()
        .setCustomId("image")
        .setMinValues(1)
        .setMaxValues(1)
        .setRequired(true);
    const imageUploadLabel = new LabelBuilder()
        .setLabel("New image")
        .setFileUploadComponent(imageUpload);

    modal.addLabelComponents(imageUploadLabel);

    await interaction.showModal(modal);
}

export async function handleEditImageModalSubmit(interaction: ModalSubmitInteraction) {
    const parts = interaction.customId.split(":");
    if (parts.length !== 2) {
        await interaction.reply({
            content: "Something went wrong - please try again.",
            flags: MessageFlags.Ephemeral,
        });
        return;
    }

    const id = Number(parts[1]);
    if (Number.isNaN(id)) {
        await interaction.reply({
            content: "Something went wrong - please try again.",
            flags: MessageFlags.Ephemeral,
        });
        return;
    }

    const uploadedFiles = interaction.fields.getUploadedFiles("image");
    const files = uploadedFiles ? [...uploadedFiles.values()] : [];

    if (files.length === 0) {
        await interaction.reply({
            content: "No file was uploaded — please try again.",
            flags: MessageFlags.Ephemeral,
        });
        return;
    }

    const invalid = files.some((f) => !f.contentType?.startsWith("image/"));
    if (invalid) {
        await interaction.reply({
            content: "Please only upload image files.",
            flags: MessageFlags.Ephemeral,
        });
        return;
    }

    try {
        await updateCard(id, { imageUrl: files[0].url });
    } catch (err) {
        console.error("update image error:", err);
        await interaction.reply({
            content: "Something went wrong updating the image.",
            flags: MessageFlags.Ephemeral,
        });
        return;
    }

    await interaction.reply({
        content: "Image updated.",
        flags: MessageFlags.Ephemeral,
    });
}
