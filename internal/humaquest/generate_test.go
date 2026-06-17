package humaquest

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/llm"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// ---------------------------------------------------------------------------
// Test doubles for the --generate path: a stub llm.Backend and a render-aware
// fake gate.Definition. (fakeProbe / fakeDef / covWith / coverEP /
// seedCoverSession live in cover_orchestration_test.go and cri_helpers_test.go.)
// ---------------------------------------------------------------------------

// stubBackend is a canned llm.Backend: it records the system/user prompt and
// returns the configured completion or error. CallFunc adapts the closure to
// the Backend interface (the test-injection seam).
func stubBackend(content string, err error, captured *string) llm.Backend {
	return llm.CallFunc(func(system, user string) (string, error) {
		if captured != nil {
			*captured = user
		}
		return content, err
	})
}

// renderDef is a fake Definition whose Render returns a configurable base prompt,
// so generatePrompt can be exercised in isolation. It also implements Evaluator.
type renderDef struct {
	fakeDef
	base       string
	renderErr  error
	renderCall int
}

func (d *renderDef) Render(s *quest.Session, it *quest.Item) (string, error) {
	d.renderCall++
	if d.renderErr != nil {
		return "", d.renderErr
	}
	return d.base, nil
}

var (
	_ gate.Definition = (*renderDef)(nil)
	_ gate.Evaluator  = (*renderDef)(nil)
)

// itemWithReason builds an Item whose last Attempt carries the given reason and
// whose Tries is set, so generatePrompt's retry branch (Tries>0 + lastReason) is
// exercised.
func itemWithReason(tries int, reason string) *quest.Item {
	it := &quest.Item{Key: "x", State: quest.TODO, Tries: tries}
	if reason != "" {
		it.Log = []quest.Attempt{{Reason: reason}}
	}
	return it
}

// genItemSetup chdirs to a fresh temp dir (no manifest → HurlDir defaults to
// "hurl") and returns a session + the single TODO item plus the conventional
// hurl path for the endpoint.
func genItemSetup(t *testing.T, ep scanner.Endpoint) (*quest.Session, *quest.Item, string) {
	t.Helper()
	dir := t.TempDir()
	chdir(t, dir)
	sessPath := seedCoverSession(t, dir, ep)
	s, err := quest.Load(sessPath)
	if err != nil {
		t.Fatal(err)
	}
	return s, s.Items[0], filepath.Join("hurl", filepath.Base(runner.HurlFileName(&ep, "hurl")))
}

// ---------------------------------------------------------------------------
// sanitizeHurl
// ---------------------------------------------------------------------------

