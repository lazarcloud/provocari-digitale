<script>
  import "./global.css"
  import { Navbar } from "$ui"
  import { onMount } from "svelte"
  let height = -70
  onMount(() => {
    setTimeout(() => {
      height = -5
    }, 100)
  })
</script>

<img src="/riftBig.png" alt="rift" style="top:{height}%;" />
<Navbar />
<slot />

<style>
  img {
    width: clamp(150px, 70%, 600px);
    height: auto;
    position: absolute;
    z-index: -1;
    left: 50%;
    transform: translateX(-50%);
    transition: top 1.25s ease-in-out;
    opacity: 0.5;
  }
</style>
