<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import * as taskApi from '@/api/task'
import type { Task } from '@/api/task'
import CommentSection from './CommentSection.vue'
import { useTaskStore } from '@/stores/task'
import { useAuthStore } from '@/stores/auth'

const props = defineProps<{
  taskId: number
}>()

const emit = defineEmits<{
  close: []
  updated: []
}>()

const task = ref<Task | null>(null)
const loading = ref(true)
const editMode = ref(false)
const authStore = useAuthStore()

const editForm = ref({
  title: '',
  description: '',
  priority: '',
  status: '',
  assignee_id: null as number | null,
  deadline: '',
})

async function fetchTask() {
  loading.value = true
  try {
    const res = await taskApi.getTask(props.taskId)
    const tdata = res.data
    task.value = tdata
    // Populate edit form
    editForm.value = {
      title: tdata.title,
      description: tdata.description,
      priority: tdata.priority,
      status: tdata.status,
      assignee_id: tdata.assignee_id,
      deadline: tdata.deadline ? tdata.deadline.substring(0, 10) : '',
    }
  } catch (err) {
    console.error('Failed to fetch task:', err)
  } finally {
    loading.value = false
  }
}

async function saveTask() {
  if (!task.value) return
  try {
    await taskApi.updateTask(task.value.id, {
      title: editForm.value.title,
      description: editForm.value.description,
      priority: editForm.value.priority,
      assignee_id: editForm.value.assignee_id,
      deadline: editForm.value.deadline ? new Date(editForm.value.deadline).toISOString() : null,
      version: task.value.version,
    })
    editMode.value = false
    await fetchTask()
    emit('updated')
  } catch (err) {
    console.error('Failed to update task:', err)
    alert('Failed to update task. Please try again.')
  }
}

async function changeStatus(status: string) {
  if (!task.value) return
  try {
    await taskApi.updateTaskStatus(task.value.id, { status })
    await fetchTask()
    emit('updated')
  } catch (err) {
    console.error('Failed to update status:', err)
  }
}

async function deleteCurrentTask() {
  if (!task.value) return
  if (!confirm('Delete this task?')) return
  try {
    await taskApi.deleteTask(task.value.id)
    emit('updated')
  } catch (err) {
    console.error('Failed to delete task:', err)
  }
}

function priorityLabel(p: string): string {
  return p
}

const statusFlow: Record<string, string> = {
  Todo: 'InProgress',
  InProgress: 'Review',
  Review: 'Testing',
  Testing: 'Done',
}

const nextStatus = computed(() => {
  if (!task.value) return null
  return statusFlow[task.value.status] || null
})

const statusLabels: Record<string, string> = {
  Todo: 'To Do',
  InProgress: 'In Progress',
  Review: 'Review',
  Testing: 'Testing',
  Done: 'Done',
}

watch(() => props.taskId, fetchTask, { immediate: true })
</script>

<template>
  <div class="modal-overlay" @click.self="emit('close')">
    <div class="modal-content task-detail-modal">
      <div v-if="loading" class="flex-center" style="padding: 40px">
        <div class="spinner"></div>
      </div>

      <template v-else-if="task">
        <div class="modal-header">
          <h3 v-if="!editMode">{{ task.title }}</h3>
          <input v-else v-model="editForm.title" class="form-input" style="flex:1" />
          <button class="modal-close" @click="emit('close')">&#x2715;</button>
        </div>

        <div class="task-body">
          <!-- Status bar -->
          <div class="task-status-bar">
            <span class="form-label" style="margin-bottom:0">Status:</span>
            <span class="tag" :class="`priority-${(task.status || '').toLowerCase()}`">{{ statusLabels[task.status] || task.status }}</span>
            <button
              v-if="nextStatus"
              class="btn btn-sm btn-primary"
              style="margin-left:12px"
              @click="changeStatus(nextStatus)"
            >
              Move to {{ statusLabels[nextStatus] }}
            </button>
          </div>

          <!-- Edit toggle -->
          <div class="flex-between mt-2">
            <div class="flex gap-1">
              <span class="tag" :class="`priority-${(task.priority || '').toLowerCase()}`">
                {{ priorityLabel(task.priority) }}
              </span>
              <span v-if="task.assignee_name" class="tag tag-info">
                {{ task.assignee_name }}
              </span>
            </div>
            <div class="flex gap-1">
              <button v-if="!editMode" class="btn btn-flat btn-sm" @click="editMode = true">
                Edit
              </button>
              <button
                v-if="authStore.user?.role === 'pm' || authStore.user?.role === 'admin'"
                class="btn btn-flat btn-sm"
                style="color: #ea4335"
                @click="deleteCurrentTask"
              >
                Delete
              </button>
            </div>
          </div>

          <!-- Description -->
          <div class="task-section mt-2">
            <div class="form-label">Description</div>
            <p v-if="!editMode" class="task-description">{{ task.description || 'No description' }}</p>
            <textarea
              v-else
              v-model="editForm.description"
              class="form-input"
              rows="4"
              placeholder="Task description..."
            ></textarea>
          </div>

          <!-- Edit fields -->
          <div v-if="editMode" class="edit-fields mt-2">
            <div class="form-group">
              <label class="form-label">Priority</label>
              <select v-model="editForm.priority" class="form-input" style="height:40px">
                <option value="Low">Low</option>
                <option value="Medium">Medium</option>
                <option value="High">High</option>
                <option value="Critical">Critical</option>
              </select>
            </div>

            <div class="form-group">
              <label class="form-label">Deadline</label>
              <input v-model="editForm.deadline" type="date" class="form-input" style="height:40px" />
            </div>

            <div class="flex gap-1 mt-2">
              <button class="btn btn-primary btn-sm" @click="saveTask">Save Changes</button>
              <button class="btn btn-flat btn-sm" @click="editMode = false">Cancel</button>
            </div>
          </div>

          <!-- Meta info -->
          <div class="task-meta mt-2" v-if="!editMode">
            <span v-if="task.deadline" class="form-label" style="display:inline">
              Deadline: {{ new Date(task.deadline).toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' }) }}
            </span>
            <span class="form-label" style="display:inline; color:#9aa0a6">
              Version: {{ task.version }}
            </span>
          </div>

          <!-- Attachments placeholder -->
          <div class="task-section mt-2">
            <div class="form-label">Attachments</div>
            <p style="font-size:13px;color:#9aa0a6">Attachments will be shown here</p>
          </div>

          <!-- Comments -->
          <div class="task-section mt-3">
            <div class="form-label">Comments</div>
            <CommentSection type="task" :targetId="task.id" />
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.task-detail-modal {
  max-width: 680px;
}

.task-body {
  overflow-y: auto;
}

.task-status-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 0;
}

.status-actions {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.task-description {
  font-size: 14px;
  color: #3c4043;
  line-height: 1.6;
  white-space: pre-wrap;
}

.task-section {
  border-top: 1px solid #f1f3f4;
  padding-top: 14px;
}

.task-meta {
  display: flex;
  gap: 16px;
}
</style>
