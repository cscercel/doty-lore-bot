import dotenv from "dotenv";
dotenv.config();


type Config = {
    db_url: string | undefined;
    db_token: string | undefined;
    discord_token: string | undefined,
    guild_id: string | undefined,
}

export const config: Config = {
    db_url: process.env.TURSO_DATABASE_URL,
    db_token: process.env.TURSO_AUTH_TOKEN,
    discord_token: process.env.DISCORD_TOKEN,
    guild_id: process.env.DISCORD_GUILD_ID,
}
