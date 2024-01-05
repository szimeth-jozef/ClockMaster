<script lang="ts">
    import BaseLayout from '../components/BaseLayout.svelte'

    export let params: any
</script>

<BaseLayout>
    WorkItem {params.id}
</BaseLayout>

<style>
</style>