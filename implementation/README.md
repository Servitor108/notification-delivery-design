# Notification Delivery — Go Reference Implementation

This directory contains runnable Go code that complements the parent system-design proposal. It is a deliberately small vertical slice, not a claim that Kafka, SQS, PostgreSQL, provider callbacks, or a production DLQ are implemented here.

## Included

- Go HTTP API to enqueue a notification and query its status.
- Idempotency-key handling with conflict detection for mismatched requests.
- A repository interface and concurrency-safe in-memory reference repository.
- Bounded concurrent worker loop and provider adapter interface.
- HTTP provider adapter with request idempotency header.
- Retry classification for timeouts, throttling, and server errors.
- Exponential backoff with full jitter.
- Unit tests for idempotency, successful delivery, and backoff.

The in-memory repository is intentionally replaceable. A production adapter should implement the same contract with PostgreSQL transactions and durable state; queue adapters can then publish outbox records to SQS. Kafka ingestion belongs in a separate adapter/process.

## Run

Requires Go 1.22+.

```bash
cd implementation
go test -race ./...
go run ./cmd/api
```

## API

Create a notification:

```bash
curl -X POST http://localhost:8080/v1/notifications \
  -H 'Content-Type: application/json' \
  -H 'Idempotency-Key: order-42-email-receipt' \
  -d '{"recipient_ref":"customer-42","channel":"email","template_id":"order-receipt-v1","payload":{"order_id":"42"}}'
```

Fetch status using the returned ID:

```bash
curl http://localhost:8080/v1/notifications/NOTIFICATION_ID
```

The demo provider logs acceptance and does not send a real email/SMS/WhatsApp message. Do not put real recipient PII into the demo service.

## Important limitations

- State is in memory and is lost on restart; there is no durable queue in this slice.
- Demo provider is a stub; the HTTP provider adapter is an extension point and requires a real provider contract.
- Provider callback reconciliation, PostgreSQL/outbox, Kafka, SQS, DLQ/replay, metrics, tracing, auth, and tenant isolation remain future integration work.
- This is portfolio sample code, not production-ready notification infrastructure.
