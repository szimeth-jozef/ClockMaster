<script lang="ts">
    import { onDestroy } from 'svelte';
    import BaseLayout from '../components/BaseLayout.svelte'
    import StatusBadge from '../components/StatusBadge.svelte'
    import MarkDoneModal from '../components/MarkDoneModal.svelte'
    import DeleteWorkItemModal from '../components/DeleteWorkItemModal.svelte'
    import { runningWorkItem } from '../stores/RunningWorkItemStore'
    import {
        Heading,
        Button,
        Table,
        TableHead,
        TableBody,
        TableHeadCell,
        TableBodyCell,
        TableBodyRow,
        Dropdown,
        DropdownItem,
        Card,
        ButtonGroup,
        Skeleton,
        P
    } from 'flowbite-svelte'
    import { DotsVerticalOutline } from 'flowbite-svelte-icons'
    import toast from 'svelte-french-toast'
    import { push } from 'svelte-spa-router'
    import type { WorkDay, WorkItem } from '../types/models.type'
    import { formatDateFromString, formatDateTimeFromString, formatDuration, formatTimer, nanosecondsToTime } from '../utils/datetime'
    import { getWorkDays, getWorkItem } from '../utils/requests'
    import { WorkItemStatus } from '../utils/workitemstatus'

    export let params: any

    const currentWorkItemID = parseInt(params.id)

    let isCurrentRunning = $runningWorkItem.workItem && $runningWorkItem.workItem.id === currentWorkItemID
    let deletePopupModal = false
    let makeDonePopupModal = false
    let workDays: WorkDay[] = []
    let workItem: WorkItem | null = null

    $: {
        workItem = workItem
        getWorkDays(currentWorkItemID)
            .then(data => {
                workDays = data
                console.log(data)
            })
    }

    getWorkItem(currentWorkItemID)
        .then(data => {
            workItem = data
            console.log(data)
        })

    const stopTimer = async () => {
        const resp = await fetch('http://localhost:8080/api/workitem/stop', {
            method: 'PATCH'
        })

        if (!resp.ok) {
            toast.error(`Failed to stop work item`, { duration: 3000 })
            return
        }

        getWorkDays(currentWorkItemID)
            .then(data => {
                workDays = data
                console.log(data)
            })

        getWorkItem(currentWorkItemID)
            .then(data => {
                workItem = data
                console.log(data)
            })

        const wiId = $runningWorkItem.workItem?.id
        runningWorkItem.stop()
        toast.success(`Stopped timer with ID: ${wiId || 'unknown'}`, { duration: 3000 })
    }

    const deleteWorkItem = async () => {
        const resp = await fetch(`http://localhost:8080/api/workitem/${currentWorkItemID}`, {
            method: 'DELETE'
        })

        if (!resp.ok) {
            toast.error(`Failed to delete work item`, { duration: 3000 })
            return
        }

        toast.success(`Deleted work item with ID: ${currentWorkItemID}`, { duration: 3000 })
        push('/')
    }

    const unsubscribeRunningWI = runningWorkItem.subscribe(value => {
        if (!value.isRunning) {
            isCurrentRunning = false
            return
        }

        isCurrentRunning = value.workItem?.id === currentWorkItemID
    })

    const markDoneAction = () => {
        if (isCurrentRunning) {
            toast.error(`Cannot mark as done while timer is running`, { duration: 3000 })
            return
        }

        makeDonePopupModal = true
    }

    onDestroy(() => {
        unsubscribeRunningWI()
    })
</script>

<BaseLayout hideRunningWorkItem={isCurrentRunning}>

    <Card class="max-w-none">
        {#if isCurrentRunning}
            <Heading level="1" class="text-center mb-6">
                {formatTimer($runningWorkItem.timer || { hours: 0, minutes: 0, seconds: 0 })}
            </Heading>
        {:else}
            <Heading level="1" class="text-center mb-6">
                -- : -- : --
            </Heading>
        {/if}

        {#if workItem}
            <div class="grid gap-2 grid-cols-1 justify-items-center">
                <Heading level="2" class="text-center mb-2" tag="h3">{workItem.name}</Heading>
                <P>#{workItem.id}</P>
                <StatusBadge statusCode={workItem.status} />
                <P>Total Time: {formatDuration(nanosecondsToTime(workItem.totalTimeNanoseconds))}</P>
            </div>
        {:else}
            <Skeleton class="max-w-none" />
        {/if}
    </Card>

    <ButtonGroup class="flex justify-end space-x-px mt-4">
        {#if isCurrentRunning}
            <Button color="primary" on:click={stopTimer}>Stop</Button>
        {:else}
            <Button color="primary">Start</Button>
        {/if}
        <Button color="primary" on:click={markDoneAction} disabled={!workItem || workItem.status === WorkItemStatus.Done || isCurrentRunning}>Mark As Done</Button>
        <Button color="red" on:click={() => (deletePopupModal = true)}>Delete</Button>
    </ButtonGroup>

    <Table shadow divClass="mt-4">
        <TableHead>
            <TableHeadCell>ID</TableHeadCell>
            <TableHeadCell>Created</TableHeadCell>
            <TableHeadCell>Last Started</TableHeadCell>
            <TableHeadCell>Total Time</TableHeadCell>
            <TableHeadCell>
                <span class="sr-only">Controls</span>
            </TableHeadCell>
        </TableHead>
        <TableBody>
            {#each workDays as workDay}
                <TableBodyRow>
                    <TableBodyCell>{workDay.ID}</TableBodyCell>
                    <TableBodyCell>{formatDateFromString(workDay.CreatedAt)}</TableBodyCell>
                    <TableBodyCell>{workDay.LastStartedAt === null ? '-' : formatDateTimeFromString(workDay.LastStartedAt)}</TableBodyCell>
                    <TableBodyCell>
                        {formatDuration(nanosecondsToTime(workDay.TotalDuration))}
                    </TableBodyCell>
                    <TableBodyCell>
                        <DotsVerticalOutline class="dots-menu text-primary-600 dark:text-primary-500 cursor-pointer ml-auto" />
                        <Dropdown triggeredBy=".dots-menu">
                            <DropdownItem>Edit</DropdownItem>
                            <DropdownItem class="text-red-500">Delete</DropdownItem>
                        </Dropdown>
                    </TableBodyCell>
                </TableBodyRow>
            {/each}
        </TableBody>
    </Table>

    <DeleteWorkItemModal bind:isOpen={deletePopupModal} workItemId={currentWorkItemID} afterDelete={() => push('/')} />

    <MarkDoneModal bind:isOpen={makeDonePopupModal} bind:workItem={workItem} />
</BaseLayout>

<style>
</style>