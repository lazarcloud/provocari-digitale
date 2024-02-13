<script>
  import { page } from "$app/stores"

  let path = $page.url.password
  $: path = $page.url.password

  let show = false

  function gotoElement(id = "") {
    show = false

    if (id == "") return console.error("No id provided")
    const element = document.getElementById(id)

    if (!element) return console.error("Element not found")

    element.scrollIntoView({
      block: "start",
      behavior: "smooth",
      inline: "nearest",
    })
  }
</script>

<nav>
  <div>
    <a href="/"><img src="/logo.png" alt="logo" /></a>
    <!-- <div class="line"></div>
      <a href="/"><img src="/logo.png" alt="logo" /></a> -->
  </div>

  <div class="burger">
    <a
      href="#menu"
      on:click|preventDefault={() => {
        show = !show
      }}
      ><img
        class={show ? "opened" : "closed"}
        src="/burger.png"
        alt="logo"
        height="auto"
      /></a
    >
  </div>

  <div class="content {show ? 'show' : ''}">
    <a
      href="#home"
      class={path == "/" ? "active" : ""}
      on:click|preventDefault={() => {
        gotoElement("home")
      }}>Probleme</a
    >
    <a
      href="#projects"
      class={path == "" ? "active" : ""}
      on:click|preventDefault={() => {
        gotoElement("projects")
      }}>Soluții</a
    >
    <a
      href="#contact"
      class={path == "" ? "active" : ""}
      on:click|preventDefault={() => {
        gotoElement("contact")
      }}>Descriere</a
    >
  </div>
  <img class="rift" src="/rift.png" alt="rift" />
</nav>

<style>
  img {
    height: 32px;
    width: auto;
  }
  .line {
    width: 1px;
    height: 16px;
    background-color: var(--white);
    margin: 0 1rem;
  }
  nav {
    background-color: var(--bg);
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 2rem;
    width: 100%;
    position: fixed;
    z-index: 1;
  }
  nav {
    position: relative;
    border-bottom: 1px solid var(--text);
  }
  .rift {
    position: absolute;
    bottom: -43%;
    right: 12.5rem;
    z-index: -1;
  }
  div,
  a {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 1rem;
  }
  a {
    text-decoration: none;
    /* color: inherit; */
    color: var(--accent);
  }
  /* a.active {
      color: var(--accent);
    } */

  .burger {
    display: none;
  }
  .burger img {
    rotate: 0;
    transition: rotate 0.3s ease-in-out;
  }
  @media only screen and (max-width: 860px) {
    .rift {
      display: none;
    }
    nav {
      padding: 1rem 2rem;
      margin-bottom: 7rem;
    }
    img {
      height: 32px;
    }
    .content {
      width: 0;
      height: 100vh;
      background-color: var(--bg);
      display: flex;
      flex-direction: column;
      position: absolute;
      border-radius: 0 0 1rem 1rem;
      top: 0;
      right: 0;
      overflow: hidden;
      transition: width 0.3s ease-in-out;
    }
    .content a {
      font-size: 1.75rem;
      padding: 1rem 2rem;
      width: 100%;
      justify-content: right;
    }
    .content.show {
      width: 100%;
    }
    .burger {
      display: block;
      z-index: 10;
    }
    .burger img.opened {
      rotate: 90deg;
    }
  }
</style>
