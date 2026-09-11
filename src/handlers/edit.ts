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
import { getCardByName, updateCard } from "../db/queries/cards";
import { Card } from "../db/schema";
import { isUniqueConstraintErr } from "./create";


export async function handleEditStart(interaction: ChatInputCommandInteraction) {
    const name = interaction.options.getString("name", true);

    const existing = await getCardByName(name);
    if (!existing) {
        await interaction.reply({
            content: `No card found named "${name}".`,
            ephemeral: true,
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

    const imageOptionSelect = new StringSelectMenuBuilder()
        .setCustomId("image_option")
        .setPlaceholder("Select an option...")
        .setRequired(true)
        .addOptions(
            { label: "Keep existing image", value: "keep", default: true },
            { label: "Delete image", value: "delete" },
            { label: "Change image", value: "change" },
        );
    const imageOptionLabel = new LabelBuilder()
        .setLabel("image")
        .setStringSelectMenuComponent(imageOptionSelect);

    const imageUpload = new FileUploadBuilder()
        .setCustomId("image")
        .setMinValues(0)
        .setMaxValues(1)
        .setRequired(false);
    const imageUploadLabel = new LabelBuilder()
        .setLabel("New image (only used if \"Change image\" is selected)")
        .setFileUploadComponent(imageUpload);

    modal.addLabelComponents(
        nameLabel,
        typeLabel,
        summaryLabel,
        bodyLabel,
        tagsLabel,
        imageOptionLabel,
        imageUploadLabel,
    );

    await interaction.showModal(modal);
}

export async function handleEditSubmit(interaction: ModalSubmitInteraction) {
    const parts = interaction.customId.split(":");
    if (parts.length !== 2) {
        await interaction.reply({
            content: "Something went wrong - please try editing again.",
            ephemeral: true,
        });
        return;
    }

    const id = Number(parts[1]);
    if (Number.isNaN(id)) {
        await interaction.reply({
            content: "Something went wrong - please try editing again.",
            ephemeral: true,
        });
        return;
    }

    const newName = interaction.fields.getTextInputValue("name").trim();
    const type = interaction.fields.getTextInputValue("type");
    const summary = interaction.fields.getTextInputValue("summary");

    const rawBody = interaction.fields.getTextInputValue("body");
    const body = rawBody !== "" ? rawBody : null;

    const rawTags = interaction.fields.getTextInputValue("tags");
    const tags = rawTags !== ""
        ? rawTags.split(",").map((t) => t.trim()).filter((t) => t !== "")
        : [];
    const tagsJSON = tags.length > 0 ? JSON.stringify(tags) : null;

    const imageOption = interaction.fields.getStringSelectValues("image_option")[0];
    if (!["keep", "delete", "change"].includes(imageOption)) {
        await interaction.reply({
            content: "Image option must be keep, delete or change",
            ephemeral: true,
        });
        return;
    }

    let imageUrl: string | null | undefined;

    if (imageOption === "delete") {
        imageUrl = null;
    } else if (imageOption === "change") {
        const uploadedFiles = interaction.fields.getUploadedFiles("image");
        const files = uploadedFiles ? [...uploadedFiles.values()] : [];

        if (files.length === 0) {
            await interaction.reply({
                content: "You selected \"Change image\" but didn't upload a file — please try again.",
                ephemeral: true,
            });
            return;
        }

        const invalid = files.some((f) => !f.contentType?.startsWith("image/"));
        if (invalid) {
            await interaction.reply({
                content: "Please only upload image files.",
                ephemeral: true,
            });
            return;
        }

        imageUrl = files[0].url;
    }

    const conflict = await getCardByName(newName);
    if (conflict && conflict.id !== id) {
        await interaction.reply({
            content: `A card named "${newName} already exists - pick a different name`,
            ephemeral: true,
        });
        return;
    }

    const params: Partial<Omit<Card, "id" | "updatedAt">> = {
        name: newName,
        type: type,
        summary: summary,
        body: body,
        tags: tagsJSON,
        ...(imageUrl !== undefined ? { imageUrl } : {}),
    };

    let updated
    try {
        updated = await updateCard(id, params);
    } catch (err) {
        if (isUniqueConstraintErr(err)) {
            await interaction.reply({
                content: `A card named "${newName}" already exists.`,
                ephemeral: true,
            });
            return;
        }
        console.error("update card error:", err);
        await interaction.reply({
            content: "Something went wrong updating card.",
            ephemeral: true,
        });
        return;
    }

    console.log(`card edited: id=${updated.id} name=${updated.name} by=${interaction.user.username}`);

    await interaction.reply({
        content: `**${updated.name}** updated.`
    });
}
