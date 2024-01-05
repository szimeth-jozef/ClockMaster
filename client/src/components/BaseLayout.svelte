<script lang="ts">
    import Navbar from './Navbar.svelte'
    import { Card, P, Heading, Button } from 'flowbite-svelte'
    import { runningWorkItem } from '../stores/RunningWorkItemStore'
    import toast from 'svelte-french-toast'

    export let hideRunningWorkItem: boolean = false

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
</script>

<header>
    <Navbar />
</header>
<main class="container mx-auto mt-4">
    <slot />
</main>
{#if $runningWorkItem.isRunning && !hideRunningWorkItem}
    <Card class="fixed bottom-4 left-4">
        <P class="text-center" weight="bold">{$runningWorkItem.workItem?.period.year}-{$runningWorkItem.workItem?.period.month}</P>
        <P class="text-center" weight="semibold">#{$runningWorkItem.workItem?.id}</P>
        <P class="text-center">{$runningWorkItem.workItem?.name}</P>
        <Heading level="3" class="text-center mb-4 min-w-64">
            {formatTimerSegment($runningWorkItem.timer?.hours || 0)}:{formatTimerSegment($runningWorkItem.timer?.minutes || 0)}:{formatTimerSegment($runningWorkItem.timer?.seconds || 0)}
        </Heading>
        <Button class="w-full" on:click={stopTimer}>Stop</Button>
    </Card>
{/if}
<footer></footer>

<style>
</style>