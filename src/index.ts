import express from "express";
import { router } from "./router";
import { appConfig } from "./config";

const app = express();

app.use(express.json());

app.get("/healthcheck", (_req, res) => {
  res.json({ status: "UP" });
});

app.use("/api/v1", router);

app.listen(appConfig.port, () => {
  console.log(`Server running on port ${appConfig.port}`);
});
