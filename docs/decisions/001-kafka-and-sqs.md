# ADR 001: Use Kafka for domain events and SQS for delivery tasks

- **Status:** Proposed
- **Date:** 2026-10-31
- **Decision owners:** Platform engineering

## Context

The platform receives business events and turns some of them into external notifications. We need to distribute domain events to independent consumers while separately buffering provider-specific delivery attempts. Kafka and Amazon SQS can both move messages, but their operational models and primary abstractions differ.

## Decision

For a platform with multiple independent event consumers and a distinct notification delivery workload:

- Use **Kafka** as the domain-event stream. It provides retained ordered partitions, consumer groups, and replay by offset.
- Use **SQS** as the delivery-task queue. It provides managed work distribution, visibility timeout, delayed retry patterns, and dead-letter redrive.
- Persist notification intent and an outbox record transactionally before publishing the SQS task.

This is a conditional architecture choice, not a rule that every notification system needs two brokers.

## Why not use only Kafka?

Kafka can also distribute work among a consumer group and can support retry topics and dead-letter topics. A Kafka-only implementation is viable when the team already operates Kafka confidently and wants one streaming platform. However, queue-specific retry, visibility, and managed redrive behavior must then be designed and operated explicitly.

## Why not use only SQS?

If notification delivery is the only asynchronous consumer, SQS alone is likely simpler and cheaper to operate. Business services can publish a notification command directly to SQS, provided the system does not need a retained domain-event stream for other independent consumers or replay.

## Consequences

### Benefits
- Clear separation between business facts and provider-specific work.
- Independent scaling and failure isolation for event consumption and delivery.
- Kafka replay can rebuild notification intents when the event consumer changes, subject to deduplication and retention.
- SQS buffers provider outages and smooths delivery bursts.

### Costs
- Two messaging systems increase infrastructure, monitoring, permissions, and operational burden.
- The outbox publisher and idempotency handling are mandatory complexity.
- Cross-system delivery is not an atomic exactly-once transaction; duplicate publication is possible and must be safe.

## Revisit when

- There is only one consumer and no event-stream replay requirement: evaluate SQS-only.
- The organization standardizes on Kafka and already has mature retry/DLQ tooling: evaluate Kafka-only.
- Workload, latency, ordering, data residency, or cost requirements materially change.
