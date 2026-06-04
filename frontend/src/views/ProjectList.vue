<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import { useAuthStore } from '@/stores/auth'
import * as projectApi from '@/api/project'
import AppHeader from '@/components/AppHeader.vue'

const router = useRouter()
const projectStore = useProjectStore()
const authStore = useAuthStore()
const loading = ref(false)

const showCreateForm = ref(false)
const newProjectName = ref('')
const newProjectDesc = ref('')
const createError = ref('')
const createLoading = ref(false)

const username = authStore.user?.username || 'User'

async function loadProjects() {
  await projectStore.fetchProjects()
}

async function handleCreateProject() {
  if (!newProjectName.value.trim()) {
    createError.value = 'Project name is required'
    return
  }
  createLoading.value = true
  createError.value = ''
  try {
    await projectApi.createProject({
      name: newProjectName.value.trim(),
      description: newProjectDesc.value.trim(),
    })
    showCreateForm.value = false
    newProjectName.value = ''
    newProjectDesc.value = ''
    await loadProjects()
  } catch (err: any) {
    createError.value = err.response?.data?.message || 'Failed to create project'
  } finally {
    createLoading.value = false
  }
}

function goToProject(id: number) {
  router.push(`/projects/${id}`)
}

const statusLabels: Record<string, string> = {
  draft: 'Draft',
  submitted: 'Submitted',
  approved: 'Approved',
  Developing: 'In Progress',
  completed: 'Completed',
  archived: 'Archived',
}

const statusTagClass: Record<string, string> = {
  draft: 'tag',
  submitted: 'tag tag-info',
  approved: 'tag tag-success',
  Developing: 'tag tag-warning',
  completed: 'tag tag-success',
  archived: 'tag',
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

onMounted(loadProjects)
</script>

<template>
  <div class="project-list-page">
    <AppHeader :username="username" />

    <div class="page-content">
      <div class="page-toolbar">
        <h2>Projects</h2>
        <button
          v-if="authStore.user?.role === 'pm' || authStore.user?.role === 'admin'"
          class="btn btn-primary"
          @click="showCreateForm = !showCreateForm"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round">
            <line x1="12" y1="5" x2="12" y2="19" />
            <line x1="5" y1="12" x2="19" y2="12" />
          </svg>
          New Project
        </button>
      </div>

      <!-- Create form -->
      <div v-if="showCreateForm" class="card mb-3">
        <h3 style="font-size:16px;font-weight:600;margin-bottom:16px">Create New Project</h3>
        <div v-if="createError" class="form-error mb-2">{{ createError }}</div>
        <div class="form-group">
          <label class="form-label">Project Name</label>
          <input v-model="newProjectName" class="form-input" placeholder="Enter project name" />
        </div>
        <div class="form-group">
          <label class="form-label">Description</label>
          <textarea v-model="newProjectDesc" class="form-input" placeholder="Enter description (optional)" rows="3"></textarea>
        </div>
        <div class="flex gap-1">
          <button class="btn btn-primary" @click="handleCreateProject" :disabled="createLoading">
            Create
          </button>
          <button class="btn btn-flat" @click="showCreateForm = false">Cancel</button>
        </div>
      </div>

      <!-- Project grid -->
      <div v-if="projectStore.loading" class="flex-center" style="padding:40px">
        <div class="spinner"></div>
      </div>

      <div v-else class="project-grid">
        <div
          v-for="project in projectStore.projects"
          :key="project.id"
          class="card project-card"
          @click="goToProject(project.id)"
        >
          <div class="flex-between mb-1">
            <h3 class="project-card-name">{{ project.name }}</h3>
            <span :class="statusTagClass[project.status] || 'tag'">
              {{ statusLabels[project.status] || project.status }}
            </span>
          </div>
          <p class="project-card-desc">{{ project.description || 'No description' }}</p>
          <div class="project-card-meta">
            <span>{{ project.members?.length || 0 }} members</span>
            <span>Created {{ formatDate(project.created_at) }}</span>
          </div>
        </div>

        <div v-if="projectStore.projects.length === 0" class="card text-center" style="padding:60px 40px">
          <p style="font-size:14px;color:#5f6368;margin-bottom:16px">No projects yet</p>
          <button
            v-if="authStore.user?.role === 'pm' || authStore.user?.role === 'admin'"
            class="btn btn-primary"
            @click="showCreateForm = true"
          >
            Create your first project
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.project-list-page {
  min-height: 100vh;
  background: #f6f8fc;
}

.page-content {
  max-width: 960px;
  margin: 0 auto;
  padding: 24px;
}

.page-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
}

.page-toolbar h2 {
  font-size: 22px;
  font-weight: 600;
  color: #202124;
}

.project-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.project-card {
  cursor: pointer;
  transition: box-shadow 0.2s, transform 0.15s;
  padding: 20px;
}

.project-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.project-card-name {
  font-size: 16px;
  font-weight: 600;
  color: #202124;
}

.project-card-desc {
  font-size: 13px;
  color: #5f6368;
  line-height: 1.5;
  margin-bottom: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.project-card-meta {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: #9aa0a6;
}
</style>
