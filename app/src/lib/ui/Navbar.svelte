<script>
  import { page } from "$app/stores"
  import { refresh, userData } from "$lib"

  let url = $page.url.pathname
  $: url = $page.url.pathname
  var right = 0
  $: if (url) {
    right = 0
    if (url == "/") right = 15.7
    if (stringContains(url, "/problems/")) right = 10.5
    if (url == "/login/" || url == "/register/") right = 3.25
  }

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

  function stringContains(string = "", substring = "") {
    return string.includes(substring)
  }
</script>

<nav>
  <div>
    <a style="opacity:1;" href="/"><img src="/logo.png" alt="logo" /></a>
    <!-- <div class="line"></div>
      <a href="/"><img src="/logo.png" alt="logo" /></a> -->
  </div>

  <div class="burger">
    <a
      href="#menu"
      style="opacity:1;"
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
      href="/"
      class={url == "/" ? "active" : ""}
      on:click={() => (show = false)}
      >Acasă
    </a>
    <a
      href="/problems"
      class={stringContains(url, "/problems/") ? "active" : ""}
      on:click={() => (show = false)}
      >Probleme
    </a>
    {#if $refresh != ""}
      <a
        href="#logout"
        class={url == "/logout/" ? "active" : ""}
        on:click={() => {
          show = false
          refresh.set("")
          userData.set("")
        }}
        >Ieși din cont
      </a>
    {:else}
      <a
        href="/login"
        class={url == "/login/" || url == "/register/" ? "active" : ""}
        on:click={() => (show = false)}
        >Loghează-te
      </a>
    {/if}
  </div>
  <img class="rift" src="/rift.png" alt="rift" style="right:{right}rem;" />
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
    width: 4rem;
    position: absolute;
    bottom: -43%;
    z-index: -1;
    transition: right 0.6s ease-out;
  }
  div,
  a {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 1.5rem;
  }
  a {
    text-decoration: none;
    color: var(--accent);
    text-wrap: nowrap;
  }
  a.active {
    opacity: 1;
  }

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
