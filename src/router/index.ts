import { Router } from "express";
import { studentHandler } from "../handler/student";

export const router = Router();

router.post("/students", studentHandler.create);
router.get("/students", studentHandler.getAll);
router.get("/students/:id", studentHandler.getById);
router.patch("/students/:id", studentHandler.patch);
router.delete("/students/:id", studentHandler.delete);
