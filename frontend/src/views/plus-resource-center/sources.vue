<template>
  <div class="resource-page">
    <ResourceHeader title="插件源" description="配置当前应用的目录来源、同步状态和签名策略。" action-text="添加插件源" @action="openCreate" />
    <div class="source-banner"><div class="source-orb"><ArtSvgIcon icon="ri:database-2-line" /></div><div><strong>已配置 {{ sources.length }} 个插件源</strong><p>插件源会先经过 Plus 校验并缓存清单，客户端不直接管理源地址。</p></div><ElTag :type="errorCount ? 'warning' : 'success'" effect="light">{{ errorCount ? `${errorCount} 个源异常` : '全部正常' }}</ElTag></div>
    <ElCard shadow="never" class="table-card" v-loading="loading">
      <div v-if="!sources.length && !loading" class="empty"><ArtSvgIcon icon="ri:database-2-line" /><strong>还没有插件源</strong><span>添加一个公开 JSON 清单或 Git 仓库地址开始使用。</span><ElButton type="primary" @click="openCreate">添加插件源</ElButton></div>
      <div v-for="source in sources" :key="source.id" class="source-row"><div class="source-mark"><ArtSvgIcon icon="ri:server-line" /></div><div class="source-main"><b>{{ source.name || '未命名源' }}</b><span>{{ source.url }}</span></div><div class="source-detail"><span>状态</span><strong :class="source.state">{{ stateText(source.state) }}</strong></div><ElButton link :loading="refreshingId === source.id" @click="refreshSource(source)"><ArtSvgIcon icon="ri:refresh-line" />刷新</ElButton><ElButton link type="primary" @click="editSource(source)">编辑</ElButton><ElButton link type="danger" @click="removeSource(source)">删除</ElButton></div>
    </ElCard>
    <ElCard shadow="never" class="security-card"><div class="security-icon"><ArtSvgIcon icon="ri:shield-check-line" /></div><div><b>源完整性保护</b><p>添加时会校验清单格式；后续将接入 SHA-256 包校验和私有源签名验证。</p></div><ElTag type="info" effect="plain">基础校验已启用</ElTag></ElCard>
    <ElDialog v-model="dialog" :title="editing ? '编辑插件源' : '添加插件源'" width="520px"><ElForm ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="saveSource"><ElFormItem label="源名称" prop="name"><ElInput v-model="form.name" placeholder="例如：authproPlus 官方源" /></ElFormItem><ElFormItem label="清单或仓库地址" prop="url"><ElInput v-model="form.url" placeholder="https://example.com/index.json" /></ElFormItem><div class="form-tip"><ArtSvgIcon icon="ri:information-line" />支持 HTTPS JSON 清单或可公开访问的 Git 仓库；保存时会立即校验。</div></ElForm><template #footer><ElButton @click="dialog = false">取消</ElButton><ElButton type="primary" :loading="saving" @click="saveSource">{{ editing ? '保存修改' : '添加并校验' }}</ElButton></template></ElDialog>
  </div>
