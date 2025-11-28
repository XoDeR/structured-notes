export interface APIResult<Data> {
  status: 'success' | 'error';
  message: string;
  result?: Data;
}

export interface FetchOptions {
  id: string;
}

function customFetch(route: string, method: string, body: object) {
  const { API } = useApi();
  if (route.endsWith('/')) route = route.slice(0, -1);
  return fetch(`${API}/${route}`, {
    method: method,
    body: method === 'GET' || method === 'DELETE' ? null : body instanceof FormData ? body : JSON.stringify(body),
    headers: body instanceof FormData ? {} : { 'Content-Type': 'application/json; charset=UTF-8' },
    credentials: 'include',
  });
}

let refreshPromise: Promise<void> | null = null;

async function refreshAccessToken(): Promise<void> {
  if (!refreshPromise) {
    refreshPromise = (async () => {
      console.log('[AUTH] Refreshing access token...');
      const res = await customFetch('auth/refresh', 'POST', {});
      const data = await res.json();

      if (!res.ok || data.status !== 'success') {
        console.warn('[AUTH] Refresh failed, logging out.');

        useUserStore().postLogout();

        window.location.replace('/login');
        throw new Error('Refresh token invalid');
      }

      console.log('[AUTH] Access token refreshed.');
    })();

    refreshPromise.finally(() => (refreshPromise = null));
  }

  return refreshPromise;
}

export async function makeRequest<T>(route: string, method: string, body: object): Promise<APIResult<T>> {
  try {
    const response = await customFetch(route, method, body);
    const data = await response.json();

    if (response.ok && data.status === 'success') {
      return data;
    }

    if (response.status === 401 && (data.message === 'Bad access token.' || data.message === 'Missing token cookies.')) {
      try {
        await refreshAccessToken();

        const retry = await customFetch(route, method, body);
        const retryData = await retry.json();
        return retryData;
      } catch {
        return { status: 'error', message: 'Authentication failed.' };
      }
    }

    return data;
  } catch (err) {
    console.error('[API ERROR]', err);
    return { status: 'error', message: String(err) };
  }
}



