<script setup lang="ts">
import { computed } from 'vue'
import type { Task } from '@/api/task'

const props = defineProps<{
  task: Task
}>()

defineEmits<{
  click: [id: number]
}>()

const priorityClass = computed(() => `priority-${(props.task.priority || '').toLowerCase()}`)

const priorityLabel = computed(() => props.task.priority || '')

const avatarLetter = computed(() => {
  if (props.task.assignee_name) {
    return props.task.assignee_name[0].toUpperCase()
  }
  return '?'
})

const formattedDeadline = computed(() => {
  if (!props.task.deadline) return ''
  const d = new Date(props.task.deadline)
  const now = new Date()
  const diffDays = Math.ceil((d.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  if (diffDays < 0) return 'Overdue'
  if (diffDays === 0) return 'Today'
  if (diffDays === 1) return 'Tomorrow'
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
})

const deadlineClass = computed(() => {
  if (!props.task.deadline) return ''
  const d = new Date(props.task.deadline)
  const now = new Date()
  const diffDays = Math.ceil((d.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
  if (diffDays < 0) return 'deadline-overdue'
  if (diffDays <= 1) return 'deadline-soon'
  return ''
})
</script>

<template>
  <div
    class="kanban-card"
    @click="$emit('click', task.id)"
  >
    <div class="card-top">
      <span class="tag" :class="priorityClass">{{ priorityLabel }}</span>
    </div>
    <h4 class="card-title">{{ task.title }}</h4>
    <div class="card-bottom">
      <div class="card-meta">
        <div class="avatar avatar-sm" :title="task.assignee_name">
          {{ avatarLetter }}
        </div>
        <span v-if="task.deadline" class="card-deadline" :class="deadlineClass">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
            <line x1="16" y1="2" x2="16" y2="6" />
            <line x1="8" y1="2" x2="8" y2="6" />
            <line x1="3" y1="10" x2="21" y2="10" />
          </svg>
          {{ formattedDeadline }}
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.kanban-card {
  background: #fff;
  border-radius: 12px;
  padding: 14px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  cursor: pointer;
  transition: box-shadow 0.2s, transform 0.15s;
}

.kanban-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
  transform: translateY(-1px);
}

.card-top {
  margin-bottom: 8px;
}

.card-title {
  font-size: 14px;
  font-weight: 500;
  color: #202124;
  line-height: 1.4;
  margin-bottom: 10px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-bottom {
  display: flex;
  align-items: center;
}

.card-meta {
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-deadline {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #5f6368;
}

.deadline-soon {
  color: #e65100;
  font-weight: 500;
}

.deadline-overdue {
  color: #ea4335;
  font-weight: 500;
}
</style>
