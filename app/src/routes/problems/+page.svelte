<script>
  import { fetchAPI } from "$lib"
</script>

<div class="container">
  <h1>Probleme</h1>
  <div class="content">
    {#await fetchAPI("/api/problems")}
      <p>loading ...</p>
    {:then data}
      {#each data.problems as { id, title, max_memory, max_time, description, owner_email }}
        <div class="pb">
          <h2>{title}</h2>
          <a href={`/problems/${id}`}>Deschide</a>
        </div>
      {/each}
    {:catch error}
      <p>{error.message}</p>
    {/await}
  </div>
</div>

<style>
  .pb {
    border-bottom: 1px solid white;
  }
  h1 {
    font-size: clamp(2rem, 10vw, 4rem);
    text-align: center;
  }
  h2 {
    font-size: clamp(1.5rem, 5vw, 2rem);
  }
  .content {
    width: clamp(400px, 80%, 700px);
    margin: 0 auto;
  }
</style>
