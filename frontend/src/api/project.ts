import api from './index'

export interface Project {
  id: number
  name: string
  description: string
  status: string
  version: number
  creator_id: number
  created_at: string
  updated_at: string
  members?: ProjectMember[]
}

export interface ProjectMember {
  id: number
  user_id: number
  username: string
  email: string
  role: string
  avatar_url: string
}

export interface ProjectListResponse {
  list: Project[]
  total: number
  page: number
  page_size: number
}

export function createProject(data: { name: string; description: string }) {
  return api.post<Project>('/api/projects', data)
}

export function getProjects(params?: { page?: number; page_size?: number }) {
  return api.get<ProjectListResponse>('/api/projects', { params })
}

export function getProject(id: number | string) {
  return api.get<Project>(`/api/projects/${id}`)
}

export function updateProject(id: number | string, data: { name: string; description: string; version: number }) {
  return api.put(`/api/projects/${id}`, data)
}

export function submitProject(id: number | string) {
  return api.post(`/api/projects/${id}/submit`)
}

export function approveProject(id: number | string) {
  return api.post(`/api/projects/${id}/approve`)
}

export function startProject(id: number | string) {
  return api.post(`/api/projects/${id}/start`)
}

export function completeProject(id: number | string) {
  return api.post(`/api/projects/${id}/complete`)
}

export function archiveProject(id: number | string) {
  return api.post(`/api/projects/${id}/archive`)
}

export function addMember(projectId: number | string, data: { user_id: number }) {
  return api.post(`/api/projects/${projectId}/members`, data)
}

export function removeMember(projectId: number | string, userId: number) {
  return api.delete(`/api/projects/${projectId}/members/${userId}`)
}

export function getMembers(projectId: number | string) {
  return api.get<ProjectMember[]>(`/api/projects/${projectId}/members`)
}
