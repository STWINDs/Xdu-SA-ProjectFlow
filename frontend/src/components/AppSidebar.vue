<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project'

// Props can be event handlers from parent if needed
const emit = defineEmits<{
  selectProject: [id: number]
  newProject: []
}>()

const projectStore = useProjectStore()
const route = useRoute()
const router = useRouter()

const currentProjectId = computed(() => {
  const id = route.params.id
  if (id) return Number(id)
  // If on dashboard without a specific project, use first project
  if (projectStore.projects.length > 0 && route.name === 'Dashboard') {
    return projectStore.projects[0].id
  }
  return null
})

function selectProject(id: number) {
  router.push(`/projects/${id}`)
  emit('selectProject', id)
}

function handleNewProject() {
  router.push('/projects')
  emit('newProject')
}

function getStatusColor(status: string): string {
  const colors: Record<string, string> = {
    draft: '#9aa0a6',
    submitted: '#1a73e8',
    approved: '#1e8e3e',
    Developing: '#f9ab00',
    completed: '#1e8e3e',
    archived: '#5f6368',
  }
  return colors[status] || '#9aa0a6'
}

const projectMenuItems = computed(() =>
  projectStore.projects.map((p) => ({
    id: p.id,
    label: p.name,
    active: currentProjectId.value === p.id,
  }))
)
</script>

<template>
  <aside class="app-sidebar">
    <nav class="sidebar-nav">
      <router-link to="/" class="sidebar-item" :class="{ active: route.name === 'Dashboard' }">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
          <rect x="3" y="3" width="7" height="7" />
          <rect x="14" y="3" width="7" height="7" />
          <rect x="3" y="14" width="7" height="7" />
          <rect x="14" y="14" width="7" height="7" />
        </svg>
        <span>Dashboard</span>
      </router-link>

      <div class="sidebar-section-label">Projects</div>

      <router-link
        v-for="p in projectStore.projects"
        :key="p.id"
        :to="`/projects/${p.id}`"
        class="sidebar-item"
        :class="{ active: currentProjectId === p.id }"
      >
        <span class="project-dot" :style="{ background: getStatusColor(p.status) }"></span>
        <span class="project-name">{{ p.name }}</span>
      </router-link>

      <router-link to="/projects" class="sidebar-item" :class="{ active: route.name === 'ProjectList' }">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
          <circle cx="12" cy="12" r="10" />
          <line x1="12" y1="8" x2="12" y2="16" />
          <line x1="8" y1="12" x2="16" y2="12" />
        </svg>
        <span>All Projects</span>
      </router-link>
    </nav>

    <div class="sidebar-footer">
      <button class="btn btn-primary btn-block" @click="handleNewProject">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
          <line x1="12" y1="5" x2="12" y2="19" />
          <line x1="5" y1="12" x2="19" y2="12" />
        </svg>
        New Project
      </button>
    </div>
  </aside>
</template>

<style scoped>
.app-sidebar {
  width: 260px;
  min-width: 260px;
  height: calc(100vh - 64px);
  background: #fff;
  border-right: 1px solid #e8eaed;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  position: sticky;
  top: 64px;
}

.sidebar-nav {
  flex: 1;
  padding: 16px 12px;
}

.sidebar-section-label {
  font-size: 11px;
  font-weight: 600;
  color: #9aa0a6;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  padding: 16px 12px 8px;
}

.sidebar-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  margin-bottom: 2px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  color: #3c4043;
  text-decoration: none;
  transition: background 0.15s, color 0.15s;
  cursor: pointer;
}

.sidebar-item:hover {
  background: #f1f3f4;
  text-decoration: none;
}

.sidebar-item.active {
  background: #e8f0fe;
  color: #1a73e8;
}

.project-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.project-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sidebar-footer {
  padding: 16px 12px;
  border-top: 1px solid #e8eaed;
}

@media (max-width: 768px) {
  .app-sidebar {
    width: 200px;
    min-width: 200px;
  }
}
</style>
