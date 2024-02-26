<script>
  import { fetchAPIAuth } from "$lib"
  import { invalidateAll } from "$app/navigation"
  import { onMount } from "svelte"

  export let data
  let id = data.id
  $: id = data.id
  let pb = data.pb
  $: pb = data.pb
  let group = data.pb.group
  $: group = data.pb.group

  console.log(pb)

  let translations = {
    NULL: "Așteptare",
    waiting: "Așteptare",
    running: "Rulare",
    finished: "Terminat",
    compiling: "Compilare",
  }

  function getRomanian(message = "", translations = {}) {
    if (translations.hasOwnProperty(message)) {
      return translations[message]
    }
    return message
  }

  // TO DO: clean up this repeated fetching
  onMount(async () => {
    setInterval(async () => {
      await invalidateAll()
    }, 1000)
  })
</script>

<div class="container">
  <h1>Soluție {id}</h1>
  <h2>Detaliile soluției</h2>
  <table>
    <tr>
      <td>ID problemă</td>
      <td>{group.problem_id}</td>
    </tr>
    <tr>
      <td>Scor</td>
      <td
        >{group.final_score == "NULL"
          ? "0"
          : group.final_score}/{group.max_score}</td
      >
    </tr>
    <tr>
      <td>Număr teste</td>
      <td>{group.test_count}</td>
    </tr>
    <tr>
      <td>Status</td>
      <td>{getRomanian(group.status, translations) || "Anulat"}</td>
    </tr>
  </table>
  <h2>Codul sursă</h2>
  <p>{group.source}</p>
  <h2>Rezultatele testelor</h2>
  <table>
    <tr>
      <th>#</th>
      <th>Status</th>
      <th>Verdict</th>
      <th>Memorie necesară</th>
      <th>Timp necesar</th>
    </tr>
    {#each pb.results as result, index}
      <tr>
        <td>{index}</td>
        <td>{getRomanian(result.status, translations)}</td>
        {#if "Așteptare" == getRomanian(result.status, translations)}
          <td>Așteptare</td>
        {:else}
          <td>{result.correct ? "Corect" : "Greșit"}</td>
        {/if}
        <td>{result.max_memory}{result.max_memory != "" ? "Kb" : ""}</td>
        <td>{result.time_taken}</td>
      </tr>
    {/each}
  </table>
</div>

<style>
  div.container {
    margin: 10px;
    padding: 10px;
  }
  table {
    width: 100%;
    border-collapse: collapse;
    max-height: 300px;
    overflow-y: auto;
  }
  th,
  td {
    border: 1px solid black;
    padding: 8px;
    text-align: left;
  }
  span {
    opacity: 0.6;
  }
  textarea {
    width: 100%;
    color: white;
    background-color: transparent;
  }
  h1 {
    font-size: clamp(2rem, 10vw, 4rem);
    text-align: center;
  }
  h2 {
    font-size: clamp(1.5rem, 5vw, 2rem);
  }
  form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    width: 100%;
    padding: 0;
  }
</style>
