<script lang="ts">
    import Navbar from './Navbar.svelte'
    import { Toast } from 'flowbite-svelte'
</script>

<header>
    <Navbar />
</header>
<main class="container mx-auto mt-4">
    <div class="absolute">
        <Toast>
            This is a toast
        </Toast>
        <Toast>
            This is another toast
        </Toast>
        <Toast>
            Oooh, a third toast
        </Toast>
    </div>
    <slot />
</main>
<footer></footer>

<style>
</style>