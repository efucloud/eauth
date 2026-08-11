# EAuth（易认证）

[中文](./README.md) | [English](./README.en.md)

EAuth 是一个面向企业场景的统一认证系统，定位类似 Auth0 / Dex 的认证能力平台。  
它既可以作为业务系统的 OIDC 认证中心，也支持对接外部第三方认证平台，统一账号登录体验与认证策略。
## 项目背景

在多系统、多组织、多终端并存的环境下，认证体系往往面临以下问题：

- 用户账号体系分散，登录入口不统一
- 第三方登录接入成本高，协议差异大
- 安全能力（MFA、人脸、设备侧风控）难以集中管理
- 企业内部系统与外部身份源之间缺少统一桥接层

EAuth 的目标是提供一套可扩展、可落地的认证基础设施，降低认证接入成本，并提升整体安全性。

## 核心能力

- OIDC 认证系统能力（统一签发与验证 Token）
- 第三方 OIDC / OAuth 登录接入能力（见 `pkg/providers`）
- 人脸识别认证能力
- MFA（多因素认证）能力
- 浏览器指纹能力（用于设备侧识别与风控场景）


## 第三方认证扩展

项目通过 `pkg/providers` 实现第三方认证连接器，可按 provider 进行扩展。  
当前代码结构已支持按类别注册和分发 provider，便于新增国内外认证平台。

## 目录说明

- `cmd/`：启动入口
- `web/`：前端源码
- `pkg/apis/`：API 层
- `pkg/services/`：业务逻辑
- `pkg/providers/`：第三方认证接入实现
- `pkg/embeds/web/`：前端构建后的嵌入静态资源目录
- `config/`：本地运行配置
- `docs/`：部署与文档

## 快速开始（本地）

### 环境要求

- Go `1.26.x`
- Node.js `20.x`
- Yarn `1.x`
- MySQL `8.x` 或兼容版本

### 启动步骤

1. 启动本地 MySQL，并创建 `eauth` 数据库

如果本机没有 MySQL，可以直接用 Docker：

```bash
docker run -d \
  --name eauth-mysql \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=EfuCloud \
  -e MYSQL_DATABASE=eauth \
  mysql:8
```

2. 检查并按需修改 `config/config.yaml`

本地最少需要确认这些字段：

- `serverAddress: http://localhost:9001`
- `mysql.host: localhost:3306`
- `mysql.user: root`
- `mysql.password: EfuCloud`
- `mysql.dbname: eauth`

3. 构建并嵌入前端静态资源

```bash
./scripts/build-web-embed.sh
```

也可以使用：

```bash
make embed-web
```

4. 启动统一服务

```bash
go run ./cmd -c ./config/config.yaml
```

默认服务端口为 `9001`。
启动后，后端会同时提供 API 和嵌入后的前端页面。

### 启动后访问

- 首页：`http://127.0.0.1:9001/`
- 健康检查：`http://127.0.0.1:9001/api/health`
- OpenAPI：`http://127.0.0.1:9001/api/v1/swagger.json`
- OIDC Metadata：`http://127.0.0.1:9001/api/.well-known/openid-configuration`

### 默认管理员账号

服务首次启动时会自动初始化管理员：

- 用户名：`admin`
- 密码：`EfuCloud`

## 本地开发说明

### 后端开发

如果只调试后端 API，可以直接启动：

```bash
go run ./cmd -c ./config/config.yaml
```

如果改了前端代码并希望后端继续直接提供页面，需要重新执行：

```bash
./scripts/build-web-embed.sh
```

### 前后端一体化本地预览

如果需要通过后端直接访问前端页面，需要先执行：

```bash
./scripts/build-web-embed.sh
```

这个脚本会完成三件事：

- 在 `web/` 下安装前端依赖
- 构建前端产物到 `web/dist`
- 将构建结果同步到 `pkg/embeds/web/`，供 Go 服务直接嵌入并对外提供

适用场景：

- 首次拉取仓库后，需要让本地后端带上前端页面一起启动
- 修改了 `web/` 下的前端代码后，需要重新生成嵌入资源
- 在构建本地镜像或发布镜像前，需要确保嵌入资源是最新的

前端代码有变更时，重新执行一次即可：

```bash
./scripts/build-web-embed.sh
go run ./cmd -c ./config/config.yaml
```

### 前端单独开发

如果只想调试前端页面，也可以单独启动前端开发服务器：

```bash
cd web
yarn install --frozen-lockfile
yarn start:dev
```

默认会启动前端开发服务。联调时请确保后端也已经在本地 `9001` 端口运行。

### 前端源码位置

- 前端源码目录：`web/`
- 前端构建输出目录：`web/dist/`
- 后端嵌入资源目录：`pkg/embeds/web/`

## 配置文件说明

EAuth 默认通过 `-c` 参数加载 YAML 配置文件，示例路径：

- 本地运行：`config/config.yaml`
- Kubernetes：`/efucloud/config/config.yaml`（由 `docs/backend.yaml` 中的 Secret 挂载）

推荐从 `config/config.yaml` 复制一份按环境维护（如 `config/config.prod.yaml`），启动时显式指定：

```bash
go run ./cmd -c ./config/config.prod.yaml
```

### 关键配置项

- `serverAddress`：系统外网访问地址（用于 OAuth/OIDC 回调拼接）
- `tokenPeriod`：Token 有效期（小时）
- `uploadPath`：上传文件目录（例如头像）
- `loginConfig.mfa`：是否开启 MFA
- `loginConfig.faceRecognition`：是否开启人脸识别
- `mysql.*`：数据库连接参数
- `email.*`：邮件服务参数（验证码、通知等）
- `logConfig.*`：日志输出与滚动策略

### 前端集成说明

- 前端源码位于 `web/`
- 生产镜像和本地统一服务都通过 `pkg/embeds/web/` 提供前端页面
- `pkg/embeds/web/` 下的文件属于构建产物，仓库只保留占位文件 `.ignore`

### 示例（最小可用）

```yaml
# serverAddress 为统一服务对外地址，用于 .well-known/openid-configuration 等元数据生成
serverAddress: "http://localhost:9001"
tokenPeriod: 16
uploadPath: "./uploads"

loginConfig:
  mfa: true
  faceRecognition: true

mysql:
  host: "127.0.0.1:3306"
  user: "root"
  password: "CHANGE_ME"
  dbname: "eauth"
  charset: "utf8mb4"
  loc: "Asia/Shanghai"
  defaultStringSize: 0
  disableDatetimePrecision: false
  dontSupportRenameColumn: false
  dontSupportRenameIndex: false
  skipInitializeWithVersion: false

email:
  smtpServer: "smtp.qq.com"
  smtpPort: 465
  username: "noreply@example.com"
  password: "CHANGE_ME"
```

## Kubernetes 部署

已提供集群部署示例：

- 部署清单：
  - `docs/namespace.yaml`
  - `docs/mysql.yaml`
  - `docs/backend.yaml`：统一服务 Deployment + Service
  - `docs/frontend.yaml`：统一服务 Ingress（可选）
- 使用说明：`docs/README.md`
