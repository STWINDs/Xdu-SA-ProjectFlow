<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useProjectStore } from '@/stores/project'
import { useAuthStore } from '@/stores/auth'
import * as projectApi from '@/api/project'
import type { Project, ProjectMember } from '@/api/project'
import AppHeader from '@/components/AppHeader.vue'
import CommentSection from '@/components/CommentSection.vue'
import KanbanBoard from '@/components/KanbanBoard.vue'

const route = useRoute()
const router = useRouter()
const projectStore = useProjectStore()
const authStore = useAuthStore()

const projectId = computed(() => Number(route.params.id))
const project = ref<Project | null>(null)
const members = ref<ProjectMember[]>([])
const loading = ref(true)
const error = ref('')

const showAddMember = ref(false)
const newMemberId = ref<number | null>(null)
const addMemberError = ref('')
const addMemberLoading = ref(false)

const showEditForm = ref(false)
const editName = ref('')
const editDesc = ref('')
const editError = ref('')
const editLoading = ref(false)

const showKanban = ref(false)

const username = authStore.user?.username || 'User'

const statusLabels: Record<string, string> = {
  draft: 'Draft',
  submitted: 'Submitted',
  approved: 'Approved',
  Developing: 'In Progress',
  completed: 'Completed',
  archived: 'Archived',
}

const isPM = computed(() => authStore.user?.role === 'pm' || authStore.user?.role === 'admin')

async function loadProject() {
  loading.value = true
  error.value = ''
  try {
    const data = await projectStore.fetchProjectDetail(projectId.value)
    project.value = data
    await fetchMembers()
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Failed to load project'
  } finally {
    loading.value = false
  }
}

async function fetchMembers() {
  try {
    const res = await projectApi.getMembers(projectId.value)
    members.value = res.data
  } catch (err) {
    console.error('Failed to fetch members:', err)
  }
}

async function handleAction(action: string) {
  try {
    switch (action) {
      case 'submit':
        await projectApi.submitProject(projectId.value)
        break
      case 'approve':
        await projectApi.approveProject(projectId.value)
        break
      case 'start':
        await projectApi.startProject(projectId.value)
        break
      case 'complete':
        await projectApi.completeProject(projectId.value)
        break
      case 'archive':
        await projectApi.archiveProject(projectId.value)
        break
    }
    await loadProject()
  } catch (err: any) {
    alert(err.response?.data?.message || 'Action failed')
  }
}

function openEdit() {
  if (!project.value) return
  editName.value = project.value.name
  editDesc.value = project.value.description
  showEditForm.value = true
}

async function handleEdit() {
  if (!project.value || !editName.value.trim()) return
  editLoading.value = true
  editError.value = ''
  try {
    await projectApi.updateProject(projectId.value, {
      name: editName.value.trim(),
      description: editDesc.value.trim(),
      version: project.value.version,
    })
    showEditForm.value = false
    await loadProject()
  } catch (err: any) {
    editError.value = err.response?.data?.message || 'Failed to update project'
  } finally {
    editLoading.value = false
  }
}

async function handleAddMember() {
  if (!newMemberId.value) {
    addMemberError.value = 'Please enter a user ID'
    return
  }
  addMemberLoading.value = true
  addMemberError.value = ''
  try {
    await projectApi.addMember(projectId.value, { user_id: newMemberId.value })
    newMemberId.value = null
    showAddMember.value = false
    await fetchMembers()
  } catch (err: any) {
    addMemberError.value = err.response?.data?.message || 'Failed to add member'
  } finally {
    addMemberLoading.value = false
  }
}

async function handleRemoveMember(userId: number) {
  if (!confirm('Remove this member?')) return
  try {
    await projectApi.removeMember(projectId.value, userId)
    await fetchMembers()
  } catch (err: any) {
    alert(err.response?.data?.message || 'Failed to remove member')
  }
}

onMounted(loadProject)
</script>

