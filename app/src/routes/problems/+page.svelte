<script>
  import { fetchAPI } from "$lib"
</script>

<div class="container">
  <h1>Probleme</h1>
</div>
{#await fetchAPI("/api/problems")}
  <p>loading ...</p>
{:then data}
  {#each data.problems as { id, title, max_memory, max_time, description, owner_email }}
    <p>{id}</p>
    <p>Problema {title} de {owner_email}</p>
    <p>{max_memory}Mb și {max_time}s</p>
    <p>{description}</p>

    {#await fetchAPI(`/api/problems/${id}/testscount`)}
      <p>Teste: loading...</p>
    {:then tests}
      <p>Test: {tests.count}</p>
    {:catch error}
      <p>{error.message}</p>
    {/await}
  {/each}
{:catch error}
  <p>{error.message}</p>
{/await}
