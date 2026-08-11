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
- `pkg/apis/`: API layer
- `pkg/services/`: business logic
- `pkg/providers/`: third-party auth integrations
- `config/`: local runtime configuration
- `docs/`: deployment manifests and docs

## Quick Start (Local)

1. Prepare MySQL and update `config/config.yaml`.
2. Start the service:

```bash
go run ./cmd -c ./config/config.yaml
```

Default service port: `9001`.

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
- `docs/backend.yaml`: unified service Deployment + Service
- `docs/frontend.yaml`: unified service Ingress (optional)

Deployment guide: `docs/README.en.md`.
