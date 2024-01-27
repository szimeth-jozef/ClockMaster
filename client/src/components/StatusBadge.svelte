<script lang="ts">
    import { Badge } from 'flowbite-svelte'

    export let statusCode: number
</script>

{#if statusCode === 0}
    <Badge rounded color="dark">ToDo</Badge>
{:else if statusCode === 1}
    <Badge rounded color="blue">InProgress</Badge>
{:else if statusCode === 2}
    <Badge rounded color="green">Done</Badge>
{:else}
    <Badge rounded color="red">Unknown</Badge>
{/if}