import { z } from "zod";

export const CreateStudentSchema = z.object({
  firstName: z.string().max(50),
  lastName: z.string().max(50),
  age: z.number().int().positive().max(149),
  gender: z.enum(["Male", "Female", "Other"]),
  email: z.string().email().max(255),
  phone: z.string().max(15),
  class: z.enum(["10th", "11th", "12th"]),
  rank: z
    .string()
    .length(1)
    .regex(/^[A-F]$/)
    .optional(),
  addressLine1: z.string().max(100),
  addressLine2: z.string().max(100).optional(),
  city: z.string().max(50),
  state: z.string().max(50),
  pincode: z.string().length(6).regex(/^\d+$/),
});

export const UpdateStudentSchema = CreateStudentSchema.partial();

export type CreateStudentInput = z.infer<typeof CreateStudentSchema>;
export type UpdateStudentInput = z.infer<typeof UpdateStudentSchema>;
