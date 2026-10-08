# Catalog Decision Engine Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Run a deterministic authentication and catalog setup, then repeatedly score catalog actions, update mental state, and either select a movie or end the session.

**Architecture:** `Session` owns the sequential workflow and mutable session state. A pure `decision` package prepares candidates, calculates appeal and action scores, chooses the highest score, and applies bounded mental-state transitions; platform capabilities remain responsible only for HTTP calls. Each concurrent session receives its own RNG so candidate order, attraction, and tie-breaking are reproducible without sharing mutable random state.

**Tech Stack:** Go 1.24, standard library `math/rand/v2`, `net/http`, `net/http/httptest`, and `testing`.

**Spec:** `docs/decision_engine.md`

## Global Constraints

- Authentication is deterministic: new users register then log in; returning users only log in.
- Catalog decisions start only after login and one successful `ListMovies` call.
- The available decisions are exactly `SELECT_MOVIE`, `CONTINUE_BROWSING`, and `LEAVE_SESSION`.
- The engine chooses the highest score; randomness is used only for candidate order, attraction, and score ties.
- Personality is immutable; mental-state values are always clamped to `[0,1]`.
- A session evaluates an eligible movie no more than once and terminates at selection, leave, exhaustion, the browse limit, context cancellation/deadline, or error.
- Decision code performs no HTTP calls; capability code contains no behavioral decisions.
- Every session owns its own `*rand.Rand`; RNG instances are not shared between session goroutines.
- No real media is streamed in this milestone.

## Review Focus

- Empty or entirely invalid catalogs must log out normally without division by zero or a fabricated selection; Task 5 adds the integration test.
- `MaxBrowseCount <= 0` must never produce `NaN` fatigue or an infinite loop; Tasks 4 and 5 add unit and integration tests.
- An authenticated catalog failure must still execute logout exactly once; Task 5 adds the ordering/error test.
- Equal or nearly equal action scores must use seeded tie resolution instead of declaration order; Task 4 adds the reproducibility test.
- Repeated action transitions at boundary inputs (`0`, `1`, and out-of-range legacy values) must return finite values inside `[0,1]`; Task 2 adds table and repetition tests.

---

## File Structure

- Create `simulator/internal/decision/action.go`: catalog action and decision result types.
- Create `simulator/internal/decision/mental_state.go`: bounded state transitions and load-failure transition.
- Create `simulator/internal/decision/catalog.go`: engine, contexts, candidate preparation, fatigue, appeal, scoring, and tie resolution.
- Create `simulator/internal/decision/mental_state_test.go`: transition and invariant tests.
- Create `simulator/internal/decision/catalog_test.go`: candidate, appeal, score, and action-selection tests.
- Create `simulator/internal/platform/httpclient/httpclient_test.go`: request, authentication-header, decode, and HTTP-error tests.
- Modify `simulator/internal/platform/httpclient/httpclient.go`: implement GET and authenticated request behavior.
- Modify `simulator/internal/platform/session/session.go`: initialize dependencies and run the catalog loop.
- Create `simulator/internal/platform/session/session_test.go`: sequential flow and termination tests with fakes.
- Modify `simulator/internal/orchastrator/orchastrator.go`: create one RNG per session before launching its goroutine.
- Create `simulator/internal/orchastrator/orchastrator_test.go`: per-session RNG reproducibility and independence test.

### Task 1: Functional Authenticated HTTP Client

**Files:**
- Modify: `simulator/internal/platform/httpclient/httpclient.go`
- Test: `simulator/internal/platform/httpclient/httpclient_test.go`

**Interfaces:**
- Consumes: existing `Client{BaseURL, HTTPClient, AccessToken, TokenType}`.
- Produces: `func (c *Client) Get(path string, out any) error` and existing `Post` behavior with an authorization header whenever `AccessToken` is non-empty.

- [ ] **Step 1: Write failing GET request tests**

Add `TestGetDecodesJSONAndSendsAuthorization` using `httptest.Server`. Assert method `GET`, path `/movies`, header `Authorization: bearer test-token`, and decoded JSON. Add `TestGetReturnsErrorForHTTPFailure` and assert a `404` response returns a non-nil error containing `404 Not Found`.

- [ ] **Step 2: Run the GET tests and verify RED**

Run: `cd simulator && go test ./internal/platform/httpclient -run 'TestGet' -v`

Expected: FAIL because `Get` returns without making a request or decoding the response.

