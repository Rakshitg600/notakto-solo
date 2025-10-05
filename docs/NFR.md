# Non-Functional Requirements (NFRs)

This document defines the system’s non-functional requirements. These guide architecture, implementation, and operations to ensure the service is production-ready.

---

## 1. Performance
- **Latency**:  
  - p50: < 100 ms  
  - p95: < 250 ms  
  - p99: < 500 ms (excluding network latency)
- **Throughput**:  
  - Minimum 500 requests/second sustained.  
  - System must degrade gracefully under overload (e.g., return 429s).

## 2. Scalability
- Horizontal scaling through stateless app servers.  
- Session store must support **10k+ concurrent sessions** in memory; Redis should scale further.  
- DB connection pool tuned for concurrency (target: 100 active connections per instance).

## 3. Availability & Reliability
- **SLO**: 99.9% monthly availability.  
- **RTO (Recovery Time Objective)**: ≤ 15 minutes.  
- **RPO (Recovery Point Objective)**: ≤ 1 minute (via WAL shipping or PITR backups for Postgres).  
- Graceful degradation: if Redis fails, fall back to in-memory with reduced durability.

## 4. Security
- All traffic via HTTPS (TLS 1.2+).  
- Authentication via Firebase tokens (short-lived).  
- Rate limiting at IP and UID levels to mitigate abuse.  
- DB user with least privilege.  
- Secrets stored in secret manager, never in source code.  
- Audit logging for sensitive operations.

## 5. Maintainability
- Modular folder structure (`routes/`, `handlers/`, `sessions/`, `fxns/`).  
- Clear separation of concerns: middleware, business logic, persistence.  
- CI pipeline with lint, tests, and security scans.  
- Well-documented configs and APIs.

## 6. Observability
- **Metrics**: request latency, RPS, error rates, active sessions, DB pool usage, rate limiter stats.  
- **Tracing**: all requests instrumented via OpenTelemetry with correlation IDs.  
- **Logging**: structured JSON logs including request-id, uid, error codes.

## 7. Monitoring & Alerting
- Alerts on: error rate > 1%, latency p95 > 250ms, DB connection exhaustion, Redis unavailability, session churn spikes.  
- Dashboards in Grafana for real-time monitoring.

## 8. Data Retention
- Session state: in-memory/Redis TTL = 15 minutes.  
- Session events: persisted in Postgres for **90 days**.  
- Logs: retained for **30 days** in logging backend.  
- Metrics & traces: retained per monitoring system (e.g., 14–30 days).

## 9. Compliance & Privacy
- UIDs treated as non-PII identifiers; no raw tokens or personal info logged.  
- Logs sanitized to avoid leaking secrets.  
- Regular dependency scanning for vulnerabilities (e.g., govulncheck).

---

## Summary

The service is designed to be **low-latency, horizontally scalable, secure, and observable**. It prioritizes developer productivity (via maintainable structure) and production readiness (via monitoring, alerting, and resilience targets).
