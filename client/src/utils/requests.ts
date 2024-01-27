import type { WorkDay, WorkItem } from "../types/models.type"

const BASE_URL = 'http://localhost:8080/api'
const ARTIFICIAL_DELAY = 0

const artificialDelay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

export const getWorkDays = async (workItemId: number): Promise<WorkDay[]> => {
    const resp = await fetch(`${BASE_URL}/workitem/${workItemId}/workday`)
    await artificialDelay(ARTIFICIAL_DELAY)

    if (resp.status !== 200) {
        return []
    }

    const data = await resp.json() as WorkDay[]
    return data
}

export const getWorkItem = async (id: number): Promise<WorkItem|null> => {
    const resp = await fetch(`${BASE_URL}/workitem/${id}`)
    await artificialDelay(ARTIFICIAL_DELAY)

    if (resp.status !== 200) {
        return null
    }

    const data = await resp.json() as WorkItem
    return data
}