- [ ] **Step 3: Implement GET and shared request helpers**

In `httpclient.go`, implement:

```go
func (c *Client) Get(path string, out any) error
func (c *Client) addAuthorization(req *http.Request)
func decodeResponse(resp *http.Response, out any) error
```

`addAuthorization` must omit the header for an empty access token and otherwise use `TokenType + " " + AccessToken`. `decodeResponse` must reject status codes `>= 400`, return without decoding when `out == nil`, and decode JSON otherwise.

- [ ] **Step 4: Run GET tests and verify GREEN**

Run: `cd simulator && go test ./internal/platform/httpclient -run 'TestGet' -v`

Expected: PASS.

- [ ] **Step 5: Write the failing POST authorization regression test**

Add `TestPostSendsAuthorizationWhenLoggedIn`. Assert the server receives `Authorization: bearer test-token` and the JSON request/response still round-trips.

- [ ] **Step 6: Run the POST test and verify RED**

Run: `cd simulator && go test ./internal/platform/httpclient -run 'TestPostSendsAuthorization' -v`

Expected: FAIL because `Post` currently sets only `Content-Type`.

- [ ] **Step 7: Reuse the helpers from POST**

Call `addAuthorization` before sending the POST request and replace duplicated response handling with `decodeResponse`.

- [ ] **Step 8: Run the package tests**

Run: `cd simulator && go test ./internal/platform/httpclient -v`

Expected: PASS.

- [ ] **Step 9: Commit**

```bash
git add simulator/internal/platform/httpclient/httpclient.go simulator/internal/platform/httpclient/httpclient_test.go
git commit -m "feat: complete simulator HTTP client"
```

### Task 2: Catalog Actions and Bounded Mental-State Transitions

**Files:**
- Create: `simulator/internal/decision/action.go`
- Create: `simulator/internal/decision/mental_state.go`
- Test: `simulator/internal/decision/mental_state_test.go`

**Interfaces:**
- Consumes: `user.MentalState`.
- Produces: `Action`, action constants, `Decision`, `ApplyMentalState(*user.MentalState, Action)`, and `ApplyMovieLoadFailure(*user.MentalState)`.

- [ ] **Step 1: Write failing action-transition tests**

Add tests named `TestApplyMentalStateSelectMovie`, `TestApplyMentalStateContinueBrowsing`, and `TestApplyMentalStateLeaveSession`. Assert the exact increase/decrease directions and rates from `docs/decision_engine.md` using the diminishing formulas, rather than only asserting that values changed.

- [ ] **Step 2: Write the failing invariant tests**

Add `TestApplyMentalStateAlwaysClampsFields` with boundary and legacy out-of-range inputs, and `TestRepeatedMentalStateTransitionsRemainBounded` that applies each action 1,000 times and asserts every field is finite and in `[0,1]`.

- [ ] **Step 3: Run transition tests and verify RED**

Run: `cd simulator && go test ./internal/decision -run 'TestApply|TestRepeated' -v`

Expected: FAIL because the package and transition functions do not exist.

- [ ] **Step 4: Define action and result types**

In `action.go`, define:

```go
type Action string

const (
	ActionSelectMovie      Action = "SELECT_MOVIE"
	ActionContinueBrowsing Action = "CONTINUE_BROWSING"
	ActionLeaveSession      Action = "LEAVE_SESSION"
)

type Decision struct {
	Action        Action
	SelectScore   float64
	ContinueScore float64
	LeaveScore    float64
}
```

- [ ] **Step 5: Implement bounded transitions**

In `mental_state.go`, implement unexported `clamp01`, `increase`, and `decrease` helpers plus:

```go
func ApplyMentalState(state *user.MentalState, action Action)
func ApplyMovieLoadFailure(state *user.MentalState)
```

Use the exact rates in the spec. `ApplyMovieLoadFailure` increases frustration by `0.10` and decreases satisfaction by `0.10`; clamp every field before returning, including unchanged fields inherited from legacy data.

- [ ] **Step 6: Run transition tests and verify GREEN**

