import { fetchAPI, fetchAPIAuth, refresh, userData } from "$lib"
import { get } from "svelte/store"

export async function load({ params }) {
  const data = await fetchAPI(`/api/problems/${params.id}`)

  let returnData = { id: params.id, pb: data }

  if (get(refresh) != "") {
    const solves = await fetchAPIAuth(`/api/solve/${params.id}`)
    console.log("solves: ", solves)
    returnData.solves = solves.solves

    const maxScore = await fetchAPIAuth(`/api/solve/max_score/${params.id}`)
    console.log("maxScore: ", maxScore)
    returnData.bestScore = maxScore.max_score
  }

  return returnData
}
