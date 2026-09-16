<template>
  <div class="resource-center page-content">
    <section class="hero-panel">
      <div>
        <div class="eyebrow"><span class="eyebrow-dot"></span> AUTHPROPLUS CONTROL PLANE</div>
        <h1>资源控制中心</h1>
        <p>按应用统一管理插件源、首页模板与广告策略，客户端只从 Plus 获取已授权内容。</p>
      </div>
      <div class="hero-actions">
        <span class="sync-state"><i></i> 配置预览模式</span>
        <ElButton type="primary" :icon="Plus" @click="openCreate">新建资源</ElButton>
      </div>
    </section>

    <section class="control-strip">
      <div class="app-context">
        <span class="label">当前应用</span>
        <ElSelect v-model="activeApp" class="app-select" placeholder="选择应用">
          <ElOption v-for="app in apps" :key="app.value" :label="app.label" :value="app.value" />
        </ElSelect>
        <span class="scope-badge">APP SCOPE · {{ activeApp }}</span>
      </div>
      <div class="quick-stats">
        <div><strong>12</strong><span>已发布插件</span></div>
        <div><strong>04</strong><span>首页模板</span></div>
        <div><strong>08</strong><span>投放广告</span></div>
        <div><strong>03</strong><span>付费插件</span></div>
      </div>
    </section>

    <ElCard class="workspace-card" shadow="never">
      <ElTabs v-model="activeTab" class="resource-tabs">
        <ElTabPane name="plugins">
          <template #label><span class="tab-label"><ArtSvgIcon icon="ri:extension-line" />插件目录</span></template>
          <div class="tab-toolbar">
            <div><h2>插件目录</h2><p>管理当前应用可分发的插件、版本和商业属性。</p></div>
            <div class="toolbar-actions"><ElInput v-model="pluginSearch" clearable placeholder="搜索插件" class="search-input"><template #prefix><ArtSvgIcon icon="ri:search-line" /></template></ElInput><ElButton @click="showSource = true" plain>管理插件源</ElButton><ElButton type="primary" :icon="Plus" @click="openCreate">发布插件</ElButton></div>
          </div>
          <div class="resource-grid">
            <article v-for="plugin in filteredPlugins" :key="plugin.key" class="resource-card">
              <div class="card-top"><div class="resource-icon" :class="plugin.color"><ArtSvgIcon :icon="plugin.icon" /></div><ElTag :type="plugin.paid ? 'warning' : 'success'" effect="light" size="small">{{ plugin.paid ? '付费插件' : '免费' }}</ElTag></div>
              <h3>{{ plugin.name }}</h3><p>{{ plugin.description }}</p>
              <div class="card-meta"><span>v{{ plugin.version }}</span><span>{{ plugin.category }}</span><span>{{ plugin.paid ? '¥' + plugin.price : '开放使用' }}</span></div>
              <div class="card-footer"><span class="published"><i></i>{{ plugin.published ? '已发布' : '草稿' }}</span><ElButton link type="primary" @click="openCreate">编辑</ElButton></div>
            </article>
          </div>
        </ElTabPane>

        <ElTabPane name="templates">
          <template #label><span class="tab-label"><ArtSvgIcon icon="ri:layout-4-line" />首页模板</span></template>
          <div class="tab-toolbar"><div><h2>首页模板</h2><p>为当前应用选择登录页和用户入口的展示模板。</p></div><ElButton type="primary" :icon="Plus" @click="openCreate">上传模板</ElButton></div>
          <div class="template-grid"><article v-for="template in templates" :key="template.name" class="template-card"><div class="template-preview" :class="template.tone"><div class="preview-bar"></div><div class="preview-content"><span></span><span></span><span></span></div><strong>{{ template.short }}</strong></div><div class="template-info"><div><h3>{{ template.name }}</h3><p>v{{ template.version }} · Schema {{ template.schema }}</p></div><ElTag :type="template.active ? 'success' : 'info'" effect="light">{{ template.active ? '当前启用' : '可启用' }}</ElTag></div><ElButton class="template-action" :type="template.active ? 'default' : 'primary'" plain @click="openCreate">{{ template.active ? '配置模板' : '启用模板' }}</ElButton></article></div>
        </ElTabPane>

        <ElTabPane name="ads">
          <template #label><span class="tab-label"><ArtSvgIcon icon="ri:advertisement-line" />广告策略</span></template>
          <div class="tab-toolbar"><div><h2>广告策略</h2><p>广告仅对当前应用生效，可按广告位、时间和权重投放。</p></div><ElButton type="primary" :icon="Plus" @click="openCreate">创建广告</ElButton></div>
          <ElTable :data="ads" class="resource-table" stripe><ElTableColumn prop="title" label="广告内容" min-width="240"><template #default="{ row }"><div class="ad-title"><span class="ad-thumb" :class="row.tone"></span><div><strong>{{ row.title }}</strong><small>{{ row.description }}</small></div></div></template></ElTableColumn><ElTableColumn prop="position" label="广告位" width="140"><template #default="{ row }"><ElTag effect="plain">{{ row.position }}</ElTag></template></ElTableColumn><ElTableColumn prop="period" label="投放周期" min-width="220" /><ElTableColumn prop="weight" label="权重" width="90" /><ElTableColumn label="状态" width="100"><template #default="{ row }"><span class="status-text"><i :class="row.active ? 'on' : ''"></i>{{ row.active ? '投放中' : '已暂停' }}</span></template></ElTableColumn><ElTableColumn label="操作" width="120" fixed="right"><template #default><ElButton link type="primary" @click="openCreate">编辑</ElButton></template></ElTableColumn></ElTable>
        </ElTabPane>

        <ElTabPane name="sources">
          <template #label><span class="tab-label"><ArtSvgIcon icon="ri:database-2-line" />插件源</span></template>
          <div class="tab-toolbar"><div><h2>插件源</h2><p>配置当前应用可读取的目录来源，后续支持私有源和签名校验。</p></div><ElButton type="primary" :icon="Plus" @click="openCreate">添加插件源</ElButton></div>
          <div class="source-list"><div v-for="source in sources" :key="source.name" class="source-row"><div class="source-mark"><ArtSvgIcon icon="ri:server-line" /></div><div class="source-main"><strong>{{ source.name }}</strong><span>{{ source.url }}</span></div><ElTag type="success" effect="light">连接正常</ElTag><span class="source-updated">同步于 {{ source.updated }}</span><ElButton link type="primary" @click="openCreate">配置</ElButton></div></div>
        </ElTabPane>
      </ElTabs>
    </ElCard>

    <ElDialog v-model="dialogVisible" title="资源编辑器" width="520px"><div class="dialog-placeholder"><div class="placeholder-icon"><ArtSvgIcon icon="ri:tools-line" /></div><h3>页面原型已就绪</h3><p>这里将接入真实的新增、编辑、发布和授权接口。当前只展示交互入口，不会写入生产数据。</p><ElTag type="info">下一阶段接入服务端 API</ElTag></div><template #footer><ElButton @click="dialogVisible = false">关闭</ElButton></template></ElDialog>
    <ElDialog v-model="showSource" title="插件源管理" width="620px"><div class="source-dialog-note"><ArtSvgIcon icon="ri:information-line" /><span>插件源按当前应用隔离。后续可在这里配置公开源、私有源和签名密钥。</span></div><div class="source-list compact"><div v-for="source in sources" :key="source.name" class="source-row"><div class="source-mark"><ArtSvgIcon icon="ri:server-line" /></div><div class="source-main"><strong>{{ source.name }}</strong><span>{{ source.url }}</span></div><ElTag type="success" effect="light">启用</ElTag></div></div><template #footer><ElButton @click="showSource = false">完成</ElButton></template></ElDialog>
  </div>
