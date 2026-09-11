import { defineConfig } from 'drizzle-kit';
import { config } from "./src/config";


export default defineConfig({
    out: './drizzle',
    schema: './src/db/schema.ts',
    dialect: 'turso',
    dbCredentials: {
        url: config.db_url || "",
        authToken: config.db_token || "",
    },
});
