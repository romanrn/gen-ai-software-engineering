# Bug Context: sec#1 — Exposed API Secret

## Summary

The API secret used to authenticate requests is visible in the source code repository.

## Wrong Behavior

Any developer with read access to the repository can obtain the key and make authenticated requests.
The secret cannot be rotated without modifying and redeploying the application.

## Intent

Demonstrates OWASP A07:2021 — secret exposure through version control.

## Severity

HIGH — credentials exposed in version control.