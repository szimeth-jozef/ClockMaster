import type { Period, Timer } from '../types/models.type'

const LOCALE = 'en-UK'

export const nanosecondsToTime = (nanoseconds: number) => {
    const inSeconds = nanoseconds / 1000000000
    const hours = Math.floor(inSeconds / 3600)
    const minutes = Math.floor((inSeconds - (hours * 3600)) / 60)
    const seconds = Math.floor(inSeconds - (hours * 3600) - (minutes * 60))

    return { hours, minutes, seconds }
}

export const formatTimer = (timer: Timer) => {
    const { hours, minutes, seconds } = timer
    return `${pad(hours)}:${pad(minutes)}:${pad(seconds)}`
}

export const formatDuration = (timer: Timer) => {
    const { hours, minutes, seconds } = timer
    return `${pad(hours)}h ${pad(minutes)}m ${pad(seconds)}s`
}

export const formatPeriod = (period: Period) => {
    const { month, year } = period
    return `${year}/${pad(month)}`
}

export const formatDateFromString = (date: string) => {
    return formatDate(new Date(date))
}

export const formatDateTimeFromString = (date: string) => {
    return formatDateTime(new Date(date))
}

export const formatDate = (date: Date) => {
    const options: Intl.DateTimeFormatOptions = { dateStyle: 'long' }
    return date.toLocaleDateString(LOCALE, options)
}

export const formatDateTime = (date: Date) => {
    const options: Intl.DateTimeFormatOptions = { dateStyle: 'long', timeStyle: 'medium' }
    return date.toLocaleString(LOCALE, options)
}

const pad = (n: number) => n.toString().padStart(2, '0')