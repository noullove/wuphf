# Autonomy Coaching Internalization

Date: 2026-04-15
Status: in progress
Scope: turn the strongest eval-time coaching into durable system behavior and tighten the eval harness so both blueprint-backed and blank-slate runs are judged on real business usefulness rather than synthetic proof artifacts.

## Why

The recent consulting and from-scratch evals proved something important and uncomfortable at the same time:
- runtime plumbing is much better than before
- real external writes are possible
- the system still needs too much human coaching at the exact moment where a business should stop planning and start operating

The main failure mode is no longer "cannot call Slack or Notion." The main failure mode is:
- over-planning
- evidence-driven behavior instead of value-driven behavior
- synthetic `proof` / `marker` / `test` outputs that satisfy durability checks without looking like a real business workflow

That means the next step is not another loose eval pass. The next step is to internalize the coaching into core behavior and into the eval harness itself.

## Issue Inventory

### 1. Behavior quality is still too coach-dependent

Observed symptoms:
- the office writes kickoff briefs, proof packets, and review shells before taking the smallest allowed real action
- the office explains blockers better than it resolves them
- blank-slate runs drift into safe internal artifacts instead of strong business choices
- the system still needs a human to say "stop proving, start operating"

Required product response:
- prefer the smallest useful live action over extra planning once the action contract is known
- reward business-semantic outputs, not eval-semantic outputs
- recover by pivoting to the next safe live action instead of collapsing into documentation

### 2. External-action tasks can still degrade into synthetic outputs

Observed symptoms:
- tasks requiring Slack / Notion / Drive can still be satisfied by repo markdown if the task wording is loose enough
- "proof artifact" style language invites manufactured outputs
- the eval itself allowed real connector writes that still looked like demo artifacts

Required product response:
- external-action tasks must declare that repo docs are not a valid substitute
- completion must require broker-side evidence of the external action when the task requires one
- prompts should ban `proof`, `marker`, `test`, and similar framing for live business work unless the task is explicitly a test

### 3. Recovery after live failure is still weak

Observed symptoms:
- after the first Slack failure, the system needed coaching to keep moving
- the system often stalls on one failed target rather than retrying once and pivoting
- headless specialist retries were too brittle for office-mode external-action tasks

Required product response:
- retry the smallest live action once when the failure is plausibly recoverable
- then pivot to the next allowed live artifact or integration instead of blocking the whole loop
- keep the same task alive where possible so the run feels continuous

### 4. Headless runtime isolation is incomplete in the base branch

Observed symptoms:
- headless Codex inherited broken user-global `~/.agents` skills
- global `HOME` leakage made autonomous behavior depend on user-machine state
- isolated runtimes also need enough auth/config copied in to use connected systems predictably

Required product response:
- force headless runs onto an isolated runtime home
- copy only the minimal auth/config needed for the connected toolchain
- make that isolation the default, not an eval-only workaround

### 5. Blank-slate invention is real but still too guided

Observed symptoms:
- the from-scratch run completed a first additive loop, but it still required meaningful steering
- the blank-slate flow still tends to choose safe consulting/evidence loops
- without stronger constraints, the system optimizes for "show autonomous structure" rather than "produce a credible business move"

Required product response:
- keep the true blank-slate launch path
- strengthen directive synthesis and task framing so the first loop is commercially legible
- make invented businesses produce a real operating move, not an autonomy demo

### 6. The eval harness still rewards the wrong thing

Observed symptoms:
- "durable evidence" is over-weighted
- real writes can still look manufactured and pass
- the harness does not clearly fail a run that produces synthetic business semantics

Required product response:
- both blueprint and from-scratch evals must prohibit synthetic proof language
- success must require business-useful outputs
- the coach role should focus on quality and recovery, not hand-authoring the business semantics for the system

## Fix Tracks

We will attack the remaining gap in parallel across three tracks.

### Track A: Runtime and recovery

Goal:
- make the system resilient enough that live external work can continue after a normal failure without a full restart

Focus:
- headless home/config isolation
- office external-action retry behavior
- failure-aware pivoting for Slack / Notion / Drive style tasks
- durable evidence checks for live-action tasks

Likely files:
- `internal/team/headless_codex.go`
- `internal/team/headless_codex_test.go`
- `internal/team/task_pipeline.go`
- `internal/team/task_pipeline_test.go`

### Track B: Agent behavior and prompt contract

Goal:
- make "operate like a business" the default behavior instead of something the coach has to repeatedly insist on

Focus:
- live-action bias in lead/specialist prompts
- explicit ban on substituting repo docs for live tasks
- stronger business-semantic task framing
- value-first artifact expectations

Likely files:
- `internal/team/launcher.go`
- `internal/team/launcher_test.go`
- `internal/team/broker.go`
- `internal/team/broker_test.go`

### Track C: Eval harness and blank-slate path

Goal:
- make the eval itself require real business workflows and keep the from-scratch path honest

Focus:
- true blank-slate launch support
- no defaulting into saved pack-era assumptions
- onboarding and Web UI phrasing that avoids synthetic-proof framing
- eval docs and replay criteria that fail manufactured outputs

Likely files:
- `cmd/wuphf/main.go`
- `web/index.html`
- `internal/team/broker_onboarding.go`
- `docs/specs/*`
- `docs/evals/*`

## Execution Plan

1. Integrate the strongest worker fixes into one autonomy branch.
2. Tighten core behavior so coached moves become default moves.
3. Tighten the eval harness so synthetic artifacts are explicit failures.
4. Run focused tests on the changed paths.
5. Restart two eval tracks:
   - blueprint-backed: `multi-agent-workflow-consulting`
   - blank-slate: `--from-scratch`
6. Judge both runs against business-semantic criteria, not just connector evidence.

## Success Criteria

This slice is complete only when all of the following are true:
- live external tasks no longer pass with repo-doc substitutes
- the office prefers a smallest allowed business action over extra planning once the action contract is known
- the system recovers from the first normal external failure without manual re-architecture
- blank-slate launch starts cleanly without hidden pack-era fallback
- the eval harness rejects synthetic `proof` / `marker` / `test` business outputs for live business runs
- both the consulting and from-scratch tracks can be rerun under the new rules and judged on real business usefulness

## Out Of Scope For This Slice

Not in scope for this specific internalization pass:
- proving a fully self-sustaining business with zero coaching forever
- broad provider comparison across all model backends
- cleaning every remaining historical pack-era compatibility shim

The target here is narrower and more practical:
- internalize the coaching that mattered
- tighten the eval
- then rerun the two high-value autonomy tracks under a materially higher bar
