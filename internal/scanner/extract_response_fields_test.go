package scanner

import (
	"sort"
	"testing"
)

func TestExtractResponseFields_FlatProperties(t *testing.T) {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"success": map[string]interface{}{"type": "boolean"},
			"message": map[string]interface{}{"type": "string"},
		},
	}

	fields := extractResponseFields(schema, nil)
	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fields))
	}

	sort.Slice(fields, func(i, j int) bool { return fields[i].Path < fields[j].Path })
	if fields[0].Path != "$.message" || fields[0].Type != "string" {
		t.Fatalf("unexpected field 0: %+v", fields[0])
	}
	if fields[1].Path != "$.success" || fields[1].Type != "boolean" {
		t.Fatalf("unexpected field 1: %+v", fields[1])
	}
}

func TestExtractResponseFields_NestedProperties(t *testing.T) {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"success": map[string]interface{}{"type": "boolean"},
			"data": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id":    map[string]interface{}{"type": "integer"},
					"email": map[string]interface{}{"type": "string"},
				},
			},
		},
	}

	fields := extractResponseFields(schema, nil)
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}

	byPath := make(map[string]string)
	for _, f := range fields {
		byPath[f.Path] = f.Type
	}

	if byPath["$.success"] != "boolean" {
		t.Fatal("expected $.success — boolean")
	}
	if byPath["$.data.id"] != "integer" {
		t.Fatal("expected $.data.id — integer")
	}
	if byPath["$.data.email"] != "string" {
		t.Fatal("expected $.data.email — string")
	}
}

func TestExtractResponseFields_WithRef(t *testing.T) {
	schema := map[string]interface{}{
		"$ref": "#/components/schemas/LoginResponse",
	}
	components := map[string]interface{}{
		"LoginResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"token": map[string]interface{}{"type": "string"},
			},
		},
	}

	fields := extractResponseFields(schema, components)
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}
	if fields[0].Path != "$.token" || fields[0].Type != "string" {
		t.Fatalf("unexpected field: %+v", fields[0])
	}
}

func TestExtractResponseFields_NoProperties(t *testing.T) {
	schema := map[string]interface{}{
		"type": "string",
	}

	fields := extractResponseFields(schema, nil)
	if fields != nil {
		t.Fatalf("expected nil, got %v", fields)
	}
}

func TestExtractResponseFields_NestedRef(t *testing.T) {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"data": map[string]interface{}{
				"$ref": "#/components/schemas/UserData",
			},
		},
	}
	components := map[string]interface{}{
		"UserData": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id":   map[string]interface{}{"type": "integer"},
				"name": map[string]interface{}{"type": "string"},
			},
		},
	}

	fields := extractResponseFields(schema, components)
	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fields))
	}

	byPath := make(map[string]string)
	for _, f := range fields {
		byPath[f.Path] = f.Type
	}

	if byPath["$.data.id"] != "integer" {
		t.Fatal("expected $.data.id — integer")
	}
	if byPath["$.data.name"] != "string" {
		t.Fatal("expected $.data.name — string")
	}
}
