<script lang="ts">
    import Navbar from '../components/Navbar.svelte'

    const getWorkItems = async () => {
        const resp = await fetch('http://localhost:8080/api/workitem')
        const data = await resp.json()
        return data
    }

    getWorkItems().then(data => {
        console.log(data)
    })
</script>

<header>
    <Navbar />
</header>

<main>
</main>

<footer></footer>

<style>
</style>