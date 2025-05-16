import { get, writable } from "svelte/store"
import Cookies from "js-cookie"
export const refresh = writable(Cookies.get("refresh") || "")
refresh.subscribe((value) => {
  Cookies.set("refresh", value)
})

export const userData = writable(Cookies.get("userData") || "")
userData.subscribe((value) => {
  Cookies.set("userData", value)
})

var baseUrl = "http://localhost:8080"
var token =
  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoiYWNjZXNzIiwiaXNzdWVyIjoibGF6YXIiLCJpc3N1ZWRBdCI6MTc0NzI5MDAzNSwiZXhwaXJlc0F0IjoxNzgzMjkwMDM1LCJ1c2VySWQiOiIiLCJhY2Nlc3NSb2xlIjoicHVibGljIn0.fM3XEO_mAtq9LDAPWXeZP89OvmQWzzQnfwyMG57zNvo"
export async function fetchAPI(path = "", options = {}) {
  options.headers = {
    ...options.headers,
    "Content-Type": "application/json",
    Authorization: "Bearer " + token,
  }
  const res = await fetch(baseUrl + path, options)
  // return res.text()
  const json = await res.json()
  if (json.error) alert(json.error)
  console.log(json)
  return json
}
export async function fetchAPIAuth(path = "", options = {}) {
  options.headers = {
    ...options.headers,
    "Content-Type": "application/json",
    Authorization: "Bearer " + get(refresh),
  }
  const res = await fetch(baseUrl + path, options)
  // return res.text()
  const json = await res.json()
  if (json.error) alert(json.error)
  console.log(json)
  return json
}
