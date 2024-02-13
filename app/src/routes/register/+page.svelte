<script>
</script>

<div class="container">
  <section>
    <h1>Fă-ți cont</h1>

    <form>
      <input type="text" placeholder="Email" />
      <input type="password" placeholder="Parolă" />
      <input type="password" placeholder="Confirmă Parola" />
      <button>Fă-ți cont</button>
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
