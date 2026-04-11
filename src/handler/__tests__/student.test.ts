import { describe, it, expect, beforeEach, vi } from "vitest";
import request from "supertest";
import express from "express";
import { studentHandler } from "../student";
import { router } from "../../router";

const app = express();
app.use(express.json());
app.use("/api/v1", router);

vi.mock("../../repository/student", () => ({
  studentRepository: {
    create: vi.fn(),
    findAll: vi.fn(),
    findById: vi.fn(),
    update: vi.fn(),
    delete: vi.fn(),
  },
}));

import { studentRepository } from "../../repository/student";

const mockStudent = {
  id: 1,
  firstName: "John",
  lastName: "Doe",
  age: 15,
  gender: "Male" as const,
  email: "john@example.com",
  phone: "1234567890",
  class: "10th" as const,
  rank: "A",
  addressLine1: "123 Main St",
  addressLine2: null,
  city: "NYC",
  state: "NY",
  pincode: "123456",
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
  deletedAt: null,
};

describe("Student Handler", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("POST /api/v1/students", () => {
    it("should create a student", async () => {
      vi.mocked(studentRepository.create).mockResolvedValue(mockStudent);

      const res = await request(app).post("/api/v1/students").send({
        firstName: "John",
        lastName: "Doe",
        age: 15,
        gender: "Male",
        email: "john@example.com",
        phone: "1234567890",
        class: "10th",
        rank: "A",
        addressLine1: "123 Main St",
        city: "NYC",
        state: "NY",
        pincode: "123456",
      });

      expect(res.status).toBe(201);
      expect(res.body.firstName).toBe("John");
    });

    it("should return 400 for invalid data", async () => {
      const res = await request(app)
        .post("/api/v1/students")
        .send({ firstName: "John" });

      expect(res.status).toBe(400);
    });
  });

  describe("GET /api/v1/students", () => {
    it("should return all students", async () => {
      vi.mocked(studentRepository.findAll).mockResolvedValue([mockStudent]);

      const res = await request(app).get("/api/v1/students");

      expect(res.status).toBe(200);
      expect(res.body).toHaveLength(1);
    });

    it("should return 500 on error", async () => {
      vi.mocked(studentRepository.findAll).mockRejectedValue(
        new Error("db error"),
      );

      const res = await request(app).get("/api/v1/students");

      expect(res.status).toBe(500);
    });
  });

  describe("GET /api/v1/students/:id", () => {
    it("should return a student", async () => {
      vi.mocked(studentRepository.findById).mockResolvedValue(mockStudent);

      const res = await request(app).get("/api/v1/students/1");

      expect(res.status).toBe(200);
      expect(res.body.firstName).toBe("John");
    });

    it("should return 400 for invalid id", async () => {
      const res = await request(app).get("/api/v1/students/abc");

      expect(res.status).toBe(400);
      expect(res.body).toEqual({ error: "invalid student ID" });
    });

    it("should return 404 when student not found", async () => {
      vi.mocked(studentRepository.findById).mockResolvedValue(null);

      const res = await request(app).get("/api/v1/students/999");

      expect(res.status).toBe(404);
    });
  });

  describe("PATCH /api/v1/students/:id", () => {
    it("should update a student", async () => {
      vi.mocked(studentRepository.update).mockResolvedValue(1);

      const res = await request(app)
        .patch("/api/v1/students/1")
        .send({ firstName: "Jane" });

      expect(res.status).toBe(200);
      expect(res.body).toEqual({ message: "student updated successfully" });
    });

    it("should return 400 for invalid id", async () => {
      const res = await request(app)
        .patch("/api/v1/students/abc")
        .send({ firstName: "Jane" });

      expect(res.status).toBe(400);
    });

    it("should return 400 for empty updates", async () => {
      const res = await request(app).patch("/api/v1/students/1").send({});

      expect(res.status).toBe(400);
      expect(res.body).toEqual({ error: "no fields to update" });
    });

    it("should return 404 when student not found", async () => {
      vi.mocked(studentRepository.update).mockResolvedValue(0);

      const res = await request(app)
        .patch("/api/v1/students/999")
        .send({ firstName: "Jane" });

      expect(res.status).toBe(404);
    });
  });

  describe("DELETE /api/v1/students/:id", () => {
    it("should delete a student", async () => {
      vi.mocked(studentRepository.delete).mockResolvedValue(1);

      const res = await request(app).delete("/api/v1/students/1");

      expect(res.status).toBe(200);
      expect(res.body).toEqual({ message: "student deleted successfully" });
    });

    it("should return 400 for invalid id", async () => {
      const res = await request(app).delete("/api/v1/students/abc");

      expect(res.status).toBe(400);
    });

    it("should return 404 when student not found", async () => {
      vi.mocked(studentRepository.delete).mockResolvedValue(0);

      const res = await request(app).delete("/api/v1/students/999");

      expect(res.status).toBe(404);
    });
  });
});
