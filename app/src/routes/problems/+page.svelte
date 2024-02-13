<script>
  import { fetchAPI } from "$lib"
</script>

<div class="container">
  <h1>Probleme</h1>
</div>
{#await fetchAPI("/api/problems") then data}
  {#each data.problems as { id, title, max_memory, max_time, description, owner_email }}
    <p>{id}</p>
    <p>Problema {title} de {owner_email}</p>
    <p>{max_memory}Mb și {max_time}s</p>
    <p>{description}</p>
  {/each}
{:catch error}
  <p>{error.message}</p>
{/await}
