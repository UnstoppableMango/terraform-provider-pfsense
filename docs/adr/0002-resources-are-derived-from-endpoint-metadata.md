# Resources are derived from the endpoint metadata in the spec, named after their path

The provider's resource list is derived from the pinned REST API spec (v2.8.1) rather than written by hand ([#49](https://github.com/UnstoppableMango/terraform-provider-pfsense/issues/49)).
Each operation's description in that spec carries a fixed `**Details:**` block naming its **Endpoint type** (Singular or Plural), **Associated model**, **Parent model** and **Required packages**, and the rule reads that block rather than guessing from the path text.
Decided in [#48](https://github.com/UnstoppableMango/terraform-provider-pfsense/issues/48); how each resource is identified is [ADR 0001](0001-resources-are-identified-by-key-not-api-id.md).
Every path and what it maps to is listed in [docs/resource-paths.md](../resource-paths.md).

## Decision

- A **Singular endpoint** and a **Plural endpoint** are paired by their Associated model, not by path. Plural paths don't follow one spelling: `/haproxy/backend/errorfiles` pairs with `/haproxy/backend/error_file`, `/routing/gateway/group/priorities` with `/priority`, and `/system/restapi/access_list` is the plural of `/system/restapi/access_list/entry`.
- `gen-config` fails when an operation's description lacks the block, so a re-pinned spec that drops it breaks the build instead of silently producing no resources.
- A model whose Singular endpoint has GET, POST, PATCH and DELETE and which has a Plural endpoint becomes a resource: Create is the singular POST, Read the singular GET, Update the singular PATCH, Delete the singular DELETE, and the Lookup goes to the plural GET. Plural PUT (replace the whole section) and plural DELETE are never used.
- The resource name is the singular path with `/api/v2/` dropped and `/` turned into `_`: `firewall_rule`, `firewall_nat_port_forward`, `services_dhcp_server_static_mapping`. A child resource's name comes from its own path like any other, and its Parent model makes its Key include the parent id (ADR 0001).
- A model whose Singular endpoint has only GET and PATCH is a **Singleton** resource: Create and Update both PATCH, Destroy removes it from state without calling the API, and its Terraform id is a fixed value. `/system/restapi/version` is excluded, because a PATCH there upgrades the REST API package.
- Models that need a package are generated like any other; their documentation names the package, and the API's error at apply time reports a missing one.
- Everything else is excluded, with the reason recorded in the table: **Actions** (POST-only endpoints such as reboot, ping, certificate generation), the `/apply` endpoints (when to apply is [#59](https://github.com/UnstoppableMango/terraform-provider-pfsense/issues/59)), GraphQL and JWT, PATCH-only `/services/dhcp_server/backend`, read-only state such as status, logs and diagnostics, and the overrides below.
- No data sources are generated in this effort. The table marks read-only paths as data-source candidates.

## Overrides

| Model | Path | Why it is excluded |
| --- | --- | --- |
| RESTAPIVersion | `/system/restapi/version` | GET and PATCH like a Singleton, but a PATCH upgrades the package |
| ACMECertificateAction | `/services/acme/certificate/action` | full CRUD, no Plural endpoint to Lookup on |
| ACMECertificateDomain | `/services/acme/certificate/domain` | full CRUD, no Plural endpoint to Lookup on |
| BINDZoneRecord | `/services/bind/zone/record` | full CRUD, no Plural endpoint to Lookup on |
| CertificateRevocationListRevokedCertificate | `/system/crl/revoked_certificate` | full CRUD, no Plural endpoint to Lookup on |
| Package | `/system/package` | GET, POST and DELETE but no PATCH; would be a replace-only resource |
| RESTAPIKey | `/auth/key` | POST and DELETE only; nothing to Read |

Models that pass the rule but have neither a declared Key nor a field that accepts a Marker are not generated (ADR 0001); the table lists them as "no Key yet". There were 44 at this pin until [#62](https://github.com/UnstoppableMango/terraform-provider-pfsense/issues/62) gave 24 of them a Key from `unique` or unique-together fields and 5 a Marker in `description`, leaving 15.

## Considered options

- Pairing by path text (singular plus `s`): rejected, three plural paths don't follow it.
- Naming after the model in snake case (`dhcp_server_static_mapping`): rejected, it needs an acronym table to avoid `i_psec_phase1`, `ha_proxy_backend` and `free_radiusmac`, and the name no longer says which URL backs it.
- Dropping the `services_`, `system_` and `vpn_` prefixes: rejected, it's an ad hoc list that would grow with every re-pin and invites collisions.
- Leaving Singletons out: rejected, settings like hostname, DNS and SSH are a large part of what people configure on a firewall.
- Skipping package-gated models: rejected, the error at apply time is clear enough and the package list would have to be maintained.

## Consequences

- Names are long (`services_dns_resolver_host_override_alias`) and change if pfSense moves a path, which re-pinning would surface as a rename.
- Destroying a Singleton leaves the setting as it was last applied, not at its default.
- Re-pinning the spec can add, drop or rename resources with no code change; the table has to be regenerated alongside it.
