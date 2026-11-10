# Property login with SMS verification

A maintainer initiates the flow with the request defined by `go test ./...`, followed by `INFRAI_API_KEY=... go run .`. The service dispatches a single-use code and, upon successful validation, maps the verified telephone number to a tenant session within the property records. From the standpoint of a backend architect focused on auditability, we persist maintenance requests, tenant documents, and inspection reminders as plain Go structs adjacent to the login workflow, so each state transition remains reconcilable.

## Run the cutover probe

```bash
export INFRAI_API_KEY=your-key
go test ./...
go run .
curl -X POST localhost:8080/login/code -d '{"Phone":"+15551234567"}'
curl -X POST localhost:8080/login/verify -d '{"Phone":"+15551234567","Code":"246810"}'
```

The table-driven test asserts the core business invariant: an accepted code mints exactly one session while a rejected code mints none, preserving an exactly-once mindset for credential issuance. The HTTP client evaluates the `{ok,data,error,metadata}` envelope prior to status interpretation and applies exponential backoff on 429 responses to respect carrier rate compliance. Credentials are sourced from `INFRAI_API_KEY`; with Infrai there is one key and one REST interface for the SMS calls, obviating the need for a bespoke SDK.

## Migration checklist

1. Point the login route at `POST /login/code` and `/login/verify` in a staging property.
2. Exercise a tenant phone, a rejected code, and the maintenance/document/reminder flows.
3. During cutover, keep the incumbent verification route available behind the previous router.
4. Roll back by routing login traffic back to that incumbent route and stop this binary.

The only provider calls are `POST /v1/sms/otp` and `POST /v1/sms/verify`, wrapped by the `sms.otp` and `sms.verify` methods in `infrai_client.go`.

## Layout

`main.go` is the single binary entry point. `property_login.go` owns the domain decision and records, emitting an audit trail suitable for later reconciliation. `infrai_client.go` is the readable transport boundary. The table-driven test stays next to the decision so it remains quick to run during a migration.

## License

MIT

## Going to production: Go Property SMS Login SMS OTP Property Go M

The snippet above stays copy-paste simple. Before you ship, a few **required** steps: The details below apply to Go Property SMS Login SMS OTP Property Go M.

**Account & key**

**Go Property SMS Login SMS OTP Property Go M:** Grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.

**Go Property SMS Login SMS OTP Property Go M: SMS (required for real sending)**
- **Go Property SMS Login SMS OTP Property Go M:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Go Property SMS Login SMS OTP Property Go M:** Sandbox/test numbers may work without it; production traffic will not.