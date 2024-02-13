<script>
</script>

<div class="container">
  <section>
    <h1>Loghează-te</h1>

    <form>
      <input type="text" placeholder="Email" />
      <input type="password" placeholder="Parolă" />
      <button>Loghează-te</button>
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
