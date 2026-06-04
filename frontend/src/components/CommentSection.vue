<script setup lang="ts">
import { ref, onMounted } from 'vue'
import * as commentApi from '@/api/comment'
import type { Comment } from '@/api/comment'
import { useAuthStore } from '@/stores/auth'

const props = defineProps<{
  type: 'project' | 'task'
  targetId: number
}>()

const comments = ref<Comment[]>([])
const newComment = ref('')
const replyContent = ref<Record<number, string>>({})
const replyOpen = ref<Record<number, boolean>>({})
const loading = ref(false)
const authStore = useAuthStore()

async function fetchComments() {
  loading.value = true
  try {
    if (props.type === 'project') {
      const res = await commentApi.getProjectComments(props.targetId)
      const pdata = res.data
      comments.value = pdata.list || pdata
    } else {
      const res = await commentApi.getTaskComments(props.targetId)
      const data: any = res.data
      comments.value = Array.isArray(data) ? data : (data.list || [])
    }
  } catch (err) {
    console.error('Failed to fetch comments:', err)
  } finally {
    loading.value = false
  }
}

async function submitComment() {
  if (!newComment.value.trim()) return
  try {
    if (props.type === 'project') {
      await commentApi.createProjectComment(props.targetId, { content: newComment.value })
    } else {
      await commentApi.createTaskComment(props.targetId, { content: newComment.value })
    }
    newComment.value = ''
    await fetchComments()
  } catch (err) {
    console.error('Failed to submit comment:', err)
  }
}

async function submitReply(commentId: number) {
  const content = replyContent.value[commentId]
  if (!content?.trim()) return
  try {
    await commentApi.replyToComment(commentId, { content })
    replyContent.value[commentId] = ''
    replyOpen.value[commentId] = false
    await fetchComments()
  } catch (err) {
    console.error('Failed to submit reply:', err)
  }
}

function toggleReply(commentId: number) {
  replyOpen.value[commentId] = !replyOpen.value[commentId]
}

async function handleDelete(commentId: number) {
  if (!confirm('Delete this comment?')) return
  try {
    await commentApi.deleteComment(commentId)
    await fetchComments()
  } catch (err) {
    console.error('Failed to delete comment:', err)
  }
}

function timeAgo(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  if (diffMins < 1) return 'Just now'
  if (diffMins < 60) return `${diffMins}m ago`
  const diffHours = Math.floor(diffMins / 60)
  if (diffHours < 24) return `${diffHours}h ago`
  const diffDays = Math.floor(diffHours / 24)
  if (diffDays < 30) return `${diffDays}d ago`
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

onMounted(fetchComments)
</script>

<template>
  <div class="comment-section">
    <div class="comment-input-area">
      <textarea
        v-model="newComment"
        class="form-input comment-textarea"
        placeholder="Write a comment..."
        rows="3"
      ></textarea>
      <button
        class="btn btn-primary btn-sm mt-1"
        :disabled="!newComment.trim()"
        @click="submitComment"
      >
        Comment
      </button>
    </div>

    <div v-if="loading" class="comment-loading text-center mt-2">
      <div class="spinner"></div>
    </div>

    <div v-else class="comment-list mt-2">
      <div v-for="comment in comments" :key="comment.id" class="comment-item">
        <div class="comment-avatar avatar avatar-sm">
          {{ (comment.username || 'U')[0].toUpperCase() }}
        </div>
        <div class="comment-body">
          <div class="comment-header">
            <span class="comment-author">{{ comment.username }}</span>
            <span class="comment-time">{{ timeAgo(comment.created_at) }}</span>
          </div>
          <div class="comment-content">{{ comment.content }}</div>
          <div class="comment-actions">
            <button class="btn btn-flat btn-sm" @click="toggleReply(comment.id)">
              Reply
            </button>
            <button
              v-if="authStore.user?.id === comment.user_id"
              class="btn btn-flat btn-sm"
              style="color: #ea4335"
              @click="handleDelete(comment.id)"
            >
              Delete
            </button>
          </div>

          <!-- Reply input -->
          <div v-if="replyOpen[comment.id]" class="reply-input-area mt-1">
            <textarea
              v-model="replyContent[comment.id]"
              class="form-input comment-textarea"
              placeholder="Write a reply..."
              rows="2"
            ></textarea>
            <div class="flex gap-1 mt-1">
              <button
                class="btn btn-primary btn-sm"
                :disabled="!replyContent[comment.id]?.trim()"
                @click="submitReply(comment.id)"
              >
                Reply
              </button>
              <button class="btn btn-flat btn-sm" @click="toggleReply(comment.id)">
                Cancel
              </button>
            </div>
          </div>

          <!-- Replies -->
          <div v-if="comment.replies && comment.replies.length > 0" class="replies mt-2">
            <div v-for="reply in comment.replies" :key="reply.id" class="comment-item reply-item">
              <div class="comment-avatar avatar avatar-sm">
                {{ (reply.username || 'U')[0].toUpperCase() }}
              </div>
              <div class="comment-body">
                <div class="comment-header">
                  <span class="comment-author">{{ reply.username }}</span>
                  <span class="comment-time">{{ timeAgo(reply.created_at) }}</span>
                </div>
                <div class="comment-content">{{ reply.content }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="comments.length === 0" class="text-center mt-2" style="color: #9aa0a6; font-size: 13px;">
        No comments yet. Be the first to comment.
      </div>
    </div>
  </div>
</template>

<style scoped>
.comment-section {
  margin-top: 8px;
}

.comment-textarea {
  min-height: 60px;
  font-size: 13px;
}

.comment-input-area {
  margin-bottom: 8px;
}

.comment-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.comment-item {
  display: flex;
  gap: 10px;
  padding: 10px 0;
  border-bottom: 1px solid #f1f3f4;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-avatar {
  flex-shrink: 0;
}

.comment-body {
  flex: 1;
  min-width: 0;
}

.comment-header {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 4px;
}

.comment-author {
  font-size: 13px;
  font-weight: 600;
  color: #202124;
}

.comment-time {
  font-size: 11px;
  color: #9aa0a6;
}

.comment-content {
  font-size: 13px;
  color: #3c4043;
  line-height: 1.5;
  white-space: pre-wrap;
  word-break: break-word;
}

.comment-actions {
  margin-top: 4px;
}

.replies {
  margin-left: 4px;
  padding-left: 16px;
  border-left: 2px solid #e8eaed;
}

.reply-item {
  border-bottom: 1px solid #f1f3f4;
  padding: 8px 0;
}

.reply-input-area {
  background: #f6f8fc;
  border-radius: 8px;
  padding: 10px;
}
</style>
