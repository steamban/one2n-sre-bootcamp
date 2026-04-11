import dotenv from "dotenv";
import { z } from "zod";

dotenv.config();

const ConfigSchema = z.object({
  PORT: z.string().default("8080"),
  DATABASE_URL: z.string(),
});

const config = ConfigSchema.parse(process.env);

export const appConfig = {
  port: config.PORT,
  databaseUrl: config.DATABASE_URL,
};
