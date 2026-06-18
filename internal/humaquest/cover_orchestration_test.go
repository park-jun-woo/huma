package humaquest

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/gate"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// errSentinel is a canned error used to assert error propagation through the
// cover orchestration.
var errSentinel = errors.New("sentinel")

// ---------------------------------------------------------------------------
// Test doubles: a fake coverageProbe and fake gate.Definition(s) so the cover
// orchestration (runCover / coverItem) is exercised with NO real server.
// ---------------------------------------------------------------------------

// fakeProbe is a canned coverageProbe: it counts calls and returns the configured
// CoverageResult/error from each method, replacing the real liveProbe.
type fakeProbe struct {
	upErr      error
	resetErr   error
	measureErr error
	downErr    error
	cov        *adapter.CoverageResult
	res        *runner.Result

	upCalls      int
	resetCalls   int
	measureCalls int
	downCalls    int
	measuredEPs  []scanner.Endpoint
}

func (p *fakeProbe) Up() error    { p.upCalls++; return p.upErr }
func (p *fakeProbe) Reset() error { p.resetCalls++; return p.resetErr }
func (p *fakeProbe) Measure(ep scanner.Endpoint) (*adapter.CoverageResult, *runner.Result, error) {
	p.measureCalls++
	p.measuredEPs = append(p.measuredEPs, ep)
	return p.cov, p.res, p.measureErr
}
func (p *fakeProbe) Down() error { p.downCalls++; return p.downErr }

var _ coverageProbe = (*fakeProbe)(nil)

// fakeDef is a fake gate.Definition that ALSO implements gate.Evaluator, so it
// drives both the Prepare short-circuit branch and the Evaluate dispatch with
// canned values. It records the Context handed to Evaluate so tests can assert
// the injected Grounds["coverage"].
type fakeDef struct {
	short     *quest.Verdict
	prepErr   error
	verdict   quest.Verdict
	evalCalls int
	lastCtx   gate.Context
}

func (d *fakeDef) Seed(args []string) ([]*quest.Item, error)               { return nil, nil }
func (d *fakeDef) Render(s *quest.Session, it *quest.Item) (string, error) { return "", nil }
func (d *fakeDef) Rules() []gate.Rule                                      { return nil }
func (d *fakeDef) Prepare(s *quest.Session, it *quest.Item, raw []byte) (gate.Context, *quest.Verdict, error) {
	if d.prepErr != nil {
		return gate.Context{}, nil, d.prepErr
	}
	if d.short != nil {
		return gate.Context{}, d.short, nil
	}
	return gate.Context{Item: it}, nil, nil
}
func (d *fakeDef) Evaluate(ctx gate.Context) quest.Verdict {
	d.evalCalls++
	d.lastCtx = ctx
	return d.verdict
}

var (
	_ gate.Definition = (*fakeDef)(nil)
	_ gate.Evaluator  = (*fakeDef)(nil)
)

// nonEvalDef implements only gate.Definition (NOT gate.Evaluator), to drive the
// evaluate() fallback branch (gate.Evaluate over the Rules() catalog).
type nonEvalDef struct {
	rules []gate.Rule
}

func (d nonEvalDef) Seed(args []string) ([]*quest.Item, error)               { return nil, nil }
func (d nonEvalDef) Render(s *quest.Session, it *quest.Item) (string, error) { return "", nil }
func (d nonEvalDef) Rules() []gate.Rule                                      { return d.rules }
func (d nonEvalDef) Prepare(s *quest.Session, it *quest.Item, raw []byte) (gate.Context, *quest.Verdict, error) {
	return gate.Context{}, nil, nil
}

var _ gate.Definition = nonEvalDef{}

// seedCoverSession writes a session file with one TODO item per endpoint, each
// carrying a payloadState. Returns the session path.
func seedCoverSession(t *testing.T, dir string, eps ...scanner.Endpoint) string {
	t.Helper()
	s := quest.New()
	for _, ep := range eps {
		it := &quest.Item{Key: ep.ID, State: quest.TODO}
		if err := it.SetPayload(payloadState{Endpoint: ep}); err != nil {
			t.Fatalf("SetPayload: %v", err)
		}
		s.Items = append(s.Items, it)
	}
	path := filepath.Join(dir, "session.json")
	if err := s.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return path
}

