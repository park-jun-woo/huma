package scanner

import (
	"encoding/json"
	"testing"
)

func TestParseOpenAPI(t *testing.T) {
	data := []byte(`
openapi: '3.0.0'
info:
  title: Test API
paths:
  /api/v1/auth/login:
    post:
      operationId: Login
      x-source-file: internal/api/auth/handler.go
      x-source-line: 50
      responses:
        200:
          description: OK
        400:
          description: Bad Request
  /api/v1/buildings/{buildingId}:
    get:
      operationId: GetBuilding
    delete:
      x-source-file: internal/api/building/handler.go
  /api/v1/users:
    get:
      operationId: ListUsers
    post:
      operationId: CreateUser
`)

	endpoints, err := parseOpenAPI(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(endpoints) != 5 {
		t.Fatalf("expected 5 endpoints, got %d", len(endpoints))
	}

	byKey := make(map[string]Endpoint)
	for _, ep := range endpoints {
		byKey[ep.Method+" "+ep.Path] = ep
	}

	// Check Login endpoint
	login, ok := byKey["POST /api/v1/auth/login"]
	if !ok {
		t.Fatal("missing POST /api/v1/auth/login")
	}
	if login.Handler != "Login" {
		t.Fatalf("expected Login, got %s", login.Handler)
	}
	if login.Source != "internal/api/auth/handler.go" {
		t.Fatalf("expected source file, got %s", login.Source)
	}
	if login.Line != 50 {
		t.Fatalf("expected line 50, got %d", login.Line)
	}
	if len(login.Responses) == 0 {
		t.Fatal("expected non-empty Responses for Login")
	}
	type respEntry struct {
		Status int `json:"status"`
	}
	var respEntries []respEntry
	if err := json.Unmarshal(login.Responses, &respEntries); err != nil {
		t.Fatalf("unmarshal responses: %v", err)
	}
	if len(respEntries) != 2 {
		t.Fatalf("expected 2 response entries, got %d", len(respEntries))
	}
	if respEntries[0].Status != 200 || respEntries[1].Status != 400 {
		t.Fatalf("unexpected response statuses: %+v", respEntries)
	}

	// Check path param conversion
	getBuilding, ok := byKey["GET /api/v1/buildings/:buildingId"]
	if !ok {
		t.Fatal("missing GET /api/v1/buildings/:buildingId — path param conversion failed")
	}
	if getBuilding.Handler != "GetBuilding" {
		t.Fatalf("expected GetBuilding, got %s", getBuilding.Handler)
	}

	// Check endpoint without responses has nil Responses
	if getBuilding.Responses != nil {
		t.Fatalf("expected nil Responses for GetBuilding, got %s", string(getBuilding.Responses))
	}

	// Check auto-generated operationId (missing operationId)
	delBuilding, ok := byKey["DELETE /api/v1/buildings/:buildingId"]
	if !ok {
		t.Fatal("missing DELETE /api/v1/buildings/:buildingId")
	}
	if delBuilding.Handler != "delete_api_v1_buildings_buildingId" {
		t.Fatalf("expected auto-generated operationId, got %s", delBuilding.Handler)
	}
	if delBuilding.Source != "internal/api/building/handler.go" {
		t.Fatalf("expected source file, got %s", delBuilding.Source)
	}

	// Check ID is generated
	for _, ep := range endpoints {
		if ep.ID == "" {
			t.Fatalf("expected non-empty ID for %s %s", ep.Method, ep.Path)
		}
	}
}

func TestParseOpenAPI_WithResponseSchema(t *testing.T) {
	data := []byte(`
openapi: '3.0.0'
info:
  title: Test API
components:
  schemas:
    LoginResponse:
      type: object
      properties:
        success:
          type: boolean
        data:
          $ref: '#/components/schemas/UserData'
    UserData:
      type: object
      properties:
        id:
          type: integer
        email:
          type: string
paths:
  /api/v1/auth/login:
    post:
      operationId: Login
      responses:
        200:
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/LoginResponse'
        400:
          description: Bad Request
`)

	endpoints, err := parseOpenAPI(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(endpoints))
	}

	ep := endpoints[0]
	if len(ep.ResponseFields) == 0 {
		t.Fatal("expected non-empty ResponseFields")
	}

	byPath := make(map[string]string)
	for _, f := range ep.ResponseFields {
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

func TestParseOpenAPI_NoResponseSchema(t *testing.T) {
	data := []byte(`
openapi: '3.0.0'
info:
  title: Test API
paths:
  /api/v1/users:
    get:
      operationId: ListUsers
      responses:
        200:
          description: OK
`)

	endpoints, err := parseOpenAPI(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(endpoints) != 1 {
		t.Fatalf("expected 1 endpoint, got %d", len(endpoints))
	}

	if len(endpoints[0].ResponseFields) != 0 {
		t.Fatalf("expected no ResponseFields, got %d", len(endpoints[0].ResponseFields))
	}
}

func TestParseOpenAPI_MissingPaths(t *testing.T) {
	data := []byte(`
openapi: '3.0.0'
info:
  title: Test API
`)

	_, err := parseOpenAPI(data)
	if err == nil {
		t.Fatal("expected error for missing paths")
	}
}
