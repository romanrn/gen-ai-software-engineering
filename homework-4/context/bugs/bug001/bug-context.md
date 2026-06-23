# Bug Context: bug#1 — Wrong Timezone

## Summary

`GET /time` returns local server time instead of UTC.

## Wrong Behavior

The endpoint responds with the server's local timezone offset instead of UTC.
For example, on a server in UTC+3 the returned time is 3 hours ahead of actual UTC.

## Intent

Demonstrates a common timezone bug that is hard to detect without cross-timezone testing.

## Severity

Medium — incorrect data returned to clients, hard to detect without cross-timezone testing.