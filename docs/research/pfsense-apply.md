# Which pfSense REST API writes need a separate apply step?

Research for [#56](https://github.com/UnstoppableMango/terraform-provider-pfsense/issues/56), part of map #54.

## Question

When does a write through REST API v2.8.1 take effect on the running firewall, and what does the provider have to call to make it happen?

## Answer

- A write to a `config_path` model **only edits `config.xml`** unless the model sets `always_apply`; the service is reloaded only when the request asks for it (`"apply": true`) or a separate `POST /api/v2/<section>/apply` is made.
- Every model on the must-work list (firewall rule, firewall alias, NAT port forward, DHCP static mapping) is staged by default, and every one of them **can apply in the same request**: `"apply": true` in the POST/PATCH JSON body, or `?apply=true` on DELETE. The spec documents this only partly (see below).
- `PUT` (replace all) and plural `DELETE` (delete many / all) **never apply**, whatever the request says.
- Applying is **asynchronous by default** (`async` defaults to `true`): the call returns once a background process is spawned. Send `"async": false` to block until the reload finishes.
- `GET` on an apply endpoint reports pending state, but **it is not reliable for the firewall section**: firewall rules never mark anything pending, and `POST /firewall/apply` never clears the firewall pending flags it reports.
- Sections are not isolated: the interface, routing, DNS forwarder and IPsec apply paths all reload the packet filter, which picks up every staged firewall/NAT/alias change, and the interface apply also applies and clears pending static routes.

Recommendation for the provider: on create/update send `"apply": true` and `"async": false`; on single-object delete send `?apply=true&async=false`. Do not rely on `GET /…/apply` to detect drift in firewall rules.

## Sources

Highest trust first. All package paths are relative to `pfSense-pkg-RESTAPI/files/usr/local/pkg/RESTAPI/` in [pfrest/pfSense-pkg-RESTAPI](https://github.com/pfrest/pfSense-pkg-RESTAPI) at tag `v2.8.1` (commit `4190c8ca40c73e5705e01f227a029b61dedf7b7f`).

1. Package source at `v2.8.1`: `Core/Model.inc`, `Core/Endpoint.inc`, `Core/Dispatcher.inc`, `Core/ModelSet.inc`, `Core/Globals.inc`, `ContentHandlers/URLContentHandler.inc`, `Schemas/OpenAPISchema.inc`, `Models/*Apply.inc`, `Dispatchers/*ApplyDispatcher.inc`, and the must-work models.
2. The pinned spec: `https://github.com/pfrest/pfSense-pkg-RESTAPI/releases/download/v2.8.1/openapi.json` (`info.version` = `v2.8.1`).
3. The package's own docs at the same tag: `docs/COMMON_CONTROL_PARAMETERS.md`, `docs/SWAGGER_AND_OPENAPI.md` (published at pfrest.org).
4. pfSense CE `master` (commit `9363ac5b`, 2026-03-31) for `filter_configure()`; no 2.8 release branch is public, so line numbers there may drift from the shipped 2.8.x.

Everything marked **Verified** was read in source; **Inference** marks conclusions drawn from source that were not exercised on a running firewall.

## 1. How the Model base class stages and applies

**Verified.**

- `$subsystem` names the pfSense "dirty" flag a model marks when it writes config; it is cleared after `apply()` runs (`Core/Model.inc:274-282`). `write_config` marks it (`Core/Model.inc:772-775`).
- `$always_apply` forces `apply()` after every create, update, delete and replace (`Core/Model.inc:284-287`).
- `create()`: `_create()` writes config, marks the subsystem dirty, runs `pre_apply_create()`, then runs `apply_create()` and clears the dirty flag **only if `$apply or $this->always_apply`** (`Core/Model.inc:2252-2290`, condition at `2276`).
- `update()` is the same (`Core/Model.inc:2344-2385`, condition at `2375`), as are `replace_all()` (`2400`, condition at `2492`) and `delete()` (`2532`, condition at `2582`).
- `apply_create/update/delete/replace_all` default to `apply()` (`Core/Model.inc:1804-1878`); the base `apply()` is empty (`1813`).
- The dirty flag is cleared by `Model` immediately after `apply_*()` returns (e.g. `Core/Model.inc:2278`), even when that call only spawned a background process.

### Where `apply` comes from in a request

**Verified.**

- `Endpoint::post()` calls `create(apply: $this->request_data['apply'] === true)` (`Core/Endpoint.inc:1226`).
- `Endpoint::patch()` calls `update(apply: $this->request_data['apply'] === true, …)` (`Core/Endpoint.inc:1250-1254`).
- `Endpoint::delete()` on a singular endpoint calls `delete(apply: $this->request_data['apply'] === true)` (`Core/Endpoint.inc:1280`).
- `Endpoint::put()` calls `replace_all(data: …)` with no `apply` argument, so it defaults to `false` (`Core/Endpoint.inc:1232-1233`).
- Plural DELETE calls `delete_all()` / `delete_many()` (`Core/Endpoint.inc:1260-1276`), which call `ModelSet::delete()`, which calls each model's `delete()` with no `apply` (`Core/ModelSet.inc:158-162`, `run_model_method` at `174-191`).
- The comparison is strict `=== true`. Query strings are decoded by `URLContentHandler::infer_type`, which turns the literal `true` into a boolean (`ContentHandlers/URLContentHandler.inc:46-63`), so `?apply=true` works on DELETE; a JSON string `"true"` in a body would not.
- `apply` and `async` are reserved field names (`Core/Globals.inc:11-12`), so they never collide with a model field.

So: **POST, PATCH and singular DELETE can apply in the same request; PUT and plural DELETE cannot.** A model with `always_apply` applies on every write regardless.

## 2. Synchronous or asynchronous

**Verified.**

- `Endpoint::$async` defaults to `true` (`Core/Endpoint.inc:303-308`); a request can override it with a boolean `async` field (`validate_async`, `Core/Endpoint.inc:773-787`), and the endpoint copies it onto the model (`Core/Endpoint.inc:1123`). Models built in code also default to async (`Core/Model.inc:427-428`).
- Every section's `apply()` spawns a Dispatcher with `async: $this->async`. `Dispatcher::spawn_process()` with `async` true starts `dispatch.sh` under `nohup` in the background and returns the PID; with `async` false it calls `process()` in the request (`Core/Dispatcher.inc:209-263`, branch at `227`, sync call at `257`). The background run is capped by `timeout` = 300 s (`Core/Dispatcher.inc:57`) and at most `max_queue` = 10 queued processes, beyond which the call fails with `503 DISPATCHER_TOO_MANY_QUEUED_PROCESSES` (`Core/Dispatcher.inc:63`, `216-224`).
- The firewall reload is doubly asynchronous by default: `FirewallApplyDispatcher` calls `filter_configure()` when async and `filter_configure_sync()` otherwise (`Dispatchers/FirewallApplyDispatcher.inc:15-19`), and pfSense's `filter_configure()` only sends a `filter reload` event to `check_reload_status` (pfSense `src/etc/inc/filter.inc:798-811`). With `async: false` the rules are loaded before the response is sent.
- The package docs say the same: `async` defaults to `true`, and the recommended way to wait is to poll the apply endpoint with GET (`docs/COMMON_CONTROL_PARAMETERS.md:45-57`).

**Inference:** with the default `async: true`, a successful 200 means only "reload queued"; a failing reload (e.g. a pf syntax error) is not reported to the caller. The provider should send `async: false` so that `apply` errors surface in the response and a subsequent read observes the applied state.

## 3. The 10 apply endpoints

**Verified.** Each `Endpoints/*ApplyEndpoint.inc` allows only `GET` and `POST`, and each `Models/*Apply.inc` sets `always_apply = true`, overrides `_create()` to do nothing, and exposes a read-only `applied` boolean from an internal callable. A POST therefore just runs `apply()`; the request body carries nothing (apart from `async`).

| Endpoint | Model (`Models/`) | Pending when (`get_applied_status`) | Dispatcher work (`Dispatchers/`) |
| --- | --- | --- | --- |
| `/api/v2/firewall/apply` | `FirewallApply.inc` | any of `aliases`, `natconf`, `filter`, `shaper` dirty (`:17`, `:53-67`); also returns `pending_subsystems` | `filter_configure()` / `filter_configure_sync()`; clears nothing (`FirewallApplyDispatcher.inc:15-19`) |
| `/api/v2/firewall/virtual_ip/apply` | `VirtualIPApply.inc` | `/tmp/.firewall_virtual_ip.apply` exists (`:38-39`) | `applyVIP()` (`VirtualIPApplyDispatcher.inc:16-18`) |
| `/api/v2/interface/apply` | `InterfaceApply.inc` | `{tmp}/.interfaces.apply` exists; returns `pending_interfaces` (`:50-68`) | reconfigures each pending interface, then `filter_configure()`, `system_routing_configure()` if `staticroutes` dirty, clears `staticroutes` and `interfaces` (`InterfaceApplyDispatcher.inc:14-82`) |
| `/api/v2/routing/apply` | `RoutingApply.inc` | `staticroutes` dirty (`:37-38`) | runs `.system_routes.apply`, `system_routing_configure()`, `system_resolvconf_generate()`, `filter_configure()`, `filter reload` event, resyncs dirty gateway groups; clears `staticroutes` (`RoutingApplyDispatcher.inc:15-49`) |
| `/api/v2/services/dhcp_server/apply` | `DHCPServerApply.inc` | `dhcpd` dirty (`:36-37`) | `dhcp_apply_changes()`; clears `dhcpd` (`DHCPServerApplyDispatcher.inc:15-18`) |
| `/api/v2/services/dns_forwarder/apply` | `DNSForwarderApply.inc` | `hosts` dirty (`:38-39`) | `services_dnsmasq_configure()`, `filter_configure()`, resolv.conf, dhcpleases; clears `hosts` (`DNSForwarderApplyDispatcher.inc:11-20`) |
| `/api/v2/services/dns_resolver/apply` | `DNSResolverApply.inc` | `unbound` dirty (`:38-39`) | `services_unbound_configure()`; clears `unbound` only on success (`DNSResolverApplyDispatcher.inc:14-22`) |
| `/api/v2/services/haproxy/apply` | `HAProxyApply.inc` | `/var/run/haproxy.conf.dirty` exists (`:40-41`) | `haproxy_check_and_run()`; removes the file on success (`HAProxyApplyDispatcher.inc:19-26`) |
| `/api/v2/vpn/ipsec/apply` | `IPsecApply.inc` | `ipsec` dirty (`:39-40`) | `ipsec_configure()`, package hook, `filter_configure()`; clears `ipsec` (`IPsecApplyDispatcher.inc:14-24`) |
| `/api/v2/vpn/wireguard/apply` | `WireGuardApply.inc` | `wireguard` dirty (`:39-40`) | `wg_toggle_wireguard()`, syncs pending tunnels; clears `wireguard` (`WireGuardApplyDispatcher.inc:19-30`) |

### What they return

**Verified:** both GET and POST return the apply model's representation, `{"applied": bool}` plus `pending_subsystems` (firewall) or `pending_interfaces` (interface); the spec's `FirewallApply` schema has exactly `applied` and `pending_subsystems`. `FirewallApply::apply()` re-reads its status after spawning (`Models/FirewallApply.inc:76-79`); the others return the status read when the model was constructed, before the apply ran.

**Inference:** the `applied` value in a POST response is not a confirmation. Re-GET after a synchronous apply if confirmation is needed.

## 4. Which writes are staged and which apply immediately

**Verified.** The spec's operation descriptions carry an `**Applies immediately**` line, computed by `Model::does_apply_immediately()`: `null` if the model has no `config_path`, `true` if `always_apply`, otherwise `false` (`Core/Model.inc:1315-1328`; rendered at `Schemas/OpenAPISchema.inc:141-146`, `163`). Across the 450 write operations in the pinned spec: 163 say `Yes`, 240 say `No`, 47 say `Not Applicable`.

`No` is a superset of "needs an apply endpoint": it also covers config-only models with no reload at all (e.g. `AuthServer`, `Certificate`, `EmailNotificationSettings`, `RESTAPIAccessListEntry`). Cross-referencing the models' `$subsystem` and `apply()` gives the staged groups, keyed by the apply endpoint that finishes them:

| Apply endpoint | Staged models (default `apply: false`) |
| --- | --- |
| `firewall/apply` | `FirewallRule` (no subsystem), `FirewallAlias` (`aliases`), `PortForward`, `OutboundNATMapping`, `OneToOneNATMapping` (`natconf`), `OutboundNATMode`, `FirewallAdvancedSettings` (no subsystem), `TrafficShaper`, `TrafficShaperQueue` (`shaper`) |
| `firewall/virtual_ip/apply` | `VirtualIP` (`vip`, plus `.firewall_virtual_ip.apply`) |
| `interface/apply` | `NetworkInterface` (`interfaces`, plus `.interfaces.apply`) |
| `routing/apply` | `StaticRoute`, `RoutingGateway`, `RoutingGatewayGroup`, `RoutingGatewayGroupPriority`, `DefaultGateway` (`staticroutes`) |
| `services/dhcp_server/apply` | `DHCPServer`, `DHCPServerAddressPool`, `DHCPServerCustomOption`, `DHCPServerStaticMapping` (`dhcpd`) |
| `services/dns_forwarder/apply` | `DNSForwarderHostOverride`, `DNSForwarderHostOverrideAlias` (`hosts`) |
| `services/dns_resolver/apply` | `DNSResolverSettings`, `DNSResolverHostOverride`, `DNSResolverHostOverrideAlias`, `DNSResolverDomainOverride`, `DNSResolverAccessList`, `DNSResolverAccessListNetwork` (`unbound`) |
| `services/haproxy/apply` | all `HAProxy*` models except `HAProxyApply` (touch `/var/run/haproxy.conf.dirty`) |
| `vpn/ipsec/apply` | `IPsecPhase1`, `IPsecPhase1Encryption`, `IPsecPhase2`, `IPsecPhase2Encryption` (`ipsec`) |
| `vpn/wireguard/apply` | `WireGuardSettings`, `WireGuardTunnel`, `WireGuardPeer` (`wireguard`), `WireGuardTunnelAddress`, `WireGuardPeerAllowedIP` (no subsystem) |
| none | `TrafficShaperLimiter`, `TrafficShaperLimiterQueue`, `TrafficShaperLimiterBandwidth` (own `TrafficShaperApplyDispatcher`, apply only via `apply: true`), `SystemTunable` (`sysctl`, no `apply()`) |

Models that apply on every write (`Applies immediately: Yes`, i.e. `always_apply`) include users and groups, system hostname/DNS/timezone/console/SSH/webGUI/log settings, NTP, cron, certificates and CAs generated or signed, VLAN/LAGG/GRE/bridge/group interfaces, OpenVPN servers/clients/overrides, the DHCP relay and backend, firewall schedules and state-table size, the REST API's own settings and keys, and the ACME, BIND and FreeRADIUS package models.

### The must-work list

**Verified**, from the model files at the tag:

| Model | Endpoint | `subsystem` | `apply()` | Same-request apply |
| --- | --- | --- | --- | --- |
| Firewall rule | `/api/v2/firewall/rule` | none (`Models/FirewallRule.inc:61-62`) | `FirewallApplyDispatcher` (`:524-526`) | yes, POST/PATCH body or DELETE query |
| Firewall alias | `/api/v2/firewall/alias` | `aliases` (`Models/FirewallAlias.inc:24-26`) | `FirewallApplyDispatcher` (`:165-167`) | yes |
| NAT port forward | `/api/v2/firewall/nat/port_forward` | `natconf` (`Models/PortForward.inc:46-48`) | `FirewallApplyDispatcher` (`:313-315`) | yes |
| DHCP static mapping | `/api/v2/services/dhcp_server/static_mapping` | `dhcpd` (`Models/DHCPServerStaticMapping.inc:41-43`) | `DHCPServerApplyDispatcher` (`:272-274`) | yes |

A port forward with an `associated_rule_id` also creates, updates or deletes its linked firewall rule in `pre_apply_*()`, which runs whether or not an apply was requested, and that inner rule write never applies itself (`Models/PortForward.inc:227-311`). **Inference:** the port forward's own apply reloads the whole filter, so the linked rule goes live with it.

### What the spec does and does not say about `apply`

**Verified:**

- No request-body schema in the spec has an `apply` or `async` property (checked the `FirewallRule`, `FirewallAlias`, `PortForward` and `DHCPServerStaticMapping` component schemas and the POST/PATCH bodies), and the string `async` appears nowhere in the spec.
- An `apply` query parameter is added to DELETE only when the model has a `$subsystem` (`Schemas/OpenAPISchema.inc:295-307`). So it is listed on DELETE `firewall/alias`, `firewall/nat/port_forward` and `services/dhcp_server/static_mapping`, but **not** on DELETE `firewall/rule`, even though `Endpoint::delete()` honours it there.
- The provider's checked-in `gen/schema.json` therefore has no `apply` attribute anywhere.

**Inference:** the generated provider cannot request an apply today. Getting it will need either a spec patch in `cmd/patch-openapi` (add `apply`/`async` to the write bodies and the DELETE parameters) or a hook in `cmd/patch-provider` that sets them on every write. Exposing `apply` as a user-visible attribute would be wrong for Terraform, since an unapplied resource is not the state that was declared.

## 5. How pending changes can be read

**Verified.** `GET /api/v2/<section>/apply` returns `applied: false` while that section's flag or file is set (table in section 3). The package docs recommend this for polling async applies (`docs/COMMON_CONTROL_PARAMETERS.md:53-55`).

Two gaps make it unreliable for the firewall section:

- `FirewallRule` sets no `$subsystem` (`Models/FirewallRule.inc`, no `subsystem` assignment), so a staged rule write marks nothing dirty and `GET /firewall/apply` keeps reporting `applied: true`. pfSense's own GUI marks `filter` dirty on rule edits (`src/usr/local/www/firewall_rules.php:161`, `187`, `203`, …), and `FirewallApply` checks for `filter` (`Models/FirewallApply.inc:17`), but nothing in the REST package ever sets it.
- `POST /firewall/apply` never clears `aliases`, `natconf` or `shaper`: the dispatcher only reloads the filter (`Dispatchers/FirewallApplyDispatcher.inc:15-19`), `FirewallApply` has no `$subsystem` so `Model::create()` clears nothing, and pfSense's `filter_configure()` does not clear flags either (the GUI clears them itself, e.g. `firewall_aliases.php:47-50`, `firewall_rules.php:150-152`). The grep for `clear_subsystem_dirty` across the package finds no call for those three subsystems. Only a per-object write with `apply: true` clears them, and then only its own subsystem (`Core/Model.inc:2278`).

**Inference:** after a staged alias or NAT change, `GET /firewall/apply` keeps reporting `applied: false` with that subsystem in `pending_subsystems` even after a successful `POST /firewall/apply`, until some alias (or NAT) write is made with `apply: true`. This was not reproduced on a running firewall.

## 6. Does one section's apply touch another's pending changes?

**Verified** from the dispatchers:

- Interface apply calls `filter_configure()`, and runs `system_routing_configure()` and clears `staticroutes` if routes are pending (`Dispatchers/InterfaceApplyDispatcher.inc:68-73`), so it applies pending routing changes too.
- Routing (`RoutingApplyDispatcher.inc:35-37`), DNS forwarder (`DNSForwarderApplyDispatcher.inc:16`) and IPsec (`IPsecApplyDispatcher.inc:18`) applies all call `filter_configure()`.
- Any firewall-family model's apply (rule, alias, NAT, shaper) runs the same `FirewallApplyDispatcher`, so it reloads the whole ruleset.

**Inference:** `filter_configure_sync()` rebuilds the ruleset from the whole of `config.xml` (e.g. it iterates `aliases/alias`, pfSense `filter.inc:1227`), so any filter reload makes every staged rule, alias and NAT change live, and does so without clearing their dirty flags. The reverse is not true: a firewall apply does not touch interfaces, routes, DHCP, DNS, VPN or HAProxy.

Consequences for the provider:

- Applying per resource is safe but repeats a full filter reload for every firewall-family resource in a plan; there is no per-object apply.
- Staging and then calling one `POST /firewall/apply` at the end would be cheaper, but Terraform's plugin protocol has no end-of-apply hook, and a later unrelated apply (interface, routing, IPsec) would push the staged firewall changes out anyway.
- One resource's `apply: true` will also push out whatever another client staged in the same section.

## Open questions

- Whether `dhcp_apply_changes()` in pfSense 2.8 also reloads DNS registration (Unbound/dnsmasq) was not checked.
- Section 5's stale `pending_subsystems` and section 6's cross-section side effects should be confirmed on a real pfSense 2.8 with REST API v2.8.1 before the provider depends on them.