Run: `cd simulator && go test ./internal/decision -run 'TestApply|TestRepeated' -v`

Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add simulator/internal/decision/action.go simulator/internal/decision/mental_state.go simulator/internal/decision/mental_state_test.go
git commit -m "feat: add catalog actions and mental state transitions"
```

### Task 3: Candidate Preparation and Appeal

**Files:**
- Create: `simulator/internal/decision/catalog.go`
- Test: `simulator/internal/decision/catalog_test.go`

**Interfaces:**
- Consumes: Task 2 `Action`/`Decision`, `catalog.MovieResponse`, `user.Personality`, and `user.MentalState`.
- Produces: `Engine`, `New(*rand.Rand)`, `PrepareCandidates`, `CalculateAppeal`, `AppealContext`, and `MaxDurationSeconds`.

- [ ] **Step 1: Write failing candidate-preparation tests**

Add `TestPrepareCandidatesFiltersInvalidAndDuplicateMovies`, asserting only unique movies with positive IDs, `READY` status, and positive duration remain. Add `TestPrepareCandidatesUsesReproducibleShuffle`, asserting two engines with identical PCG seeds return the same order and that the output contains each eligible candidate exactly once.

- [ ] **Step 2: Write failing appeal tests**

Add `TestCalculateAppealUsesExactWeights`, `TestCalculateAppealClampsResult`, `TestCalculateAppealUsesExplorationForUnseenMovie`, and `TestCalculateAppealUsesConsistencyForWatchedMovie`. Use a fixed RNG and assert the formula `0.50*durationFit + 0.30*noveltyFit + 0.20*attraction`. Add `TestNewPanicsForNilRNG` and assert the panic text is `decision: nil rng`.

- [ ] **Step 3: Run preparation and appeal tests and verify RED**

Run: `cd simulator && go test ./internal/decision -run 'TestPrepareCandidates|TestCalculateAppeal' -v`

Expected: FAIL because the engine and appeal functions do not exist.

- [ ] **Step 4: Define the engine and appeal context**

In `catalog.go`, define:

```go
type Engine struct {
	rng *rand.Rand
}

type AppealContext struct {
	Personality        user.Personality
	MentalState        user.MentalState
	Candidate          catalog.MovieResponse
	MaxDurationSeconds int
	PreviouslyWatched  bool
}

func New(rng *rand.Rand) *Engine
func (e *Engine) PrepareCandidates(movies []catalog.MovieResponse) []catalog.MovieResponse
func MaxDurationSeconds(candidates []catalog.MovieResponse) int
func (e *Engine) CalculateAppeal(context AppealContext) float64
```

`New` must panic with `decision: nil rng` when passed a nil RNG. `PrepareCandidates` returns a new slice and must not mutate the API response slice.

- [ ] **Step 5: Implement candidate filtering, shuffling, and appeal**

Use `e.rng.Shuffle` for candidates and `e.rng.Float64()` for attraction. When `MaxDurationSeconds <= 0`, use neutral duration fit `0.5` rather than dividing by zero. Use novelty `Consistency` for watched candidates and `Exploration` for unseen candidates; callers pass `PreviouslyWatched=false` until history integration is enabled.

- [ ] **Step 6: Run preparation and appeal tests and verify GREEN**

Run: `cd simulator && go test ./internal/decision -run 'TestPrepareCandidates|TestCalculateAppeal' -v`

Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add simulator/internal/decision/catalog.go simulator/internal/decision/catalog_test.go
git commit -m "feat: prepare and evaluate catalog candidates"
```

### Task 4: Catalog Scoring and Action Selection

**Files:**
- Modify: `simulator/internal/decision/catalog.go`
- Modify: `simulator/internal/decision/catalog_test.go`

**Interfaces:**
- Consumes: Task 3 `Engine` and a caller-calculated candidate appeal.
- Produces: `CatalogContext`, `BrowsingFatigue`, and `(*Engine).Decide(CatalogContext) Decision`.

- [ ] **Step 1: Write failing exact-score tests**

Add table tests for select, continue, and leave using the exact coefficients in `docs/decision_engine.md`. Assert scores to a floating-point tolerance of `1e-12`.

- [ ] **Step 2: Write failing behavioral tests**

Add tests proving strong appeal/impulsiveness favors select, curiosity/exploration favors continue, frustration/low energy favors leave, and increasing fatigue lowers continue while raising leave.

- [ ] **Step 3: Write failing termination-safety and tie tests**

Add `TestBrowsingFatigueHandlesNonPositiveMaximum` and assert it returns `1` without `NaN`. Add `TestDecideBreaksTiesReproducibly` using two engines with identical seeds and exact/epsilon ties, and assert both produce the same action sequence rather than always choosing the first declared action.

- [ ] **Step 4: Run scoring tests and verify RED**

Run: `cd simulator && go test ./internal/decision -run 'Test.*Score|TestBrowsingFatigue|TestDecide' -v`

