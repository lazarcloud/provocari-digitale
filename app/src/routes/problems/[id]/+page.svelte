<script>
  import { fetchAPIAuth } from "$lib"
  import { invalidateAll } from "$app/navigation"

  export let data
  let id = data.id
  let pb = data.pb
  let solution = ""
  function encodeToBase64(str) {
    return btoa(unescape(encodeURIComponent(str)))
  }
  let index = 0
</script>

<h1>Problem {id}</h1>
{JSON.stringify(pb)}

<h2>My solves</h2>

{#key index}
  {#await fetchAPIAuth(`/api/solve/${id}`)}
    <p>loading...</p>
  {:then solves}
    {#each solves.solves as solve}
      <p>{JSON.stringify(solve)}</p>
    {/each}
  {:catch error}
    <p>{error.message}</p>
  {/await}
{/key}

<form
  on:submit|preventDefault={async () => {
    console.log(solution)
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
  <label for="solution">Solution:</label>
  <textarea
    id="solution"
    name="solution"
    bind:value={solution}
    rows="4"
    cols="50"
  ></textarea>
  <button type="submit">Submit</button>
</form>
