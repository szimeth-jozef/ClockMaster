<script lang="ts">
    import {
        Button,
        Modal,
        P,
        NumberInput
    } from 'flowbite-svelte'
    import { formatDuration, getSuggestedHours, nanosecondsToTime } from '../utils/datetime'
    import type { WorkItem } from '../types/models.type';
    import { getWorkDays, getWorkItem } from '../utils/requests'
    import toast from 'svelte-french-toast'

    export let isOpen: boolean = false
    export let workItem: WorkItem|null = null

    interface MarkDoneWorkDay {
        workDayId: number,
        createdAt: string
        totalDurationInNanoseconds: number,
        roundedDurationInHours: number,
    }

    let markDoneWorkDays: MarkDoneWorkDay[] = []

    const getMarkDoneWorkDays = async (workItemId: number) => {
        const workDays = await getWorkDays(workItemId)
        markDoneWorkDays = workDays.map(wd => ({
            workDayId: wd.ID,
            createdAt: wd.CreatedAt,
            totalDurationInNanoseconds: wd.TotalDuration,
            roundedDurationInHours: getSuggestedHours(wd.TotalDuration)
        }))
    }

    const markAsDone = async (workItemId: number) => {
        console.log(markDoneWorkDays)

        const payload = {
            workdays: markDoneWorkDays.map(wd => ({
                workday_id: wd.workDayId,
                rounded_duration_in_hours: wd.roundedDurationInHours
            }))
        }

        const resp = await fetch(`http://localhost:8080/api/workitem/${workItemId}/done`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(payload)
        })

        if (!resp.ok) {
            toast.error(`Failed to mark work item as done`, { duration: 3000 })
            return
        }

        toast.success(`Marked work item as done`, { duration: 3000 })

        // update ui
        workItem = await getWorkItem(workItemId)
    }

    $: if (workItem) {
        console.log("hello update")
        getMarkDoneWorkDays(workItem.id)
    }
</script>

<Modal title="Mark Work Item As Done" bind:open={isOpen} size="md" autoclose>
    <P class="mb-5">
        To mark as Done, you must first round up total time <span class="font-bold">{formatDuration(nanosecondsToTime(workItem?.totalTimeNanoseconds || 0))}</span> (suggested: {getSuggestedHours(workItem?.totalTimeNanoseconds || 0)}h)
    </P>
    {#each markDoneWorkDays as workDay}
        <div>
            {formatDuration(nanosecondsToTime(workDay.totalDurationInNanoseconds))}
            <NumberInput bind:value={workDay.roundedDurationInHours} />
        </div>
    {/each}
    <svelte:fragment slot="footer">
        <Button on:click={() => workItem && markAsDone(workItem.id)} class="me-2">Mark Done</Button>
        <Button color="alternative">Cancel</Button>
    </svelte:fragment>
</Modal>