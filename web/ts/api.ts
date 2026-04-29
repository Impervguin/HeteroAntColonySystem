import { validateDto } from "./utils/validate.js"

import { ParseTSPResponse, GetTSPResponse } from "./models/tsp.js"
import { HacoRunRequest } from "./models/haco/request.js"
import { HacoRunDetailsResponse } from "./models/haco/response.js"

// =======================

async function safeJson(res: Response) {
  if (!res.ok) {
    const text = await res.text()
    throw new Error(`HTTP ${res.status}: ${text}`)
  }
  return res.json()
}

function getApiBase(): string {
  return window.APP_CONFIG.apiBase ?? "http://localhost:8080/api/v1"
}

export async function parseTSP(file: File): Promise<ParseTSPResponse> {
  const formData = new FormData()
  formData.append("file", file)

  const res = await fetch(`${getApiBase()}/tsp/parse`, {
    method: "POST",
    body: formData
  })

  return validateDto(ParseTSPResponse, await safeJson(res))
}

export async function getTSP(filename: string): Promise<GetTSPResponse> {
  const res = await fetch(`${getApiBase()}/tsp/${encodeURIComponent(filename)}`)
  return validateDto(GetTSPResponse, await safeJson(res))
}

export async function runHacoDetails(request: HacoRunRequest): Promise<HacoRunDetailsResponse> {
  const res = await fetch(`${getApiBase()}/haco/run/details`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(request)
  })

  return validateDto(HacoRunDetailsResponse, await safeJson(res))
}