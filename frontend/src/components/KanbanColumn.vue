<script setup lang="ts">
import type { KanbanColumn, Task } from '@/api/task'
import KanbanCard from './KanbanCard.vue'

defineProps<{
  column: KanbanColumn
  title: string
}>()

const emit = defineEmits<{
  taskClick: [id: number]
}>()
</script>

<template>
  <div class="kanban-column">
    <div class="column-header">
      <div class="column-title-row">
        <span class="column-title">{{ title }}</span>
        <span class="column-count">{{ column.tasks.length }}</span>
      </div>
    </div>

    <div class="column-body">
      <KanbanCard
        v-for="task in column.tasks"
        :key="task.id"
        :task="task"
        @click="emit('taskClick', task.id)"
      />

      <div v-if="column.tasks.length === 0" class="column-empty">
        No tasks
      </div>
    </div>
  </div>
</template>

<style scoped>
.kanban-column {
  min-width: 300px;
  width: 300px;
  max-width: 340px;
  background: #eef2ff;
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.column-header {
  padding: 16px 16px 8px;
  flex-shrink: 0;
}

.column-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.column-title {
  font-size: 14px;
  font-weight: 600;
  color: #202124;
}

.column-count {
  font-size: 12px;
  font-weight: 500;
  color: #5f6368;
  background: #fff;
  padding: 1px 8px;
  border-radius: 10px;
}

.column-body {
  flex: 1;
  overflow-y: auto;
  padding: 4px 12px 12px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.column-empty {
  padding: 24px 0;
  text-align: center;
  font-size: 13px;
  color: #9aa0a6;
}
</style>