func TestSanitizeHurl(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"clean", "GET /x\nHTTP 200\n", "GET /x\nHTTP 200"},
		{"trim_only", "   GET /x\nHTTP 200   ", "GET /x\nHTTP 200"},
		{"fenced_with_lang", "here:\n```hurl\nGET /x\nHTTP 200\n```\ndone", "GET /x\nHTTP 200"},
		{"fenced_no_lang", "```\nGET /x\nHTTP 200\n```", "GET /x\nHTTP 200"},
		{"fence_no_closing", "```hurl\nGET /x\nHTTP 200", "GET /x\nHTTP 200"},
		{"fence_no_newline", "```", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := sanitizeHurl(c.in); got != c.want {
				t.Errorf("sanitizeHurl(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// stashHurl / writeHurl / restoreHurl
// ---------------------------------------------------------------------------

func TestStashHurl_NoPriorFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "absent.hurl")
	st, err := stashHurl(path)
	if err != nil {
		t.Fatalf("stashHurl: %v", err)
	}
	if st.existed {
		t.Error("existed = true for an absent file, want false")
	}
	if st.path != path {
		t.Errorf("path = %q, want %q", st.path, path)
	}
}

func TestStashHurl_ExistingFile(t *testing.T) {
	path := writeTempFile(t, "prior.hurl", "PRIOR CONTENT\n")
	st, err := stashHurl(path)
	if err != nil {
		t.Fatalf("stashHurl: %v", err)
	}
	if !st.existed {
		t.Fatal("existed = false for a present file, want true")
	}
	if string(st.content) != "PRIOR CONTENT\n" {
		t.Errorf("content = %q, want PRIOR CONTENT", st.content)
	}
}

func TestStashHurl_ReadError(t *testing.T) {
	// A directory at the path is not IsNotExist → ReadFile returns a genuine error.
	dir := t.TempDir()
	if _, err := stashHurl(dir); err == nil {
		t.Fatal("want read error when path is a directory")
	}
}

func TestWriteHurl_CreatesParentDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "deep", "out.hurl")
	if err := writeHurl(path, "BODY\n"); err != nil {
		t.Fatalf("writeHurl: %v", err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "BODY\n" {
		t.Errorf("content = %q, want BODY", b)
	}
}

func TestWriteHurl_MkdirError(t *testing.T) {
	blocker := writeTempFile(t, "blocker", "x")
	if err := writeHurl(filepath.Join(blocker, "out.hurl"), "y"); err == nil {
		t.Fatal("want mkdir error when parent path is a file")
	}
}

// Stash an existing file → overwrite → restore brings back the original bytes.
func TestRestoreHurl_RestoresPriorContent(t *testing.T) {
	path := writeTempFile(t, "x.hurl", "ORIGINAL\n")
	st, err := stashHurl(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeHurl(path, "GENERATED\n"); err != nil {
		t.Fatal(err)
	}
	if err := restoreHurl(st); err != nil {
		t.Fatalf("restoreHurl: %v", err)
	}
	b, _ := os.ReadFile(path)
	if string(b) != "ORIGINAL\n" {
		t.Errorf("restored content = %q, want ORIGINAL (byte-identical)", b)
	}
}

// Stash when no prior file → write new → restore deletes it (no-file state).
func TestRestoreHurl_DeletesWhenNoPrior(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new.hurl")
	st, err := stashHurl(path) // existed=false
	if err != nil {
		t.Fatal(err)
	}
	if err := writeHurl(path, "GENERATED\n"); err != nil {
		t.Fatal(err)
	}
	if err := restoreHurl(st); err != nil {
		t.Fatalf("restoreHurl: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("file should be removed, stat err = %v", err)
	}
}

// existed=false + already-absent file → Remove's IsNotExist is swallowed.
func TestRestoreHurl_NoPriorAlreadyAbsent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ghost.hurl")
	if err := restoreHurl(hurlStash{path: path, existed: false}); err != nil {
		t.Errorf("restoreHurl over absent file should be nil, got %v", err)
	}
}

// existed=false but the path is a NON-empty directory → os.Remove fails with a
// genuine (non-IsNotExist) error, which must propagate.
func TestRestoreHurl_NoPrior_RemoveError(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "busy")
	if err := os.MkdirAll(filepath.Join(target, "child"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := restoreHurl(hurlStash{path: target, existed: false}); err == nil {
		t.Fatal("want Remove error for a non-empty directory")
	}
}

// existed=true but the path is a directory → os.WriteFile fails, propagating.
func TestRestoreHurl_WriteBackError(t *testing.T) {
	dir := t.TempDir()
	if err := restoreHurl(hurlStash{path: dir, content: []byte("x"), existed: true}); err == nil {
		t.Fatal("want WriteFile error when the path is a directory")
	}
}

// ---------------------------------------------------------------------------
// generatePrompt
// ---------------------------------------------------------------------------

func TestGeneratePrompt_FirstAttempt_NoFeedback(t *testing.T) {
	def := &renderDef{base: "BASE PROMPT"}
	it := itemWithReason(0, "ignored reason")
	got, err := generatePrompt(def, nil, it)
	if err != nil {
		t.Fatalf("generatePrompt: %v", err)
	}
	if got != "BASE PROMPT" {
		t.Errorf("Tries==0 prompt = %q, want just the base", got)
	}
	if strings.Contains(got, "Retry") {
		t.Error("first attempt must not carry the retry preamble")
	}
}

func TestGeneratePrompt_Retry_AppendsCoachingAndLastReason(t *testing.T) {
	def := &renderDef{base: "BASE"}
	it := itemWithReason(1, "C-03: body too shallow")
	got, err := generatePrompt(def, nil, it)
	if err != nil {
		t.Fatalf("generatePrompt: %v", err)
	}
	if !strings.HasPrefix(got, "BASE") {
		t.Errorf("prompt should start with the base, got %q", got)
	}
	if !strings.Contains(got, ruleSystem["C-03"]) {
		t.Error("retry prompt missing the C-03 static coaching preamble")
	}
	if !strings.Contains(got, "Previous verdict:") || !strings.Contains(got, "body too shallow") {
		t.Errorf("retry prompt missing lastReason feedback: %q", got)
	}
}

func TestGeneratePrompt_Retry_NoLogNoFeedbackLine(t *testing.T) {
	def := &renderDef{base: "BASE"}
	it := itemWithReason(2, "") // Tries>0 but empty log → lastReason == ""
	got, err := generatePrompt(def, nil, it)
	if err != nil {
		t.Fatalf("generatePrompt: %v", err)
	}
	if !strings.Contains(got, ruleSystem["C-03"]) {
		t.Error("retry preamble should still be present")
	}
	if strings.Contains(got, "Previous verdict:") {
		t.Error("no log → must not append a Previous verdict line")
	}
}

func TestGeneratePrompt_RenderError(t *testing.T) {
	def := &renderDef{renderErr: errSentinel}
	if _, err := generatePrompt(def, nil, itemWithReason(0, "")); err != errSentinel {
		t.Fatalf("err = %v, want sentinel from Render", err)
	}
}

// ---------------------------------------------------------------------------
// runtimeFailVerdict
// ---------------------------------------------------------------------------

func TestRuntimeFailVerdict(t *testing.T) {
	ep := scanner.Endpoint{Method: "GET", Path: "/api/v1/users"}
	v := runtimeFailVerdict(ep, "assert status==200 but got 500")
	if v.Outcome != quest.OutFail {
		t.Errorf("Outcome = %s, want FAIL", v.Outcome)
	}
	if v.RootCause != "H-03" {
		t.Errorf("RootCause = %q, want H-03", v.RootCause)
	}
	if len(v.Facts) != 1 {
		t.Fatalf("Facts = %d, want 1", len(v.Facts))
	}
	f := v.Facts[0]
	if f.Where != "GET /api/v1/users" {
		t.Errorf("Where = %q, want endpoint key", f.Where)
	}
	if !strings.Contains(f.Actual, "assert status==200 but got 500") {
		t.Errorf("Fact.Actual must embed the runtime feedback, got %q", f.Actual)
	}
	if !strings.Contains(v.Feedback, ruleSystem["H-03"]) {
		t.Error("Feedback missing the H-03 static preamble")
	}
}

func TestRuntimeFailVerdict_Truncates(t *testing.T) {
	long := strings.Repeat("z", 2000)
	v := runtimeFailVerdict(scanner.Endpoint{Method: "POST", Path: "/x"}, long)
	if !strings.Contains(v.Facts[0].Actual, "…(truncated)") {
		t.Error("over-long feedback should be truncated in Fact.Actual")
	}
	if len(v.Facts[0].Actual) > 1600 {
		t.Errorf("Fact.Actual not capped: len=%d", len(v.Facts[0].Actual))
	}
}

// ---------------------------------------------------------------------------
// parseFailVerdict
// ---------------------------------------------------------------------------

func TestParseFailVerdict(t *testing.T) {
	ep := scanner.Endpoint{Method: "GET", Path: "/api/v1/users"}
	v := parseFailVerdict(ep, errSentinel)
	if v.Outcome != quest.OutFail || v.RootCause != "H-03" {
		t.Errorf("verdict = %+v, want FAIL/H-03", v)
	}
	if len(v.Facts) != 1 || v.Facts[0].Where != "GET /api/v1/users" {
		t.Fatalf("Facts wrong: %+v", v.Facts)
	}
	if !strings.Contains(v.Facts[0].Actual, errSentinel.Error()) {
		t.Errorf("Fact.Actual must carry the parse error, got %q", v.Facts[0].Actual)
	}
	if !strings.Contains(v.Feedback, "ENTIRE .hurl") {
		t.Errorf("Feedback should instruct emitting the whole file, got %q", v.Feedback)
	}
}

// ---------------------------------------------------------------------------
// persistCoverState
// ---------------------------------------------------------------------------

func TestPersistCoverState(t *testing.T) {
	ep := coverEP("a")
	it := &quest.Item{Key: ep.ID, State: quest.TODO}
	cov := covWith(10, 63.0, 1, 2)
	ps := payloadState{ImproveCount: 2}

	if err := persistCoverState(it, ep, cov, ps); err != nil {
		t.Fatalf("persistCoverState: %v", err)
	}
	var got payloadState
	if err := it.DecodePayload(&got); err != nil {
		t.Fatal(err)
	}
	if got.PrevCoverage != 63.0 {
		t.Errorf("PrevCoverage = %v, want 63.0", got.PrevCoverage)
	}
	if got.ImproveCount != 3 {
		t.Errorf("ImproveCount = %d, want 3 (prior 2 + 1)", got.ImproveCount)
	}
	if got.Endpoint.ID != ep.ID {
		t.Errorf("Endpoint not round-tripped: %+v", got.Endpoint)
	}
}

func TestPersistCoverState_NilCoverageZeroPercent(t *testing.T) {
	it := &quest.Item{Key: "a", State: quest.TODO}
	if err := persistCoverState(it, coverEP("a"), nil, payloadState{}); err != nil {
		t.Fatalf("persistCoverState: %v", err)
	}
	var got payloadState
	_ = it.DecodePayload(&got)
	if got.PrevCoverage != 0 {
		t.Errorf("PrevCoverage = %v, want 0 for nil coverage", got.PrevCoverage)
	}
	if got.ImproveCount != 1 {
		t.Errorf("ImproveCount = %d, want 1", got.ImproveCount)
	}
}

// ---------------------------------------------------------------------------
// buildBackend
// ---------------------------------------------------------------------------

func TestBuildBackend_GenerateOff_NilNoError(t *testing.T) {
	b, err := buildBackend(false, "claude:sonnet")
	if err != nil {
		t.Fatalf("buildBackend: %v", err)
	}
	if b != nil {
		t.Errorf("backend = %v, want nil in manual mode", b)
	}
}

func TestBuildBackend_GenerateOn_BogusModel(t *testing.T) {
	// No colon → FromFlag rejects it.
	b, err := buildBackend(true, "no-colon-here")
	if err == nil {
		t.Fatal("want error for an unparseable --model")
	}
	if b != nil {
		t.Errorf("backend = %v, want nil on error", b)
	}
}

func TestBuildBackend_GenerateOn_Success(t *testing.T) {
	// ollama resolves to a Backend with no CLI/network needed at construction.
	b, err := buildBackend(true, "ollama:gemma")
	if err != nil {
		t.Fatalf("buildBackend: %v", err)
	}
	if b == nil {
		t.Fatal("backend = nil, want a constructed Backend")
	}
}

// ---------------------------------------------------------------------------
// evaluateGenerated
// ---------------------------------------------------------------------------

func TestEvaluateGenerated_InjectsCoverageAndEvaluates(t *testing.T) {
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutPass}}
	cov := covWith(10, 80, 1, 2)
	v, err := evaluateGenerated(def, quest.New(), &quest.Item{Key: "a"}, cov)
	if err != nil {
		t.Fatalf("evaluateGenerated: %v", err)
	}
	if v.Outcome != quest.OutPass {
		t.Errorf("verdict = %s, want PASS", v.Outcome)
	}
	if def.evalCalls != 1 {
		t.Errorf("Evaluate calls = %d, want 1", def.evalCalls)
	}
	ground := def.lastCtx.Grounds["coverage"]
	if ground == "" {
		t.Fatal("coverage not injected as a ground")
	}
	got, present := decodeCoverage(ground)
	if !present || got.Percent != 80 {
		t.Errorf("injected coverage wrong: %+v present=%v", got, present)
	}
}

func TestEvaluateGenerated_NilCoverage_NoGround(t *testing.T) {
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutReview}}
	if _, err := evaluateGenerated(def, quest.New(), &quest.Item{Key: "a"}, nil); err != nil {
		t.Fatalf("evaluateGenerated: %v", err)
	}
	if _, ok := def.lastCtx.Grounds["coverage"]; ok {
		t.Error("ground injected despite nil coverage")
	}
}

