<template>
  <div class="file-tree">
    <div
      v-for="node in tree"
      :key="node.path"
      :class="['tree-node', { 'tree-node-active': activePath === node.path }]"
    >
      <div class="tree-item" @click="handleClick(node)">
        <span class="tree-expand" @click.stop="toggleExpand(node)">
          <template v-if="node.type === 'dir'">
            {{ expandedPaths.has(node.path) ? '▼' : '▶' }}
          </template>
        </span>
        <span class="tree-icon">{{ node.type === 'dir' ? '📁' : getFileIcon(node.name) }}</span>
        <span class="tree-name">{{ node.name }}</span>
      </div>
      <div v-if="node.type === 'dir' && expandedPaths.has(node.path) && node.children" class="tree-children">
        <FileTree
          :tree="node.children"
          :active-path="activePath"
          :depth="depth + 1"
          @select="$emit('select', $event)"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

export interface FileNode {
  name: string
  path: string
  type: 'file' | 'dir'
  children?: FileNode[]
}

const props = withDefaults(defineProps<{
  tree: FileNode[]
  activePath?: string
  depth?: number
}>(), {
  depth: 0
})

const emit = defineEmits(['select'])

const expandedPaths = ref(new Set<string>())

function toggleExpand(node: FileNode) {
  if (node.type !== 'dir') return
  if (expandedPaths.value.has(node.path)) {
    expandedPaths.value.delete(node.path)
  } else {
    expandedPaths.value.add(node.path)
  }
}

function handleClick(node: FileNode) {
  if (node.type === 'dir') {
    toggleExpand(node)
  } else {
    emit('select', node)
  }
}

function getFileIcon(name: string): string {
  const ext = name.split('.').pop()?.toLowerCase() || ''
  const iconMap: Record<string, string> = {
    ts: '📘', tsx: '📘', js: '📙', jsx: '📙',
    vue: '💚', go: '🔵', py: '🐍', rs: '🦀',
    md: '📝', json: '📋', yaml: '📋', yml: '📋',
    css: '🎨', scss: '🎨', html: '🌐'
  }
  return iconMap[ext] || '📄'
}
</script>

<style scoped>
.file-tree {
  font-size: 14px;
}

.tree-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  cursor: pointer;
  border-radius: 4px;
  transition: background 0.15s;
}

.tree-item:hover {
  background: var(--bg-secondary);
}

.tree-node-active > .tree-item {
  background: var(--accent-color);
  color: #fff;
}

.tree-expand {
  width: 16px;
  font-size: 10px;
  text-align: center;
  flex-shrink: 0;
}

.tree-icon {
  font-size: 14px;
  flex-shrink: 0;
}

.tree-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tree-children {
  padding-left: 16px;
}
</style>
