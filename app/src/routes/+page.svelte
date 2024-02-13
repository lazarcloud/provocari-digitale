<script>
  import { fetchAPI } from "$lib"
</script>

<div class="container">
  <h1>
    Puntea <span>riftului</span><br />
    dintre <span>probleme</span><br />
    și <span>soluții.</span><br />
  </h1>
</div>

<!-- {#await fetchAPI("/api/problems") then data}
  {JSON.stringify(data)}
{:catch error}
  <p>{error.message}</p>
{/await} -->

<style>
  h1 {
    font-size: clamp(2rem, 10vw, 4rem);
    /* margin-top: 4em; */
  }
  span {
    font-family: "Protest Riot", sans-serif;
    font-weight: 400;
    font-style: normal;
    font-size: 1.1em;
  }
</style>
