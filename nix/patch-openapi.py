import json, copy, sys

src, dst = sys.argv[1], sys.argv[2]

with open(src) as f:
    spec = json.load(f)


def resolve_ref(spec, ref):
    parts = ref.lstrip("#/").split("/")
    node = spec
    for p in parts:
        node = node[p]
    return copy.deepcopy(node)


def flatten_schema(spec, schema):
    """Recursively resolve $ref and flatten allOf into a single object schema."""
    if not isinstance(schema, dict):
        return schema

    if "$ref" in schema:
        schema = resolve_ref(spec, schema["$ref"])

    if "allOf" in schema:
        merged = {"type": "object", "properties": {}}
        for sub in schema["allOf"]:
            sub = flatten_schema(spec, sub)
            merged["properties"].update(sub.get("properties", {}))
            if "required" in sub:
                merged.setdefault("required", [])
                merged["required"].extend(sub["required"])
        if not merged.get("required"):
            merged.pop("required", None)
        schema = merged

    if "properties" in schema:
        for k, v in schema["properties"].items():
            schema["properties"][k] = flatten_schema(spec, v)

    if "items" in schema:
        schema["items"] = flatten_schema(spec, schema["items"])

    return schema


for path_item in spec.get("paths", {}).values():
    for op in path_item.values():
        if not isinstance(op, dict):
            continue
        for content in op.get("requestBody", {}).get("content", {}).values():
            if "schema" in content:
                content["schema"] = flatten_schema(spec, content["schema"])
        for response in op.get("responses", {}).values():
            for content in response.get("content", {}).values():
                if "schema" in content:
                    content["schema"] = flatten_schema(spec, content["schema"])

with open(dst, "w") as f:
    json.dump(spec, f)
