package scanner

import "testing"

func TestResolveRef_DirectRef(t *testing.T) {
	schema := map[string]interface{}{
		"$ref": "#/components/schemas/User",
	}
	components := map[string]interface{}{
		"User": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{"type": "integer"},
			},
		},
	}

	result := resolveRef(schema, components, 10)
	if result["type"] != "object" {
		t.Fatalf("expected type=object, got %v", result["type"])
	}
}

func TestResolveRef_NestedRef(t *testing.T) {
	schema := map[string]interface{}{
		"$ref": "#/components/schemas/Wrapper",
	}
	components := map[string]interface{}{
		"Wrapper": map[string]interface{}{
			"$ref": "#/components/schemas/Inner",
		},
		"Inner": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{"type": "string"},
			},
		},
	}

	result := resolveRef(schema, components, 10)
	if result["type"] != "object" {
		t.Fatalf("expected type=object, got %v", result["type"])
	}
}

func TestResolveRef_NoRef(t *testing.T) {
	schema := map[string]interface{}{
		"type": "string",
	}

	result := resolveRef(schema, nil, 10)
	if result["type"] != "string" {
		t.Fatalf("expected type=string, got %v", result["type"])
	}
}

func TestResolveRef_DepthLimit(t *testing.T) {
	schema := map[string]interface{}{
		"$ref": "#/components/schemas/A",
	}
	components := map[string]interface{}{
		"A": map[string]interface{}{
			"$ref": "#/components/schemas/B",
		},
		"B": map[string]interface{}{
			"type": "object",
		},
	}

	result := resolveRef(schema, components, 1)
	// depth=1 resolves A, but A has another $ref and depth is now 0
	if _, hasRef := result["$ref"]; !hasRef {
		t.Fatal("expected unresolved $ref at depth limit")
	}
}

func TestResolveRef_UnknownRef(t *testing.T) {
	schema := map[string]interface{}{
		"$ref": "#/components/schemas/Missing",
	}
	components := map[string]interface{}{}

	result := resolveRef(schema, components, 10)
	if result["$ref"] != "#/components/schemas/Missing" {
		t.Fatal("expected original schema returned for unknown ref")
	}
}

func TestResolveRef_NonComponentRef(t *testing.T) {
	schema := map[string]interface{}{
		"$ref": "#/definitions/Foo",
	}

	result := resolveRef(schema, nil, 10)
	if result["$ref"] != "#/definitions/Foo" {
		t.Fatal("expected original schema for non-component ref")
	}
}
