<template>
  <span
    class="uava"
    :class="{ gray: !hasName }"
    :style="{ width: size + 'px', height: size + 'px', fontSize: Math.max(12, Math.round(size * 0.4)) + 'px' }"
  >
    <img v-if="user?.avatar" :src="user.avatar" :alt="name" />
    <template v-else>{{ letter }}</template>
  </span>
</template>

<script setup>
import { computed } from 'vue'

// 统一头像：有 avatar 显示图片，无则显示昵称/用户名首字符
const props = defineProps({
  user: { type: Object, default: null },
  size: { type: Number, default: 34 }
})

const name = computed(() => (props.user && (props.user.nickname || props.user.username)) || '')
const hasName = computed(() => !!name.value)
const letter = computed(() => (name.value ? name.value.slice(0, 1).toUpperCase() : '?'))
</script>

<style scoped>
.uava {
  border-radius: 50%; overflow: hidden; flex-shrink: 0;
  display: grid; place-items: center; line-height: 1;
  background: var(--primary-100); color: var(--primary-700); font-weight: 700;
}
.uava.gray { background: #e5e7eb; color: #6b7280; }
.uava img { width: 100%; height: 100%; object-fit: cover; display: block; }
</style>