</template>

<script setup lang="ts">
  import { computed, ref } from 'vue'
  import { Plus } from '@element-plus/icons-vue'
  const activeTab = ref('plugins')
  const activeApp = ref('authpro')
  const pluginSearch = ref('')
  const dialogVisible = ref(false)
  const showSource = ref(false)
  const apps = [{ label: 'authpro 主应用', value: 'authpro' }, { label: '演示应用', value: 'demo' }]
  const plugins = [
    { key: 'payment', name: '支付聚合插件', description: '统一接入支付渠道与订单回调。', version: '2.4.0', category: '支付', paid: true, price: '99', published: true, icon: 'ri:secure-payment-line', color: 'violet' },
    { key: 'realname', name: '实名认证插件', description: '为应用提供可配置的实名核验能力。', version: '1.8.2', category: '身份认证', paid: true, price: '49', published: true, icon: 'ri:shield-user-line', color: 'blue' },
    { key: 'analytics', name: '运营分析', description: '查看应用授权、设备与转化数据。', version: '1.2.1', category: '运营', paid: false, price: '', published: true, icon: 'ri:bar-chart-2-line', color: 'orange' }
  ]
  const filteredPlugins = computed(() => plugins.filter((item) => item.name.includes(pluginSearch.value)))
  const templates = [{ name: 'Aurora 登录页', short: 'A', version: '1.3.0', schema: 2, active: true, tone: 'aurora' }, { name: 'Midnight 控制台', short: 'M', version: '1.0.4', schema: 2, active: false, tone: 'midnight' }, { name: 'Minimal 极简页', short: 'M', version: '0.9.2', schema: 1, active: false, tone: 'minimal' }]
  const ads = [{ title: 'Plus Pro 专业版', description: '升级授权，解锁完整能力', position: 'home-banner', period: '2026/09/01 — 2026/10/01', weight: 90, active: true, tone: 'violet' }, { title: '支付插件限时优惠', description: '订阅插件首月优惠', position: 'sidebar', period: '2026/09/10 — 长期', weight: 60, active: true, tone: 'orange' }, { title: '新版本功能预告', description: '了解 authproPlus 最新能力', position: 'popup', period: '未发布', weight: 20, active: false, tone: 'blue' }]
  const sources = [{ name: 'authproPlus 官方源', url: 'https://source.authproplus.example/api', updated: '刚刚' }, { name: '企业私有源', url: 'https://packages.example.com/authpro', updated: '12 分钟前' }]
  const openCreate = () => { dialogVisible.value = true }
