<script>
  import { goto } from "$app/navigation"
  import { fetchAPI } from "$lib"
</script>

<div class="container">
  <h1>
    Puntea <span>riftului</span><br />
    dintre <span>probleme</span><br />
    și <span>soluții.</span><br />
  </h1>
  <button on:click={() => goto("/register")}>{`{ Începe }`}</button>
</div>

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
  .container {
    display: flex;
    flex-direction: column;
    justify-content: center;
    height: 70vh;
  }
  button {
    color: black;
    background-color: white;
    border: 1px solid black;
    transition: all 0.3s ease-in-out;
    width: 160px;
  }
  button:hover {
    color: white;
    background-color: black;
  }
</style>
