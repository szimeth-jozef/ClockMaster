<script lang="ts">
    import BaseLayout from '../components/BaseLayout.svelte'
    import StatusBadge from '../components/StatusBadge.svelte'
    import type { Period, WorkItem, WorkItemResponse } from '../types/models.type'
    import { period } from '../stores/PeriodStore'
    import { runningWorkItem } from '../stores/RunningWorkItemStore'
    import toast from 'svelte-french-toast'
    import { formatDateFromString, formatDuration, formatPeriod, nanosecondsToTime } from '../utils/datetime'
    import {
        ArrowLeftOutline,
        ArrowRightOutline,
        ClipboardCheckOutline,
        PlayOutline,
        DotsVerticalOutline,
        EyeOutline,
        ReceiptOutline
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
        Select,
        Dropdown,
        DropdownItem,
        Card,
    } from 'flowbite-svelte'
    import { onDestroy } from 'svelte'
    import DeleteWorkItemModal from '../components/DeleteWorkItemModal.svelte';

    let workItems: WorkItem[] = []
    let currentPeriodTotalTime = "0h 0m 0s"
    let formModal = false
    const deleteModal = { isOpen: false, workItemId: -1 }
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

    const startWorkItem = async (workItem: WorkItem) => {
        const resp = await fetch(`http://localhost:8080/api/workitem/${workItem.id}/start`, {
            method: 'PATCH'
        })

        if (!resp.ok) {
            toast.error(`Failed to start work item with ID: ${workItem.id}`, { duration: 3000 })
            return
        }

        workItem.isRunning = true
        workItem.status = 1
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

        const durationHours = parseInt(String(formData.get('duration-hours')) || '0')
        const durationMinutes = parseInt(String(formData.get('duration-minutes')) || '0')

        const durationSeconds = (durationHours * 60 * 60) + (durationMinutes * 60)

        const payload = {
            name,
            period_month: month,
            period_year: year,
            init_total_duration_in_seconds: durationSeconds
        }

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
        })
    }

    const openDeleteModal = (workItem: WorkItem) => {
        deleteModal.isOpen = true
        deleteModal.workItemId = workItem.id
    }

    const unsubscribePeriod = period.subscribe(value => {
        getWorkItems(value).then(data => {
            const responsePeriod = data.period
            if (responsePeriod.month - 1 !== value.month || responsePeriod.year !== value.year) {
                console.error('Periods do not match')
            }

            workItems = data.workItems || []
            const totalPeriodTime = workItems.reduce((acc, curr) => acc + curr.totalTimeNanoseconds, 0)
            currentPeriodTotalTime = formatDuration(nanosecondsToTime(totalPeriodTime))
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
            const totalPeriodTime = workItems.reduce((acc, curr) => acc + curr.totalTimeNanoseconds, 0)
            currentPeriodTotalTime = formatDuration(nanosecondsToTime(totalPeriodTime))
        })
    })

    onDestroy(() => {
        unsubscribePeriod()
        unsubscribeRunningWI()
    })
</script>

<BaseLayout>
    <div class="flex justify-end">
        <Button size="md" class="mr-4" on:click={() => (formModal = true)}>
            New Work Item <ClipboardCheckOutline class="w-3.5 h-3.5 ms-2" />
        </Button>
        <Button size="md" on:click={() => alert('TODO')}>
            Create Invoce <ReceiptOutline class="w-3.5 h-3.5 ms-2" />
        </Button>
    </div>

    <div class="flex flex-row justify-evenly mt-8">
        <Card class="border-blue-500 dark:border-blue-500 text-blue-500 dark:text-blue-500 w-60">Total Time</Card>
        <Card class="border-green-500 dark:border-green-500 text-green-500 dark:text-green-500 w-60">Total Time Done</Card>
        <Card class="max-w-none w-60">Current DevOps Invoice ID</Card>
    </div>

    <div class="flex items-center justify-evenly mt-8">
        <Button size="sm" on:click={period.previousPeriod}>
            <ArrowLeftOutline class="w-3.5 h-3.5 me-2" /> Previous
        </Button>
        <div class="flex flex-col items-center w-32">
            <span class="dark:text-gray-400 font-bold">
                {getMonthName($period.month)} {$period.year}
            </span>
            <span class="dark:text-gray-400">{currentPeriodTotalTime}</span>
        </div>
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
                    <TableBodyCell>{formatDateFromString(workItem.created)}</TableBodyCell>
                    <TableBodyCell>{formatPeriod(workItem.period)}</TableBodyCell>
                    <TableBodyCell>
                        <StatusBadge statusCode={workItem.status} />
                    </TableBodyCell>
                    <TableBodyCell>{workItem.isInvoiced}</TableBodyCell>
                    <TableBodyCell>
                        {formatDuration(nanosecondsToTime(workItem.totalTimeNanoseconds))}
                    </TableBodyCell>
                    <TableBodyCell>
                        <div class="flex justify-between">
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
                            <a
                                href={`#/workitem/${workItem.id}`}
                                class="font-medium text-primary-600 dark:text-primary-500"
                                title="Edit"
                            >
                                <EyeOutline />
                            </a>
                            <DotsVerticalOutline id="dots-menu-{workItem.id}" class="text-primary-600 dark:text-primary-500 cursor-pointer" />
                            <Dropdown triggeredBy="#dots-menu-{workItem.id}">
                                <DropdownItem>Move to next</DropdownItem>
                                <DropdownItem>Move to previous</DropdownItem>
                                <DropdownItem>Mark as done</DropdownItem>
                                <DropdownItem>Edit</DropdownItem>
                                <DropdownItem class="text-red-500" on:click={() => openDeleteModal(workItem)}>
                                    Delete
                                </DropdownItem>
                            </Dropdown>
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
            <div>
                <Label>Initial total duration</Label>
                <hr class="my-2" />
                <div class="grid md:grid-cols-2 md:gap-6">
                    <div class="relative z-0 w-full group">
                        <Label class="space-y-2">
                            <span>Hours</span>
                            <Input type="number" name="duration-hours" placeholder="0" />
                        </Label>
                    </div>
                    <div class="relative z-0 w-full group">
                        <Label class="space-y-2">
                            <span>Minutes</span>
                            <Input type="number" name="duration-minutes" placeholder="0" />
                        </Label>
                    </div>
                </div>
            </div>
            <Button type="submit" class="w-full1">Create</Button>
        </form>
    </Modal>

    <DeleteWorkItemModal bind:isOpen={deleteModal.isOpen} bind:workItemId={deleteModal.workItemId} afterDelete={() => (workItems = workItems.filter((wi) => wi.id !== deleteModal.workItemId))} />
</BaseLayout>

<style>
</style>