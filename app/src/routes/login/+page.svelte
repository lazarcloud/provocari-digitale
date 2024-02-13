<script>
  import { goto } from "$app/navigation"
  import { fetchAPI, fetchAPIAuth, userData } from "$lib"
  import { refresh } from "$lib"
  let email = "",
    password = ""

  async function login() {
    const data = await fetchAPI("/api/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    })
    if (data.status == "ok") {
      const refreshToken = data.refreshToken
      refresh.set(refreshToken)
      userData.set(email.toLocaleLowerCase())
      goto("/problems")
    }
  }
</script>

<div class="container">
  <section>
    <h1>Loghează-te</h1>
    <form on:submit={login}>
      <input type="email" bind:value={email} placeholder="Email" />
      <input type="password" bind:value={password} placeholder="Parolă" />
      <input type="submit" value="Loghează-te" />
    </form>
    <p>Nu ai cont? <a href="/register">Fă-ți unul gratuit!</a></p>
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
