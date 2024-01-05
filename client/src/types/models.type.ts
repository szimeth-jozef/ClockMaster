export interface WorkItem {
    id: number,
    created: string,
    name: string,
    status: number,
    period: string,
    isInvoiced: boolean,
    totalTimeNanoseconds: number,
    isRunning: boolean,
}

export interface Period {
    month: number,
    year: number
}

export interface WorkItemResponse {
    period: Period,
    workItems: WorkItem[]
}