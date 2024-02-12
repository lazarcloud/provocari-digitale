<script>
  import { fetchAPI } from "$lib"
</script>

{#await fetchAPI("/api/problems") then data}
  {JSON.stringify(data)}
{:catch error}
  <p>{error.message}</p>
{/await}
