# Autonomy Proof Plan

Date: 2026-04-15
Status: active

## Objective

Prove two things in sequence:

1. Blueprint-backed autonomy
- `multi-agent-workflow-consulting` should run a real business-semantic loop with real external systems and correct approval behavior.

2. Blank-slate autonomy
- `--from-scratch` should stand up a business from a plain directive and execute the first real business-semantic loop without silently collapsing into an eval/demo workflow.

## Rules For Both Eval Tracks

These are now the hard rules for future runs:
- no `proof`, `marker`, `test`, or `demo` naming in business outputs unless the task itself is explicitly a test
- real external writes are allowed only within the user-approved safety envelope
- repo markdown is not a valid substitute for a live external-action task
- success is judged on business usefulness, not just connector evidence
- the coach may steer quality and unblock the system, but should not hand-author the substantive business outputs for it

## Track Order

### Track 1: Consulting blueprint

Purpose:
- establish one repeatable reference loop with real business semantics

Must prove:
- intake
- scoping
- task decomposition
- one real external action chain
- a value-bearing deliverable or client-facing/internal handoff artifact
- review and loop continuation without re-bootstrap

### Track 2: Blank-slate from-scratch run

Purpose:
- prove the system can invent and operate without template handholding

Must prove:
- true blank-slate bootstrap
- commercially legible business choice
- a first operating loop that looks like a real business action, not an autonomy demo
- the same external-action discipline as the consulting proof

## Parallel Fix Order

1. Runtime / recovery fixes
2. Prompt / behavior fixes
3. Eval harness and blank-slate fixes
4. Focused tests
5. Fresh reruns of consulting and from-scratch tracks

## Companion Specs

- `docs/specs/2026-04-15-autonomy-coaching-internalization.md`
- `docs/evals/2026-04-14-youtube-business-e2e.md`
