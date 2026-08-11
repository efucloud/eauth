# EAuth Kubernetes Deployment Guide

[中文](./README.md) | [English](./README.en.md)

This guide corresponds to the following Kubernetes manifests:

- `docs/namespace.yaml`: Namespace (`efucloud`)
- `docs/mysql.yaml`: MySQL (PVC + ConfigMap + Deployment + Service)
- `docs/backend.yaml`: Unified EAuth service (config Secret + Deployment + Service)
- `docs/frontend.yaml`: Unified-service Ingress (optional)

Default deployment namespace: `efucloud`.

## 1. Prerequisites

- Kubernetes cluster is available and Ingress NGINX is installed (if domain access is needed).
- Image pull credentials are prepared (for private registry).
- A usable StorageClass exists (example: `mysql-local-retain`).
- Update these settings before deployment:
  - `docs/mysql.yaml`: `MYSQL_ROOT_PASSWORD`, `MYSQL_DATABASE`
  - `docs/backend.yaml`: `mysql`, `email`, and `serverAddress` under `Secret.stringData.config.yaml`
  - `docs/frontend.yaml`: domain, IngressClass, and TLS secret

## 2. Recommended Apply Order

```bash
kubectl apply -f docs/namespace.yaml
kubectl apply -f docs/mysql.yaml
kubectl apply -f docs/backend.yaml
kubectl apply -f docs/frontend.yaml
kubectl port-forward -n efucloud svc/eauth 9001:80
```

You can also apply all manifests at once:

```bash
kubectl apply -f docs/
```

## 3. Status Checks

```bash
kubectl -n efucloud get pods,svc,deploy,ingress
```

Backend logs:

```bash
kubectl -n efucloud logs -f deploy/eauth
```

## 4. Connectivity Verification

Unified service access inside the cluster:

- `http://eauth.efucloud.svc.cluster.local`

After local port-forward:

- Console root: `http://127.0.0.1:9001/`
- Health check: `http://127.0.0.1:9001/api/health`
- OpenAPI: `http://127.0.0.1:9001/api/v1/swagger.json`
- OIDC metadata: `http://127.0.0.1:9001/api/.well-known/openid-configuration`

Health check endpoints:

- `GET /api/health`
- `GET /api/v1/swagger.json`
- `GET /metrics`

If you apply the optional Ingress in `docs/frontend.yaml`, the frontend can also be accessed directly from the root path:

- `GET /`
- `GET /.well-known/openid-configuration` (compatibility routing handled by Ingress)

## 5. Production Recommendations

- Service config is currently managed with `Secret.stringData`; consider external secret management (for example SealedSecrets or External Secrets).
- For MySQL, prefer a dedicated instance or StatefulSet and enable backup policies.
- For `uploads`, use PVC or object storage to avoid data loss during pod rescheduling.
- Enforce TLS on Ingress and add observability (access logs, error logs, alerting).