func TestEvaluateGenerated_ShortCircuit(t *testing.T) {
	short := &quest.Verdict{Outcome: quest.OutSkip, Feedback: "exempt"}
	def := &fakeDef{short: short}
	v, err := evaluateGenerated(def, quest.New(), &quest.Item{Key: "a"}, covWith(1, 1, 1))
	if err != nil {
		t.Fatalf("evaluateGenerated: %v", err)
	}
	if v.Outcome != quest.OutSkip {
		t.Errorf("verdict = %s, want the short-circuit SKIP", v.Outcome)
	}
	if def.evalCalls != 0 {
		t.Error("Evaluate must not run on a short-circuit")
	}
}

func TestEvaluateGenerated_PrepareError(t *testing.T) {
	def := &fakeDef{prepErr: errSentinel}
	if _, err := evaluateGenerated(def, quest.New(), &quest.Item{Key: "a"}, nil); err != errSentinel {
		t.Fatalf("err = %v, want sentinel from Prepare", err)
	}
}

func TestEvaluateGenerated_CoverageGroundError(t *testing.T) {
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutPass}}
	// Non-finite Percent fails json.Marshal in coverageGround.
	cov := &adapter.CoverageResult{Total: 1, Percent: math.Inf(1)}
	if _, err := evaluateGenerated(def, quest.New(), &quest.Item{Key: "a"}, cov); err == nil {
		t.Fatal("want coverageGround marshal error")
	}
	if def.evalCalls != 0 {
		t.Error("Evaluate must not run after a coverageGround error")
	}
}

