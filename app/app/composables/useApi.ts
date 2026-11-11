
export const useApi = () => {
  const config = useRuntimeConfig()

  const apiBase = config.public.apiBase

  return {
    apiBase,
  }
}
