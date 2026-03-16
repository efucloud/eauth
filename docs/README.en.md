# EAuth Kubernetes Deployment Guide

[中文](./README.md) | [English](./README.en.md)

This guide corresponds to the following Kubernetes manifests:

- `docs/namespace.yaml`: Namespace (`efucloud`)
- `docs/mysql.yaml`: MySQL (PVC + ConfigMap + Deployment + Service)
- `docs/backend.yaml`: EAuth backend (Secret + Deployment + Service)
- `docs/frontend.yaml`: EAuth frontend console (Deployment + Service + Ingress)

Default deployment namespace: `efucloud`.

## 1. Prerequisites

- Kubernetes cluster is available and Ingress NGINX is installed (if domain access is needed).
- Image pull credentials are prepared (for private registry).
- A usable StorageClass exists (example: `mysql-local-retain`).
- Update these settings before deployment:
  - `docs/mysql.yaml`: `MYSQL_ROOT_PASSWORD`, `MYSQL_DATABASE`
  - `docs/backend.yaml`: `mysql`, `email`, and `serverAddress` under `Secret.stringData.config.yaml`
  - `docs/frontend.yaml`: domain, TLS secret, frontend image version

## 2. Recommended Apply Order

```bash
kubectl apply -f docs/namespace.yaml
kubectl apply -f docs/mysql.yaml
kubectl apply -f docs/backend.yaml
kubectl apply -f docs/frontend.yaml
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

Frontend logs:

```bash
kubectl -n efucloud logs -f deploy/eauth-console
```

## 4. Connectivity Verification

Backend service access inside cluster:

- `http://eauth.efucloud.svc.cluster.local:9001`

Health check endpoints:

- `GET /api/v1/swagger.json`
- `GET /metrics`

## 5. Production Recommendations

- Backend config is currently managed with `Secret.stringData`; consider external secret management (for example SealedSecrets or External Secrets).
- For MySQL, prefer a dedicated instance or StatefulSet and enable backup policies.
- For `uploads`, use PVC or object storage to avoid data loss during pod rescheduling.
- Enforce TLS on frontend Ingress and add observability (access logs, error logs, alerting).