Expected: FAIL because scoring and selection are not implemented.

- [ ] **Step 5: Define context and scoring interface**

Add:

```go
type CatalogContext struct {
	Personality     user.Personality
	MentalState     user.MentalState
	CandidateAppeal float64
	BrowsedCount    int
	MaxBrowseCount  int
}

func BrowsingFatigue(browsedCount, maxBrowseCount int) float64
func (e *Engine) Decide(context CatalogContext) Decision
```

- [ ] **Step 6: Implement the three formulas and highest-score choice**

Copy all coefficients exactly from the spec. Resolve scores within `scoreEpsilon = 1e-9` by collecting every tied action and selecting one with `e.rng.IntN`. Do not use a proportional random roll.

- [ ] **Step 7: Run all decision tests**

Run: `cd simulator && go test ./internal/decision -v`

Expected: PASS.

- [ ] **Step 8: Commit**

```bash
git add simulator/internal/decision/catalog.go simulator/internal/decision/catalog_test.go
git commit -m "feat: score and select catalog actions"
```

### Task 5: Sequential Session Workflow and Catalog Loop

**Files:**
- Modify: `simulator/internal/platform/session/session.go`
- Create: `simulator/internal/platform/session/session_test.go`

**Interfaces:**
- Consumes: Tasks 1-4 HTTP client and decision interfaces, plus existing authentication/catalog request and response types.
- Produces: `NewSession(user *user.User, registerFirst bool, rng *rand.Rand) *Session`, `(*Session).Run(context.Context) error`, an unexported dependency-injected constructor for tests, and `(*Session).runCatalogLoop(context.Context, []catalog.MovieResponse) (*catalog.MovieDetailsResponse, error)`.

- [ ] **Step 1: Write failing sequential-order tests**

Create small fake authentication and catalog capabilities that append operation names to a shared slice. Add `TestRunRegistersThenLogsInThenListsCatalog`, asserting a new user produces `register, login, list, logout`; add the returning-user case asserting `login, list, logout`.

- [ ] **Step 2: Write failing loop-action tests**

Add tests for: select calls `GetMovie` once; continue advances to a different candidate; leave calls no detail endpoint; every candidate is evaluated at most once; exhausted candidates return without selection; the configured browse limit stops iteration; and a cancelled/deadline-exceeded context terminates the loop and still logs out.

- [ ] **Step 3: Write failing error and empty-input tests**

Assert authentication failure prevents listing, list failure after login still logs out exactly once, an empty/invalid catalog logs out normally, and detail failure applies the failure transition before advancing to the next candidate.

- [ ] **Step 4: Run session tests and verify RED**

Run: `cd simulator && go test ./internal/platform/session -v`

Expected: FAIL because `Session` does not initialize dependencies, accept an RNG, or run the decision loop.

- [ ] **Step 5: Introduce narrow dependency interfaces**

In `session.go`, define unexported interfaces matching the existing concrete methods:

```go
type authenticator interface {
	Register(*httpclient.Client, authentication.RegisterRequest) (*authentication.RegisterResponse, error)
	Login(*httpclient.Client, authentication.LoginRequest) (*authentication.LoginResponse, error)
	Logout(*httpclient.Client) error
}

type catalogBrowser interface {
	ListMovies(*httpclient.Client) ([]catalog.MovieResponse, error)
	GetMovie(*httpclient.Client, int) (*catalog.MovieDetailsResponse, error)
}

type catalogDecider interface {
	PrepareCandidates([]catalog.MovieResponse) []catalog.MovieResponse
	CalculateAppeal(decision.AppealContext) float64
	Decide(decision.CatalogContext) decision.Decision
}
```

Add an unexported test constructor that accepts the client, these interfaces, and a configurable maximum browse count. The public constructor creates the HTTP client, real authentication/catalog capabilities, `decision.New(rng)`, and uses the default browse count. The fake decider in tests returns a prescribed action sequence so loop tests do not depend on score calibration.

- [ ] **Step 6: Change the public session constructor**

Implement:

```go
func NewSession(
	u *user.User,
	registerFirst bool,
	rng *rand.Rand,
) *Session
```

Store `decision.New(rng)` in the session. Use `const defaultMaxBrowseCount = 20`.

Expose cancellation through:

```go
func (s *Session) Run(ctx context.Context) error
```

Store the configured browse count on `Session`; the dependency-injected test constructor may set it to `0` or another small value.

