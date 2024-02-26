import { fetchAPI, fetchAPIAuth } from "$lib"

export async function load({ params }) {
  const data = await fetchAPIAuth(`/api/solve/progress/${params.id}`)
  console.log("progress:", data)
  return { id: params.id, pb: data }
}
