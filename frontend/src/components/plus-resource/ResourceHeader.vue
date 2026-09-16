<template>
  <div class="resource-head">
    <div>
      <div class="eyebrow"><span></span> AUTHPROPLUS CONTROL PLANE</div>
      <h1>{{ title }}</h1>
      <p>{{ description }}</p>
    </div>
    <div class="head-actions">
      <span class="scope">当前应用</span>
      <ElSelect v-model="activeApp" size="large" class="app-select">
        <ElOption v-for="app in apps" :key="app.value" :label="app.label" :value="app.value" />
      </ElSelect>
      <ElButton v-if="actionText" type="primary" :icon="Plus" @click="$emit('action')">{{ actionText }}</ElButton>
    </div>
  </div>
</template>
<script setup lang="ts">
  import { Plus } from '@element-plus/icons-vue'
  import { ref } from 'vue'
  defineProps<{ title: string; description: string; actionText?: string }>()
  defineEmits<{ action: [] }>()
  const activeApp = ref('authpro')
  const apps = [{ label: 'authpro 主应用', value: 'authpro' }, { label: '演示应用', value: 'demo' }]
</script>
<style scoped>
  .resource-head { display:flex; justify-content:space-between; align-items:flex-end; gap:24px; padding:28px 32px; background:linear-gradient(120deg,#1c2541,#2d3c67 55%,#6a55a3); border-radius:18px; color:#fff; box-shadow:0 14px 34px rgba(39,49,92,.16); }
  .eyebrow { color:#bfc9e8; font-size:11px; letter-spacing:1.5px; font-weight:700; }.eyebrow span { display:inline-block; width:7px; height:7px; margin-right:8px; border-radius:50%; background:#8ee7c0; }
  h1 { margin:9px 0 6px; font-size:27px; letter-spacing:-.6px; }.resource-head p { margin:0; color:#c3cce1; font-size:13px; }.head-actions { display:flex; align-items:center; gap:10px; }.scope { color:#d5dcf0; font-size:12px; }.app-select { width:170px; }.app-select :deep(.el-input__wrapper) { background:rgba(255,255,255,.12); box-shadow:0 0 0 1px rgba(255,255,255,.22) inset; }.app-select :deep(.el-input__inner) { color:#fff; }
  @media(max-width:760px){.resource-head{align-items:flex-start;flex-direction:column;padding:22px}.head-actions{width:100%;flex-wrap:wrap}.app-select{flex:1}}
</style>
