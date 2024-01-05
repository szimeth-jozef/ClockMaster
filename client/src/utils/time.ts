export const nanosecondsToTime = (nanoseconds: number) => {
    const inSeconds = nanoseconds / 1000000000
    const hours = Math.floor(inSeconds / 3600)
    const minutes = Math.floor((inSeconds - (hours * 3600)) / 60)
    const seconds = Math.floor(inSeconds - (hours * 3600) - (minutes * 60))

    return { hours, minutes, seconds }
}