func coverEP(id string) scanner.Endpoint {
	return scanner.Endpoint{ID: id, Method: "GET", Path: "/api/v1/" + id}
}

// countLines returns the number of non-empty lines in a file (0 if absent).
func countLines(t *testing.T, path string) int {
	t.Helper()
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return 0
	}
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	n := 0
	for _, ln := range strings.Split(strings.TrimRight(string(b), "\n"), "\n") {
		if strings.TrimSpace(ln) != "" {
			n++
		}
	}
	return n
}

// ---------------------------------------------------------------------------
// runCover
// ---------------------------------------------------------------------------

func TestRunCover_PassItem_TerminalAndExported(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	outPath := filepath.Join(dir, "out", "results.jsonl")

	probe := &fakeProbe{cov: covWith(10, 90, 1)}
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutPass, Feedback: "covered"}}

	var buf bytes.Buffer
	if err := runCover(def, probe, nil, 0, sessPath, outPath, &buf); err != nil {
		t.Fatalf("runCover: %v", err)
	}

	// Up once, Down deferred once.
	if probe.upCalls != 1 || probe.downCalls != 1 {
		t.Errorf("Up=%d Down=%d, want 1/1", probe.upCalls, probe.downCalls)
	}

	// Item locked PASS in the saved session.
	s, err := quest.Load(sessPath)
	if err != nil {
		t.Fatal(err)
	}
	if s.Items[0].State != quest.PASS {
		t.Errorf("item state = %s, want PASS", s.Items[0].State)
	}
	// Exported exactly once.
	if n := countLines(t, outPath); n != 1 {
		t.Errorf("JSONL lines = %d, want 1", n)
	}
	if !strings.Contains(buf.String(), "PASS") {
		t.Errorf("render missing PASS: %q", buf.String())
	}
}

func TestRunCover_ImproveRetriesUntilDONE(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	outPath := filepath.Join(dir, "results.jsonl")

	// OutFail is a retryable attempt: Apply locks DONE at MaxTries (3).
	probe := &fakeProbe{cov: covWith(10, 50, 1)}
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutFail, RootCause: "R1"}}

	if err := runCover(def, probe, nil, 0, sessPath, outPath, &bytes.Buffer{}); err != nil {
		t.Fatalf("runCover: %v", err)
	}

	// Re-measured each retry: MaxTries attempts.
	if probe.measureCalls != quest.MaxTries {
		t.Errorf("Measure calls = %d, want %d (one per retry)", probe.measureCalls, quest.MaxTries)
	}

	s, err := quest.Load(sessPath)
	if err != nil {
		t.Fatal(err)
	}
	if s.Items[0].State != quest.DONE {
		t.Errorf("item state = %s, want DONE", s.Items[0].State)
	}
	if s.Items[0].Tries != quest.MaxTries {
		t.Errorf("Tries = %d, want %d", s.Items[0].Tries, quest.MaxTries)
	}
	// ImproveCount persisted across retries.
	var ps payloadState
	if err := s.Items[0].DecodePayload(&ps); err != nil {
		t.Fatal(err)
	}
	if ps.ImproveCount != quest.MaxTries {
		t.Errorf("ImproveCount = %d, want %d", ps.ImproveCount, quest.MaxTries)
	}
	// Terminal item exported once.
	if n := countLines(t, outPath); n != 1 {
		t.Errorf("JSONL lines = %d, want 1", n)
	}
}

func TestRunCover_PrepareShortCircuit(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	outPath := filepath.Join(dir, "results.jsonl")

	probe := &fakeProbe{cov: covWith(10, 90, 1)}
	// Prepare returns a SKIP short verdict; Evaluate must NOT be called.
	def := &fakeDef{short: &quest.Verdict{Outcome: quest.OutSkip, Feedback: "all exempt"}}

	if err := runCover(def, probe, nil, 0, sessPath, outPath, &bytes.Buffer{}); err != nil {
		t.Fatalf("runCover: %v", err)
	}
	if def.evalCalls != 0 {
		t.Errorf("Evaluate called %d times on short-circuit, want 0", def.evalCalls)
	}
	s, _ := quest.Load(sessPath)
	if s.Items[0].State != quest.SKIPPED {
		t.Errorf("item state = %s, want SKIPPED", s.Items[0].State)
	}
}

