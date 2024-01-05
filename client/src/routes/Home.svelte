<script lang="ts">
    import BaseLayout from '../components/BaseLayout.svelte'
    import type { Period, WorkItem, WorkItemResponse } from '../types/models.type'
    import { period } from '../stores/PeriodStore'
    import { runningWorkItem } from '../stores/RunningWorkItemStore'
    import toast from 'svelte-french-toast'
    import { nanosecondsToTime } from '../utils/time'
    import {
        ArrowLeftOutline,
        ArrowRightOutline,
        ClipboardCheckOutline,
        PenOutline,
        PlayOutline,
    } from 'flowbite-svelte-icons'
    import {
        Button,
        Table,
        TableBody,
        TableBodyCell,
        TableBodyRow,
        TableHead,
        TableHeadCell,
        Modal,
        Label,
        Input,
        Select
    } from 'flowbite-svelte'
    import { onDestroy } from 'svelte'

    let workItems: WorkItem[] = []
    let formModal = false
    let createFormSelectedMonth: number = new Date().getMonth() + 1
    let createFormSelectedYear: number = new Date().getFullYear()

    const months = [
        { name: 'January', value: 1 },
        { name: 'February', value: 2 },
        { name: 'March', value: 3 },
        { name: 'April', value: 4 },
        { name: 'May', value: 5 },
        { name: 'June', value: 6 },
        { name: 'July', value: 7 },
        { name: 'August', value: 8 },
        { name: 'September', value: 9 },
        { name: 'October', value: 10 },
        { name: 'November', value: 11 },
        { name: 'December', value: 12 }
    ]
    const years = [
        { name: "2023", value: 2023 },
        { name: "2024", value: 2024 },
        { name: "2025", value: 2025 },
        { name: "2026", value: 2026 },
    ]

    const getWorkItems = async (period: Period): Promise<WorkItemResponse> => {
        const url = new URL('http://localhost:8080/api/workitem')

        const month = period.month + 1
        const year = period.year

        url.searchParams.append('month', month.toString())
        url.searchParams.append('year', year.toString())
        const resp = await fetch(url.toString())
        const data = await resp.json() as WorkItemResponse
        return data
    }

    const getMonthName = (monthNumber: number) => {
        const date = new Date()
        date.setMonth(monthNumber)
        return date.toLocaleString('default', { month: 'long' })
    }

    const workItemStatusToString = (status: number) => {
        switch (status) {
            case 0:
                return 'ToDo'
            case 1:
                return 'InProgress'
            case 2:
                return 'Done'
            default:
                return 'Unknown'
        }
    }

    const formatTime = (time: {hours: number, minutes: number, seconds: number}) => {
        return `${time.hours}h ${time.minutes}m ${time.seconds}s`
    }

    const startWorkItem = async (workItem: WorkItem) => {
        const resp = await fetch(`http://localhost:8080/api/workitem/${workItem.id}/start`, {
            method: 'PATCH'
        })

        if (!resp.ok) {
            toast.error(`Failed to start work item with ID: ${workItem.id}`, { duration: 3000 })
            return
        }

        workItem.isRunning = true
        workItems = [...workItems]
        runningWorkItem.start(workItem)
        toast.success(`Started work item with ID: ${workItem.id}`, { duration: 3000 })
    }

    const stopWorkItem = async (workItem: WorkItem) => {
        const resp = await fetch('http://localhost:8080/api/workitem/stop', {
            method: 'PATCH'
        })

        if (!resp.ok) {
            toast.error(`Failed to stop work item`, { duration: 3000 })
            return
        }

        runningWorkItem.stop()
        toast.success(`Stopped work item with ID: ${workItem.id}`, { duration: 3000 })
    }

    const createWorkItem = async (event: SubmitEvent) => {
        event.preventDefault()
        const formData = new FormData(event.target as HTMLFormElement)

        const name = formData.get('name')
        const month = parseInt(String(formData.get('month-period')))
        const year = parseInt(String(formData.get('year-period')))
        const durationSeconds = parseInt(String(formData.get('duration-minutes')) || '0') * 60

        const payload = {
            name,
            period_month: month,
            period_year: year,
            init_total_duration_in_seconds: durationSeconds
        }

        console.log(payload)

        const resp = await fetch('http://localhost:8080/api/workitem', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        })

        const data = await resp.json() as WorkItem

        formModal = false
        toast.success(`Work item created with ID: ${data.id}`, { duration: 3000 })

        if (!($period.year === year || $period.month === month)) {
            return
        }

        const p = $period
        getWorkItems(p).then(data => {
            const responsePeriod = data.period
            if (responsePeriod.month - 1 !== p.month || responsePeriod.year !== p.year) {
                console.error('Periods do not match')
            }

            workItems = data.workItems || []
            console.log(data)
        })
    }

    const unsubscribePeriod = period.subscribe(value => {
        console.log(value)

        getWorkItems(value).then(data => {
            const responsePeriod = data.period
            if (responsePeriod.month - 1 !== value.month || responsePeriod.year !== value.year) {
                console.error('Periods do not match')
            }

            workItems = data.workItems || []
            console.log(data)
        })
    })

    const unsubscribeRunningWI = runningWorkItem.subscribe(value => {
        if (value.isRunning) {
            return
        }

        const year = value.workItem?.period.year || 0
        const month = value.workItem?.period.month || 0

        const p = $period
        if (!(year === p.year && month - 1 === p.month)) {
            return
        }

        getWorkItems(p).then(data => {
            const responsePeriod = data.period
            if (responsePeriod.month - 1 !== p.month || responsePeriod.year !== p.year) {
                console.error('Periods do not match')
            }

            workItems = data.workItems || []
            console.log(data)
        })
    })

    onDestroy(() => {
        unsubscribePeriod()
        unsubscribeRunningWI()
    })
