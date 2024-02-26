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
      <td>Scor</td>
      <td
        >{group.final_score == "NULL"
          ? "0"
          : group.final_score}/{group.max_score}</td
      >
    </tr>
    <tr>
      <td>Memorie maximă</td>
      <td>Kb</td>
    </tr>
  </table>
  <p>
    Scor {group.final_score == "NULL"
      ? "0"
      : group.final_score}/{group.max_score}
  </p>
  <p>ID Problemă: {group.problem_id}</p>
  <p>Număr teste: {group.test_count}</p>
  <p>Status: {group.status || "finished"}</p>
  <p>{group.source}</p>
  {#each pb.results as result, index}
    <div>
      <h2>Test {index}: {getRomanian(result.status, translations)}</h2>
      <p>{result.correct ? "Corect" : "Greșit"}</p>
      <p>Memorie necesară: {result.max_memory}Kb</p>
      <p>Timp necesar: {result.time_taken}</p>
    </div>
  {/each}
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
