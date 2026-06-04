import api from './index'

export interface Task {
  id: number
  title: string
  description: string
  priority: string
  status: string
  project_id: number
  assignee_id: number | null
  assignee_name: string
  deadline: string | null
  version: number
  created_at: string
  updated_at: string
}

export interface TaskListResponse {
  list: Task[]
  total: number
  page: number
  page_size: number
}

export interface KanbanResponse {
  columns: KanbanColumn[]
}

export interface KanbanColumn {
  status: string
  title: string
  tasks: Task[]
}

export function createTask(projectId: number | string, data: {
  title: string
  description: string
  priority: string
  project_id: number
  assignee_id?: number
  deadline?: string
}) {
  return api.post<Task>(`/api/projects/${projectId}/tasks`, data)
}

export function getTasks(projectId: number | string, params?: {
  status?: string
  priority?: string
  assignee_id?: number
  page?: number
  page_size?: number
}) {
  return api.get<TaskListResponse>(`/api/projects/${projectId}/tasks`, { params })
}

export function getTask(id: number | string) {
  return api.get<Task>(`/api/tasks/${id}`)
}

export function updateTask(id: number | string, data: {
  title?: string
  description?: string
  priority?: string
  assignee_id?: number | null
  deadline?: string | null
  version: number
}) {
  return api.put(`/api/tasks/${id}`, data)
}

export function deleteTask(id: number | string) {
  return api.delete(`/api/tasks/${id}`)
}

export function assignTask(id: number | string, data: { assignee_id: number }) {
  return api.put(`/api/tasks/${id}/assign`, data)
}

export function updateTaskStatus(id: number | string, data: { status: string }) {
  return api.put(`/api/tasks/${id}/status`, data)
}

export function getKanban(projectId: number | string, params?: { assignee_id?: number; priority?: string }) {
  return api.get<KanbanResponse>(`/api/projects/${projectId}/kanban`, { params })
}
