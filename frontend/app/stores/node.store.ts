import { Collection } from './collection';
import type { DbNode, Node, Permission } from './interfaces';

export const useNodeStore = defineStore('node', {
  state: () => ({
    nodes: new Collection<string, Node>(),
    publicNodes: new Collection<string, Node>(),
    allTags: [] as string[],
    isFetching: false,
  }),
});