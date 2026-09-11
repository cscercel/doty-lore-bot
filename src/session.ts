import { Client, Events, Interaction } from "discord.js";
import { handleCreateStart, handleCreateSubmit } from "./handlers/create";
import { handleDelete } from "./handlers/delete";
import { handleEditStart, handleEditSubmit } from "./handlers/edit";
import { handleList } from "./handlers/list";
import { handleViewAutoComplete, handleView } from "./handlers/view";


export function registerHandlers(client: Client) {
    client.on(Events.InteractionCreate, async (interaction: Interaction) => {
        // Slash commands
        if (interaction.isChatInputCommand()) {
            if (interaction.commandName !== "lore") return;

            const sub = interaction.options.getSubcommand(false);
            if (!sub) return;

            switch (sub) {
                case "create":
                    await handleCreateStart(interaction);
                    break;
                case "view":
                    await handleView(interaction);
                    break;
                case "list":
                    await handleList(interaction);
                    break;
                case "edit":
                    await handleEditStart(interaction);
                    break;
                case "delete":
                    await handleDelete(interaction);
                    break;
            }
            return;
        }

        // Autocomplete
        if (interaction.isAutocomplete()) {
            if (interaction.commandName !== "lore") return;

            const sub = interaction.options.getSubcommand(false);
            if (!sub) return;

            if (sub === "view" || sub === "edit" || sub === "delete") {
                await handleViewAutoComplete(interaction);
            }
            return;
        }

        // Modal submissions
        if (interaction.isModalSubmit()) {
            const customId = interaction.customId;

            if (customId.startsWith("lore_create_modal:")) {
                await handleCreateSubmit(interaction);
            } else if (customId.startsWith("lore_edit_modal:")) {
                await handleEditSubmit(interaction);
            }
            return;
        }
    });
}
