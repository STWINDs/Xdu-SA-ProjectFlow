import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as projectApi from '@/api/project'
import type { Project, ProjectMember } from '@/api/project'

export const useProjectStore = defineStore('project', () => {
  const projects = ref<Project[]>([])
  const currentProject = ref<Project | null>(null)
  const members = ref<ProjectMember[]>([])
  const total = ref(0)
  const loading = ref(false)

  async function fetchProjects(page = 1, pageSize = 20) {
    loading.value = true
    try {
      const res = await projectApi.getProjects({ page, page_size: pageSize })
      const pdata = res.data
      projects.value = pdata.list
      total.value = pdata.total
    } catch (err) {
      console.error('Failed to fetch projects:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchProjectDetail(id: number | string) {
    loading.value = true
    try {
      const res = await projectApi.getProject(id)
      const detail = res.data
      currentProject.value = detail
      return detail
    } catch (err) {
      console.error('Failed to fetch project detail:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchMembers(projectId: number | string) {
    try {
      const res = await projectApi.getMembers(projectId)
      members.value = res.data
    } catch (err) {
      console.error('Failed to fetch members:', err)
    }
  }

  return {
    projects,
    currentProject,
    members,
    total,
    loading,
    fetchProjects,
    fetchProjectDetail,
    fetchMembers,
  }
})