func TestRunCover_ProbeUpError_NoLoopNoDown(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))

	probe := &fakeProbe{upErr: errSentinel}
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutPass}}

	err := runCover(def, probe, nil, 0, sessPath, filepath.Join(dir, "o.jsonl"), &bytes.Buffer{})
	if err != errSentinel {
		t.Fatalf("err = %v, want sentinel", err)
	}
	if probe.measureCalls != 0 {
		t.Errorf("Measure called %d times after Up error, want 0", probe.measureCalls)
	}
	// Down is deferred only AFTER a successful Up.
	if probe.downCalls != 0 {
		t.Errorf("Down called %d times after Up error, want 0", probe.downCalls)
	}
}

func TestRunCover_EmptySession(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir) // no items
	probe := &fakeProbe{}
	def := &fakeDef{}

	if err := runCover(def, probe, nil, 0, sessPath, filepath.Join(dir, "o.jsonl"), &bytes.Buffer{}); err != nil {
		t.Fatalf("runCover: %v", err)
	}
	if probe.measureCalls != 0 {
		t.Errorf("Measure called on empty session: %d", probe.measureCalls)
	}
	// Up/Down still ran (compile once / teardown).
	if probe.upCalls != 1 || probe.downCalls != 1 {
		t.Errorf("Up=%d Down=%d, want 1/1", probe.upCalls, probe.downCalls)
	}
}

func TestRunCover_LoadSessionError(t *testing.T) {
	probe := &fakeProbe{}
	def := &fakeDef{}
	err := runCover(def, probe, nil, 0, filepath.Join(t.TempDir(), "missing.json"), "o.jsonl", &bytes.Buffer{})
	if err == nil {
		t.Fatal("want error for missing session")
	}
	if probe.upCalls != 0 {
		t.Error("Up should not run when session load fails")
	}
}

func TestRunCover_CoverItemError_Propagates(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	probe := &fakeProbe{measureErr: errSentinel, cov: covWith(1, 1, 1)}
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutPass}}

	err := runCover(def, probe, nil, 0, sessPath, filepath.Join(dir, "o.jsonl"), &bytes.Buffer{})
	if err != errSentinel {
		t.Fatalf("err = %v, want sentinel from Measure", err)
	}
	// Down still runs (deferred after successful Up).
	if probe.downCalls != 1 {
		t.Errorf("Down=%d, want 1 (deferred)", probe.downCalls)
	}
}

func TestRunCover_SinkCreateError(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	// outPath's parent is a regular file → newJSONLSink's MkdirAll fails.
	blocker := filepath.Join(dir, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	probe := &fakeProbe{cov: covWith(1, 1, 1)}
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutPass}}

	err := runCover(def, probe, nil, 0, sessPath, filepath.Join(blocker, "out.jsonl"), &bytes.Buffer{})
	if err == nil {
		t.Fatal("want sink-create error")
	}
	// Up succeeded, so Down is deferred.
	if probe.downCalls != 1 {
		t.Errorf("Down=%d, want 1", probe.downCalls)
	}
	if probe.measureCalls != 0 {
		t.Errorf("Measure ran despite sink-create failure: %d", probe.measureCalls)
	}
}

func TestRunCover_CommitVerdictError(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	// Point --out at an existing DIRECTORY: newJSONLSink succeeds (parent exists)
	// but the first terminal item's Export → Emit(open dir) fails inside
	// commitVerdict, exercising runCover's commit-error branch.
	outDir := filepath.Join(dir, "outdir")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		t.Fatal(err)
	}
	probe := &fakeProbe{cov: covWith(1, 1, 1)}
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutPass}}

	err := runCover(def, probe, nil, 0, sessPath, outDir, &bytes.Buffer{})
	if err == nil {
		t.Fatal("want commitVerdict export error")
	}
}

// ---------------------------------------------------------------------------
// coverItem
// ---------------------------------------------------------------------------

