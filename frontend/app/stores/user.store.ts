import { defineStore } from 'pinia'
import type { User, PublicUser, ConnectionLog } from './interfaces';

export const useUserStore = defineStore('user', {
  state: () => ({
    user: undefined as User | undefined,
    users: {} as Record<string, PublicUser>,
    lastConnection: null as ConnectionLog | null,
    currentFetching: [] as string[], // List of currently fetching user ids to avoid duplicate requests
  }),
  getters: {
    getById: state => (id: string) => state.users[id],
    search: state => (query: string) => {
      const lowerQuery = query.toLowerCase();
      return Object.values(state.users).filter(
        u => u.username?.toLowerCase().includes(lowerQuery) || u.email?.toLowerCase().includes(lowerQuery) || u.id === query,
      );
    },
  },
  actions: {
    async login(username: string, password: string) { },
    async register(user: Omit<User, 'id' | 'created_timestamp' | 'updated_timestamp'>): Promise<boolean> { return false; },
    async postLogout() {
      this.user = undefined;
      useUserStore().clear();
      if (import.meta.client) localStorage.removeItem('isLoggedIn');
    },
    clear() {
      this.user = undefined;
      this.users = {};
      this.lastConnection = null;
      this.currentFetching = [];
    },
  },
});