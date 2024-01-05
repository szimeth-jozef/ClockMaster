<script lang="ts">
    import Home from './routes/Home.svelte'
    import WorkItem from './routes/WorkItem.svelte'
    import Invoices from './routes/Invoices.svelte'
    import Router from 'svelte-spa-router'
    import { Toaster } from 'svelte-french-toast'
    import {runningWorkItem} from './stores/RunningWorkItemStore'
    import type { StatusResponse } from './types/models.type'


    const routes = {
        '/': Home,
        '/workitem/:id': WorkItem,
        '/invoices': Invoices
        // "*": NotFound
    }

    const getStatus = async () => {
        const resp = await fetch('http://localhost:8080/api/workitem/status')

        if (!resp.ok) {
            return
        }

        const data = await resp.json() as StatusResponse

        console.log(data)

        if (data.isRunning) {
            runningWorkItem.start(data.workItem, data.deltaDurationNanoseconds)
        }
    }

    getStatus()

    const beforeUnload = (e: BeforeUnloadEvent) => {
        if (!$runningWorkItem.isRunning) {
            return
        }

        e.preventDefault()
        e.returnValue = ''
        return '...'
    }
</script>

<Router {routes} />
<Toaster position="bottom-right" />
<svelte:window on:beforeunload={beforeUnload} />

<style>
</style>