</script>

<style scoped>
  .resource-center { --ink: #172033; --muted: #7d879b; --line: #e9edf4; padding: 24px; color: var(--ink); background: #f6f8fc; min-height: 100%; }
  .hero-panel { display:flex; justify-content:space-between; align-items:flex-end; gap:24px; padding:30px 34px; border-radius:20px; color:#fff; background:linear-gradient(120deg,#1c2541 0%,#27345a 54%,#674e9b 100%); box-shadow:0 16px 42px rgba(32,44,86,.18); }
  .eyebrow { color:#b8c4e7; font-size:11px; letter-spacing:1.5px; font-weight:700; }.eyebrow-dot { display:inline-block;width:7px;height:7px;border-radius:50%;background:#8ee7c0;margin-right:8px;box-shadow:0 0 0 4px rgba(142,231,192,.14); }
  h1 { margin:10px 0 6px; font-size:30px; letter-spacing:-.8px; }.hero-panel p { margin:0; color:#bdc7df; font-size:13px; }.hero-actions { display:flex;align-items:center;gap:18px; }.sync-state { color:#cbd5ee;font-size:12px;white-space:nowrap; }.sync-state i { display:inline-block;width:6px;height:6px;border-radius:50%;background:#8ee7c0;margin-right:7px; }
  .control-strip { display:flex;justify-content:space-between;align-items:center;margin:18px 0;padding:15px 20px;background:#fff;border:1px solid var(--line);border-radius:14px; }.app-context,.quick-stats { display:flex;align-items:center;gap:14px; }.label { color:var(--muted);font-size:12px; }.app-select { width:180px; }.scope-badge { font:10px ui-monospace,monospace;color:#6673a5;background:#f0f2ff;padding:6px 9px;border-radius:5px; }.quick-stats { gap:28px; }.quick-stats div { display:flex;flex-direction:column;gap:2px; }.quick-stats strong { font-size:18px; }.quick-stats span { color:var(--muted);font-size:11px; }
  .workspace-card { border:0;border-radius:16px; }.resource-tabs :deep(.el-tabs__header) { margin-bottom:26px; }.tab-label { display:flex;align-items:center;gap:7px;font-size:13px; }.tab-label :deep(.art-svg-icon) { font-size:16px; }.tab-toolbar { display:flex;justify-content:space-between;align-items:flex-end;gap:16px;margin-bottom:22px; }.tab-toolbar h2 { margin:0 0 5px;font-size:18px; }.tab-toolbar p { margin:0;color:var(--muted);font-size:12px; }.toolbar-actions { display:flex;gap:10px; }.search-input { width:200px; }
  .resource-grid { display:grid;grid-template-columns:repeat(3,1fr);gap:16px; }.resource-card { padding:20px;border:1px solid var(--line);border-radius:13px;transition:.2s;background:#fff; }.resource-card:hover { border-color:#b7b5ed;box-shadow:0 10px 24px rgba(35,45,90,.08);transform:translateY(-2px); }.card-top,.card-footer,.template-info { display:flex;justify-content:space-between;align-items:center; }.resource-icon { display:grid;place-items:center;width:42px;height:42px;border-radius:11px;font-size:22px; }.resource-icon.violet { color:#7354c8;background:#f1ecff; }.resource-icon.blue { color:#3c83d7;background:#eaf4ff; }.resource-icon.orange { color:#d97729;background:#fff1e5; }.resource-card h3 { margin:17px 0 7px;font-size:15px; }.resource-card p { height:36px;margin:0;color:var(--muted);font-size:12px;line-height:18px; }.card-meta { display:flex;gap:12px;margin:18px 0 15px;color:var(--muted);font-size:11px; }.card-footer { padding-top:12px;border-top:1px solid #f0f2f6; }.published,.status-text { color:#6c778d;font-size:11px; }.published i,.status-text i { display:inline-block;width:6px;height:6px;border-radius:50%;background:#b4bdca;margin-right:6px; }.published i,.status-text i.on { background:#47bd88; }
  .template-grid { display:grid;grid-template-columns:repeat(3,1fr);gap:18px; }.template-card { padding:14px;border:1px solid var(--line);border-radius:13px; }.template-preview { height:150px;border-radius:9px;padding:13px;color:#fff;position:relative;overflow:hidden; }.template-preview.aurora { background:linear-gradient(140deg,#7162d4,#d78fb1); }.template-preview.midnight { background:linear-gradient(140deg,#17233e,#3d6791); }.template-preview.minimal { background:linear-gradient(140deg,#cbd3df,#f4b38e); }.preview-bar { height:8px;width:42px;border-radius:4px;background:rgba(255,255,255,.7); }.preview-content { display:flex;flex-direction:column;gap:8px;margin-top:46px; }.preview-content span { display:block;height:7px;border-radius:4px;background:rgba(255,255,255,.75); }.preview-content span:nth-child(1){width:62%}.preview-content span:nth-child(2){width:84%;opacity:.55}.preview-content span:nth-child(3){width:35%;opacity:.35}.template-preview strong { position:absolute;bottom:14px;right:15px;font-size:28px;opacity:.5; }.template-info { padding:15px 2px 12px; }.template-info h3 { margin:0 0 5px;font-size:14px; }.template-info p { margin:0;color:var(--muted);font-size:11px; }.template-action { width:100%; }
  .resource-table { border:1px solid var(--line);border-radius:10px;overflow:hidden; }.ad-title { display:flex;align-items:center;gap:11px; }.ad-title small { display:block;color:var(--muted);font-size:11px;margin-top:4px; }.ad-thumb { width:34px;height:34px;border-radius:8px; }.ad-thumb.violet { background:linear-gradient(135deg,#7c67d9,#c19ee5); }.ad-thumb.orange { background:linear-gradient(135deg,#ec9957,#f6d2a0); }.ad-thumb.blue { background:linear-gradient(135deg,#68a5dc,#afd3f0); }
  .source-list { display:flex;flex-direction:column;border:1px solid var(--line);border-radius:10px;overflow:hidden; }.source-list.compact { margin-top:16px; }.source-row { display:flex;align-items:center;gap:15px;padding:17px 20px;border-bottom:1px solid #f0f2f6; }.source-row:last-child { border-bottom:0; }.source-mark { display:grid;place-items:center;width:36px;height:36px;border-radius:9px;color:#6673a5;background:#eef0ff;font-size:18px; }.source-main { flex:1;display:flex;flex-direction:column;gap:4px; }.source-main strong { font-size:13px; }.source-main span,.source-updated { color:var(--muted);font-size:11px; }.source-updated { margin-right:10px; }.source-dialog-note { display:flex;gap:8px;padding:12px 14px;color:#63709b;background:#f2f3ff;border-radius:8px;font-size:12px; }.dialog-placeholder { text-align:center;padding:18px 30px 30px; }.placeholder-icon { display:grid;place-items:center;width:60px;height:60px;margin:auto;border-radius:17px;color:#6656c8;background:#f0edff;font-size:28px; }.dialog-placeholder h3 { margin:16px 0 8px; }.dialog-placeholder p { color:var(--muted);font-size:13px;line-height:21px; }
  @media (max-width: 1000px) { .hero-panel,.control-strip,.tab-toolbar { align-items:flex-start;flex-direction:column; }.quick-stats { width:100%;justify-content:space-between; }.resource-grid,.template-grid { grid-template-columns:1fr 1fr; } }. 
  @media (max-width: 640px) { .resource-center { padding:12px; }.resource-grid,.template-grid { grid-template-columns:1fr; }.toolbar-actions { width:100%;flex-wrap:wrap; }.search-input { width:100%; }.hero-actions { width:100%;justify-content:space-between; }.control-strip { padding:14px; }.app-context { flex-wrap:wrap; } }
</style>