- [ ] **Step 7: Implement deterministic setup and reliable logout**

In `Run`, check `ctx.Err()`, register when required, log in, install the logout `defer` immediately after successful login, call `ListMovies` exactly once, then call `runCatalogLoop(ctx, movies)`. Remove the fixed five-second sleep. If a movie is selected, retain the local result as the handoff point for the future playback loop and log only its ID/title. Cancellation or deadline expiry returns `ctx.Err()` after deferred logout.

- [ ] **Step 8: Implement the catalog loop**

Prepare candidates once and calculate maximum duration once. Before each iteration, return `ctx.Err()` when cancelled. For each candidate up to `min(len(candidates), s.maxBrowseCount)`, calculate appeal with `PreviouslyWatched=false`, call `Decide`, log state/scores/action, and then:

- `SELECT_MOVIE`: apply its state transition, call `GetMovie`, return it on success; on failure apply `ApplyMovieLoadFailure` and advance.
- `CONTINUE_BROWSING`: apply its transition and advance.
- `LEAVE_SESSION`: apply its transition and return no movie.

When the catalog is empty, invalid, exhausted, or capped, apply the leave transition once and return no movie. Do not place authentication or `ListMovies` inside the loop.

- [ ] **Step 9: Run session tests and verify GREEN**

Run: `cd simulator && go test ./internal/platform/session -v`

Expected: PASS.

- [ ] **Step 10: Commit**

```bash
git add simulator/internal/platform/session/session.go simulator/internal/platform/session/session_test.go
git commit -m "feat: run catalog decisions in user sessions"
```

### Task 6: Per-Session RNG Provisioning

**Files:**
- Modify: `simulator/internal/orchastrator/orchastrator.go`
- Create: `simulator/internal/orchastrator/orchastrator_test.go`

**Interfaces:**
- Consumes: Task 5 `session.NewSession(..., *rand.Rand)`.
- Produces: an orchestrator-owned session seed generator and `(*Orchastrator).nextSessionRNG() *rand.Rand`.

- [ ] **Step 1: Write the failing RNG test**

Add `TestNextSessionRNGIsReproducibleAndIndependent`. Construct two orchestrators with equal seed-generator state, assert corresponding session RNGs produce equal sequences, and assert the first and second session RNGs from one orchestrator do not produce the same sequence.

- [ ] **Step 2: Run the RNG test and verify RED**

Run: `cd simulator && go test ./internal/orchastrator -run TestNextSessionRNG -v`

Expected: FAIL because the seed generator/helper does not exist.

- [ ] **Step 3: Implement per-session RNG creation**

Add `sessionSeedGenerator *rand.Rand` to `Orchastrator`, initialize it with a PCG source separate from the population and mental-state generators, and implement `nextSessionRNG` by drawing two `uint64` seeds and creating a new PCG-backed RNG.

- [ ] **Step 4: Pass the RNG before starting each goroutine**

Call `nextSessionRNG` in the orchestrator loop before `go func`, capture it as a goroutine argument, and pass it to `session.NewSession`. Pass `context.Background()` to `Session.Run`; a future orchestrator shutdown can replace it with a cancellable parent context. Never call the shared seed generator from inside concurrent session goroutines.

- [ ] **Step 5: Run orchestrator and full simulator tests**

Run: `cd simulator && go test ./internal/orchastrator -v && go test ./...`

Expected: PASS with no race-prone shared session RNG.

- [ ] **Step 6: Run formatting and static checks**

Run: `cd simulator && gofmt -w internal/decision internal/platform/httpclient internal/platform/session internal/orchastrator && go vet ./... && go test ./...`

Expected: all commands exit `0`.

- [ ] **Step 7: Commit**

```bash
git add simulator/internal/orchastrator/orchastrator.go simulator/internal/orchastrator/orchastrator_test.go
git commit -m "feat: give each simulator session an independent rng"
```

## Final Verification

- [ ] Run: `cd simulator && gofmt -w internal/decision internal/platform/httpclient internal/platform/session internal/orchastrator`
- [ ] Run: `cd simulator && go vet ./...`
- [ ] Run: `cd simulator && go test ./...`
- [ ] Confirm logs contain no password or access token.
- [ ] Confirm `ListMovies` appears once per successfully authenticated session.
- [ ] Confirm every catalog decision log includes candidate ID, appeal, three scores, selected action, and mental state before/after.
- [ ] Confirm no production code was added to playback as part of this plan.
