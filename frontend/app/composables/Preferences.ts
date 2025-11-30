export const DEFAULT_PREFERENCES = {
  datatableItemsCount: 10 as number,
};

export type Preferences = {
  [K in keyof typeof DEFAULT_PREFERENCES]: (typeof DEFAULT_PREFERENCES)[K];
};
export type PreferenceKey = keyof Preferences;

const preferences = reactive<Preferences>(loadPreferences());

function loadPreferences(): Preferences {
  try {
    const raw = JSON.parse(localStorage.getItem('preferences') || '{}') as Partial<Preferences>;
    return {
      ...DEFAULT_PREFERENCES,
      ...raw,
    };
  } catch {
    return { ...DEFAULT_PREFERENCES };
  }
}

function savePreferences() {
  localStorage.setItem('preferences', JSON.stringify(preferences));
}

export function usePreferences() {
  function get<K extends PreferenceKey>(key: K) {
    return computed<Preferences[K]>({
      get: () => preferences[key] as Preferences[K],
      set: val => {
        preferences[key] = val;
        savePreferences();
      },
    });
  }

  function set<K extends PreferenceKey>(key: K, value: Preferences[K]) {
    preferences[key] = value;
    savePreferences();
  }

  function reset() {
    Object.assign(preferences, DEFAULT_PREFERENCES);
    savePreferences();
  }

  const all = preferences;

  return {
    get,
    set,
    reset,
    all,
  };
}