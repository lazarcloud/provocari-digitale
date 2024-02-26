<script>
  import { fetchAPIAuth, refresh } from "$lib"
  import { invalidateAll } from "$app/navigation"
  import { onMount } from "svelte"

  export let data
  let id = data.id
  let pb = data.pb
  let solution = ""
  function encodeToBase64(str) {
    return btoa(unescape(encodeURIComponent(str)))
  }
  function formatTimeFromUnix(unix) {
    const date = new Date(unix)
    return date.toLocaleString()
  }
  let index = 0
  // TO DO: optimize this
  let interval
  onMount(async () => {
    interval = setInterval(async () => {
      await invalidateAll()
    }, 5000)

    return () => clearInterval(interval)
  })
</script>

<div class="container">
  <h1>{pb.title}</h1>
  <p>creată de <span>{pb.owner_email}</span></p>
  <p>{pb.description}</p>
  <p>max_memory: {pb.max_memory}Mb și max_time: {pb.max_time}s</p>
  {#if pb.uses_standard_io}
    <p>Se folosește standard input/output</p>
  {:else}
    <p>Se folosesc fișiere</p>
    <p>IN: {pb.input_file_name} OUT: {pb.output_file_name}</p>
  {/if}
  <!-- {JSON.stringify(data)} -->
  {#if $refresh != ""}
    {#key index}
      {#if data.solves.length}
        <h2>Rezolvările mele</h2>
        <table>
          <tr>
            <th>ID Test</th>
            <th>Scor Final</th>
            <th>Scor Maxim</th>
            <th>Nr. teste</th>
            <th>Data</th>
          </tr>
          {#each data.solves as test}
            <tr>
              <td><a href={`/solves/${test.id}`}>{test.id}</a></td>
              <td>{test.final_score == "NULL" ? "0" : test.final_score}</td>
              <td>{test.max_score}</td>
              <td>{test.test_count}</td>
              <td>{formatTimeFromUnix(test.created_at * 1000)}</td>
            </tr>
          {/each}
        </table>
      {:else}<h2>Nu ai rezolvat această problemă</h2>{/if}
    {/key}

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
        index++
      }}
    >
      <label for="solution">Soluție:</label>
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
</style>
