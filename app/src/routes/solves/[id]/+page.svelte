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