// ---------------------------------------------------------------------------
// measureGenerated
// ---------------------------------------------------------------------------

func TestMeasureGenerated_RuntimePass_Evaluates(t *testing.T) {
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutPass}}
	probe := &fakeProbe{cov: covWith(10, 90, 1), res: &runner.Result{Pass: true}}
	v, cov, err := measureGenerated(def, probe, quest.New(), &quest.Item{Key: "a"}, coverEP("a"))
	if err != nil {
		t.Fatalf("measureGenerated: %v", err)
	}
	if v.Outcome != quest.OutPass {
		t.Errorf("verdict = %s, want PASS", v.Outcome)
	}
	if def.evalCalls != 1 {
		t.Errorf("Evaluate calls = %d, want 1 (runtime passed)", def.evalCalls)
	}
	if cov == nil || cov.Percent != 90 {
		t.Errorf("cov = %+v, want the measured coverage", cov)
	}
	if probe.resetCalls != 1 {
		t.Errorf("Reset calls = %d, want 1", probe.resetCalls)
	}
}

func TestMeasureGenerated_RuntimeFail_ForcesRuntimeFailVerdict(t *testing.T) {
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutPass}}
	probe := &fakeProbe{
		cov: covWith(10, 90, 1),
		res: &runner.Result{Pass: false, Feedback: "assertion failed at line 3"},
	}
	v, cov, err := measureGenerated(def, probe, quest.New(), &quest.Item{Key: "a"}, coverEP("a"))
	if err != nil {
		t.Fatalf("measureGenerated: %v", err)
	}
	if v.Outcome != quest.OutFail || v.RootCause != "H-03" {
		t.Errorf("verdict = %+v, want forced FAIL/H-03 on runtime failure", v)
	}
	if def.evalCalls != 0 {
		t.Error("Evaluate must NOT run when the hurl failed at runtime (§4 hole closure)")
	}
	if !strings.Contains(v.Facts[0].Actual, "assertion failed at line 3") {
		t.Errorf("runtime feedback not threaded into the verdict: %q", v.Facts[0].Actual)
	}
	if cov == nil {
		t.Error("measured coverage should still be returned for persistence")
	}
}

