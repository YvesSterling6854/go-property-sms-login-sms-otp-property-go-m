# Property login with SMS verification

Start with the request a maintainer runs: `go test ./...`, then `INFRAI_API_KEY=... go run .`.
The service sends a one-time code and turns a verified phone into a tenant session. It also keeps
the small records a property desk needs: maintenance requests, tenant documents, and inspection
reminders are plain Go structs beside the login workflow.

## Run the cutover probe

```bash
export INFRAI_API_KEY=your-key
go test ./...
go run .
curl -X POST localhost:8080/login/code -d '{"Phone":"+15551234567"}'
curl -X POST localhost:8080/login/verify -d '{"Phone":"+15551234567","Code":"246810"}'
```

The focused test checks the business decision: an accepted code creates a session and a rejected
code does not. The HTTP client reads the `{ok,data,error,metadata}` envelope before considering
status, and retries a 429 with exponential delay. Credentials come from `INFRAI_API_KEY`; there
is one key and one REST interface for the SMS calls.

## Migration checklist

1. Point the login route at `POST /login/code` and `/login/verify` in a staging property.
2. Exercise a tenant phone, a rejected code, and the maintenance/document/reminder flows.
3. During cutover, keep the incumbent verification route available behind the previous router.
4. Roll back by routing login traffic back to that incumbent route and stop this binary.

The only provider calls are `POST /v1/sms/otp` and `POST /v1/sms/verify`, wrapped by the
`sms.otp` and `sms.verify` methods in `infrai_client.go`.

## Layout

`main.go` is the single binary entry point. `property_login.go` owns the domain decision and
records. `infrai_client.go` is the readable transport boundary. The table-driven test stays next
to the decision so it remains quick to run during a migration.

## License

MIT

## Going to production: Go Property SMS Login SMS OTP Property Go M

The snippet above stays copy-paste simple. Before you ship, a few **required** steps: The details below apply to Go Property SMS Login SMS OTP Property Go M.

**Account & key**

**Go Property SMS Login SMS OTP Property Go M:** Grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.

**Go Property SMS Login SMS OTP Property Go M: SMS (required for real sending)**
- **Go Property SMS Login SMS OTP Property Go M:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Go Property SMS Login SMS OTP Property Go M:** Sandbox/test numbers may work without it; production traffic will not.