func TestCoverItem_MeasuredBranch_InjectsCoverageAndPersists(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	s, _ := quest.Load(sessPath)
	it := s.Items[0]

	cov := covWith(10, 72.5, 1, 2)
	probe := &fakeProbe{cov: cov}
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutFail}}

	v, err := coverItem(def, probe, s, it, nil)
	if err != nil {
		t.Fatalf("coverItem: %v", err)
	}
	if v.Outcome != quest.OutFail {
		t.Errorf("verdict = %s, want FAIL", v.Outcome)
	}
	// Reset then Measure happened.
	if probe.resetCalls != 1 || probe.measureCalls != 1 {
		t.Errorf("Reset=%d Measure=%d, want 1/1", probe.resetCalls, probe.measureCalls)
	}
	// coverage injected as a ground before Evaluate.
	ground := def.lastCtx.Grounds["coverage"]
	if ground == "" {
		t.Fatal("Grounds[coverage] not injected")
	}
	gotCov, present := decodeCoverage(ground)
	if !present || gotCov == nil || gotCov.Percent != 72.5 {
		t.Errorf("injected coverage round-trip wrong: %+v present=%v", gotCov, present)
	}
	// payloadState persisted: PrevCoverage + ImproveCount++.
	var ps payloadState
	if err := it.DecodePayload(&ps); err != nil {
		t.Fatal(err)
	}
	if ps.PrevCoverage != 72.5 {
		t.Errorf("PrevCoverage = %v, want 72.5", ps.PrevCoverage)
	}
	if ps.ImproveCount != 1 {
		t.Errorf("ImproveCount = %d, want 1", ps.ImproveCount)
	}
}

func TestCoverItem_NoLiveSignal_NoGroundInjected(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	s, _ := quest.Load(sessPath)
	it := s.Items[0]

	probe := &fakeProbe{cov: nil} // no live signal
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutReview}}

	if _, err := coverItem(def, probe, s, it, nil); err != nil {
		t.Fatalf("coverItem: %v", err)
	}
	if _, ok := def.lastCtx.Grounds["coverage"]; ok {
		t.Error("coverage ground injected despite nil CoverageResult")
	}
	var ps payloadState
	_ = it.DecodePayload(&ps)
	if ps.PrevCoverage != 0 {
		t.Errorf("PrevCoverage = %v, want 0 (no signal)", ps.PrevCoverage)
	}
}

func TestCoverItem_ShortCircuit_ReturnsShortNoPersist(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	s, _ := quest.Load(sessPath)
	it := s.Items[0]

	probe := &fakeProbe{cov: covWith(1, 99, 1)}
	short := &quest.Verdict{Outcome: quest.OutSkip, Feedback: "exempt"}
	def := &fakeDef{short: short}

	v, err := coverItem(def, probe, s, it, nil)
	if err != nil {
		t.Fatalf("coverItem: %v", err)
	}
	if v.Outcome != quest.OutSkip {
		t.Errorf("verdict = %s, want SKIPPED short", v.Outcome)
	}
	if def.evalCalls != 0 {
		t.Error("Evaluate called on short-circuit")
	}
	// short-circuit returns before SetPayload → ImproveCount stays 0.
	var ps payloadState
	_ = it.DecodePayload(&ps)
	if ps.ImproveCount != 0 {
		t.Errorf("ImproveCount = %d, want 0 (short-circuit skips persist)", ps.ImproveCount)
	}
}

func TestCoverItem_CoverageGroundError(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	s, _ := quest.Load(sessPath)
	// A non-finite Percent makes json.Marshal (coverageGround) fail, driving the
	// defensive error branch before Evaluate.
	cov := &adapter.CoverageResult{Total: 1, Percent: math.Inf(1)}
	probe := &fakeProbe{cov: cov}
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutPass}}
	if _, err := coverItem(def, probe, s, s.Items[0], nil); err == nil {
		t.Fatal("want coverageGround marshal error on non-finite Percent")
	}
	if def.evalCalls != 0 {
		t.Error("Evaluate must not run after a coverageGround error")
	}
}

