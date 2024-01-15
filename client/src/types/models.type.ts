export interface WorkItem {
    id: number,
    created: string,
    name: string,
    status: number,
    period: Period,
    isInvoiced: boolean,
    totalTimeNanoseconds: number,
    isRunning: boolean,
}

export interface Period {
    month: number,
    year: number
}

export interface Timer {
    hours: number,
    minutes: number,
    seconds: number
}

export interface WorkItemResponse {
    period: Period,
    workItems: WorkItem[]
}

export interface StatusResponse {
    isRunning: boolean,
    deltaDurationNanoseconds: number,
    workItem: WorkItem
}

export interface WorkDay {
    ID: number,
    CreatedAt: string,
    UpdatedAt: string,
    WorkItemID: number,
    LastStartedAt: string,
    TotalDuration: number
}