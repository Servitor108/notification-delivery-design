# Runbook: Retry and dead-letter queue operations

## Purpose

Use this runbook to investigate delayed notification delivery, repeated provider failures, and messages moved to the SQS dead-letter queue (DLQ). Follow the team's incident process and access controls. Do not replay messages blindly.

## 1. Initial triage

1. Check notification delivery dashboards: queue depth, oldest-message age, worker concurrency, provider error rate, and DLQ depth.
2. Check Kafka consumer lag and outbox backlog to locate the stage where work is accumulating.
3. Group failures by provider, channel, error class, tenant, and time window.
4. Check provider status pages and rate-limit responses.
5. Confirm whether messages are still within their business expiry window.

## 2. Classify the failure

| Failure class | Typical action |
|---|---|
| Throttling / HTTP 429 | Reduce concurrency, respect Retry-After, retry with jitter |
| Provider 5xx / transient network issue | Check provider health; bounded exponential retry |
| Timeout with unknown provider outcome | Reconcile by provider idempotency key/reference before retry |
| Invalid recipient / malformed request | Do not retry unchanged; correct source data or mark terminal |
| Unsupported template / configuration | Fix template or routing configuration, then selectively replay |
| Expired notification | Do not replay unless product owner explicitly approves |

## 3. DLQ investigation

For each sample message:
- Capture task ID, notification ID, correlation ID, receive count, and failure classification.
- Inspect metadata only; avoid copying recipient details or message content into tickets or chat.
- Find the corresponding notification and attempt history in the database.
- Determine whether a prior provider attempt may already have succeeded.
- Identify whether the root cause is systemic or isolated.

## 4. Safe replay checklist

Before replay:
- [ ] Root cause is understood and mitigated.
- [ ] Message is not expired, or an explicit product decision authorizes replay.
- [ ] Notification state permits another attempt.
- [ ] Provider status/idempotency lookup has been checked for ambiguous outcomes.
- [ ] Duplicate-send risk is understood and accepted.
- [ ] Replay is limited to a reviewed cohort or bounded batch.
- [ ] Operator, time, reason, and count are audit logged.

Replay through the approved operational tool or controlled redrive workflow. Never manually edit production queue payloads or bypass notification state checks.

## 5. Recovery verification

- Confirm DLQ depth is decreasing as expected.
- Confirm queue age and delivery latency recover.
- Watch provider throttling and worker error rates.
- Verify delivered/failed state reconciliation from callbacks.
- Stop replay if duplicate suppression, provider error rate, or queue age regresses.

## 6. Escalation

Escalate to the provider owner for sustained provider errors, to platform/on-call for queue or consumer backlog, and to the product owner for expiry exceptions or user-impact decisions.

## 7. Follow-up

Record root cause, affected time range, notification count, customer impact, replay decision, and prevention actions. Add a test for any newly discovered failure mode.
