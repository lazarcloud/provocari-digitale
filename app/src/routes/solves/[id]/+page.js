import { fetchAPI, fetchAPIAuth } from "$lib"

export async function load({ params }) {
  const data = await fetchAPIAuth(`/api/solve/progress/${params.id}`)

  return { id: params.id, pb: data }
}