func TestMeasureGenerated_NilResult_Evaluates(t *testing.T) {
	// A nil *runner.Result (e.g. no live hurl signal) is not a runtime failure.
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutReview}}
	probe := &fakeProbe{cov: nil, res: nil}
	v, _, err := measureGenerated(def, probe, quest.New(), &quest.Item{Key: "a"}, coverEP("a"))
	if err != nil {
		t.Fatalf("measureGenerated: %v", err)
	}
	if def.evalCalls != 1 || v.Outcome != quest.OutReview {
		t.Errorf("nil result should fall through to evaluate, got %+v evalCalls=%d", v, def.evalCalls)
	}
}

func TestMeasureGenerated_ResetError(t *testing.T) {
	probe := &fakeProbe{resetErr: errSentinel}
	if _, _, err := measureGenerated(&fakeDef{}, probe, quest.New(), &quest.Item{Key: "a"}, coverEP("a")); err != errSentinel {
		t.Fatalf("err = %v, want sentinel from Reset", err)
	}
}

func TestMeasureGenerated_MeasureError(t *testing.T) {
	probe := &fakeProbe{measureErr: errSentinel}
	if _, _, err := measureGenerated(&fakeDef{}, probe, quest.New(), &quest.Item{Key: "a"}, coverEP("a")); err != errSentinel {
		t.Fatalf("err = %v, want sentinel from Measure", err)
	}
}

