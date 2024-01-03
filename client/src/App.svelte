<script lang="ts">
    import Home from './routes/Home.svelte'
    import WorkItem from './routes/WorkItem.svelte'
    import Invoices from './routes/Invoices.svelte'
    import Router from 'svelte-spa-router'


    const routes = {
        '/': Home,
        '/workitem/:id': WorkItem,
        '/invoices': Invoices
        // "*": NotFound
    }
</script>

<Router {routes} />

<style>
</style>
