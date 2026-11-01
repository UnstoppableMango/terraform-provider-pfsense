# Object ids in pfSense REST API v2.8.1

Research for [#55](https://github.com/UnstoppableMango/terraform-provider-pfsense/issues/55) (map #54).
Question: what identifies an object in REST API v2.8.1, and is it stable enough to be a Terraform resource id?

## Answer

No.
For almost every plural model the `id` is the PHP array index of the object in its `config.xml` section, and it shifts whenever an earlier sibling is deleted, an object is moved with `placement`, or the model re-sorts its section on write.
The upstream docs say so outright.
A Terraform resource has to store a natural key and resolve it to the current `id` with a filtered query on the plural endpoint before every singular GET, PATCH or DELETE.
The natural keys are good for aliases (`name`) and DHCP static mappings (`parent_id` plus `mac`), weak for firewall rules (`tracker` is a Unix timestamp in seconds with no uniqueness check), and absent for NAT port forwards.

## Sources

Highest trust first:

1. Package source at tag `v2.8.1`, commit `4190c8ca40c73e5705e01f227a029b61dedf7b7f`, cloned from <https://github.com/pfrest/pfSense-pkg-RESTAPI>.
   Paths below are relative to the repo root; `RESTAPI/` stands for `pfSense-pkg-RESTAPI/files/usr/local/pkg/RESTAPI/`.
   Links use `https://github.com/pfrest/pfSense-pkg-RESTAPI/blob/v2.8.1/<path>#L<n>`.
2. pfSense CE source, `pfsense/pfsense` at `master` commit `9363ac5b8651a1c7a333180425ce7719070f95f9` (2026-03-31).
   The public repo has no 2.8 branch or tag, so this is the closest available; the functions cited are long-standing, but this is not the exact code running on a 2.8.1 box.
3. The release asset `openapi.json` from the `v2.8.1` release (`info.version` is `v2.8.1`).
4. First-party docs in the package repo at the same tag (`docs/WORKING_WITH_OBJECT_IDS.md`, `docs/QUERIES_FILTERS_AND_SORTING.md`), which are what pfrest.org renders.

Each claim below is marked **verified** (read in the source or spec) or **inference** (reasoned from what was read, not executed against a pfSense box).
Nothing here was tested against a running pfSense.

## 1. What `id` refers to

**Verified.** `id` is a property of the `Model` base class, not a Field, which is why the `FirewallRule` schema has no `id` property.

- `RESTAPI/Core/Model.inc:81-85` declares `public mixed $id`, "applies exclusively to Models with $many enabled".
- `RESTAPI/Core/Model.inc:88-93` declares `$id_type`, default `'integer'`.
- `RESTAPI/Core/Model.inc:1133-1141`, in `to_representation()`, injects `parent_id` and `id` into the response alongside the Fields; that is the only place they enter the output.

**Verified.** For a `many` model backed by a `config_path`, the id is the array key in that config section.

- `RESTAPI/Core/Model.inc:814-858`, `get_next_id()`: a new object's id is `count($conf)` of the section, or `0` when it is empty. The docblock says it "searches for the next index ID".
- `RESTAPI/Core/Model.inc:2085-2104`, `read_all()`: each object is constructed with `id: $internal_id`, the key of the array returned from the config path (`get_internal_objects()`, `Model.inc:967-969`).
- `RESTAPI/Core/Model.inc:1032-1061`, `from_internal()`: a singular read looks the id up with `array_key_exists($this->id, $internal_objects)` and returns 404 `MODEL_OBJECT_NOT_FOUND` otherwise.
- `RESTAPI/Core/Model.inc:2229-2236`, `_create()`: writes to `"{config_path}/$this->id"`. `_update()` (`2299-2324`) and `_delete()` (`2510-2516`) address the same path.
- For `FirewallRule` the section is `filter/rule` (`RESTAPI/Models/FirewallRule.inc:61`), so `id=3` is the fourth `<rule>` across *all* interfaces, not the fourth rule on one interface.

**Verified.** The upstream docs state it directly: "pfSense does not use persistent object IDs ... pfSense uses the object's array index as the object ID instead of a persistent, stored value" (`docs/WORKING_WITH_OBJECT_IDS.md:5-11`).

**Verified.** The spec agrees in shape but not always in type.

- Singular GET and DELETE on `/api/v2/firewall/rule`, `/firewall/alias`, `/firewall/nat/port_forward` take a required `id` query parameter typed `oneOf: [integer, string]`.
- PATCH request bodies wrap the model schema in `allOf` with a required `id` of type `integer`. That is hard-coded at `RESTAPI/Schemas/OpenAPISchema.inc:190-198`, regardless of the model's `id_type`.

## 2. Do ids shift?

**Verified: yes, on delete.**

- `_delete()` calls `del_config("{config_path}/$id")` (`RESTAPI/Core/Model.inc:2513`), which reaches pfSense's `array_del_path()`; that is a plain `unset($el[$vkey])` (`pfsense src/etc/inc/util.inc:4809-4831`) and leaves a hole in memory.
- `Model::write_config()` then calls pfSense `write_config()` (`RESTAPI/Core/Model.inc:766`), which serializes to XML and re-reads it with `config_read_file(true, false)` (`pfsense src/etc/inc/config.lib.inc:637-639`). XML carries no indices, so the reparsed list is renumbered from 0.
- The package relies on this: `ModelSet::delete()` deletes "in reverse order to avoid index issues when deleting objects" (`RESTAPI/Core/ModelSet.inc:158-162`).
- The docs give the example: with ids 0, 1, 2, deleting 1 turns 2 into 1 (`docs/WORKING_WITH_OBJECT_IDS.md:10-11`).

**Verified: yes, on placement.** `set_placement()` (`RESTAPI/Core/Model.inc:1928-1965`) does `array_splice` out of its index and back in at `placement`, then sets `$this->id = $this->placement`. Every object between the old and new index changes id. POST and PATCH accept `placement` for every model (`RESTAPI/Core/Endpoint.inc:1218, 1243`).

**Verified: yes, on create and update for models with `sort_by`.** `sort()` (`RESTAPI/Core/Model.inc:1884-1922`) rewrites the whole section in sorted order after every `_create()` and `_update()` (`Model.inc:2235, 2323`). Models that set `sort_by` in v2.8.1: `DHCPServerStaticMapping` (by `ipaddr`, `RESTAPI/Models/DHCPServerStaticMapping.inc:44-45`), `DNSResolverHostOverride`, `DNSResolverHostOverrideAlias`, `FirewallSchedule`, `RESTAPIAccessListEntry`, `UserGroup`. On these, creating or re-addressing any one object can renumber its siblings. A client can also pass `sort_by` on POST/PATCH to force a sort on any model (`Endpoint.inc:1221-1223, 1245-1247`).

**Inference: yes, on changes made outside the API.** Edits in the web GUI, config restores, and other API clients all rewrite the same arrays, and the API has no record of what an index used to point at.

**Verified: which models are exempt.** Only models that override the id scheme. In v2.8.1, three set `id_type = 'string'`: `DHCPServer`, `NetworkInterface`, `Table`. `DHCPServer` keys `dhcpd` by interface id (`lan`, `opt1`) and overrides `get_next_id()` to return it (`RESTAPI/Models/DHCPServer.inc:54-56, 527-529`), so its id is stable. The docs call these out as the "small number of endpoints" with persistent ids (`docs/WORKING_WITH_OBJECT_IDS.md:13-15`). Of 119 `many` models, 92 use a `config_path` and take the index scheme unless they override it; models backed by an `internal_callable` get whatever key their callable returns (`Model.inc:972-990`), so their stability depends on each callable (not surveyed).

## 3. Natural keys on the must-work list

**Verified: plural endpoints filter by any field.**

- `GET` on a plural endpoint passes all remaining request data to `Model::query()` as filters (`RESTAPI/Core/Endpoint.inc:1186-1199`); `DELETE` on a plural endpoint does the same through `delete_many()` (`Endpoint.inc:1260-1278`).
- `ModelSet::query()` (`RESTAPI/Core/ModelSet.inc:235-291`) matches `field` or `field__<filter>`, and treats `id` and `parent_id` as fields too (`ModelSet.inc:265-266`).
- The default filter `exact` compares with loose `==` unless one side is null (`RESTAPI/QueryFilters/ExactQueryFilter.inc:21-28`), so a query-string `"1700000000"` matches an integer `tracker`.
- Other filters: `startswith`, `endswith`, `contains`, `in`, `lt`, `lte`, `gt`, `gte`, `except`, `regex`, `format` (`RESTAPI/QueryFilters/`, `docs/QUERIES_FILTERS_AND_SORTING.md`).
- The spec models this as one free-form `query` parameter, `style: form`, `explode: true`, `additionalProperties: string`, on `GET /api/v2/firewall/rules` and the other plural paths.

**Verified: an unknown filter field is ignored, not rejected.** `ModelSet.inc:259-262` skips any query key that is not a property of the model. A misspelled key (`?trackr=...`) matches every object. **Inference:** a provider lookup must check that it got exactly one result and fail on zero or several, never take the first.

Per model:

| Model | Natural key | Source | Strength |
| --- | --- | --- | --- |
| `FirewallAlias` | `name` | `RESTAPI/Models/FirewallAlias.inc:29-37`: `required`, `unique`, `editable: false` | Strong. Unique across aliases, cannot be renamed. Lookup: `GET /api/v2/firewall/aliases?name=<name>`. |
| `FirewallRule` | `tracker` | `RESTAPI/Models/FirewallRule.inc:301-307`: `UnixTimeField`, `editable: false`, `read_only: true`, `auto_add_now: true`, **not** `unique` | Weak. The value is `time()`, whole seconds (`RESTAPI/Fields/UnixTimeField.inc:141-144`), and nothing checks uniqueness, so two rules created in the same second share a tracker (**inference** from the code). It is read-only and never regenerated on update. Lookup: `GET /api/v2/firewall/rules?tracker=<n>`. |
| `FirewallRule` | `associated_rule_id` | `FirewallRule.inc:308-317`: `unique`, read-only | Only for rules created as the companion of a NAT port forward; null otherwise. |
| `PortForward` | none | `RESTAPI/Models/PortForward.inc:24-42` | No unique field. `associated_rule_id` (`PortForward.inc:159-172`) is a random `get_unique_id()` only when created with `new` (`PortForward.inc:215-219`); it is `''` or `pass` otherwise. `created_time` is a seconds timestamp in `created/time`, not unique. **Inference:** the only stable handle is a key the provider itself writes, such as a marker in `descr`, or forcing `associated_rule_id: new`. |
| `DHCPServerStaticMapping` | `parent_id` + `mac` | `RESTAPI/Models/DHCPServerStaticMapping.inc:48-54`: `mac` is `required`, `unique` | Strong. Uniqueness is checked against siblings under the same parent only, because `Field::check_field_unique()` queries `parent_id: $this->context->parent_id` (`RESTAPI/Core/Field.inc:727-751`). `mac` is editable, so a changed MAC is a replace. Lookup: `GET /api/v2/services/dhcp_server/static_mappings?parent_id=lan&mac=<mac>`. |

## 4. How child models address their parent

**Verified.**

- A child model sets `parent_model_class`; `DHCPServerStaticMapping` sets `'DHCPServer'` with `config_path = 'staticmap'` (`RESTAPI/Models/DHCPServerStaticMapping.inc:40-41`).
- Its config path is `"{parent config_path}/{parent_id}/{config_path}"`, here `dhcpd/<parent_id>/staticmap` (`RESTAPI/Core/Model.inc:628-631`).
- `parent_id` is the parent's `id`, so for static mappings it is the DHCP server's interface id string (`lan`, `opt1`), which is stable.
- The child's own `id` is the index within that parent's `staticmap` list, and it is re-sorted by `ipaddr` on every write (section 2).
- `parent_id` is required whenever the parent is a `many` model (`Model.inc:1506-1514`, `1033-1039`).
- The singular endpoint takes both `parent_id` and `id` as required query parameters in the spec.
- The plural endpoint reads children from every parent (`get_internal_objects_from_all_parents()`, `Model.inc:909-928`) and filters on `parent_id` like any field.
- **Spec bug, verified:** the PATCH body schema declares `parent_id` as `integer` (`RESTAPI/Schemas/OpenAPISchema.inc:230-240`), but the static mapping's parent id is a string. The query-parameter form is correctly `oneOf: [integer, string]` (`OpenAPISchema.inc:282-291`). Generated code that trusts the body schema will mistype it.
- The docs warn that `parent_id` is no more persistent than `id` in general (`docs/WORKING_WITH_OBJECT_IDS.md:63-64`); for `DHCPServer` parents it is, because their id is the interface name.

## Implications for the provider

All of this is **inference**.

- Do not use the API `id` as the Terraform resource id, or as anything stored in state and reused across operations.
- Store the natural key as the resource id: alias `name`; static mapping `<parent_id>/<mac>`; firewall rule `tracker`, with a documented collision risk; port forward needs a provider-owned key.
- Before each singular call, resolve the key to the current `id` through the plural endpoint with an exact filter, and require exactly one match.
  There is still a race between the lookup and the write if anything else edits the same section in between; the API offers no compare-and-swap.
- A PATCH that sends `placement` or `sort_by`, or any create on a `sort_by` model, invalidates every sibling's cached id; never reuse an id across operations.
- `HashiCorp tfplugingen-openapi` will map `id` to an integer attribute from the PATCH body schema; that attribute must be treated as computed and volatile, not as identity.
