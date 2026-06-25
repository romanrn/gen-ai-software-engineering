# Bug Context: bug#2 — Incorrect Uptime Reporting

## Summary

The `/health` endpoint reports incorrect uptime values after the service has been running for a period of time.

## Wrong Behavior

After approximately 2 minutes of uptime the reported `uptime_seconds` value becomes negative and continues to decrease instead of increasing.

## Intent

Demonstrates a silent numeric bug that only manifests at runtime after a delay, making it hard to catch in short-lived tests.

## Severity

Low-Medium — service stays functional but monitoring/health checks report incorrect data.