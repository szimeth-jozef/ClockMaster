<script lang="ts">
    import { Button, Modal, P, Heading } from 'flowbite-svelte'
    import { ExclamationCircleOutline } from 'flowbite-svelte-icons'
    import toast from 'svelte-french-toast'
    import { push } from 'svelte-spa-router'

    export let isOpen: boolean = false
    export let workItemId: number = -1
    export let afterDelete: () => void = () => {}

    const deleteWorkItem = async (id: number) => {
        const resp = await fetch(`http://localhost:8080/api/workitem/${id}`, {
            method: 'DELETE'
        })

        if (!resp.ok) {
            toast.error(`Failed to delete work item`, { duration: 3000 })
            return
        }

        toast.success(`Deleted work item with ID: ${id}`, { duration: 3000 })
        afterDelete()
    }
</script>

<Modal bind:open={isOpen} size="xs" autoclose>
    <div class="text-center">
        <ExclamationCircleOutline class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
        <Heading tag="h5" class="mb-5">Are you sure you want to delete work item with ID {workItemId}?</Heading>
        <Button on:click={() => deleteWorkItem(workItemId)} color="red" class="me-2">Yes, I'm sure</Button>
        <Button color="alternative">No, cancel</Button>
    </div>
</Modal>