</script>

<BaseLayout>
    <div class="flex justify-end">
        <Button size="md" on:click={() => (formModal = true)}>
            New Work Item <ClipboardCheckOutline class="w-3.5 h-3.5 ms-2" />
        </Button>
    </div>

    <div class="flex items-center justify-evenly mt-8">
        <Button size="sm" on:click={period.previousPeriod}>
            <ArrowLeftOutline class="w-3.5 h-3.5 me-2" /> Previous
        </Button>
        <span class="dark:text-gray-400 font-bold w-32">
            {getMonthName($period.month)} {$period.year}
        </span>
        <Button size="sm" on:click={period.nextPeriod}>
            Next <ArrowRightOutline class="w-3.5 h-3.5 ms-2" />
        </Button>
    </div>

    <Table shadow divClass="mt-4">
        <TableHead>
            <TableHeadCell>ID</TableHeadCell>
            <TableHeadCell>Name</TableHeadCell>
            <TableHeadCell>Created</TableHeadCell>
            <TableHeadCell>Period</TableHeadCell>
            <TableHeadCell>Status</TableHeadCell>
            <TableHeadCell>Invoiced</TableHeadCell>
            <TableHeadCell>Total time</TableHeadCell>
            <TableHeadCell>
                <span class="sr-only">Controls</span>
            </TableHeadCell>
        </TableHead>
        <TableBody>
            {#each workItems as workItem}
                <TableBodyRow>
                    <TableBodyCell>{workItem.id}</TableBodyCell>
                    <TableBodyCell>{workItem.name}</TableBodyCell>
                    <TableBodyCell>{workItem.created}</TableBodyCell>
                    <TableBodyCell>{workItem.period.year}-{workItem.period.month}</TableBodyCell>
                    <TableBodyCell>{workItemStatusToString(workItem.status)}</TableBodyCell>
                    <TableBodyCell>{workItem.isInvoiced}</TableBodyCell>
                    <TableBodyCell>
                        {formatTime(nanosecondsToTime(workItem.totalTimeNanoseconds))}
                    </TableBodyCell>
                    <TableBodyCell>
                        <div class="flex justify-between">
                            <a
                                href={`#/workitem/${workItem.id}`}
                                class="font-medium text-primary-600 dark:text-primary-500"
                                title="Edit"
                            >
                                <PenOutline />
                            </a>
                            {#if workItem.isRunning}
                                <div
                                    on:click={() => stopWorkItem(workItem)}
                                    on:keypress={() => stop}
                                    class="w-5 h-5 bg-primary-600 dark:bg-primary-500 rounded-sm"
                                    role="button"
                                    tabindex="0"
                                    title="Stop">
                                </div>
                            {:else}
                                <div
                                    on:click={() => startWorkItem(workItem)}
                                    on:keypress={() => startWorkItem(workItem)}
                                    class="font-medium text-primary-600 dark:text-primary-500"
                                    role="button"
                                    tabindex="0"
                                    title="Start"
                                >
                                    <PlayOutline />
                                </div>
                            {/if}
                        </div>
                    </TableBodyCell>
                </TableBodyRow>
            {/each}
        </TableBody>
    </Table>

    <Modal bind:open={formModal} size="xs" autoclose={false} class="w-full">
        <form class="flex flex-col space-y-6" on:submit={createWorkItem}>
            <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">Create new work item</h3>
            <Label class="space-y-2">
                <span>Work item name</span>
                <Input type="text" name="name" placeholder="Fix bug..." required />
            </Label>
            <div class="grid md:grid-cols-2 md:gap-6">
                <div class="relative z-0 w-full group">
                    <Label class="space-y-2">
                        <span>Month of period</span>
                        <Select items={months} name="month-period" bind:value={createFormSelectedMonth} placeholder="" />
                    </Label>
                </div>
                <div class="relative z-0 w-full group">
                    <Label class="space-y-2">
                        <span>Year of period</span>
                        <Select items={years} name="year-period" bind:value={createFormSelectedYear} placeholder="" />
                    </Label>
                </div>
            </div>
            <Label class="space-y-2">
                <span>Initial total duration in minutes</span>
                <Input type="number" name="duration-minutes" placeholder="0" />
            </Label>
            <Button type="submit" class="w-full1">Create</Button>
        </form>
    </Modal>
</BaseLayout>

<style>
</style>