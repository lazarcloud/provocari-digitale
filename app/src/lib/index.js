import { get, writable } from "svelte/store"
export const refresh = writable("")
export const userData = writable("")

var baseUrl = "http://localhost:8080"
var token =
  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0eXBlIjoiYWNjZXNzIiwiaXNzdWVyIjoibGF6YXIiLCJpc3N1ZWRBdCI6MTcwNzc2OTg4OCwiZXhwaXJlc0F0IjoxNzQzNzY5ODg4LCJ1c2VySWQiOiIiLCJhY2Nlc3NSb2xlIjoicHVibGljIn0.YEXte-T7r6ax6Hopzt6KYoesj7Ia-caf7dlWrHGnVEs"
export async function fetchAPI(path, options) {
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
export async function fetchAPIAuth(path, options) {
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
