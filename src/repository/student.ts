import { prisma } from "../db";
import type { Student } from "@prisma/client";

export const studentRepository = {
  async create(
    data: Omit<Student, "id" | "createdAt" | "updatedAt" | "deletedAt">,
  ): Promise<Student> {
    return prisma.student.create({ data });
  },

  async findAll(): Promise<Student[]> {
    return prisma.student.findMany({
      where: { deletedAt: null },
    });
  },

  async findById(id: number): Promise<Student | null> {
    return prisma.student.findFirst({
      where: { id, deletedAt: null },
    });
  },

  async update(
    id: number,
    data: Partial<
      Omit<Student, "id" | "createdAt" | "updatedAt" | "deletedAt">
    >,
  ): Promise<number> {
    const result = await prisma.student.updateMany({
      where: { id, deletedAt: null },
      data: { ...data, updatedAt: new Date() },
    });
    return result.count;
  },

  async delete(id: number): Promise<number> {
    const result = await prisma.student.updateMany({
      where: { id, deletedAt: null },
      data: { deletedAt: new Date() },
    });
    return result.count;
  },
};
