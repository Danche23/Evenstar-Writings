<template>
  <div>
    <div class="admin-page-head">
      <h2>用户管理</h2>
    </div>

    <!-- 搜索 -->
    <div class="toolbar-row">
      <el-input v-model="keyword" placeholder="搜索用户名 / 昵称 / 邮箱…" clearable style="width: 260px"
        @keyup.enter="reload" @clear="reload" />
      <button class="btn btn-primary btn-sm" @click="reload">搜索</button>
    </div>

    <div class="admin-card">
      <el-table :data="list" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="用户" min-width="150">
          <template #default="{ row }">
            <div class="u-cell">
              <el-avatar :size="30" :src="row.avatar || undefined">{{ (row.nickname || row.username || '?').slice(0, 1) }}</el-avatar>
              <div>
                <div class="u-name">{{ row.nickname || row.username }}</div>
                <div class="faint small">@{{ row.username }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="180" show-overflow-tooltip />
        <el-table-column label="角色" width="100">
          <template #default="{ row }">
            <span class="badge" :class="row.role === 1 ? 'badge-primary' : 'badge-gray'">
              {{ row.role === 1 ? '管理员' : '用户' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <span class="badge" :class="row.status === 1 ? 'badge-success' : 'badge-danger'">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="注册时间" width="130">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <template v-if="row.id === store.user?.id">
              <span class="faint small">（自己）</span>
            </template>
            <template v-else>
              <span class="op-link" @click="openEdit(row)">编辑</span>
              <el-popconfirm
                :title="row.status === 1 ? '禁用后该用户将无法登录，确定？' : '确定解禁该用户？'"
                width="230"
                @confirm="toggleStatus(row)"
              >
                <template #reference>
                  <span class="op-link">{{ row.status === 1 ? '禁用' : '解禁' }}</span>
                </template>
              </el-popconfirm>
              <el-popconfirm title="删除后该用户数据将软删除，确定？" width="220" @confirm="remove(row)">
                <template #reference><span class="op-link danger">删除</span></template>
              </el-popconfirm>
            </template>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="pager">
      <el-pagination
        layout="prev, pager, next, total"
        :total="total"
        :page-size="query.page_size"
        :current-page="query.page"
        @current-change="(p) => { query.page = p; load() }"
      />
    </div>

    <el-dialog v-model="editVisible" title="编辑用户资料" width="420px">
      <el-form label-width="56px">
        <el-form-item label="昵称">
          <el-input v-model="editForm.nickname" maxlength="50" placeholder="展示昵称" />
        </el-form-item>
        <el-form-item label="简介">
          <el-input v-model="editForm.bio" type="textarea" :rows="3" maxlength="500" placeholder="博主简介，将展示在前台关于卡片" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" :disabled="editSaving" @click="saveEdit">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { adminListUsers, adminUpdateUserStatus, adminDeleteUser, adminUpdateUser } from '@/api/user'
import { useUserStore } from '@/stores/user'
import { formatDateTime } from '@/utils/format'

const store = useUserStore()
const list = ref([])
const total = ref(0)
const loading = ref(false)
const keyword = ref('')
const query = reactive({ page: 1, page_size: 10 })

async function load() {
  loading.value = true
  try {
    const params = { page: query.page, page_size: query.page_size }
    if (keyword.value.trim()) params.keyword = keyword.value.trim()
    const data = await adminListUsers(params)
    list.value = data.list || []
    total.value = Number(data.total || 0)
  } catch (e) { /* 拦截器提示 */ } finally {
    loading.value = false
  }
}

function reload() {
  query.page = 1
  load()
}

async function toggleStatus(row) {
  try {
    await adminUpdateUserStatus(row.id, row.status === 1 ? 2 : 1)
    ElMessage.success(row.status === 1 ? '已禁用' : '已解禁')
    load()
  } catch (e) { /* 拦截器提示 */ }
}

async function remove(row) {
  try {
    await adminDeleteUser(row.id)
    ElMessage.success('已删除')
    load()
  } catch (e) { /* 拦截器提示 */ }
}

const editVisible = ref(false)
const editSaving = ref(false)
const editForm = reactive({ id: null, nickname: '', bio: '' })

function openEdit(row) {
  editForm.id = row.id
  editForm.nickname = row.nickname || ''
  editForm.bio = row.bio || ''
  editVisible.value = true
}

async function saveEdit() {
  editSaving.value = true
  try {
    await adminUpdateUser(editForm.id, {
      nickname: editForm.nickname || null,
      bio: editForm.bio || null
    })
    ElMessage.success('已保存')
    editVisible.value = false
    load()
  } catch (e) { /* 拦截器提示 */ } finally {
    editSaving.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar-row { display: flex; align-items: center; gap: 10px; margin-bottom: 16px; }
.u-cell { display: flex; align-items: center; gap: 10px; }
.u-name { font-weight: 500; font-size: 14px; }
.op-link { color: var(--primary); font-size: 13px; margin-right: 12px; cursor: pointer; }
.op-link.danger { color: var(--danger); }
.pager { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
