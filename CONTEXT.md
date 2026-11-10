# pfSense provider

Terraform provider for pfSense, generated from the pfSense REST API's OpenAPI spec.

## Language

**API id**:
The `id` the pfSense REST API returns for an object. For most resources it is the object's position in its config section, which shifts when siblings are deleted, moved or re-sorted, so it identifies the object only for the instant it was read. A few resources instead have a stable string API id, such as an interface name like `lan`.
_Avoid_: object id, resource id

**Key**:
The field or fields that identify one object of a resource across operations, and the value the provider stores as the Terraform resource id.
_Avoid_: natural key, primary key

**Lookup**:
Resolving a **Key** to the object's current **API id** through a filtered query on the resource's plural endpoint, which must match exactly one object.
_Avoid_: search, find

**Marker**:
A provider-written tag in an object's description that serves as its **Key** when the resource has no field that pfSense keeps unique.
_Avoid_: tag, label

**Apply**:
Making changes saved in pfSense's configuration take effect on the running system, such as reloading the packet filter.
_Avoid_: commit, reload, sync

**Staged write**:
A change saved to pfSense's configuration but not yet applied, so the running system still behaves as before.
_Avoid_: pending change, dirty config

**Create-only field**:
A field pfSense accepts when an object is created but refuses to change afterwards, so changing it means replacing the object.
_Avoid_: immutable field, non-editable field

**Singular endpoint**:
A REST API path that operates on one object of a model, such as `/firewall/rule`.
_Avoid_: item endpoint, detail endpoint

**Plural endpoint**:
A REST API path that operates on every object of a model at once, and the one a **Lookup** queries, such as `/firewall/rules`. Its spelling isn't always the singular plus `s`.
_Avoid_: list endpoint, collection endpoint

**Singleton**:
A model with exactly one object and no create or delete, such as the system hostname, managed as a resource that only updates it.
_Avoid_: settings resource, global

**Action**:
An endpoint that does something rather than storing an object, such as reboot, ping, or certificate generation; it never becomes a resource.
_Avoid_: command, operation

## Relationships

- Every operation on an existing object starts with a **Lookup** of its **Key**; a positional **API id** is never stored or reused.
- A stable string **API id** is itself the resource's **Key**.
- A resource's **Key** is one or more of its fields, or a **Marker** when no field is unique.
- A child resource's **Key** includes its parent's id.
- A model that can be created, read, updated and deleted through its **Singular endpoint**, and has a **Plural endpoint**, becomes a resource; the two are paired by model, not by path.
- A resource is named after its **Singular endpoint**.
- Changing a **Key** field or a **Create-only field** replaces the object.
- Every write the provider makes is followed by its **Apply** in the same request; the provider never leaves a **Staged write** behind.

## Flagged ambiguities

- "id" means the **API id** in the REST API and the **Key** in Terraform state; they are the same value only for resources whose **API id** is a stable string.
