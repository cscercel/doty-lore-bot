import { Client, Events, GatewayIntentBits } from "discord.js";
import { config } from "./config";
import { registerHandlers } from "./session";


// Create a new client instance
const client = new Client({ intents: [GatewayIntentBits.Guilds] });
registerHandlers(client);

client.once(Events.ClientReady, (readyClient) => {
    console.log(`Ready! Logged in as ${readyClient.user.tag}`);
});

client.login(config.discord_token);
