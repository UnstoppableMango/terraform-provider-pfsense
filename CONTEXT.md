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

## Relationships

- Every operation on an existing object starts with a **Lookup** of its **Key**; a positional **API id** is never stored or reused.
- A stable string **API id** is itself the resource's **Key**.
- A resource's **Key** is one or more of its fields, or a **Marker** when no field is unique.
- A child resource's **Key** includes its parent's id.

## Flagged ambiguities

- "id" means the **API id** in the REST API and the **Key** in Terraform state; they are the same value only for resources whose **API id** is a stable string.
