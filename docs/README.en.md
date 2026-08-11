# EAuth Kubernetes Deployment Guide

[中文](./README.md) | [English](./README.en.md)

The deployment model is now split by service:

- `docs/namespace.yaml`: namespace (`efucloud`)
- `docs/mysql.yaml`: MySQL (PVC + ConfigMap + Deployment + Service)
- `docs/backend.yaml`: EAuth backend (Secret + Deployment + Service)
- `docs/frontend.yaml`: EAuth frontend console (Secret + Deployment + Service + Ingress)

Default deployment namespace: `efucloud`.

## 1. Prerequisites

- Kubernetes is available.
- Install Ingress NGINX if public domain access is required.
- Prepare image pull credentials if you use a private registry.
- Update these files before deployment:
  - `docs/mysql.yaml`: `MYSQL_ROOT_PASSWORD`, `MYSQL_DATABASE`
  - `docs/backend.yaml`: `mysql`, `email`, and `serverAddress` in `Secret.stringData.config.yaml`
  - `docs/frontend.yaml`: frontend image, domain, TLS settings, and Nginx proxy settings

`serverAddress` should be the public frontend URL, not the backend Service address.

## 2. Recommended Apply Order

```bash
kubectl apply -f docs/namespace.yaml
kubectl apply -f docs/mysql.yaml
kubectl apply -f docs/backend.yaml
kubectl apply -f docs/frontend.yaml
```

Once environment-specific values are updated, you can also apply everything at once:

```bash
kubectl apply -f docs/
```

## 3. Access Model

- Backend Service: `eauth.efucloud.svc.cluster.local:9001`
- Frontend Service: `eauth-console.efucloud.svc.cluster.local:80`
- Public entrypoint: the Ingress defined in `docs/frontend.yaml`

For local verification, use port-forwarding:

```bash
kubectl -n efucloud port-forward svc/eauth 9001:9001
kubectl -n efucloud port-forward svc/eauth-console 8000:80
```

## 4. Status Checks

```bash
kubectl -n efucloud get pods,svc,deploy,ingress
kubectl -n efucloud logs -f deploy/eauth
kubectl -n efucloud logs -f deploy/eauth-console
```

## 5. Connectivity Verification

- Frontend home: `http://127.0.0.1:8000`
- Backend health check: `http://127.0.0.1:9001/api/health`
- OpenAPI: `http://127.0.0.1:9001/api/v1/swagger.json`
- OIDC metadata: access `/.well-known/openid-configuration` through the frontend domain

## 6. Production Recommendations

- Backend config is currently managed through `Secret.stringData`; consider SealedSecrets or External Secrets.
- Use PVC or object storage for `uploads` to avoid data loss after pod rescheduling.
- Prefer a dedicated MySQL instance or StatefulSet with backups enabled.
- Enforce TLS on the frontend Ingress and restrict CORS to trusted domains.
