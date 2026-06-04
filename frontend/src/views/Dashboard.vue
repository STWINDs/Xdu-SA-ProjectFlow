<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useProjectStore } from '@/stores/project'
import { useAuthStore } from '@/stores/auth'
import AppHeader from '@/components/AppHeader.vue'
import AppSidebar from '@/components/AppSidebar.vue'
import KanbanBoard from '@/components/KanbanBoard.vue'

const projectStore = useProjectStore()
const authStore = useAuthStore()
const selectedProjectId = ref<number | null>(null)
const sidebarVisible = ref(true)

const username = computed(() => authStore.user?.username || 'User')
const avatarUrl = computed(() => authStore.user?.avatar_url || '')

onMounted(async () => {
  await projectStore.fetchProjects()
  if (projectStore.projects.length > 0) {
    selectedProjectId.value = projectStore.projects[0].id
  }
})

function handleSelectProject(id: number) {
  selectedProjectId.value = id
}
</script>

<template>
  <div class="dashboard">
    <AppHeader :username="username" :avatarUrl="avatarUrl" />

    <div class="dashboard-body">
      <AppSidebar @select-project="handleSelectProject" />

      <main class="dashboard-main">
        <KanbanBoard
          v-if="selectedProjectId"
          :key="selectedProjectId"
          :projectId="selectedProjectId"
        />
        <div v-else class="dashboard-empty">
          <div class="card" style="text-align:center;padding:60px 40px;">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="#9aa0a6" stroke-width="1.5" stroke-linecap="round" style="margin-bottom:16px">
              <rect x="3" y="3" width="7" height="7" />
              <rect x="14" y="3" width="7" height="7" />
              <rect x="3" y="14" width="7" height="7" />
              <rect x="14" y="14" width="7" height="7" />
            </svg>
            <h3 style="font-size:16px;font-weight:600;color:#202124;margin-bottom:8px">No Project Selected</h3>
            <p style="font-size:13px;color:#5f6368;margin-bottom:20px">
              Select a project from the sidebar or create a new one to get started.
            </p>
            <router-link to="/projects" class="btn btn-primary">Browse Projects</router-link>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<style scoped>
.dashboard {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f6f8fc;
}

.dashboard-body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.dashboard-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.dashboard-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}
</style>