// ---------------------------------------------------------------------------
// generateItem (orchestration: stub backend + fake probe, NO real LLM/server)
// ---------------------------------------------------------------------------

func TestGenerateItem_Exempt_ShortCircuitsSkip(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	writeUnreachable(t, dir,
		"- endpoint: GET /api/v1/users\n  status: 200\n  reason: dead\n  evidence: e\n")
	ep := scanner.Endpoint{ID: "e", Method: "GET", Path: "/api/v1/users",
		Responses: json.RawMessage(`[{"status":200,"line":1}]`)}

	var called string
	backend := stubBackend("should not be called", nil, &called)
	v, err := generateItem(&renderDef{}, &fakeProbe{}, quest.New(), &quest.Item{Key: "e"}, ep,
		payloadState{}, backend)
	if err != nil {
		t.Fatalf("generateItem: %v", err)
	}
	if v.Outcome != quest.OutSkip {
		t.Errorf("verdict = %s, want SKIP for a fully-exempt endpoint", v.Outcome)
	}
	if called != "" {
		t.Error("backend must not be called on the exempt short-circuit")
	}
}

func TestGenerateItem_HappyPath_PassKeepsGeneratedFile(t *testing.T) {
	ep := scanner.Endpoint{ID: "u", Method: "GET", Path: "/api/v1/users"}
	s, it, hurlPath := genItemSetup(t, ep)

	var capturedPrompt string
	backend := stubBackend(validHurl, nil, &capturedPrompt)
	probe := &fakeProbe{cov: covWith(10, 95, 1), res: &runner.Result{Pass: true}}
	def := &renderDef{base: "AUTHORING", fakeDef: fakeDef{verdict: quest.Verdict{Outcome: quest.OutPass}}}

	v, err := generateItem(def, probe, s, it, ep, payloadState{}, backend)
	if err != nil {
		t.Fatalf("generateItem: %v", err)
	}
	if v.Outcome != quest.OutPass {
		t.Errorf("verdict = %s, want PASS", v.Outcome)
	}
	if def.evalCalls != 1 {
		t.Errorf("Evaluate calls = %d, want 1", def.evalCalls)
	}
	// Generated file kept on PASS, byte-identical to the sanitized completion.
	b, err := os.ReadFile(hurlPath)
	if err != nil {
		t.Fatalf("generated file not kept: %v", err)
	}
	if string(b) != sanitizeHurl(validHurl) {
		t.Errorf("kept file content = %q, want sanitized completion", b)
	}
	if !strings.Contains(capturedPrompt, "AUTHORING") {
		t.Errorf("backend got prompt %q, want it built from Render", capturedPrompt)
	}
	// payloadState persisted.
	var ps payloadState
	_ = it.DecodePayload(&ps)
	if ps.ImproveCount != 1 || ps.PrevCoverage != 95 {
		t.Errorf("payload not persisted: %+v", ps)
	}
}

