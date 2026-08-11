# EAuth 集群部署说明

[中文](./README.md) | [English](./README.en.md)

本文档对应以下 Kubernetes 清单：

- `docs/namespace.yaml`：命名空间（`efucloud`）
- `docs/mysql.yaml`：MySQL（PVC + ConfigMap + Deployment + Service）
- `docs/backend.yaml`：EAuth 单服务部署（配置 Secret + Deployment + Service）
- `docs/frontend.yaml`：EAuth 单服务 Ingress（可选）

默认部署命名空间：`efucloud`。

## 1. 部署前准备

- 已安装 Kubernetes 与 Ingress NGINX（如需域名访问）
- 已准备镜像拉取凭据（如私有仓库）
- 已有可用 StorageClass（示例中 `mysql-local-retain`）
- 部署前请先修改以下配置：
  - `docs/mysql.yaml`：`MYSQL_ROOT_PASSWORD`、`MYSQL_DATABASE`
  - `docs/backend.yaml`：`Secret.stringData.config.yaml` 中的 `mysql` / `email` / `serverAddress`
  - `docs/frontend.yaml`：域名、IngressClass、TLS 证书 secret

## 2. 推荐部署顺序

```bash
kubectl apply -f docs/namespace.yaml
kubectl apply -f docs/mysql.yaml
kubectl apply -f docs/backend.yaml
kubectl apply -f docs/frontend.yaml
kubectl port-forward -n efucloud svc/eauth 9001:80
```

也可以一次性执行：

```bash
kubectl apply -f docs/
```

## 3. 状态检查

```bash
kubectl -n efucloud get pods,svc,deploy,ingress
```

查看后端日志：

```bash
kubectl -n efucloud logs -f deploy/eauth
```

## 4. 联通验证

集群内统一访问地址：

- `http://eauth.efucloud.svc.cluster.local`

本地端口转发后访问：

- 控制台首页：`http://127.0.0.1:9001/`
- 健康检查：`http://127.0.0.1:9001/api/health`
- OpenAPI：`http://127.0.0.1:9001/api/v1/swagger.json`
- OIDC Metadata：`http://127.0.0.1:9001/api/.well-known/openid-configuration`

健康检查接口：

- `GET /api/health`
- `GET /api/v1/swagger.json`
- `GET /metrics`

如果使用 `docs/frontend.yaml` 中的 Ingress，还可以通过根路径访问前端页面：

- `GET /`
- `GET /.well-known/openid-configuration`（Ingress 做了兼容转发）

## 5. 生产环境建议

- 当前服务配置采用 `Secret.stringData`，建议对接外部密钥管理（如 SealedSecrets / External Secrets）
- MySQL 建议使用独立实例或 StatefulSet，并启用备份策略
- `uploads` 建议使用 PVC 或对象存储，避免 Pod 重建导致数据丢失
- Ingress 建议强制 TLS，并配置可观测性（访问日志/错误日志/告警）
