package scanner

import "testing"

func TestExtractFieldsFromOp_200Schema(t *testing.T) {
	op := map[string]interface{}{
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"success": map[string]interface{}{"type": "boolean"},
							},
						},
					},
				},
			},
		},
	}

	fields := extractFieldsFromOp(op, nil)
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}
	if fields[0].Path != "$.success" || fields[0].Type != "boolean" {
		t.Fatalf("unexpected field: %+v", fields[0])
	}
}

func TestExtractFieldsFromOp_WithRef(t *testing.T) {
	op := map[string]interface{}{
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": map[string]interface{}{
							"$ref": "#/components/schemas/LoginResponse",
						},
					},
				},
			},
		},
	}
	components := map[string]interface{}{
		"LoginResponse": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"token": map[string]interface{}{"type": "string"},
				"data": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id": map[string]interface{}{"type": "integer"},
					},
				},
			},
		},
	}

	fields := extractFieldsFromOp(op, components)
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
	if byPath["$.token"] != "string" {
		t.Fatal("expected $.token — string")
	}
}

func TestExtractFieldsFromOp_NoResponses(t *testing.T) {
	op := map[string]interface{}{
		"operationId": "Test",
	}

	fields := extractFieldsFromOp(op, nil)
	if fields != nil {
		t.Fatalf("expected nil, got %v", fields)
	}
}

func TestExtractFieldsFromOp_NoContent(t *testing.T) {
	op := map[string]interface{}{
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"description": "OK",
			},
		},
	}

	fields := extractFieldsFromOp(op, nil)
	if fields != nil {
		t.Fatalf("expected nil, got %v", fields)
	}
}

func TestExtractFieldsFromOp_201Status(t *testing.T) {
	op := map[string]interface{}{
		"responses": map[string]interface{}{
			"201": map[string]interface{}{
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"id": map[string]interface{}{"type": "integer"},
							},
						},
					},
				},
			},
		},
	}

	fields := extractFieldsFromOp(op, nil)
	if len(fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(fields))
	}
	if fields[0].Path != "$.id" {
		t.Fatalf("expected $.id, got %s", fields[0].Path)
	}
}

func TestExtractFieldsFromOp_SortedOutput(t *testing.T) {
	op := map[string]interface{}{
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"z_field": map[string]interface{}{"type": "string"},
								"a_field": map[string]interface{}{"type": "integer"},
							},
						},
					},
				},
			},
		},
	}

	fields := extractFieldsFromOp(op, nil)
	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fields))
	}
	if fields[0].Path != "$.a_field" {
		t.Fatalf("expected sorted: $.a_field first, got %s", fields[0].Path)
	}
	if fields[1].Path != "$.z_field" {
		t.Fatalf("expected sorted: $.z_field second, got %s", fields[1].Path)
	}
}
