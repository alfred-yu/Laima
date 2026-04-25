<template>
  <div class="repo-code">
    <div class="code-header">
      <div class="branch-selector">
        <Skeleton v-if="branchesLoading" type="text" width="100px" />
        <select v-else v-model="selectedBranch" @change="loadFiles">
          <option v-for="branch in branches" :key="branch" :value="branch">{{ branch }}</option>
        </select>
      </div>
      <div class="code-actions">
        <button class="btn">上传文件</button>
        <button class="btn">新建文件</button>
        <button class="btn">上传文件夹</button>
      </div>
    </div>
    <div class="code-content">
      <div class="file-tree">
        <div class="file-tree-header">
          <h3>文件</h3>
        </div>
        <div v-if="filesLoading" class="loading-state">
          <Skeleton type="list" :count="8" />
        </div>
        <div v-else-if="filesError" class="error-state">
          {{ filesError }}
        </div>
        <ul v-else class="file-list">
          <li 
            v-for="file in files" 
            :key="file.path" 
            class="file-item"
            :class="{ active: currentFile.path === file.path }"
            @click="selectFile(file)"
          >
            <span class="file-icon">{{ file.type === 'dir' ? '📁' : getFileIcon(file.name) }}</span>
            <span class="file-name">{{ file.name }}</span>
          </li>
        </ul>
      </div>
      <div class="file-content">
        <div v-if="fileContentLoading" class="loading-state">
          <Skeleton type="text" :count="20" />
        </div>
        <div v-else-if="fileContentError" class="error-state">
          {{ fileContentError }}
        </div>
        <template v-else>
          <div class="file-header">
            <h2>{{ currentFile.name }}</h2>
            <div class="file-actions">
              <button class="btn">编辑</button>
              <button class="btn">删除</button>
              <button class="btn">历史</button>
            </div>
          </div>
          <div class="code-editor">
            <pre>{{ currentFile.content }}</pre>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { Skeleton } from '../components'
import { repoApi } from '../services/api'

const route = useRoute()
const owner = computed(() => route.params.owner as string)
const repo = computed(() => route.params.repo as string)

const selectedBranch = ref('main')
const branches = ref<string[]>([])
const files = ref<any[]>([])
const currentFile = ref<any>({ name: '', path: '', content: '', type: 'file' })

const branchesLoading = ref(true)
const filesLoading = ref(true)
const fileContentLoading = ref(false)
const branchesError = ref('')
const filesError = ref('')
const fileContentError = ref('')

const loadBranches = async () => {
  try {
    branchesLoading.value = true
    branchesError.value = ''
    const response = await repoApi.listBranches(owner.value, repo.value)
    const branchList = (response as any).items || []
    branches.value = branchList.length > 0 ? branchList : ['main']
    selectedBranch.value = branches.value[0]
  } catch (err: any) {
    branchesError.value = err.message || '加载分支失败'
    branches.value = ['main']
  } finally {
    branchesLoading.value = false
  }
}

const loadFiles = async () => {
  try {
    filesLoading.value = true
    filesError.value = ''
    // TODO: 当 API 支持文件列表接口后替换
    // const response = await repoApi.listFiles(owner.value, repo.value, selectedBranch.value)
    // files.value = response.items || []
    files.value = []
  } catch (err: any) {
    filesError.value = err.message || '加载文件失败'
  } finally {
    filesLoading.value = false
  }
}

const selectFile = async (file: any) => {
  if (file.type === 'dir') {
    return
  }
  
  currentFile.value = file
  try {
    fileContentLoading.value = true
    fileContentError.value = ''
    // TODO: 当 API 支持文件内容接口后替换
    // const response = await repoApi.getFileContent(owner.value, repo.value, file.path, selectedBranch.value)
    // currentFile.value.content = response.content
    currentFile.value.content = ''
  } catch (err: any) {
    fileContentError.value = err.message || '加载文件内容失败'
  } finally {
    fileContentLoading.value = false
  }
}

const getFileIcon = (name: string): string => {
  const ext = name.split('.').pop()?.toLowerCase() || ''
  const iconMap: Record<string, string> = {
    ts: '📘', tsx: '📘', js: '📙', jsx: '📙',
    vue: '💚', go: '🔵', py: '🐍', rs: '🦀',
    md: '📝', json: '📋', yaml: '📋', yml: '📋',
    css: '🎨', scss: '🎨', html: '🌐'
  }
  return iconMap[ext] || '📄'
}

onMounted(() => {
  loadBranches()
  loadFiles()
})
</script>

<style scoped>
.repo-code {
  padding: 20px;
}

.code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.branch-selector select {
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  background: var(--bg-primary);
  color: var(--text-primary);
  transition: border-color 0.3s, background-color 0.3s, color 0.3s;
}

.code-actions {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 6px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.3s;
}

.btn:hover {
  background: var(--bg-secondary);
}

.code-content {
  display: flex;
  gap: 20px;
  min-height: 600px;
}

.file-tree {
  width: 300px;
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 20px;
  transition: background-color 0.3s;
}

.file-tree-header h3 {
  margin: 0 0 20px 0;
  font-size: 16px;
  color: var(--text-primary);
  transition: color 0.3s;
}

.loading-state,
.error-state {
  padding: 20px 0;
  color: var(--text-secondary);
}

.error-state {
  color: var(--danger-color);
}

.file-list {
  list-style: none;
  padding: 0;
}

.file-item {
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  display: flex;
  align-items: center;
}

.file-item:hover {
  background: var(--bg-primary);
}

.file-item.active {
  background: var(--accent-color);
  color: #fff;
}

.file-icon {
  margin-right: 8px;
  font-size: 16px;
}

.file-name {
  font-size: 14px;
  color: var(--text-primary);
  transition: color 0.3s;
}

.file-content {
  flex: 1;
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: background-color 0.3s, box-shadow 0.3s;
}

.file-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--border-color);
  transition: border-color 0.3s;
}

.file-header h2 {
  margin: 0;
  font-size: 18px;
  color: var(--text-primary);
  transition: color 0.3s;
}

.file-actions {
  display: flex;
  gap: 10px;
}

.code-editor {
  background: var(--bg-primary);
  border-radius: 4px;
  padding: 20px;
  overflow: auto;
  max-height: 500px;
  transition: background-color 0.3s;
}

.code-editor pre {
  margin: 0;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
  color: var(--text-primary);
  transition: color 0.3s;
}
</style>
