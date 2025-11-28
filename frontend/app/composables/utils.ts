export function useApi() {
  const config = useRuntimeConfig();
  const API = `${config.public.baseApi}/api`;
  return { API };
}