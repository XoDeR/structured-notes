import { makeRequest, type FetchOptions } from '~/helpers/apiClient';
import { Collection } from './collection';
import type { DbNode, Node, Permission } from './interfaces';

export const useNodeStore = defineStore('node', {
  state: () => ({
    nodes: new Collection<string, Node>(),
    publicNodes: new Collection<string, Node>(),
    allTags: [] as string[],
    isFetching: false,
  }),
  getters: {
    media: state => state.nodes.filter(d => d.role === 4),
  },
  actions: {
    recomputeTags() {
      const tags = new Set<string>();
      this.nodes.forEach(node => {
        if (node.tags) {
          parseTags(node.tags).forEach(tag => tags.add(tag));
        }
        if (node.parent_id && !this.nodes.get(node.parent_id)) {
          node.parent_id = '';
        }
      });
      this.allTags = Array.from(tags).sort();
    },

    async fetch<T extends FetchOptions>(opts?: T): Promise<'id' extends keyof T ? Node : Collection<string, Node>> {
      if (opts?.id && !this.nodes.get(opts.id)?.partial) {
        // retrieve from the store
        return this.nodes.get(opts.id) as 'id' extends keyof T ? Node : Collection<string, Node>;
      }

      console.log(`[store/node] Fetching node(s) with options: ${JSON.stringify(opts)}`);

      if (!this.nodes.size) {
        this.isFetching = true;
      }

      const request = await makeRequest(`nodes/@me/${opts?.id || ''}`, 'GET', {});

      this.isFetching = false;

      if (request.status == 'success') {
        if (opts?.id) {
          // Single node
          const result = request.result as { node: DbNode; permissions: Permission[] };
          const n = this.nodes.get(opts.id);
          let shared = false;
          if (n) shared = n.shared;
          const updatedNode: Node = { ...(result.node as DbNode), partial: false, shared: shared, permissions: result.permissions };
          this.nodes.set(opts.id, updatedNode);
          return updatedNode as 'id' extends keyof T ? Node : Collection<string, Node>;
        } else {
          // List of nodes
          for (const node of request.result as DbNode[]) {
            this.nodes.set(node.id, { ...node, partial: true, shared: false, permissions: [] });
          }
          this.recomputeTags();
          return this.nodes as 'id' extends keyof T ? Node : Collection<string, Node>;
        }
      } else throw request;
    },
  }
});

function parseTags(tags: string): string[] {
  if (typeof tags === 'string') {
    return tags
      .split(',')
      .map(tag => tag.trim())
      .filter(Boolean);
  }
  return [];
}