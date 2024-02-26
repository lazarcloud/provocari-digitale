import { fetchAPI, fetchAPIAuth, refresh, userData } from "$lib"
import { get } from "svelte/store"

export async function load({ params }) {
  const data = await fetchAPI(`/api/problems/${params.id}`)

  let returnData = { id: params.id, pb: data }

  if (get(refresh) != "") {
    const solves = await fetchAPIAuth(`/api/solve/${params.id}`)
    console.log("solves: ", solves)
    returnData.solves = solves.solves
  }

  return returnData
}
