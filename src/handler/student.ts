import { Request, Response } from "express";
import { studentRepository } from "../repository/student";
import { CreateStudentSchema, UpdateStudentSchema } from "./student.schema";

export const studentHandler = {
  async create(req: Request, res: Response) {
    const parsed = CreateStudentSchema.safeParse(req.body);
    if (!parsed.success) {
      res.status(400).json({ error: parsed.error.errors });
      return;
    }

    try {
      const student = await studentRepository.create({
        ...parsed.data,
        rank: parsed.data.rank ?? null,
        addressLine2: parsed.data.addressLine2 ?? null,
      });
      res.status(201).json(student);
    } catch (error) {
      console.error("failed to create student", error);
      res.status(500).json({ error: "failed to create student" });
    }
  },

  async getAll(_req: Request, res: Response) {
    try {
      const students = await studentRepository.findAll();
      res.json(students);
    } catch (error) {
      console.error("failed to retrieve students", error);
      res.status(500).json({ error: "failed to retrieve students" });
    }
  },

  async getById(req: Request, res: Response) {
    const id = parseInt(req.params.id, 10);
    if (isNaN(id)) {
      res.status(400).json({ error: "invalid student ID" });
      return;
    }

    try {
      const student = await studentRepository.findById(id);
      if (!student) {
        res.status(404).json({ error: "student not found" });
        return;
      }
      res.json(student);
    } catch (error) {
      console.error("failed to get student by ID", { id, error });
      res.status(500).json({ error: "failed to get student" });
    }
  },

  async patch(req: Request, res: Response) {
    const id = parseInt(req.params.id, 10);
    if (isNaN(id)) {
      res.status(400).json({ error: "invalid student ID" });
      return;
    }

    const parsed = UpdateStudentSchema.safeParse(req.body);
    if (!parsed.success) {
      res.status(400).json({ error: parsed.error.errors });
      return;
    }

    if (Object.keys(parsed.data).length === 0) {
      res.status(400).json({ error: "no fields to update" });
      return;
    }

    try {
      const count = await studentRepository.update(id, parsed.data);
      if (count === 0) {
        res.status(404).json({ error: "student not found" });
        return;
      }
      res.json({ message: "student updated successfully" });
    } catch (error) {
      console.error("failed to patch student", { id, error });
      res.status(500).json({ error: "failed to update student" });
    }
  },

  async delete(req: Request, res: Response) {
    const id = parseInt(req.params.id, 10);
    if (isNaN(id)) {
      res.status(400).json({ error: "invalid student ID" });
      return;
    }

    try {
      const count = await studentRepository.delete(id);
      if (count === 0) {
        res.status(404).json({ error: "student not found" });
        return;
      }
      res.json({ message: "student deleted successfully" });
    } catch (error) {
      console.error("failed to delete student", { id, error });
      res.status(500).json({ error: "failed to delete student" });
    }
  },
};
