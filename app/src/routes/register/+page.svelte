<script>
  import { goto } from "$app/navigation"
  import { fetchAPI } from "$lib"
  let email = "",
    password = "",
    confirmPassword = ""

  async function register() {
    if (password !== confirmPassword) return alert("Parolele nu coincid")
    const data = await fetchAPI("/api/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, password, confirmPassword }),
    })
    if (data.status == "ok") goto("/login")
  }
</script>

<div class="container">
  <section>
    <h1>Fă-ți cont</h1>

    <form on:submit={register}>
      <input type="email" bind:value={email} placeholder="Email" />
      <input type="password" bind:value={password} placeholder="Parolă" />
      <input
        type="password"
        bind:value={confirmPassword}
        placeholder="Confirmă Parola"
      />
      <input type="submit" value="Fă-ți cont" />
    </form>
    <p>Ai cont? <a href="/login">Loghează-te!</a></p>
  </section>
</div>

<style>
  h1 {
    border-bottom: 1px solid var(--text);
  }
  section {
    width: clamp(300px, 80%, 500px);
    display: flex;
    flex-direction: column;
    align-items: center;
    margin: 0 auto;
  }
</style>
