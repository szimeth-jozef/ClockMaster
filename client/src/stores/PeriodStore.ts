import { writable } from "svelte/store"
import type { Period } from "../types/models.type"

const createPeriodStore = () => {
    const { subscribe, update } = writable<Period>({
        month: new Date().getMonth(),
        year: new Date().getFullYear()
    })

    const nextPeriod = () => {
        update(period => {
            if (period.month === 11) {
                return {
                    month: 0,
                    year: period.year + 1
                }
            }

            return {
                month: period.month + 1,
                year: period.year
            }
        })
    }

    const previousPeriod = () => {
        update(period => {
            if (period.month === 0) {
                return {
                    month: 11,
                    year: period.year - 1
                }
            }

            return {
                month: period.month - 1,
                year: period.year
            }
        })
    }

    return {
        subscribe,
        nextPeriod,
        previousPeriod
    }

}

export const period = createPeriodStore()