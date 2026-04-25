<template>
  <div class="diff-viewer">
    <div v-for="(file, index) in files" :key="index" class="diff-file">
      <div class="diff-file-header">
        <span class="diff-file-name">{{ file.name }}</span>
        <span :class="['diff-file-status', `status-${file.status}`]">{{ file.status }}</span>
      </div>
      <div class="diff-content">
        <div v-for="(line, lineIndex) in file.lines" :key="lineIndex" :class="['diff-line', `diff-line-${line.type}`]">
          <span class="line-number">{{ line.oldNum || '' }}</span>
          <span class="line-number">{{ line.newNum || '' }}</span>
          <span class="line-prefix">{{ linePrefix(line.type) }}</span>
          <span class="line-content">{{ line.content }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
export interface DiffLine {
  type: 'added' | 'removed' | 'context'
  oldNum?: number
  newNum?: number
  content: string
}

export interface DiffFile {
  name: string
  status: 'added' | 'modified' | 'removed'
  lines: DiffLine[]
}

defineProps<{
  files: DiffFile[]
}>()

function linePrefix(type: string): string {
  switch (type) {
    case 'added': return '+'
    case 'removed': return '-'
    default: return ' '
  }
}
</script>

<style scoped>
.diff-viewer {
  border: 1px solid var(--border-color);
  border-radius: 6px;
  overflow: hidden;
}

.diff-file + .diff-file {
  border-top: 1px solid var(--border-color);
}

.diff-file-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
}

.diff-file-name {
  font-family: monospace;
  font-size: 13px;
  color: var(--text-primary);
}

.diff-file-status {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 10px;
  text-transform: uppercase;
}

.status-added { background: var(--success-color); color: #fff; }
.status-modified { background: var(--warning-color); color: var(--text-primary); }
.status-removed { background: var(--danger-color); color: #fff; }

.diff-content {
  font-family: monospace;
  font-size: 13px;
  overflow-x: auto;
}

.diff-line {
  display: flex;
  min-height: 20px;
  line-height: 20px;
}

.line-number {
  min-width: 40px;
  padding: 0 8px;
  text-align: right;
  color: var(--text-secondary);
  background: var(--bg-secondary);
  user-select: none;
  flex-shrink: 0;
}

.line-prefix {
  min-width: 20px;
  text-align: center;
  flex-shrink: 0;
}

.line-content {
  padding: 0 8px;
  white-space: pre;
  flex: 1;
}

.diff-line-added { background: rgba(46, 160, 67, 0.15); }
.diff-line-added .line-prefix { color: var(--success-color); }

.diff-line-removed { background: rgba(248, 81, 73, 0.15); }
.diff-line-removed .line-prefix { color: var(--danger-color); }

.diff-line-context .line-prefix { color: var(--text-secondary); }
</style>
