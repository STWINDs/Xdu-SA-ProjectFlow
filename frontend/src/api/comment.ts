import api from './index'

export interface Comment {
  id: number
  content: string
  user_id: number
  username: string
  avatar_url: string
  parent_id: number | null
  replies: Comment[]
  created_at: string
  updated_at: string
}

export interface CommentListResponse {
  list: Comment[]
  total: number
  page: number
  page_size: number
}

export function getProjectComments(projectId: number | string, params?: { page?: number; page_size?: number }) {
  return api.get<CommentListResponse>(`/api/projects/${projectId}/comments`, { params })
}

export function createProjectComment(projectId: number | string, data: { content: string }) {
  return api.post(`/api/projects/${projectId}/comments`, data)
}

export function getTaskComments(taskId: number | string) {
  return api.get<Comment[]>(`/api/tasks/${taskId}/comments`)
}

export function createTaskComment(taskId: number | string, data: { content: string }) {
  return api.post(`/api/tasks/${taskId}/comments`, data)
}

export function replyToComment(commentId: number | string, data: { content: string }) {
  return api.post(`/api/comments/${commentId}/reply`, data)
}

export function deleteComment(id: number | string) {
  return api.delete(`/api/comments/${id}`)
}