</template>
<script setup lang="ts">
  import { computed, onMounted, reactive, ref } from 'vue'
  import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
  import ResourceHeader from '@/components/plus-resource/ResourceHeader.vue'
  import { fetchAddPluginSource, fetchDeletePluginSource, fetchPluginList, fetchRefreshPluginSource, type PluginSource } from '@/api/system-manage'
  const sources = ref<PluginSource[]>([])
  const loading = ref(false)
  const saving = ref(false)
  const refreshingId = ref<number | null>(null)
  const dialog = ref(false)
  const editing = ref(false)
  const editingSource = ref<PluginSource | null>(null)
  const formRef = ref<FormInstance>()
  const form = reactive({ name: '', url: '' })
  const rules: FormRules = { url: [{ required: true, message: '请输入源地址', trigger: 'blur' }, { type: 'url', message: '请输入有效的 URL', trigger: 'blur' }] }
  const errorCount = computed(() => sources.value.filter((item) => item.state === 'error').length)
  const stateText = (state: PluginSource['state']) => ({ ok: '连接正常', error: '连接异常', unknown: '待检查' })[state]
  const loadSources = async () => { loading.value = true; try { const data = await fetchPluginList(); sources.value = data.sources || [] } catch (error: any) { ElMessage.error(error?.message || '插件源加载失败') } finally { loading.value = false } }
  const openCreate = () => { editing.value = false; editingSource.value = null; form.name = ''; form.url = ''; dialog.value = true }
  const editSource = (source: PluginSource) => { editing.value = true; editingSource.value = source; form.name = source.name; form.url = source.url; dialog.value = true }
  const saveSource = async () => { if (!formRef.value) return; const valid = await formRef.value.validate().catch(() => false); if (!valid) return; saving.value = true; try { if (editing.value && editingSource.value) { if (form.url !== editingSource.value.url) { await fetchAddPluginSource(form.name, form.url); await fetchDeletePluginSource(editingSource.value.id) } else if (form.name !== editingSource.value.name) { await fetchAddPluginSource(form.name, form.url) } ElMessage.success('插件源已更新') } else { await fetchAddPluginSource(form.name, form.url); ElMessage.success('插件源添加成功') } dialog.value = false; await loadSources() } catch (error: any) { ElMessage.error(error?.message || '保存插件源失败') } finally { saving.value = false } }
  const refreshSource = async (source: PluginSource) => { refreshingId.value = source.id; try { const result = await fetchRefreshPluginSource(source.id); ElMessage.success(`刷新成功：发现 ${result.plugins} 个插件、${result.homeTemplates} 个模板`); await loadSources() } catch (error: any) { ElMessage.error(error?.message || '刷新插件源失败') } finally { refreshingId.value = null } }
  const removeSource = async (source: PluginSource) => { try { await ElMessageBox.confirm(`确认删除插件源「${source.name || source.url}」？已下载的本地插件不会受影响。`, '删除插件源', { type: 'warning', confirmButtonText: '删除', cancelButtonText: '取消' }); await fetchDeletePluginSource(source.id); ElMessage.success('插件源已删除'); await loadSources() } catch (error: any) { if (error !== 'cancel' && error !== 'close') ElMessage.error(error?.message || '删除插件源失败') } }
  onMounted(loadSources)
</script>
<style scoped>
  .resource-page{padding:24px;background:#f6f8fc;min-height:100%;color:#172033}.source-banner,.security-card{display:flex;align-items:center;gap:16px;margin:18px 0;padding:20px 24px;border:1px solid #e9edf4;border-radius:14px;background:#fff}.source-banner>div:nth-child(2),.security-card>div:nth-child(2){flex:1}.source-orb,.security-icon{display:grid;place-items:center;width:42px;height:42px;border-radius:11px;color:#6556bd;background:#f0edff;font-size:21px}.source-banner strong,.security-card b{font-size:14px}.source-banner p,.security-card p{margin:6px 0 0;color:#8b94a5;font-size:12px}.table-card{border:0;border-radius:14px;min-height:180px}.source-row{display:flex;align-items:center;gap:14px;padding:19px 4px;border-bottom:1px solid #edf0f5}.source-row:last-child{border-bottom:0}.source-mark{display:grid;place-items:center;flex:0 0 36px;width:36px;height:36px;border-radius:9px;color:#6673a5;background:#eef0ff;font-size:18px}.source-main{flex:1;min-width:0;display:flex;flex-direction:column;gap:5px}.source-main span{overflow:hidden;color:#8b94a5;font-size:11px;text-overflow:ellipsis;white-space:nowrap}.source-detail{display:flex;flex-direction:column;gap:5px;min-width:75px}.source-detail span{color:#8b94a5;font-size:11px}.source-detail strong{font-size:12px}.source-detail strong.ok{color:#2d9d68}.source-detail strong.error{color:#dc625e}.source-detail strong.unknown{color:#9a7a33}.security-card{margin-top:18px}.empty{display:flex;align-items:center;flex-direction:column;gap:10px;padding:28px;color:#8b94a5;text-align:center}.empty>svg{font-size:34px;color:#7354c8}.empty strong{color:#172033}.empty span{font-size:12px}.form-tip{display:flex;gap:6px;padding:10px;color:#6673a5;background:#f2f3ff;border-radius:7px;font-size:12px;line-height:18px}.form-tip svg{flex:0 0 auto;margin-top:2px}.security-card{margin-top:18px}@media(max-width:700px){.resource-page{padding:12px}.source-banner,.security-card{align-items:flex-start;flex-wrap:wrap;padding:16px}.source-banner>div:nth-child(2),.security-card>div:nth-child(2){min-width:calc(100% - 58px)}.source-row{align-items:flex-start;flex-wrap:wrap;gap:10px}.source-main{min-width:calc(100% - 52px)}.source-detail{margin-left:50px}.source-row .el-button{margin-left:0}.source-row .el-button:nth-last-child(3){margin-left:50px}}
</style>
