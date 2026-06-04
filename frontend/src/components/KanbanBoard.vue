<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as taskApi from '@/api/task'
import type { KanbanColumn as KanbanColumnType } from '@/api/task'
import KanbanColumn from './KanbanColumn.vue'
import TaskDetailModal from './TaskDetailModal.vue'
import { useTaskStore } from '@/stores/task'

const props = defineProps<{
  projectId: number
}>()

const columns = ref<KanbanColumnType[]>([])
const loading = ref(true)
const selectedTaskId = ref<number | null>(null)
const showTaskModal = ref(false)

const assigneeFilter = ref<number | undefined>(undefined)
const priorityFilter = ref<string | undefined>(undefined)

const taskStore = useTaskStore()

async function fetchKanban() {
  try {
    const res = await taskApi.getKanban(props.projectId, {
      assignee_id: assigneeFilter.value,
      priority: priorityFilter.value,
    })
    const kdata = res.data
    columns.value = kdata.columns

    // Cache all tasks
    columns.value.forEach((col) => {
      taskStore.cacheTasks(col.tasks)
    })
  } catch (err) {
    console.error('Failed to fetch kanban:', err)
  } finally {
    loading.value = false
  }
}

function handleTaskClick(taskId: number) {
  selectedTaskId.value = taskId
  showTaskModal.value = true
}

function handleCloseModal() {
  showTaskModal.value = false
  selectedTaskId.value = null
}

function handleTaskUpdated() {
  showTaskModal.value = false
  selectedTaskId.value = null
  fetchKanban()
}

let intervalId: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  fetchKanban()
  intervalId = setInterval(fetchKanban, 5000)
})

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId)
  }
})

watch([assigneeFilter, priorityFilter], () => {
  fetchKanban()
})

function getColumnTitle(status: string): string {
  return status
}
</script>

<template>
  <div class="kanban-board">
    <div class="kanban-toolbar">
      <h2 class="kanban-title">Kanban Board</h2>
      <div class="kanban-filters">
        <select v-model="priorityFilter" class="form-input filter-select">
          <option :value="undefined">All Priorities</option>
          <option value="Low">Low</option>
          <option value="Medium">Medium</option>
          <option value="High">High</option>
          <option value="Critical">Critical</option>
        </select>
      </div>
    </div>

    <div v-if="loading && columns.length === 0" class="kanban-loading">
      <div class="spinner"></div>
      <span>Loading board...</span>
    </div>

    <div v-else class="kanban-columns">
      <KanbanColumn
        v-for="col in columns"
        :key="col.status"
        :column="col"
        :title="getColumnTitle(col.status)"
        @task-click="handleTaskClick"
      />
    </div>

    <TaskDetailModal
      v-if="showTaskModal && selectedTaskId"
      :taskId="selectedTaskId"
      @close="handleCloseModal"
      @updated="handleTaskUpdated"
    />
  </div>
</template>

<style scoped>
.kanban-board {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  height: calc(100vh - 64px);
  overflow: hidden;
}

.kanban-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  flex-shrink: 0;
}

.kanban-title {
  font-size: 18px;
  font-weight: 600;
  color: #202124;
}

.kanban-filters {
  display: flex;
  gap: 12px;
}

.filter-select {
  width: 160px;
  height: 36px;
  font-size: 13px;
  padding: 0 12px;
  background: #fff;
}

.kanban-columns {
  flex: 1;
  display: flex;
  gap: 20px;
  padding: 0 24px 24px;
  overflow-x: auto;
  overflow-y: hidden;
  min-height: 0;
}

.kanban-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  flex: 1;
  color: #5f6368;
  font-size: 14px;
}
</style>