<template>
  <div class="project-detail-page">
    <AppHeader :username="username" />

    <div class="page-content">
      <!-- Loading -->
      <div v-if="loading" class="flex-center" style="padding:60px">
        <div class="spinner"></div>
      </div>

      <!-- Error -->
      <div v-else-if="error" class="card text-center" style="padding:60px 40px">
        <p style="font-size:14px;color:#ea4335;margin-bottom:16px">{{ error }}</p>
        <button class="btn btn-primary" @click="loadProject">Retry</button>
      </div>

      <template v-else-if="project">
        <!-- Back link -->
        <router-link to="/projects" class="back-link mb-2" style="display:inline-flex;align-items:center;gap:6px;font-size:13px">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
            <line x1="19" y1="12" x2="5" y2="12" />
            <polyline points="12 19 5 12 12 5" />
          </svg>
          Back to Projects
        </router-link>

        <!-- Project header card -->
        <div class="card mb-3">
          <div v-if="!showEditForm">
            <div class="flex-between mb-1">
              <h2 style="font-size:20px;font-weight:600">{{ project.name }}</h2>
              <div class="flex gap-1">
                <span class="tag" :class="{
                  'tag-info': project.status === 'submitted',
                  'tag-success': project.status === 'approved' || project.status === 'completed',
                  'tag-warning': project.status === 'Developing',
                }">
                  {{ statusLabels[project.status] || project.status }}
                </span>
                <button v-if="isPM" class="btn btn-flat btn-sm" @click="openEdit">Edit</button>
              </div>
            </div>
            <p style="font-size:14px;color:#3c4043;line-height:1.6;margin-bottom:16px">
              {{ project.description || 'No description' }}
            </p>
          </div>

          <div v-else>
            <h3 style="font-size:16px;font-weight:600;margin-bottom:16px">Edit Project</h3>
            <div v-if="editError" class="form-error mb-2">{{ editError }}</div>
            <div class="form-group">
              <label class="form-label">Name</label>
              <input v-model="editName" class="form-input" />
            </div>
            <div class="form-group">
              <label class="form-label">Description</label>
              <textarea v-model="editDesc" class="form-input" rows="3"></textarea>
            </div>
            <div class="flex gap-1">
              <button class="btn btn-primary btn-sm" @click="handleEdit" :disabled="editLoading">Save</button>
              <button class="btn btn-flat btn-sm" @click="showEditForm = false">Cancel</button>
            </div>
          </div>

          <!-- Status actions -->
          <div class="flex gap-1 mt-2" style="flex-wrap:wrap">
            <button
              v-if="project.status === 'draft'"
              class="btn btn-primary btn-sm"
              @click="handleAction('submit')"
            >
              Submit for Approval
            </button>
            <button
              v-if="isPM && project.status === 'submitted'"
              class="btn btn-primary btn-sm"
              @click="handleAction('approve')"
            >
              Approve
            </button>
            <button
              v-if="project.status === 'approved'"
              class="btn btn-primary btn-sm"
              @click="handleAction('start')"
            >
              Start Development
            </button>
            <button
              v-if="project.status === 'Developing'"
              class="btn btn-primary btn-sm"
              @click="handleAction('complete')"
            >
              Mark Complete
            </button>
            <button
              v-if="isPM && (project.status === 'completed' || project.status === 'draft')"
              class="btn btn-secondary btn-sm"
              @click="handleAction('archive')"
            >
              Archive
            </button>

            <button
              v-if="project.status === 'Developing'"
              class="btn btn-secondary btn-sm"
              @click="showKanban = !showKanban"
            >
              {{ showKanban ? 'Hide Kanban' : 'Show Kanban' }}
            </button>
          </div>
        </div>

        <!-- Kanban Board -->
        <div v-if="showKanban && project.status === 'Developing'" class="mb-3" style="overflow-x:auto">
          <KanbanBoard :projectId="projectId" />
        </div>

        <div class="detail-grid">
          <!-- Members -->
          <div class="card">
            <div class="flex-between mb-2">
              <h3 style="font-size:15px;font-weight:600">Team Members</h3>
              <button v-if="isPM" class="btn btn-flat btn-sm" @click="showAddMember = !showAddMember">
                + Add
              </button>
            </div>

            <div v-if="showAddMember" class="mb-2" style="background:#f6f8fc;padding:12px;border-radius:8px">
              <div v-if="addMemberError" class="form-error mb-1">{{ addMemberError }}</div>
              <div class="flex gap-1">
                <input
                  v-model="newMemberId"
                  type="number"
                  class="form-input"
                  placeholder="User ID"
                  style="height:36px;flex:1"
                />
                <button class="btn btn-primary btn-sm" @click="handleAddMember" :disabled="addMemberLoading">
                  Add
                </button>
              </div>
            </div>

            <div v-if="members.length === 0" style="font-size:13px;color:#9aa0a6">
              No members yet
            </div>

            <div v-for="m in members" :key="m.id" class="member-item flex-between">
              <div class="flex gap-1">
                <div class="avatar avatar-sm">{{ m.username[0].toUpperCase() }}</div>
                <div>
                  <div style="font-size:13px;font-weight:500">{{ m.username }}</div>
                  <div style="font-size:11px;color:#9aa0a6">{{ m.email }}</div>
                </div>
              </div>
              <button
                v-if="isPM"
                class="btn btn-flat btn-sm"
                style="color:#ea4335"
                @click="handleRemoveMember(m.user_id)"
              >
                Remove
              </button>
            </div>
          </div>

          <!-- Comments -->
          <div class="card">
            <h3 style="font-size:15px;font-weight:600;margin-bottom:12px">Comments</h3>
            <CommentSection type="project" :targetId="projectId" />
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.project-detail-page {
  min-height: 100vh;
  background: #f6f8fc;
}

.page-content {
  max-width: 1100px;
  margin: 0 auto;
  padding: 24px;
}

.back-link:hover {
  text-decoration: none;
  color: #1557b0;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

@media (max-width: 768px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
}

.member-item {
  padding: 8px 0;
  border-bottom: 1px solid #f1f3f4;
}

.member-item:last-child {
  border-bottom: none;
}
</style>
