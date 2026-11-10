# Resources are identified by a Key, not by the positional pfSense API id

pfSense REST API v2.8.1 gives an object an `id` that is its position in its config section, and it shifts whenever a sibling is deleted, moved or re-sorted ([#55](https://github.com/UnstoppableMango/terraform-provider-pfsense/issues/55)).
So the provider stores a **Key** as the Terraform resource id and does a **Lookup** on the plural endpoint, requiring exactly one match, before every singular call; a positional API id is never stored or reused.
Decided in [#58](https://github.com/UnstoppableMango/terraform-provider-pfsense/issues/58).

## Decision

- Each resource's Key is declared next to its entry in `cmd/gen-config`, because the spec carries no uniqueness metadata. The declared Key is a field the REST API's model source marks `unique: true`, or the fields it lists in `unique_together_fields`, both of which the API enforces on every write. When a model has several, the Key is the one the user sets and rarely changes, usually `name`. A server-assigned field is the Key only when it never changes after create, as with OpenVPN `vpnid`. Models whose API id is a stable string (DHCP server, interface, table) use that string API id as the Key; this is the one case where state holds an API id.
- A Lookup filters with an exact match and then checks that the returned object's Key fields equal the request, because pfSense ignores an unknown filter key and would match everything. A Marker Lookup can't match exactly, because the Marker is appended to a description the import doesn't know, so it filters with `<field>__contains=[tf:<hex>]` and then checks that the Marker extracted from each result equals the Key.
- Read: zero matches removes the resource from state; several matches is an error naming the Key and the count. Update and Delete fail on zero or several.
- Key fields force replacement when changed.
- Firewall rules are keyed by `tracker`, a whole-second timestamp with no uniqueness check, so the provider spaces rule creates at least one second apart.
- Models with no unique field but a `descr` or `description` field are keyed by a **Marker**: ` [tf:<12 hex chars>]` appended to that field, stripped from its attribute and exposed as a computed `marker` attribute. NAT port forwards use this.
- Models with neither a declared Key nor a field that accepts the Marker are not generated as resources; the generator lists them in a report instead. A Marker can't go in a field whose validator rejects its brackets, as FreeRADIUS interface's `description` does; that model is keyed by its unique-together fields instead.
- Amended by [#62](https://github.com/UnstoppableMango/terraform-provider-pfsense/issues/62): Keys taken from `unique` and `unique_together_fields` in the model source, the Marker in `description`, and how a Marker is looked up.
- All writes from one provider instance go through one lock, which also serializes the Lookup-then-write sequence between resources in the same apply.
- Import takes the Key, with child resources prefixed by the parent id and multi-field Keys joined by `/` (`lan/00:11:22:33:44:55`). A Marker-keyed object that has no Marker yet is imported by its current API id and matched on every configured field until the next apply writes its Marker.

## Considered options

- Storing the API id and re-reading it: rejected, the id silently points at a different object after a sibling is deleted.
- Forcing `associated_rule_id: new` on port forwards to get a random unique id: rejected, it always creates a linked filter rule.
- A composite Key of interface, protocol and destination port for port forwards: rejected, pfSense doesn't enforce its uniqueness.

## Consequences

- A rename in the GUI of a Key field, or deleting a Marker from a description, makes the object look deleted, and the next plan recreates it.
- The Lookup-then-write race with writers outside the provider (GUI, other API clients) remains; the API has no compare-and-swap.
