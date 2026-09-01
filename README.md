# EAuth（易认证）

[中文](./README.md) | [English](./README.en.md)

EAuth 是一个面向企业场景的统一认证系统，定位类似 Auth0 / Dex 的认证能力平台。  
它既可以作为业务系统的 OIDC 认证中心，也支持对接外部第三方认证平台，统一账号登录体验与认证策略。
## 前端仓库地址：https://github.com/efucloud/eauth-console
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
- `pkg/apis/`：API 层
- `pkg/services/`：业务逻辑
- `pkg/providers/`：第三方认证接入实现
- `config/`：本地运行配置
- `docs/`：部署与文档

## 快速开始（本地）

- Go `1.26.x`
- MySQL `8.x` 或兼容版本

### 1. 启动 MySQL

如果本机没有 MySQL，可以直接用 Docker：

```bash
docker run -d \
  --name eauth-mysql \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=EfuCloud \
  -e MYSQL_DATABASE=eauth \
  mysql:8
```

### 2. 修改后端配置

本地最少需要确认这些字段：

- `serverAddress: http://localhost:8000`
- `mysql.host: localhost:3306`
- `mysql.user: root`
- `mysql.password: EfuCloud`
- `mysql.dbname: eauth`

### 3. 启动后端

```bash
go run ./cmd -c ./config/config.yaml
```

后端默认监听 `9001`。

### 4. 启动前端

```bash
git clone https://github.com/efucloud/eauth-console.git
cd eauth-console
yarn install --frozen-lockfile
yarn start:dev
```

前端默认监听 `8000`，并通过代理访问本地 `9001` 端口的后端。

### 启动后访问

- 首页：`http://127.0.0.1:8000/`
- 健康检查：`http://127.0.0.1:9001/api/health`
- OpenAPI：`http://127.0.0.1:9001/api/v1/swagger.json`
- OIDC Metadata：`http://127.0.0.1:8000/.well-known/openid-configuration`

### 默认管理员账号

服务首次启动时会自动初始化管理员：

- 用户名：`admin`
- 密码：`EfuCloud`

## 本地开发说明

### 仅运行后端

如果只调试后端 API，可以直接启动：

```bash
go run ./cmd -c ./config/config.yaml
```

此时可直接调试 API，但前端页面需要在 `eauth-console` 仓库单独启动。

### 前后端联调

建议先启动本仓库中的后端，再在前端仓库中执行：

```bash
yarn install --frozen-lockfile
yarn start:dev
```

联调时保持后端为 `9001`、前端为 `8000`。

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

### 示例（最小可用）

```yaml
# serverAddress 为前端地址，为 .well-known/openid-configuration 等元数据生成提供地址信息
serverAddress: "http://localhost:8000"
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
  - `docs/backend.yaml`
  - `docs/frontend.yaml`
- 使用说明：`docs/README.md`
