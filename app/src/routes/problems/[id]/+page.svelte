<script>
  import { fetchAPIAuth, refresh } from "$lib"
  import { invalidateAll } from "$app/navigation"
  import { onMount } from "svelte"

  export let data
  let id = data.id
  let pb = data.pb
  let solution = ""
  let currentPage = 1
  let solvesPerPage = 5 // Number of solves to display per page

  function encodeToBase64(str) {
    return btoa(unescape(encodeURIComponent(str)))
  }

  function formatTimeFromUnix(unix) {
    const date = new Date(unix)
    return date.toLocaleString()
  }

  let interval

  onMount(async () => {
    interval = setInterval(async () => {
      await invalidateAll()
    }, 5000)

    return () => clearInterval(interval)
  })

  function nextPage() {
    currentPage++
  }

  function prevPage() {
    if (currentPage > 1) {
      currentPage--
    }
  }

  function getVisibleSolves() {
    const startIndex = (currentPage - 1) * solvesPerPage
    const endIndex = startIndex + solvesPerPage
    return data.solves.slice(startIndex, endIndex)
  }

  function formatUUID4chars(uuid) {
    return uuid.slice(0, 4)
  }
</script>

<div class="container">
  <h1>{pb.title}</h1>
  <h2>Detaliile problemei</h2>
  <table>
    <!-- <tr>
      <th>Autor</th>
      <th>da</th>
    </tr> -->
    <tr>
      <td>Autor</td>
      <td>{pb.owner_email}</td>
    </tr>
    <tr>
      <td>Memorie maximă</td>
      <td>{pb.max_memory}Kb</td>
    </tr>
    <tr>
      <td>Timp maxim</td>
      <td>{pb.max_time}ms</td>
    </tr>
    {#if pb.uses_standard_io}
      <tr>
        <td>Citire / Scriere</td>
        <td>Standard IO</td>
      </tr>
    {:else}
      <tr>
        <td>Citire / Scriere</td>
        <td>IN: {pb.input_file_name} OUT: {pb.output_file_name}</td>
      </tr>
    {/if}
  </table>
  <h2>Cerință</h2>
  <p>{pb.description}</p>

  {#if $refresh != ""}
    {#key currentPage}
      <h2>
        Rezolvările mele, maxim {data.bestScore} / {data.solves[0].max_score}
      </h2>
      <table>
        <tr>
          <th>#</th>
          <th>Scor Final</th>
          <th>Scor Maxim</th>
          <th>Nr. teste</th>
          <th>Status</th>
          <th>Data</th>
        </tr>
        {#each getVisibleSolves() as test, index}
          <tr>
            <td
              ><a href={`/solves/${test.id}`}
                >{(currentPage - 1) * solvesPerPage + index + 1}</a
              ></td
            >
            <td>{test.final_score == "NULL" ? "0" : test.final_score}</td>
            <td>{test.max_score}</td>
            <td>{test.test_count}</td>
            <td><a href={`/solves/${test.id}`}>{test.status}</a></td>
            <td>{formatTimeFromUnix(test.created_at * 1000)}</td>
          </tr>
        {/each}
      </table>
      <button on:click={prevPage}>Previous</button>
      <button on:click={nextPage}>Next</button>
    {/key}
    <h2>Soluție</h2>
    <form
      on:submit|preventDefault={async () => {
        console.log(solution)
        if (solution === "") alert("Soluția nu poate fi goală")
        const data = await fetchAPIAuth(`/api/solve/submit/${id}`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ code: encodeToBase64(solution) }),
        })
        console.log(data)
        await invalidateAll()
        solution = ""
      }}
    >
      <textarea
        id="solution"
        name="solution"
        bind:value={solution}
        rows="4"
        cols="50"
      ></textarea>
      <button type="submit">Trimite</button>
    </form>
  {:else}
    <h2>
      Pentru a încărca soluții ai nevoie de un <a href="/register">cont</a>.
    </h2>
  {/if}
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