func TestGenerateItem_RuntimeFail_RestoresPriorNoEvaluate(t *testing.T) {
	ep := scanner.Endpoint{ID: "u", Method: "GET", Path: "/api/v1/users"}
	s, it, hurlPath := genItemSetup(t, ep)

	// A human-authored prior hurl exists and must survive a non-PASS attempt.
	if err := os.MkdirAll(filepath.Dir(hurlPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hurlPath, []byte("PRIOR HUMAN HURL\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	backend := stubBackend(validHurl, nil, nil)
	probe := &fakeProbe{
		cov: covWith(10, 40, 1),
		res: &runner.Result{Pass: false, Feedback: "status 500 != 200"},
	}
	def := &renderDef{fakeDef: fakeDef{verdict: quest.Verdict{Outcome: quest.OutPass}}}

	v, err := generateItem(def, probe, s, it, ep, payloadState{}, backend)
	if err != nil {
		t.Fatalf("generateItem: %v", err)
	}
	if v.Outcome != quest.OutFail || v.RootCause != "H-03" {
		t.Errorf("verdict = %+v, want forced runtime FAIL", v)
	}
	if def.evalCalls != 0 {
		t.Error("gate Evaluate must not run on a runtime failure")
	}
	// Prior asset restored byte-identical.
	b, _ := os.ReadFile(hurlPath)
	if string(b) != "PRIOR HUMAN HURL\n" {
		t.Errorf("prior hurl not restored, got %q", b)
	}
}

func TestGenerateItem_RuntimeFail_NoPrior_RemovesGenerated(t *testing.T) {
	ep := scanner.Endpoint{ID: "u", Method: "GET", Path: "/api/v1/users"}
	s, it, hurlPath := genItemSetup(t, ep)

	backend := stubBackend(validHurl, nil, nil)
	probe := &fakeProbe{cov: covWith(1, 10, 1), res: &runner.Result{Pass: false, Feedback: "boom"}}
	def := &renderDef{fakeDef: fakeDef{verdict: quest.Verdict{Outcome: quest.OutPass}}}

	if _, err := generateItem(def, probe, s, it, ep, payloadState{}, backend); err != nil {
		t.Fatalf("generateItem: %v", err)
	}
	if _, err := os.Stat(hurlPath); !os.IsNotExist(err) {
		t.Errorf("generated file should be removed (no prior), stat err = %v", err)
	}
}

func TestGenerateItem_ParseFail_RestoresNoMeasure(t *testing.T) {
	ep := scanner.Endpoint{ID: "u", Method: "GET", Path: "/api/v1/users"}
	s, it, hurlPath := genItemSetup(t, ep)

	// An over-long single line is written to disk but ParseHurlEntries cannot parse it.
	backend := stubBackend(strings.Repeat("A", 2*1024*1024), nil, nil)
	probe := &fakeProbe{cov: covWith(1, 1, 1), res: &runner.Result{Pass: true}}
	def := &renderDef{fakeDef: fakeDef{verdict: quest.Verdict{Outcome: quest.OutPass}}}

	v, err := generateItem(def, probe, s, it, ep, payloadState{}, backend)
	if err != nil {
		t.Fatalf("generateItem: %v", err)
	}
	if v.Outcome != quest.OutFail || v.RootCause != "H-03" {
		t.Errorf("verdict = %+v, want parse FAIL", v)
	}
	if !strings.Contains(v.Feedback, "ENTIRE .hurl") {
		t.Errorf("want parseFailVerdict feedback, got %q", v.Feedback)
	}
	if probe.measureCalls != 0 {
		t.Errorf("Measure ran on an unparseable hurl: %d", probe.measureCalls)
	}
	if def.evalCalls != 0 {
		t.Error("Evaluate must not run on a parse failure")
	}
	// No prior → generated discarded.
	if _, err := os.Stat(hurlPath); !os.IsNotExist(err) {
		t.Errorf("unparseable generated file should be removed, stat err = %v", err)
	}
}

func TestGenerateItem_BackendError(t *testing.T) {
	ep := scanner.Endpoint{ID: "u", Method: "GET", Path: "/api/v1/users"}
	s, it, hurlPath := genItemSetup(t, ep)
	backend := stubBackend("", errSentinel, nil)
	_, err := generateItem(&renderDef{}, &fakeProbe{}, s, it, ep, payloadState{}, backend)
	if err != errSentinel {
		t.Fatalf("err = %v, want sentinel from backend.Complete", err)
	}
	// Nothing written before the backend error.
	if _, statErr := os.Stat(hurlPath); !os.IsNotExist(statErr) {
		t.Errorf("no file should exist after a backend error, stat err = %v", statErr)
	}
}

func TestGenerateItem_RenderError(t *testing.T) {
	ep := scanner.Endpoint{ID: "u", Method: "GET", Path: "/api/v1/users"}
	s, it, _ := genItemSetup(t, ep)
	def := &renderDef{renderErr: errSentinel}
	backend := stubBackend(validHurl, nil, nil)
	if _, err := generateItem(def, &fakeProbe{}, s, it, ep, payloadState{}, backend); err != errSentinel {
		t.Fatalf("err = %v, want sentinel from generatePrompt/Render", err)
	}
}

func TestGenerateItem_MeasureError_RestoresPrior(t *testing.T) {
	ep := scanner.Endpoint{ID: "u", Method: "GET", Path: "/api/v1/users"}
	s, it, hurlPath := genItemSetup(t, ep)
	if err := os.MkdirAll(filepath.Dir(hurlPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hurlPath, []byte("PRIOR\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	backend := stubBackend(validHurl, nil, nil)
	probe := &fakeProbe{measureErr: errSentinel}
	def := &renderDef{fakeDef: fakeDef{verdict: quest.Verdict{Outcome: quest.OutPass}}}
	if _, err := generateItem(def, probe, s, it, ep, payloadState{}, backend); err != errSentinel {
		t.Fatalf("err = %v, want sentinel from Measure", err)
	}
	// Best-effort restore on the error path.
	b, _ := os.ReadFile(hurlPath)
	if string(b) != "PRIOR\n" {
		t.Errorf("prior hurl not restored on Measure error, got %q", b)
	}
}
