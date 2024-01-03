<script lang="ts">
    import {location} from 'svelte-spa-router'
    import {
        Navbar,
        NavBrand,
        NavLi,
        NavUl,
        NavHamburger,
        DarkMode,
    } from 'flowbite-svelte'

    $: activeUrl = `#${$location}`;
    let activeClass = 'text-white bg-green-700 md:bg-transparent md:text-green-700 md:dark:text-white dark:bg-green-600 md:dark:bg-transparent';
    let nonActiveClass = 'text-gray-700 hover:bg-gray-100 md:hover:bg-transparent md:border-0 md:hover:text-green-700 dark:text-gray-400 md:dark:hover:text-white dark:hover:bg-gray-700 dark:hover:text-white md:dark:hover:bg-transparent';
</script>

<header>
    <Navbar class="border-b">
        <NavBrand href="/">
            <svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" height="100" viewBox="0 0 48 48" class="me-3 h-6 sm:h-9">
                <path fill="#00acc1" d="M44,24c0,11.044-8.956,20-20,20S4,35.044,4,24S12.956,4,24,4S44,12.956,44,24z"></path>
                <path fill="#eee" d="M40,24c0,8.838-7.162,16-16,16S8,32.838,8,24S15.163,8,24,8S40,15.163,40,24z"></path>
                <path fill="#78909c" d="M23 39.95C23.332 39.97 23.662 40 24 40s.668-.03 1-.05V36h-2V39.95zM36 25h3.95C39.97 24.668 40 24.338 40 24s-.03-.668-.05-1H36V25zM8.05 25H12v-2H8.05C8.03 23.332 8 23.662 8 24S8.03 24.668 8.05 25z"></path>
                <path d="M31.664,30.27L29.934,32l-6.797-6.797l1.73-1.73L31.664,30.27z"></path>
                <path d="M27,24c0,1.656-1.344,3-3,3s-3-1.344-3-3s1.344-3,3-3S27,22.344,27,24"></path>
                <path d="M23,11h2v10.336h-2V11z"></path>
                <path fill="#00acc1" d="M25,24c0,0.551-0.449,1-1,1s-1-0.449-1-1s0.449-1,1-1S25,23.449,25,24"></path>
            </svg>
          <span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">ClockMaster</span>
        </NavBrand>
        <div class="flex md:order-2">
            <div class="md:hidden text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-700 rounded-lg text-sm p-2.5" data-collapse-toggle="mobile-menu-3" aria-controls="mobile-menu-3" aria-expanded="false">
                <DarkMode />
            </div>
            <div class="hidden relative md:block">
                <DarkMode />
            </div>
            <NavHamburger />
          </div>
        <NavUl {activeUrl} {activeClass} {nonActiveClass}>
          <NavLi href="#/">Home</NavLi>
          <NavLi href="#/invoices">Invoices</NavLi>
        </NavUl>
        <!-- <DarkMode /> -->
    </Navbar>
</header>