func TestCoverItem_DecodePayloadError(t *testing.T) {
	it := &quest.Item{Key: "bad", State: quest.TODO, Payload: json.RawMessage(`{`)}
	_, err := coverItem(&fakeDef{}, &fakeProbe{}, quest.New(), it, nil)
	if err == nil {
		t.Fatal("want DecodePayload error")
	}
}

func TestCoverItem_ResetError(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	s, _ := quest.Load(sessPath)
	probe := &fakeProbe{resetErr: errSentinel}
	if _, err := coverItem(&fakeDef{}, probe, s, s.Items[0], nil); err != errSentinel {
		t.Fatalf("err = %v, want sentinel from Reset", err)
	}
}

func TestCoverItem_MeasureError(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	s, _ := quest.Load(sessPath)
	probe := &fakeProbe{measureErr: errSentinel}
	if _, err := coverItem(&fakeDef{}, probe, s, s.Items[0], nil); err != errSentinel {
		t.Fatalf("err = %v, want sentinel from Measure", err)
	}
}

func TestCoverItem_PrepareError(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	s, _ := quest.Load(sessPath)
	probe := &fakeProbe{cov: covWith(1, 1, 1)}
	def := &fakeDef{prepErr: errSentinel}
	if _, err := coverItem(def, probe, s, s.Items[0], nil); err != errSentinel {
		t.Fatalf("err = %v, want sentinel from Prepare", err)
	}
}

// ---------------------------------------------------------------------------
// commitVerdict
// ---------------------------------------------------------------------------

func TestCommitVerdict_TerminalLocksAndExports(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	s, _ := quest.Load(sessPath)
	it := s.Items[0]
	outPath := filepath.Join(dir, "out.jsonl")
	sink, err := newJSONLSink(outPath)
	if err != nil {
		t.Fatal(err)
	}

	if err := commitVerdict(s, it, quest.Verdict{Outcome: quest.OutPass}, sink, sessPath); err != nil {
		t.Fatalf("commitVerdict: %v", err)
	}
	if it.State != quest.PASS {
		t.Errorf("state = %s, want PASS", it.State)
	}
	if !it.Emitted {
		t.Error("Emitted not set after export")
	}
	if n := countLines(t, outPath); n != 1 {
		t.Errorf("JSONL lines = %d, want 1", n)
	}
	// Persisted to disk.
	reloaded, _ := quest.Load(sessPath)
	if reloaded.Items[0].State != quest.PASS || !reloaded.Items[0].Emitted {
		t.Error("session not saved with locked+emitted item")
	}
}

func TestCommitVerdict_NonTerminalNotExported(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	s, _ := quest.Load(sessPath)
	it := s.Items[0]
	outPath := filepath.Join(dir, "out.jsonl")
	sink, _ := newJSONLSink(outPath)

	// FAIL once (tries 1 < MaxTries) → not terminal → not exported.
	if err := commitVerdict(s, it, quest.Verdict{Outcome: quest.OutFail}, sink, sessPath); err != nil {
		t.Fatalf("commitVerdict: %v", err)
	}
	if it.State != quest.TODO {
		t.Errorf("state = %s, want still TODO", it.State)
	}
	if n := countLines(t, outPath); n != 0 {
		t.Errorf("JSONL lines = %d, want 0 (non-terminal)", n)
	}
}

// failingSink is a quest.Sink whose Emit always errors, to drive commitVerdict's
// export-error branch (Export fails after Save succeeds).
type failingSink struct{}

func (failingSink) Emit(it *quest.Item) error { return errSentinel }

func TestCommitVerdict_ExportError(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	s, _ := quest.Load(sessPath)
	// Terminal verdict → Export tries to Emit the item → failingSink errors,
	// while Save (valid sessPath) succeeds, so commitVerdict returns the exportErr.
	err := commitVerdict(s, s.Items[0], quest.Verdict{Outcome: quest.OutPass}, failingSink{}, sessPath)
	if err != errSentinel {
		t.Fatalf("err = %v, want export sentinel", err)
	}
	// The pre-export Apply+Save still persisted the locked state.
	reloaded, _ := quest.Load(sessPath)
	if reloaded.Items[0].State != quest.PASS {
		t.Errorf("state = %s, want PASS persisted before export", reloaded.Items[0].State)
	}
}

