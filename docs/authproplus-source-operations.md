# authproPlus 广告、插件源与更新源操作手册

## 1. 当前可控范围

当前版本已经支持在后台管理插件源和首页模板源；广告和 authproPlus 自身更新源通过服务端配置控制。这样保持了上游项目现有的软件源协议和目录格式，后续同步作者版本时不会因为重做协议而产生冲突。

| 项目 | 当前控制方式 |
| --- | --- |
| 插件源 | 后台“插件中心/插件管理”中的“软件源管理” |
| 首页模板 | 插件源 `index.json` 的 `homeTemplates` 字段，后台启用 |
| 广告 | `AUTO_PRO_ADVERTISEMENT_URL` 指向兼容广告 API |
| authproPlus 更新源 | `AUTO_PRO_UPDATE_URL` 指向 `latest.json` |

## 2. 插件源

后台进入“插件中心”或“插件管理”，打开“软件源管理”，添加 JSON 清单地址或 Git 仓库地址。

JSON 清单格式保持上游格式：

```json
{
  "name": "authproPlus 插件源",
  "plugins": [
    {
      "id": "my-plugin",
      "category": "other",
      "name": "我的插件",
      "description": "插件描述",
      "icon": "ri:plug-line",
      "version": "1.0.0",
      "author": { "name": "你的名称", "url": "", "email": "" },
      "downloadUrl": "plugins/my-plugin.zip"
    }
  ],
  "homeTemplates": []
}
```

Git 仓库根目录必须包含 `index.json`。插件包地址可以使用仓库内相对路径。系统会校验 JSON、插件 ID、名称和版本，并缓存目录约 5 分钟。后台支持对已添加的软件源执行刷新。

## 3. 首页模板

首页模板继续放在同一份 `index.json` 的 `homeTemplates` 字段中：

```json
{
  "name": "authproPlus 主题源",
  "plugins": [],
  "homeTemplates": [
    {
      "id": "plus-home",
      "name": "Plus 首页",
      "description": "authproPlus 首页模板",
      "version": "1.0.0",
      "schemaVersion": 1,
      "sha256": "模板 JSON 文件的 SHA256",
      "templatePath": "templates/plus-home.json"
    }
  ]
}
```

模板文件使用声明式 JSON。发布后在后台刷新首页模板列表，然后执行“启用”。模板下载、SHA-256 校验和 schema 校验失败时，系统自动回退到内置默认模板。

## 4. 广告服务

当前前端从后端读取以下接口：

```text
GET /api/advertisements?position=home-banner
GET /api/advertisements?position=sidebar
GET /api/advertisements?position=popup
```

后端再代理到 `AUTO_PRO_ADVERTISEMENT_URL`。该地址需要返回：

```json
{
  "code": 200,
  "msg": "ok",
  "data": {
    "records": [
      {
        "id": "ad-001",
        "title": "Plus 广告",
        "imageUrl": "https://cdn.example.com/ad.png",
        "destinationUrl": "https://example.com",
        "position": "home-banner",
        "weight": 100,
        "startAt": "2026-09-16T00:00:00Z",
        "endAt": "2026-12-31T23:59:59Z",
        "description": "广告说明"
      }
    ]
  }
}
```

广告位只允许 `home-banner`、`sidebar` 和 `popup`。系统会按开始时间、结束时间和权重过滤，并缓存结果。广告源暂时不可用时，系统会在容忍时间内继续使用旧缓存，超过时间后返回空列表。

当前版本没有完整的广告 CRUD 后台。若要在 authproPlus 后台直接新建、编辑、上下架广告，需要后续增加广告管理插件；该插件应继续输出上面的兼容 JSON，不改变前端协议。

## 5. 更新源

服务端启动时通过 `AUTO_PRO_UPDATE_URL` 读取更新清单地址。如果未设置，则仍使用上游默认地址。因此生产环境必须显式设置：

```bash
AUTO_PRO_UPDATE_URL=https://plus.example.com/releases/latest.json
```

`latest.json` 至少需要包含：

```json
{
  "version": "1.1.7",
  "channel": "stable",
  "minVersion": "1.0.0",
  "force": false,
  "releasedAt": "2026-09-16T14:00:00Z",
  "releasesUrl": "https://plus.example.com/releases/releases.json",
  "package": {
    "os": "linux",
    "arch": "amd64",
    "fileName": "auth_pro-full-v1.1.7.tar.gz",
    "url": "https://plus.example.com/releases/auth_pro-full-v1.1.7.tar.gz",
    "sha256": "64 位 SHA256",
    "size": 12345678,
    "signature": ""
  },
  "actions": {
    "updateFrontend": true,
    "updateBackend": true,
    "restartBackend": true,
    "backupDatabase": true
  },
  "notes": ["修复问题", "增加功能"]
}
```

更新地址、软件包地址和历史版本地址必须使用 HTTPS。当前一键整包更新主要支持 Linux amd64；更新包会校验文件大小和 SHA-256。

## 6. 当前测试环境设置

当前测试服务可以使用以下配置重启：

```bash
cd /home/ubuntu/auth_pro/backend
export AUTO_PRO_DATA_DIR=/home/ubuntu/auth_pro/runtime
export AUTO_PRO_FRONTEND_DIR=/home/ubuntu/auth_pro/backend/static
export AUTO_PRO_ADVERTISEMENT_URL=https://你的广告接口/api/v1/public/advertisements
export AUTO_PRO_UPDATE_URL=https://你的域名/releases/latest.json
./auth_pro
```

正式部署时建议复制项目根目录的 `.env.authproplus.example`，改名为部署环境的配置文件，并通过 systemd、容器编排或面板安全注入环境变量。数据库密码和签名密钥不能提交到 Git 仓库。

## 7. 推荐的下一步插件化改造

广告管理应作为 Plus 插件实现。插件负责广告 CRUD、图片存储、投放规则和统计；对终端仍输出现有 `records` 格式。

更新源管理也可以作为 Plus 插件实现。插件负责保存更新频道和清单地址；核心更新器继续优先读取环境变量，并在没有环境变量时读取数据库中的插件配置。这样既保留现有部署兼容性，也能提供后台配置能力。

插件管理、广告管理和更新源管理的新增页面必须复用现有 RBAC 权限体系，不能只依靠前端隐藏按钮。
