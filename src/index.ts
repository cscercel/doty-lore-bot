import { Client, Events, GatewayIntentBits } from "discord.js";
import { config } from "./config";
import { registerHandlers } from "./session";
import { db } from "./db";


// Create a new client instance
const client = new Client({ intents: [GatewayIntentBits.Guilds] });
registerHandlers(client);

client.once(Events.ClientReady, (readyClient) => {
    console.log(`Ready! Logged in as ${readyClient.user.tag}`);
});

client.on(Events.Error, (err) => {
    console.error("client error:", err);
});

process.on("unhandledRejection", (err) => {
    console.error("unhandled rejection:", err);
});

async function shutdown() {
    console.log("Shutting down...");
    client.destroy();
    db.$client.close();
    process.exit(0);
}

process.on("SIGINT", shutdown);
process.on("SIGTERM", shutdown);

client.login(config.discord_token);