func TestCommitVerdict_SaveError(t *testing.T) {
	dir := t.TempDir()
	sessPath := seedCoverSession(t, dir, coverEP("a"))
	s, _ := quest.Load(sessPath)
	sink, _ := newJSONLSink(filepath.Join(dir, "out.jsonl"))
	// A session path whose parent does not exist → Save fails.
	badPath := filepath.Join(dir, "nope", "deep", "session.json")
	err := commitVerdict(s, s.Items[0], quest.Verdict{Outcome: quest.OutPass}, sink, badPath)
	if err == nil {
		t.Fatal("want Save error for unwritable path")
	}
}

// ---------------------------------------------------------------------------
// evaluate (dispatch)
// ---------------------------------------------------------------------------

func TestEvaluate_DispatchesToEvaluator(t *testing.T) {
	def := &fakeDef{verdict: quest.Verdict{Outcome: quest.OutReview, RootCause: "X"}}
	v := evaluate(def, gate.Context{})
	if v.Outcome != quest.OutReview || v.RootCause != "X" {
		t.Errorf("verdict = %+v, want the Evaluator's", v)
	}
	if def.evalCalls != 1 {
		t.Errorf("Evaluate calls = %d, want 1", def.evalCalls)
	}
}

func TestEvaluate_FallsBackToRules(t *testing.T) {
	// A def WITHOUT gate.Evaluator → flat gate.Evaluate over Rules().
	fired := false
	def := nonEvalDef{rules: []gate.Rule{{
		Meta: gate.RuleMeta{ID: "F-1", Level: gate.LevelFail},
		Check: func(ctx gate.Context) (bool, quest.Fact) {
			fired = true
			return true, quest.Fact{Rule: "F-1", Where: "x"}
		},
	}}}
	v := evaluate(def, gate.Context{})
	if !fired {
		t.Error("rule Check never ran → fallback path not taken")
	}
	if v.Outcome != quest.OutFail {
		t.Errorf("verdict = %s, want FAIL from fired rule", v.Outcome)
	}
}

func TestEvaluate_FallbackEmptyRulesPasses(t *testing.T) {
	v := evaluate(nonEvalDef{}, gate.Context{})
	if v.Outcome != quest.OutPass {
		t.Errorf("verdict = %s, want PASS with no rules", v.Outcome)
	}
}

// ---------------------------------------------------------------------------
// renderCoverVerdict
// ---------------------------------------------------------------------------

func TestRenderCoverVerdict_Feedback(t *testing.T) {
	var buf bytes.Buffer
	it := &quest.Item{State: quest.DONE}
	v := quest.Verdict{Outcome: quest.OutFail, Feedback: "line one\nline two\n"}
	renderCoverVerdict(&buf, "GET /x", it, v)
	out := buf.String()
	if !strings.HasPrefix(out, "GET /x -> FAIL (state DONE)\n") {
		t.Errorf("header wrong: %q", out)
	}
	if !strings.Contains(out, "  line one\n") || !strings.Contains(out, "  line two\n") {
		t.Errorf("feedback not indented: %q", out)
	}
}

func TestRenderCoverVerdict_Facts(t *testing.T) {
	var buf bytes.Buffer
	it := &quest.Item{State: quest.REVIEW}
	v := quest.Verdict{
		Outcome: quest.OutReview,
		Facts:   []quest.Fact{{Rule: "C-01", Where: "body", Expected: "a", Actual: "b"}},
	}
	renderCoverVerdict(&buf, "GET /y", it, v)
	out := buf.String()
	if !strings.Contains(out, "GET /y -> REVIEW (state REVIEW)") {
		t.Errorf("header missing: %q", out)
	}
	if !strings.Contains(out, `- C-01: body expected="a" actual="b"`) {
		t.Errorf("fact line missing: %q", out)
	}
}

// ---------------------------------------------------------------------------
// NewLoopCmd
// ---------------------------------------------------------------------------

