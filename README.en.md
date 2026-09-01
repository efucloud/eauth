# EAuth

[中文](./README.md) | [English](./README.en.md)

EAuth is a unified authentication platform for enterprise scenarios, positioned similarly to Auth0 and Dex.  
It can act as an OIDC identity provider for business systems and can also integrate with external third-party identity platforms to provide a consistent login experience and centralized authentication policy.
## Frontend Repository

https://github.com/efucloud/eauth-console

## Project Background

In multi-system, multi-organization, and multi-device environments, identity and authentication often face these challenges:

- Account systems are fragmented and login entries are inconsistent.
- Integrating third-party login providers is expensive due to protocol differences.
- Security capabilities (MFA, face recognition, device-side risk controls) are hard to manage centrally.
- Internal enterprise systems and external identity sources lack a unified bridge.

EAuth is designed to provide extensible and production-ready authentication infrastructure that reduces integration cost while improving security posture.

## Core Capabilities

- OIDC provider capabilities (unified token issuing and validation)
- Third-party OIDC/OAuth integration (see `pkg/providers`)
- Face recognition authentication
- MFA (multi-factor authentication)
- Browser fingerprint support for device identification and risk-control scenarios

## Third-Party Auth Extension

Third-party connectors are implemented under `pkg/providers` and can be extended by provider type.  
The current codebase already supports categorized provider registration and dispatch, making it straightforward to add more domestic and international identity providers.

## Directory Overview

- `cmd/`: service entrypoints
- `pkg/apis/`: API layer
- `pkg/services/`: business logic
- `pkg/providers/`: third-party auth integrations
- `config/`: local runtime configuration
- `docs/`: deployment manifests and docs

## Quick Start (Local)

- Go `1.26.x`
- MySQL `8.x` or a compatible version

### 1. Start MySQL

If MySQL is not installed locally, you can use Docker:

```bash
docker run -d \
  --name eauth-mysql \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=EfuCloud \
  -e MYSQL_DATABASE=eauth \
  mysql:8
```

### 2. Update backend config

At minimum, confirm these fields:

- `serverAddress: http://localhost:8000`
- `mysql.host: localhost:3306`
- `mysql.user: root`
- `mysql.password: EfuCloud`
- `mysql.dbname: eauth`

### 3. Start backend

```bash
go run ./cmd -c ./config/config.yaml
```

The backend listens on `9001`.

### 4. Start frontend

```bash
git clone https://github.com/efucloud/eauth-console.git
cd eauth-console
yarn install --frozen-lockfile
yarn start:dev
```

The frontend listens on `8000` and proxies API requests to the backend on `9001`.

### Access URLs

- Home page: `http://127.0.0.1:8000/`
- Health check: `http://127.0.0.1:9001/api/health`
- OpenAPI: `http://127.0.0.1:9001/api/v1/swagger.json`
- OIDC metadata: `http://127.0.0.1:8000/.well-known/openid-configuration`

### Default Admin Account

On first startup, the service creates a default administrator:

- Username: `admin`
- Password: `EfuCloud`

## Local Development

### Backend only

If you only need to work on backend APIs, you can start the service directly:

```bash
go run ./cmd -c ./config/config.yaml
```

This is enough for API development, but the UI must be started from the `eauth-console` repository.

### Full-stack local development

Start the backend in this repository first, then run in the frontend repository:

```bash
yarn install --frozen-lockfile
yarn start:dev
```

Recommended local setup: backend on `9001`, frontend on `8000`.

## Configuration

EAuth loads YAML config via the `-c` flag. Common paths:

- Local: `config/config.yaml`
- Kubernetes: `/efucloud/config/config.yaml` (mounted from Secret in `docs/backend.yaml`)

Recommended approach: copy from `config/config.yaml` for each environment (for example `config/config.prod.yaml`) and start with an explicit config file:

```bash
go run ./cmd -c ./config/config.prod.yaml
```

### Key Config Fields

- `serverAddress`: public system URL used for OAuth/OIDC callback construction
- `tokenPeriod`: token lifetime (hours)
- `uploadPath`: upload directory (for example avatars)
- `loginConfig.mfa`: enable MFA
- `loginConfig.faceRecognition`: enable face recognition
- `mysql.*`: database connection settings
- `email.*`: email service settings (verification codes, notifications)
- `logConfig.*`: logging output and rotation settings

### Minimal Example

```yaml
# serverAddress should be the frontend URL and is used for
# .well-known/openid-configuration and other OIDC metadata
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

## Kubernetes Deployment

Deployment manifests are provided:

- `docs/namespace.yaml`
- `docs/mysql.yaml`
- `docs/backend.yaml`
- `docs/frontend.yaml`

Deployment guide: `docs/README.en.md`.
