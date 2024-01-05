import { writable } from 'svelte/store'
import type { WorkItem, Timer } from '../types/models.type'
import { nanosecondsToTime } from '../utils/time'


export interface RunningWorkItemContext {
    isRunning: boolean,
    workItem?: WorkItem,
    timer?: Timer
}

const createRunningWorkItemStore = () => {
    const { subscribe, set, update } = writable<RunningWorkItemContext>({isRunning: false})
    let interval: number|null = null

    const start = (workItem: WorkItem, initNanosec: number = 0) => {
        const timer = (initNanosec <= 0) ?
            { hours: 0, minutes: 0, seconds: 0 } :
            nanosecondsToTime(initNanosec)

        const context: RunningWorkItemContext = {
            isRunning: true,
            workItem,
            timer
        }

        set(context)

        interval = setInterval(() => {
            update((ctx) => {
                if (!ctx.timer) {
                    return ctx
                }

                ctx.timer.seconds++

                if (ctx.timer.seconds === 60) {
                    ctx.timer.seconds = 0
                    ctx.timer.minutes++
                }

                if (ctx.timer.minutes === 60) {
                    ctx.timer.minutes = 0
                    ctx.timer.hours++
                }

                return ctx
            })
        }, 1000)
    }

    const stop = () => {
        if (interval !== null) {
            clearInterval(interval)
            interval = null
        }

        update(value => {
            value.isRunning = false
            return value
        })
    }

    return {
        subscribe,
        start,
        stop
    }
}

export const runningWorkItem = createRunningWorkItemStore()