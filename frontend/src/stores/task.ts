import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Task } from '@/api/task'

export const useTaskStore = defineStore('task', () => {
  const tasksCache = ref<Map<number, Task>>(new Map())

  function cacheTask(task: Task) {
    tasksCache.value.set(task.id, task)
  }

  function cacheTasks(tasks: Task[]) {
    tasks.forEach((t) => tasksCache.value.set(t.id, t))
  }

  function getTask(id: number): Task | undefined {
    return tasksCache.value.get(id)
  }

  function removeTask(id: number) {
    tasksCache.value.delete(id)
  }

  return {
    tasksCache,
    cacheTask,
    cacheTasks,
    getTask,
    removeTask,
  }
})
