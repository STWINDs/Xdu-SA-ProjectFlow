import api from './index'

export interface Attachment {
  id: number
  filename: string
  original_name: string
  file_size: number
  mime_type: string
  url: string
  user_id: number
  task_id: number | null
  project_id: number | null
  created_at: string
}

export function uploadAttachment(file: File, taskId?: number | string, projectId?: number | string) {
  const formData = new FormData()
  formData.append('file', file)
  if (taskId) formData.append('task_id', String(taskId))
  if (projectId) formData.append('project_id', String(projectId))
  return api.post<Attachment>('/api/attachments/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export function getTaskAttachments(taskId: number | string) {
  return api.get<Attachment[]>(`/api/tasks/${taskId}/attachments`)
}

export function getProjectAttachments(projectId: number | string) {
  return api.get<Attachment[]>(`/api/projects/${projectId}/attachments`)
}

export function deleteAttachment(id: number | string) {
  return api.delete(`/api/attachments/${id}`)
}
