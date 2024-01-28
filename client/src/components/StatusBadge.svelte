<script lang="ts">
    import { Badge } from 'flowbite-svelte'
    import { WorkItemStatus } from '../utils/workitemstatus';

    export let statusCode: number
</script>

{#if statusCode === WorkItemStatus.ToDo}
    <Badge rounded color="dark">ToDo</Badge>
{:else if statusCode === WorkItemStatus.InProgress}
    <Badge rounded color="blue">InProgress</Badge>
{:else if statusCode === WorkItemStatus.Done}
    <Badge rounded color="green">Done</Badge>
{:else}
    <Badge rounded color="red">Unknown</Badge>
{/if}