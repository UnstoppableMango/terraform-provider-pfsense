# Every write applies in the same request

pfSense REST API v2.8.1 only edits `config.xml` for most models, including all four on the must-work list, and the change reaches the running firewall only when the request asks for an **Apply** or a separate `POST /api/v2/<section>/apply` is made ([#56 findings](https://github.com/UnstoppableMango/terraform-provider-pfsense/blob/research/pfsense-apply/docs/research/pfsense-apply.md)).
Applies are asynchronous by default, and the spec declares neither `apply` nor `async` in any request body.
Decided in [#59](https://github.com/UnstoppableMango/terraform-provider-pfsense/issues/59); the write lock and Lookup it builds on are [ADR 0001](0001-resources-are-identified-by-key-not-api-id.md), and the `/apply` endpoints are excluded from the resource list by [ADR 0002](0002-resources-are-derived-from-endpoint-metadata.md).

## Decision

- Every Create and Update on every generated resource sends `"apply": true` and `"async": false` in its JSON body, and every Delete sends `?apply=true&async=false`. No resource is left as a **Staged write**, and there is no apply resource or provider switch.
- The flags are added by the request code that `cmd/patch-provider` generates, not to the spec in `cmd/patch-openapi`, so they never become Terraform attributes. This also covers DELETE `/firewall/rule`, where the spec omits the `apply` parameter that the server honours.
- The flags are sent to every model, including those the spec says apply immediately or where applying doesn't apply; `apply` and `async` are reserved field names, so they never collide with a model field.
- A list in `cmd/gen-config` names models whose writes keep `async: true`, because applying them synchronously breaks the request itself. It starts with `WebGUISettings`, whose apply restarts nginx. Further entries come from testing on the disposable pfSense VM ([#57](https://github.com/UnstoppableMango/terraform-provider-pfsense/issues/57)).
- Only a non-2xx response counts as a failed write.
- When a write fails, state follows what pfSense holds:
  - Create: the provider does a Lookup of the Key. If the object exists, it is saved to state along with the error, so Terraform taints it and replaces it on the next apply. If not, nothing is saved.
  - Update: the prior state is kept, so the next plan proposes the change again.
  - Delete: the state is kept, and a later Delete whose Lookup finds nothing treats the object as deleted.
- The generated client's per-request timeout is 5 minutes, matching the REST API dispatcher's 300 s cap, instead of 30 s.

## Considered options

- A user-declared apply resource such as `pfsense_firewall_apply`: rejected. Terraform has no end-of-apply hook to order it after every staged write, and any later interface, routing, DNS forwarder or IPsec apply reloads the filter and pushes staged firewall changes out anyway.
- A provider-level switch between the two: rejected, two behaviours to test for no must-work case that needs staging.
- Declaring `apply` and `async` in the spec: rejected, tfplugingen would expose them as user-settable attributes, and an unapplied resource is not the declared state.
- Reading pfSense's system notices after each write to catch a failed filter load: rejected, no REST endpoint exposes them at v2.8.1.

## Consequences

- Each write to a firewall-family resource (rule, alias, NAT, shaper) triggers a full synchronous filter reload, and the write lock runs them one at a time, so a plan touching many rules takes that many reloads.
- A ruleset that pf refuses to load is not detected. The firewall dispatcher discards `filter_configure_sync()`'s result, so the response is still 200, pfSense keeps running the previous rules, and the error only appears as a system notice in the GUI. This should be confirmed on the VM.
- One resource's apply also makes live whatever another client staged in the same section.
