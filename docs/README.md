# EAuth 集群部署说明

[中文](./README.md) | [English](./README.en.md)

当前部署模型为前后端分离：

- `docs/namespace.yaml`：命名空间（`efucloud`）
- `docs/mysql.yaml`：MySQL（PVC + ConfigMap + Deployment + Service）
- `docs/backend.yaml`：EAuth 后端（Secret + Deployment + Service）
- `docs/frontend.yaml`：EAuth 前端控制台（Secret + Deployment + Service + Ingress）

默认部署命名空间：`efucloud`。

## 1. 部署前准备

- 已安装 Kubernetes。
- 如需域名访问，请提前安装 Ingress NGINX。
- 如镜像仓库为私有仓库，请准备镜像拉取凭据。
- 请先按环境修改以下配置：
  - `docs/mysql.yaml`：`MYSQL_ROOT_PASSWORD`、`MYSQL_DATABASE`
  - `docs/backend.yaml`：`Secret.stringData.config.yaml` 中的 `mysql`、`email`、`serverAddress`
  - `docs/frontend.yaml`：前端镜像地址、域名、TLS 配置、Nginx 反向代理参数

`serverAddress` 应配置为用户访问前端控制台的公网地址，而不是后端 Service 地址。

## 2. 推荐部署顺序

```bash
kubectl apply -f docs/namespace.yaml
kubectl apply -f docs/mysql.yaml
kubectl apply -f docs/backend.yaml
kubectl apply -f docs/frontend.yaml
```

如果已经完成环境定制，也可以直接执行：

```bash
kubectl apply -f docs/
```

## 3. 访问方式

- 后端 Service：`eauth.efucloud.svc.cluster.local:9001`
- 前端 Service：`eauth-console.efucloud.svc.cluster.local:80`
- 对外访问入口：`docs/frontend.yaml` 中定义的 Ingress

本地验证可使用端口转发：

```bash
kubectl -n efucloud port-forward svc/eauth 9001:9001
kubectl -n efucloud port-forward svc/eauth-console 8000:80
```

## 4. 状态检查

```bash
kubectl -n efucloud get pods,svc,deploy,ingress
kubectl -n efucloud logs -f deploy/eauth
kubectl -n efucloud logs -f deploy/eauth-console
```

## 5. 联通验证

- 前端首页：`http://127.0.0.1:8000`
- 后端健康检查：`http://127.0.0.1:9001/api/health`
- OpenAPI：`http://127.0.0.1:9001/api/v1/swagger.json`
- OIDC 元数据：通过前端域名访问 `/.well-known/openid-configuration`

## 6. 生产环境建议

- 后端配置目前通过 `Secret.stringData` 管理，建议接入 SealedSecrets 或 External Secrets。
- `uploads` 建议改为 PVC 或对象存储，避免 Pod 重建后文件丢失。
- MySQL 建议使用独立实例或 StatefulSet，并配置备份。
- 前端 Ingress 建议启用 TLS，并根据实际域名限制 CORS。
