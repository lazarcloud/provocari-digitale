import { fetchAPI } from "$lib"

export async function load({ params }) {
  const data = await fetchAPI(`/api/problems/${params.id}`)

  return { id: params.id, pb: data }
}
