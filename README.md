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
- `pkg/apis/`：API 层
- `pkg/services/`：业务逻辑
- `pkg/providers/`：第三方认证接入实现
- `config/`：本地运行配置
- `docs/`：部署与文档

## 快速开始（本地）

1. 准备 MySQL 并修改配置文件：`config/config.yaml`
2. 启动服务：

```bash
go run ./cmd -c ./config/config.yaml
```

默认服务端口为 `9001`。

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
