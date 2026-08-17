---
title: "Deployment Validation Report"
category: deployment
validator: validator-deployment
date: 2026-03-20
---

# Deployment Fact Validation — docker-nats-setup.md

## Claims Validated

### 1. `nats:2.10-alpine` Docker Image
**Status: CONFIRMED**

The image tag `nats:2.10-alpine` exists on Docker Hub as an official image.
Docker Hub layer detail confirms: `hub.docker.com/layers/library/nats/2.10-alpine/images/sha256-c1c3f0bceb82215bc89c429d235dfd564a0154336df7611788915005d5d22bfd`

More specific patch-version tags also exist (e.g. `nats:2.10.20-alpine`, `nats:2.10.29-alpine`).

**Source:** https://hub.docker.com/_/nats/tags

---

### 2. JetStream `--js` Flag
**Status: CONFIRMED (with minor notation)**

JetStream is enabled via the `-js` flag passed to the NATS server binary. Official NATS docs show:
```bash
docker run -p 4222:4222 nats -js
```

The document uses `--js` (double dash). NATS server CLI accepts both `-js` and `--js` as equivalent; both are valid. No functional difference.

**Source:** https://docs.nats.io/running-a-nats-service/nats_docker/jetstream_docker

---

### 3. `/healthz` Endpoint on Port 8222
**Status: CONFIRMED**

NATS exposes an HTTP monitoring server on port 8222. The `/healthz` endpoint exists and returns `{ "status": "ok" }` when the server can accept connections. This is a documented, official monitoring endpoint.

Full list of monitoring endpoints on port 8222 includes: `/varz`, `/connz`, `/routez`, `/gatewayz`, `/leafz`, `/subsz`, `/accountz`, `/accstatz`, `/jsz`, `/healthz`.

**Source:** https://docs.nats.io/running-a-nats-service/nats_admin/monitoring

---

### 4. Docker Compose `condition: service_healthy`
**Status: CONFIRMED**

Docker Compose supports `condition: service_healthy` in the `depends_on` block. The app service waiting on the NATS healthcheck is valid Docker Compose syntax:

```yaml
depends_on:
  nats:
    condition: service_healthy
```

Compose waits until the dependency's healthcheck passes before starting the dependent service.

**Source:** https://docs.docker.com/compose/how-tos/startup-order/

---

### 5. Port Mapping `4229:4222` Syntax
**Status: CONFIRMED**

The Docker Compose port mapping syntax `HOST:CONTAINER` is correct. `4229:4222` maps external host port 4229 to internal container port 4222. This is standard Docker Compose syntax.

**Source:** https://docs.docker.com/reference/compose-file/services/

---

## Summary

| Claim | Status | Notes |
|-------|--------|-------|
| `nats:2.10-alpine` image exists | ✅ CONFIRMED | Valid Docker Hub official tag |
| `--js` enables JetStream | ✅ CONFIRMED | `-js` and `--js` are both accepted |
| `/healthz` on port 8222 | ✅ CONFIRMED | Returns `{"status":"ok"}` |
| `condition: service_healthy` syntax | ✅ CONFIRMED | Standard Docker Compose feature |
| Port mapping `4229:4222` syntax | ✅ CONFIRMED | Standard Docker Compose `HOST:CONTAINER` |

**All facts in `docker-nats-setup.md` are accurate.** No corrections needed.

---

## Sources
- [nats — Official Image | Docker Hub](https://hub.docker.com/_/nats)
- [nats Tags | Docker Hub](https://hub.docker.com/_/nats/tags)
- [JetStream | NATS Docs](https://docs.nats.io/running-a-nats-service/nats_docker/jetstream_docker)
- [Monitoring | NATS Docs](https://docs.nats.io/running-a-nats-service/nats_admin/monitoring)
- [Control startup order — Docker Compose](https://docs.docker.com/compose/how-tos/startup-order/)
- [Define services in Docker Compose](https://docs.docker.com/reference/compose-file/services/)
