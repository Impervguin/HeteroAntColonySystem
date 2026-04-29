import "reflect-metadata"
import { plainToInstance } from "class-transformer"
import { validateOrReject } from "class-validator"

export async function validateDto<T extends object>(
  cls: new () => T,
  data: unknown
): Promise<T> {
  const instance = plainToInstance(cls, data)
  await validateOrReject(instance)
  return instance
} 