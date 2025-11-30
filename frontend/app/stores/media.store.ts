import { makeRequest } from '~/helpers/apiClient';
import type { Node } from './interfaces';

export const useMediaStore = defineStore('media', {
  state: () => ({
    mediaFiles: [] as Node[],
    isFetching: false,
  }),
  getters: {
    getAll: state => state.mediaFiles,
    getById: state => (id: string) => state.mediaFiles.find(c => c.id == id),
  },
  actions: {
    async post(mediaFile: FormData): Promise<Node> {
      const request = await makeRequest(`media`, 'POST', mediaFile);
      if (request.status == 'success') {
        useNodeStore().nodes.set((request.result as Node).id, request.result as Node);
        return request.result as Node;
      } else throw request.message;
    },
    clear() {
      this.mediaFiles = [];
      this.isFetching = false;
    },
  },
});
