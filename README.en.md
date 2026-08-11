# EAuth

[中文](./README.md) | [English](./README.en.md)

EAuth is a unified authentication platform for enterprise scenarios, positioned similarly to Auth0 and Dex.  
It can act as an OIDC identity provider for business systems and can also integrate with external third-party identity platforms to provide a consistent login experience and centralized authentication policy.

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
- `web/`: frontend source
- `pkg/apis/`: API layer
- `pkg/services/`: business logic
- `pkg/providers/`: third-party auth integrations
- `pkg/embeds/web/`: embedded frontend build output
- `config/`: local runtime configuration
- `docs/`: deployment manifests and docs

## Quick Start (Local)

### Requirements

- Go `1.26.x`
- Node.js `20.x`
- Yarn `1.x`
- MySQL `8.x` or a compatible version

### Startup Steps

1. Start a local MySQL instance and create the `eauth` database

If MySQL is not installed locally, you can use Docker:

```bash
docker run -d \
  --name eauth-mysql \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=EfuCloud \
  -e MYSQL_DATABASE=eauth \
  mysql:8
```

2. Review and update `config/config.yaml` if needed

At minimum, confirm these fields:

- `serverAddress: http://localhost:9001`
- `mysql.host: localhost:3306`
- `mysql.user: root`
- `mysql.password: EfuCloud`
- `mysql.dbname: eauth`

3. Build and embed the frontend assets:

```bash
./scripts/build-web-embed.sh
```

Or:

```bash
make embed-web
```

4. Start the unified service:

```bash
go run ./cmd -c ./config/config.yaml
```

Default service port: `9001`.
After startup, the backend serves both the API and the embedded frontend UI.

### Access URLs

- Home page: `http://127.0.0.1:9001/`
- Health check: `http://127.0.0.1:9001/api/health`
- OpenAPI: `http://127.0.0.1:9001/api/v1/swagger.json`
- OIDC metadata: `http://127.0.0.1:9001/api/.well-known/openid-configuration`

### Default Admin Account

On first startup, the service creates a default administrator:

- Username: `admin`
- Password: `EfuCloud`

## Local Development

### Backend Development

If you only need to work on backend APIs, you can start the service directly:

```bash
go run ./cmd -c ./config/config.yaml
```

If you changed frontend code and want the backend to keep serving the latest UI, re-run:

```bash
./scripts/build-web-embed.sh
```

### Unified Local Preview

If you want the backend to serve the frontend UI locally, run:

```bash
./scripts/build-web-embed.sh
```

This script does three things:

- installs frontend dependencies under `web/`
- builds the frontend into `web/dist`
- syncs the build output into `pkg/embeds/web/` so the Go service can embed and serve it

Use it when:

- you have just cloned the repository and want the local backend to serve the UI
- you changed frontend code under `web/` and need fresh embedded assets
- you are preparing a local or release image build and need the embedded assets up to date

After frontend changes, re-run:

```bash
./scripts/build-web-embed.sh
go run ./cmd -c ./config/config.yaml
```

### Frontend-Only Development

If you only want to work on the frontend UI, you can start the frontend dev server separately:

```bash
cd web
yarn install --frozen-lockfile
yarn start:dev
```

For integrated testing, make sure the backend is also running locally on port `9001`.

### Frontend Paths

- frontend source: `web/`
- frontend build output: `web/dist/`
- backend embedded assets: `pkg/embeds/web/`

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

### Frontend Integration

- Frontend source lives under `web/`
- Production images and the local unified service serve the UI from `pkg/embeds/web/`
- Files under `pkg/embeds/web/` are generated build artifacts; the repository only keeps the `.ignore` placeholder

### Minimal Example

```yaml
# serverAddress should be the public URL of the unified service and is used for
# .well-known/openid-configuration and other OIDC metadata
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

## Kubernetes Deployment

Deployment manifests are provided:

- `docs/namespace.yaml`
- `docs/mysql.yaml`
- `docs/backend.yaml`: unified service Deployment + Service + optional Ingress

Deployment guide: `docs/README.en.md`.