func TestNewLoopCmd_Shape(t *testing.T) {
	cmd := NewLoopCmd(Def())
	if cmd.Use != "loop" {
		t.Errorf("Use = %q, want loop", cmd.Use)
	}
	if len(cmd.Aliases) != 0 {
		t.Errorf("Aliases = %v, want none (clean cut: cover removed)", cmd.Aliases)
	}
	if cmd.Short == "" {
		t.Error("Short is empty")
	}
	if cmd.RunE == nil {
		t.Fatal("RunE is nil")
	}
	if cmd.Flags().Lookup("measure-only") == nil {
		t.Error("--measure-only flag is not registered")
	}
	modelFlag := cmd.Flags().Lookup("model")
	if modelFlag == nil {
		t.Fatal("--model flag is not registered")
	}
	if modelFlag.DefValue != "ollama:gemma4:e4b" {
		t.Errorf("--model default = %q, want ollama:gemma4:e4b", modelFlag.DefValue)
	}
	if cmd.Flags().Lookup("max-items") == nil {
		t.Error("--max-items flag is not registered")
	}
	if cmd.Flags().Lookup("generate") != nil {
		t.Error("--generate flag should be removed (clean cut)")
	}
}

func TestNewLoopCmd_RunE_NoManifest(t *testing.T) {
	chdir(t, t.TempDir()) // no manifest.yaml here
	cmd := NewLoopCmd(Def())
	// Register the persistent flags the RunE reads (normally inherited from root).
	cmd.Flags().String("session", "session.json", "")
	cmd.Flags().String("out", "out.jsonl", "")

	err := cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatal("want error when no manifest is present")
	}
	if !strings.Contains(err.Error(), "manifest") {
		t.Errorf("error = %v, want it to mention manifest", err)
	}
}

func TestNewLoopCmd_RunE_MissingSessionFlag(t *testing.T) {
	chdir(t, t.TempDir())
	cmd := NewLoopCmd(Def())
	// No "session" flag registered → GetString("session") errors.
	if err := cmd.RunE(cmd, nil); err == nil {
		t.Fatal("want error when the session flag is undefined")
	}
}

func TestNewLoopCmd_RunE_MissingOutFlag(t *testing.T) {
	chdir(t, t.TempDir())
	cmd := NewLoopCmd(Def())
	cmd.Flags().String("session", "session.json", "")
	// "out" flag absent → GetString("out") errors.
	if err := cmd.RunE(cmd, nil); err == nil {
		t.Fatal("want error when the out flag is undefined")
	}
}

func TestNewLoopCmd_RunE_BackendBuildError(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	// A valid manifest gets past config.Load, so RunE reaches buildBackend.
	if err := os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte(coverManifest), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := NewLoopCmd(Def())
	cmd.Flags().String("session", filepath.Join(dir, "session.json"), "")
	cmd.Flags().String("out", filepath.Join(dir, "out.jsonl"), "")
	// Generation is ON by default; an unparseable --model (no backend:model colon)
	// makes buildBackend → llm.FromFlag fail, exercising RunE's backend-error branch.
	if err := cmd.Flags().Set("model", "bogus-no-colon"); err != nil {
		t.Fatal(err)
	}

	err := cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatal("want error when --model is unparseable for generate mode")
	}
	if !strings.Contains(err.Error(), "--model") {
		t.Errorf("error = %v, want it to mention the invalid --model", err)
	}
}

// coverManifest is a minimal yongol manifest whose build/start commands are the
// no-op `true` binary, so newLiveProbe(cfg).Up() compiles instantly and no real
// server is ever spawned (the seeded session is empty → no Measure).
const coverManifest = `apiVersion: yongol/v1
kind: Project
metadata:
  name: cover-test
backend:
  lang: go
testing:
  base_url: http://localhost:8080
  hurl_dir: hurl
  server:
    build: "true"
    start: "true"
    ready: /health
`

func TestNewLoopCmd_RunE_RunsCoverOnEmptySession(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	if err := os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte(coverManifest), 0o644); err != nil {
		t.Fatal(err)
	}
	// Empty seeded session → runCover loops zero times (no server, no Measure).
	sessPath := seedCoverSession(t, dir)

	cmd := NewLoopCmd(Def())
	cmd.Flags().String("session", sessPath, "")
	cmd.Flags().String("out", filepath.Join(dir, "out.jsonl"), "")

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("RunE over empty session: %v", err)
	}
}
