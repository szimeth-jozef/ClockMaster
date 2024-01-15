<script lang="ts">
  import { onDestroy } from 'svelte';
    import BaseLayout from '../components/BaseLayout.svelte'
    import { runningWorkItem } from '../stores/RunningWorkItemStore'
    import {
        Heading,
        Button,
        Modal,
        Table,
        TableHead,
        TableBody,
        TableHeadCell,
        TableBodyCell,
        TableBodyRow
    } from 'flowbite-svelte'
    import { ExclamationCircleOutline } from 'flowbite-svelte-icons'
    import toast from 'svelte-french-toast'
    import { push } from 'svelte-spa-router'
    import type { WorkDay } from '../types/models.type'
    import { nanosecondsToTime } from '../utils/time'

    export let params: any

    const currentWorkItemID = parseInt(params.id)

    let isCurrentRunning = $runningWorkItem.workItem && $runningWorkItem.workItem.id === currentWorkItemID
    let popupModal = false
    let workDays: WorkDay[] = []

    const getWorkDays = async () => {
        const resp = await fetch(`http://localhost:8080/api/workitem/${currentWorkItemID}/workday`)
        const data = await resp.json() as WorkDay[]
        workDays = data
        console.log(data)
    }
    getWorkDays()

    const formatTimerSegment = (segment: number) => {
        return String(segment).padStart(2, '0')
    }

    const stopTimer = async () => {
        const resp = await fetch('http://localhost:8080/api/workitem/stop', {
            method: 'PATCH'
        })

        if (!resp.ok) {
            toast.error(`Failed to stop work item`, { duration: 3000 })
            return
        }

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

    const formatTime = (time: {hours: number, minutes: number, seconds: number}) => {
        return `${time.hours}h ${time.minutes}m ${time.seconds}s`
    }

    onDestroy(() => {
        unsubscribeRunningWI()
    })
</script>

<BaseLayout hideRunningWorkItem={isCurrentRunning}>
    {#if isCurrentRunning}
        <Heading level="1" class="text-center">
            {formatTimerSegment($runningWorkItem.timer?.hours || 0)}:{formatTimerSegment($runningWorkItem.timer?.minutes || 0)}:{formatTimerSegment($runningWorkItem.timer?.seconds || 0)}
        </Heading>
    {/if}


    {#if isCurrentRunning}
        <Button on:click={stopTimer}>Stop</Button>
    {:else}
        <Button>Start</Button>
    {/if}

    <Button on:click={() => (popupModal = true)} color="red">Delete</Button>
    <Button on:click={() => null} color="green">Mark As Done</Button>

    <Table shadow divClass="mt-4">
        <TableHead>
            <TableHeadCell>ID</TableHeadCell>
            <TableHeadCell>Created</TableHeadCell>
            <TableHeadCell>LastStarted</TableHeadCell>
            <TableHeadCell>Total time</TableHeadCell>
            <TableHeadCell>
                <span class="sr-only">Controls</span>
            </TableHeadCell>
        </TableHead>
        <TableBody>
            {#each workDays as workDay}
                <TableBodyRow>
                    <TableBodyCell>{workDay.ID}</TableBodyCell>
                    <TableBodyCell>{workDay.CreatedAt}</TableBodyCell>
                    <TableBodyCell>{workDay.LastStartedAt}</TableBodyCell>
                    <TableBodyCell>
                        {formatTime(nanosecondsToTime(workDay.TotalDuration))}
                    </TableBodyCell>
                    <TableBodyCell>
                        Controls
                    </TableBodyCell>
                </TableBodyRow>
            {/each}
        </TableBody>
    </Table>

    <Modal bind:open={popupModal} size="xs" autoclose>
        <div class="text-center">
          <ExclamationCircleOutline class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
          <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">Are you sure you want to delete this work item?</h3>
          <Button on:click={deleteWorkItem} color="red" class="me-2">Yes, I'm sure</Button>
          <Button color="alternative">No, cancel</Button>
        </div>
      </Modal>
</BaseLayout>

<style>
</style>