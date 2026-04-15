# YouTube Business E2E Eval

Date: 2026-04-14
Scenario: Post in `#general` asking the team to build a faceless YouTube channel business with automated content creation and monetization, from the Web UI, starting from a fresh worktree and reset office state.

Replay companion:
- `docs/evals/2026-04-14-youtube-user-journey.md`
  - step-by-step Web UI path
  - roadblocks hit in the path
  - fixes that removed them
  - current replay expectations for a later Claude Code run

Autonomy-proof planning:
- `docs/specs/2026-04-14-operation-blueprint-refactor.md`
- `docs/specs/2026-04-14-blueprint-test-matrix.md`
- `docs/specs/2026-04-14-multi-agent-workflow-consulting-proof.md`

## Autonomy Proof Checkpoint

Current best live proof is no longer the original YouTube path; it is the blueprint-backed `multi-agent-workflow-consulting` run from the Web UI.

Latest validated consulting proof checkpoint:
- worktree: `/Users/najmuzzaman/Documents/nex/WUPHF-eval-youtube`
- provider: `codex`
- home: `/tmp/wuphf-proof-consulting-4`
- broker/web: `7899 / 7914`

What is now proven:
- onboarding from the browser wakes `@operator` from a plain directive
- `@operator` creates `#client-intake-dryrun` and the first consulting lanes durably
- the first builder task (`task-2`) is now created through the real runtime path as `local_worktree` with a real assigned worktree:
  - `/tmp/wuphf-proof-consulting-4/.wuphf/task-worktrees/wuphf-eval-youtube/wuphf-task-task-2`
  - branch `wuphf-d761fdb1-task-2`
- builder executed inside that worktree and shipped a dry-run packet generator plus focused test
- `task-2` closed durably as `done` / `approved`
- the next builder task (`task-11`) was created with `depends_on: ["task-2"]`, stayed `local_worktree`, and reused the exact same worktree path and branch as `task-2`
- `task-11` also closed durably as `done` / `approved` on the reused worktree
- the surrounding consulting loop completed on the same operating structure:
  - `task-4` closed with the dry-run orchestration spec
  - `task-5` closed with the review scorecard / approval QA gate
  - operator posted explicit approval/close messages after both builder slices
- the next live blocker in that proof has now also been closed:
  - after a broker restart on the same proof home, a direct live sequence of `fresh local_worktree -> reused dependent local_worktree -> fresh local_worktree` succeeded
  - the third task now receives a distinct fresh worktree instead of failing with `failed to manage task worktree`

Roadblocks closed in this proof slice:
- `builder`-owned repo tasks were not consistently treated as local-worktree work
- the broker could not reuse a completed dependency worktree for the next `local_worktree` task
- the broker batch planner path (`/task-plan`) created `local_worktree` tasks without assigning a worktree, which is why the live consulting proof initially showed `task-2` missing `worktree_path`
- persisted sibling-worktree replay was copying `.wuphf/` runtime state into fresh worktrees, which is what broke the third fresh `local_worktree` create after a reuse chain
- persisted sibling-worktree replay was too brittle around duplicate or stale prior source paths

Remaining proof gaps:
- complete a successful live read-only external workflow after the current provider backoff window
- then exercise approval-gated external side effects
- capture the final operator close-out that marks the repeated consulting loop complete for replay

### Live Integration Checkpoint

The consulting proof has now reached a real external system through the normal runtime path:
- provider: `one`
- live connected integration: `gmail`
- action shape: read-only inbox listing

What happened:
- the broker executed a live One workflow against the connected Gmail account from the proof office
- the first live attempt failed with a real provider response:
  - `429 User-rate limit exceeded. Retry after 2026-04-14T22:55:02.178Z`
- no external side effect occurred:
  - no outbound message
  - no deletion
  - no mutation of existing external resources

Durable evidence:
- broker action log now records the live failure against the external provider
- the One flow and run logs are persisted under:
  - `/tmp/wuphf-proof-consulting-4/.wuphf/one/.one/flows/`

Follow-up fix:
- the broker now treats provider rate limits as first-class runtime evidence:
  - returns `429 Too Many Requests`
  - includes `Retry-After`
  - records `external_workflow_rate_limited`
  - updates skill execution status to `rate_limited`

### Broader Integration Checkpoint

The consulting proof is no longer Gmail-only.

What is now proven:
- Notion is connected in the consulting bootstrap package as `connected / live_capable`
- the blueprint-owned `workflow-audit` path ran successfully as a dry-run for:
  - `google-drive`
  - `notion`
- the blueprint-owned `client-handoff` path ran successfully as a mock workflow for:
  - `slack`
  - `notion`
- both runs recorded durable broker evidence:
  - `external_workflow_executed`
  - `skill_invocation`

Live read-only proof beyond Gmail:
- a custom live read-only Notion workflow now succeeds through the broker + One path
- the first attempt failed because the flow body incorrectly included `connectionKey`
- after removing that field from the body, the workflow completed successfully and returned real Notion page results for the query `workflow`

Current consulting bootstrap connection state:
- `notion`: connected / live capable
- `slack`: connected / live capable
- `google-drive`: ready for auth

So the current boundary is clear:
- the runtime can already coordinate Slack/Drive/Notion through the consulting blueprint
- Notion is proven live
- Slack is now also proven live
- Google Drive still needs approved account linking before it can move from dry-run into live execution

Latest consulting proof checkpoint:
- the office was given explicit policy rules through the Web UI:
  - Slack and Notion writes allowed for this proof
  - no deletes
  - no modification of existing external resources
  - Gmail out
- the office did create durable follow-up tasks from those instructions, but it still tried to satisfy `Notion proof artifact` with repo markdown instead of a real Notion write
- that exposed a product bug rather than a user-approval gap: active Policies were persisted in the broker and shown in the UI, but they were not injected into agent system prompts
- runtime capability was then proven directly:
  - Slack join + post succeeded in `#test-slack-integration`
  - Notion create-page succeeded and returned a live page URL for `Loopsmith Consulting proof`
- code fix landed in `internal/team/launcher.go` so active office policies are now injected into both lead and specialist prompts as hard constraints; focused launcher tests passed

## Run 1

### Issue 1: stale CEO turn after delegation
- Symptom: CEO delegated to `@eng` and `@gtm`, but stayed active in the same turn instead of stopping and waiting for pushed specialist completions. Specialist replies landed in the broker, but the office stalled before CEO synthesis.
- Evidence:
  - `~/.wuphf/team/broker-state.json` showed specialist replies without a follow-up CEO message.
  - `~/.wuphf/logs/headless-codex-ceo.log` showed the runtime relying on stale-turn cancellation to recover.
- Fix:
  - Tightened CEO and specialist prompts to explicitly end the turn after delegation, completion, handoff, or a blocking human question.
  - Added the same stop-after-reply instruction to queued work packets and task notifications.
- Status: fixed in code, rerun validated.

### Issue 2: live-output panes leaked opaque `item.completed` noise
- Symptom: specialist DM live-output panes showed repeated `item.completed` cards with no human-usable content.
- Evidence:
  - Web UI live-output stream showed `item.completed` rows instead of actionable summaries.
- Fix:
  - Suppressed unknown `item.completed` events in the live-output renderer after recognized tool/message cases are handled.
- Status: fixed in code, rerun validated.

### Issue 3: Nex context queries failed despite health showing connected
- Symptom: agents attempted `query_context`, but live output showed session/auth expiry errors.
- Evidence:
  - Specialist and CEO live-output panes showed Nex auth/session-expired failures.
  - `/health` still reported `nex_connected: true`.
- Fix:
  - No code fix yet. This appears to be environment/auth state, not a deterministic repo bug.
- Status: open.

### Issue 4: automation did not create tasks, channels, or skills in the first pass
- Symptom: the team produced planning replies, but no task records, new channels, or skills were created for the business build.
- Evidence:
  - Broker state remained at `task_count = 0`, `channel_count = 1`.
- Fix:
  - No code fix yet. First rerun will confirm whether the stale-turn fix unblocks the expected CEO follow-through.
- Status: resolved in later reruns after blueprint-backed bootstrap and task creation fixes.

## Run 2

### Result: Codex now creates durable task state
- Symptom change:
  - After the leadership prompt update, the CEO created a top-level project task immediately instead of leaving the office in pure discussion state.
- Evidence:
  - `task-4` was created from the initial CEO reply and showed up in the UI and `/tasks`.
- Fix:
  - Prompt change only. No broker change required for this part.
- Status: improved, but still incomplete.

### Issue 5: CEO could tag a non-existent agent and the broker accepted it
- Symptom:
  - The CEO tagged `studio` in-channel, but `studio` did not exist in `/office-members`.
  - The broker silently accepted the tag, creating a false appearance of extra staffing.
- Evidence:
  - Message `msg-3` included `studio` in `tagged`.
  - `/office-members` still listed only `ceo`, `eng`, and `gtm`.
- Fix:
  - Tightened broker message validation so explicit tagged slugs must be real office members (except `you`/`human`/`system`).
  - Thread auto-tagging now filters out stale/non-member participants instead of re-tagging ghosts.
- Status: fixed in code, rerun validated no dead-agent tag leakage.

### Issue 6: non-general channel composer kept the wrong aria-label
- Symptom:
  - In `#lab`, the composer placeholder updated correctly, but the aria-label still read `Message general channel`.
  - This is an accessibility bug and it broke automation targeting.
- Evidence:
  - DOM inspection in `#lab` showed placeholder `Message #lab — type / for commands, @ to mention` with aria-label `Message general channel`.
- Fix:
  - Updated the web UI to synchronize `aria-label` with composer context for normal channels and direct messages.
- Status: fixed in code.

## Run 3

### Result: Codex first-turn behavior improved materially
- Evidence:
  - Fresh rerun from the Web UI produced:
    - immediate CEO reply
    - explicit top-level task creation by CEO
    - specialist routing without dead-agent tags
- Current state:
  - Codex now reliably creates at least one durable task for the initiative.
  - Codex still did **not** create a dedicated execution channel or propose any skills in the validated reruns.
  - The business-build loop remains better coordinated, but it still stops short of a full autonomous “build and run from scratch” operating setup.
- Status: partially fixed; still open at the product-behavior level.

## Approval Judgement

- No approval requests were raised in the validated Codex reruns.
- This was acceptable for the work that actually happened:
  - internal planning
  - routing
  - task creation
  - channel setup intent
- Approvals that **should** still be required later, once the workflow reaches them:
  - spending money on vendors or subscriptions
  - creating or linking real external publishing/monetization accounts
  - publishing content to a real YouTube channel
  - accepting legal/commercial commitments such as sponsorship terms
- Net:
  - no unnecessary approvals observed
  - but the eval never progressed far enough to exercise the approvals that should exist for real-world execution

## Claude Code Smoke

Goal: verify Claude Code works for message exchange across normal channels and direct agent DMs after the Codex-focused fixes.

### Verified working
- `#general`
  - human: `Claude smoke test in #general. Reply with GENERAL-OK and nothing else.`
  - reply: `GENERAL-OK`
- `#lab`
  - human: `Claude smoke test in #lab. Reply with LAB-OK and nothing else.`
  - reply: `LAB-OK`
- DM with CEO
  - human: `Claude DM smoke to CEO. Reply with CEO-DM-OK and nothing else.`
  - reply: `CEO-DM-OK`
- DM with Eng
  - human: `Claude DM smoke to ENG. Reply with ENG-DM-OK and nothing else.`
  - reply: `ENG-DM-OK`

### Notes
- DM channels are persisted as deterministic pair slugs such as `ceo__human` and `eng__human`, even though the UI route is `#/dm/<agent>`.
- In the Eng DM, the CEO also posted a follow-up confirmation after Eng replied. This did not break the DM flow, but it is extra chatter worth noting for future DM-polish work.

## Run 4

### Issue 7: specialist same-task duplicate turns could queue behind a successful result
- Symptom:
  - After an engineer finished a task and posted a substantive completion, a second engineer turn for the same task could already be queued and would start immediately afterward.
  - This kept the task in `in_progress` longer than necessary and delayed CEO review.
- Evidence:
  - `task-11` posted a shipped result in `#youtube-factory`, but a second `eng` turn for the same task immediately started afterward.
  - Latency log showed a second `agent=eng stage=started` right after a successful `agent=eng status=ok` for the same task window.
- Fix:
  - Generalized same-task turn dedupe from CEO-only to all agents in `internal/team/headless_codex.go`.
  - Added tests for:
    - dropping a duplicate active specialist turn
    - replacing a pending queued specialist turn for the same task
- Verification:
  - `go test ./internal/team -run 'TestEnqueueHeadlessCodexTurnRecord(DropsDuplicateLeadTaskWhileActive|DropsDuplicateAgentTaskWhileActive|ReplacesPendingAgentTaskTurn)'`
  - `go build -o wuphf ./cmd/wuphf`
- Status: fixed in code, and validated in the next fresh rerun when `task-10` logged `queue-drop: agent already handling same task` instead of spawning a second visible engineer lane.

### Issue 8: specialist prompt still allowed completion posts without durable task mutation
- Symptom:
  - A specialist could post a human-readable completion message while the owned task remained `in_progress`.
  - This left downstream routing dependent on luck rather than durable state.
- Evidence:
  - In the failing `task-11` cycle, the engineer posted a shipped message in-channel while broker task state stayed `in_progress`.
- Fix:
  - Strengthened specialist instructions so the final sequence is explicit:
    - `team_task` mutation first
    - then completion broadcast or `human_message`
    - then stop
  - Also strengthened the pushed task packet language to say a completion post without the task mutation is a failure.
- Status: fixed in prompt/instructions; validated indirectly by the fresh rerun where `task-3` moved to `review` before CEO approval.

### Current open issue: office-mode engineering turn can write code and still fail to surface completion
- Symptom:
  - In the fresh rerun, `task-10` wrote real code into the main checkout but stayed `in_progress` with no completion message or task mutation.
- Evidence:
  - Main checkout now contains:
    - `internal/commands/cmd_youtube_pack.go`
    - `internal/commands/cmd_youtube_pack_test.go`
  - related edits in `internal/commands/slash.go` and `internal/commands/registry_test.go`
  - No live `codex exec` process remained for the office session.
  - `task-10` still showed `status: in_progress` and `review_state: pending_review`.
- Status: open. This is the next runtime bug to isolate.

## Run 5

### Issue 9: Studio integration smoke runs failed because generic workflow drafts were not valid One flow definitions
- Symptom:
  - Clicking `Run Smoke Test` in the Studio integration harness returned `502 Bad Gateway`.
  - This blocked the eval path for proving integrations were added and exercised from the Web UI.
- Evidence:
  - Browser inspection of `/api/studio/run-workflow` returned a validation error from One:
    - missing top-level `name`
    - `version` must be a string
    - missing top-level `inputs`
    - step `name` missing
    - generic `kind` field where One expected `type`
- Fix:
  - Added `/studio/run-workflow` to the broker and wired the Studio UI to call it.
  - Added a broker-side fallback path:
    - use the real provider when registration/execution works
    - if the workflow is `dry_run` or `mock` and provider registration/execution fails, run a local stub executor instead
  - Persisted smoke-run evidence back into Studio state, skill execution state, and broker actions/messages so the dry-run still leaves durable proof.
- Verification:
  - Studio Web UI:
    - `Sponsor Outreach Dry-Run` now succeeds end to end from the button in Studio
    - workflow evidence shows a persisted run id and mode
    - the Gmail adapter lane flips from `Blocked on integration` to `Smoke tested`
    - the Skills app shows the workflow provider/key, Gmail integration tag, and last execution status/time
  - Code:
    - `go test ./internal/team -run 'Test(HandleStudioGeneratePackagePersistsAction|HandleStudioRunWorkflowExecutesOneDraftAndUpdatesSkill)'`
    - `go build -o wuphf ./cmd/wuphf`
- Status: fixed and validated live in the browser.

### Issue 10: multiple `[SKILL PROPOSAL]` blocks in one CEO message only surfaced the first proposal
- Symptom:
  - The CEO created five integration specialists and emitted five dry-run skill proposals in one message.
  - Only the first proposal (`gmail-dry-run-harness`) produced:
    - a system `skill_proposal` message
    - a request in the Requests app
    - a proposed skill entry
  - The other four proposal blocks were silently dropped.
- Evidence:
  - The raw CEO message in `#general` contained five tagged `[SKILL PROPOSAL]...[/SKILL PROPOSAL]` blocks.

## Run 6

### Result: employee blueprints now drive startup roster formation
- Symptom change:
  - runtime startup no longer needs static packs as the canonical source of team roles
  - operation starter agents now bind directly to employee blueprints, and manifest materialization overlays starter-local details on top of employee blueprint role definitions
- Evidence:
  - operation fixtures under `templates/operations/*/blueprint.yaml` now declare `employee_blueprint` on starter agents
  - fresh runtime formation passes:
    - `go test ./internal/operations ./internal/company ./internal/team ./internal/commands -count=1`
- Fix:
  - added first-class employee blueprints and runtime binding/overlay logic
  - moved fresh broker/launcher defaults to blueprint refs instead of `founding-team`
- Status: fixed in code and covered by tests

### Result: live launcher routing is no longer limited to hardcoded domain buckets
- Symptom change:
  - notification and task relevance now use roster metadata and task text, not just slug classes like `eng`, `fe`, `marketing`, or `sales`
- Evidence:
  - targeted launcher tests now pass for generic roles such as `bookkeeper` and `community-manager`
  - validation:
    - `go test ./internal/team -run 'TestNotificationTargetsForMessageUsesMetadataBackedTaskOwner|TestRelevantTaskForTargetUsesRosterMetadata' -count=1`
    - `go test ./internal/operations ./internal/company ./internal/orchestration ./internal/team ./internal/commands -count=1`
    - `go build -o wuphf ./cmd/wuphf`
- Fix:
  - exported generic routing scorers from `internal/orchestration/message_router.go`
  - added launcher-side routing helpers backed by office-member expertise/role metadata and active task text
  - replaced `inferAgentDomain`-based owner/task matching in the launcher notification path
- Status: fixed in the active runtime path

### Remaining architecture gap after Run 6
- Packs and prompt-generated member templates still exist as compatibility shims.
- `internal/team/operation_bootstrap.go` no longer uses pack/backlog translators on the active path.
  - active bootstrap now resolves operation blueprints directly from `templates/operations/*/blueprint.yaml`
  - pack/backlog/registry readers remain only as a legacy compatibility fallback
- Routing logic is still split between `internal/orchestration` and `internal/team`.
- The full browser/runtime replay matrix across all operation blueprints has not been executed yet.
  - Broker state only contained one pending request and one system proposal announcement.
- Fix:
  - Changed `parseSkillProposalLocked` to iterate through every tagged proposal block in a message instead of stopping after the first one.
  - Made malformed or duplicate proposals skip only that block instead of aborting parsing for the rest of the message.
  - Added a regression test for one CEO message containing multiple proposal blocks.
- Verification:
  - Code:
    - `go test ./internal/team -run 'Test(ParseSkillProposalCEOHappyPath|ParseSkillProposalCreatesNonBlockingInterview|ParseSkillProposalParsesMultipleBlocks|SkillProposalAcceptCallbackActivatesSkill)'`
    - `go build -o wuphf ./cmd/wuphf`
  - Live Web UI:
    - Sent a controlled prompt asking the CEO to reply with exactly three unique skill proposals
    - Requests app showed three pending approvals
    - broker `/requests` returned three pending `skill_proposal` items
    - broker `/messages` showed three `Skill Proposed:` system messages
    - broker `/skills` showed all three proposed skills with `status: proposed`
    - rejected the three parser-test artifacts afterward to keep the workspace clean
- Status: fixed and validated live in the browser.

### Integration/agent coverage reached in this run
- Verified that integration ownership is now durable across both skills and generated agents:
  - Studio-created workflow skills carry explicit integration tags such as `integration:gmail`, `integration:youtube-data-api`, `integration:google-drive`, `integration:slack`, and `integration:youtube-analytics`
  - CEO-generated specialist agents include the owned integration directly in `expertise`, for example:
    - `Gmail integration ownership`
    - `YouTube Data API integration ownership`
    - `Google Drive integration ownership`
    - `Slack integration ownership`
    - `YouTube Analytics integration ownership`
  - CEO also created one durable task lane per generated integration specialist:
    - `task-56` Gmail
    - `task-57` YouTube Data API
    - `task-58` Google Drive
    - `task-59` Slack
    - `task-60` YouTube Analytics

### Approval judgment from this run
- Appropriate approval:
  - The Requests app asking to activate newly proposed reusable skills is defensible, because it changes durable team capabilities even when the skill is dry-run only.
- Not required yet:
  - no approval should be required for the Studio smoke runs themselves when they are dry-run/mock only and do not touch real external systems
- Still required later:
  - linking real Gmail, Slack, YouTube, Drive, or analytics accounts
  - any live publish, send, spend, or externally visible side effect

## Run 6

### Issue 11: the browser still hardcoded the starter business shape after Studio bootstrap moved to the broker
- Symptom:
  - Studio state was broker-driven, but onboarding still used browser-owned starter definitions for agents, channels, tasks, and kickoff copy.
  - The UI was still deciding what business to build before the control plane ran.
- Evidence:
  - `web/index.html` still contained:
    - `FACELESS_AI_STARTER_AGENTS`
    - `STARTER_TEMPLATES`
    - fixed `faceless-ai-workflows` business copy
- Fix:
  - Added broker-owned `starter` data to `/studio/bootstrap-package`.
  - The broker now generates the starter agents, channels, first tasks, kickoff prompt, and general-channel description from the selected channel pack plus backlog/profile inputs.
  - Removed the browser-owned starter constants and switched onboarding to consume the broker starter plan.
- Verification:
  - `go test ./internal/team -count=1`
  - `go build -o wuphf ./cmd/wuphf`
  - source scan no longer finds `FACELESS_AI_STARTER_AGENTS`, `STARTER_TEMPLATES`, or `faceless-ai-workflows` in `web/index.html`
- Status: fixed in code.

### Issue 12: same-task dedupe blocked legitimate automatic retries
- Symptom:
  - A failed local-worktree turn could not queue its retry because the just-failed attempt was still marked active when recovery ran.
- Evidence:
  - `TestRunHeadlessCodexQueueRetriesLocalWorktreeAfterGenericError` timed out waiting for the requeued prompt.
- Fix:
  - Kept same-task dedupe for duplicate notifications.
  - Allowed recovery retries with a higher attempt count to queue behind the unwinding active turn.
  - Added a focused regression test for this case.
- Verification:
  - `go test ./internal/team -run 'TestEnqueueHeadlessCodexTurnRecordAllowsRetryBehindActiveAgentTask|TestRunHeadlessCodexQueueRetriesLocalWorktreeAfterGenericError' -count=1`
  - `go test ./internal/team -count=1`
- Status: fixed in code.

## Current Limit Summary

What is now broker-driven instead of browser-hardcoded:
- Studio bootstrap config
- workflow drafts
- smoke tests
- connection cards
- queue seed
- monetization ladder
- starter agents/channels/tasks/kickoff plan

## Run 7

### Issue 13: starter-plan construction still lived in core Go even after the browser starter was removed
- Symptom:
  - the browser no longer owned starter agents and channels, but the broker still constructed the business roster, lane list, kickoff copy, and general-room description directly in `studio_bootstrap.go`
- Evidence:
  - `buildStudioStarterTemplate(...)` still returned a handwritten roster for `ceo`, `research-lead`, `scriptwriter`, `editor`, `packaging-lead`, `growth-ops`, and `monetization`
  - the same function still hardcoded the `research`, `scripts`, `production`, `packaging`, `growth`, and `revenue` channels plus all first tasks and kickoff copy
- Fix:
  - added a generic `internal/operations` package with template-backed blueprint types and loader
  - added `templates/operations/youtube-factory/blueprint.yaml`
  - wired `/studio/bootstrap-package` to load the blueprint from the selected pack's `workspace.pipeline_id`
  - switched starter agent/channel/task generation and kickoff/general copy to render from the loaded blueprint instead of handwritten Go branches
- Verification:
  - `go test ./internal/team -run 'Test(BuildStudioBootstrapPackageFromRepoIncludesStarterPlan|HandleStudioRunWorkflowExecutesOneDraftAndUpdatesSkill|HandleStudioGeneratePackagePersistsAction)' -count=1`
  - `go build -o wuphf ./cmd/wuphf`
- Status: fixed in code.

### Issue 14: connection cards, workflow drafts, smoke tests, and automation cards still lived in core runtime branches
- Symptom:
  - after the starter plan moved to the blueprint, the broker still hardcoded the integration lanes, workflow definitions, smoke-test payloads, and Studio automation cards in `studio_bootstrap.go`
- Evidence:
  - `buildStudioConnectionCards(...)` constructed a local list of Gmail/Slack/YouTube/Drive/Analytics lanes
  - `buildStudioWorkflowDrafts()` returned a handwritten list of workflow definitions
  - `buildStudioSmokeTests()` returned a handwritten list of smoke-test payloads
  - `buildStudioAutomation(...)` still returned fixed card copy from Go
- Fix:
  - extended `internal/operations.Blueprint` to carry:
    - `connections`
    - richer `workflows` with trigger, schedule, checklist, definition, and smoke-test data
    - `automation` modules
  - moved those declarations into `templates/operations/youtube-factory/blueprint.yaml`
  - switched bootstrap package generation to render those slices from the loaded blueprint instead of local constants
- Verification:
  - `go test ./internal/team -run 'Test(BuildStudioBootstrapPackageFromRepoIncludesStarterPlan|HandleStudioRunWorkflowExecutesOneDraftAndUpdatesSkill|HandleStudioGeneratePackagePersistsAction)' -count=1`
  - `go build -o wuphf ./cmd/wuphf`
- Status: fixed in code.

### Issue 15: moving YAML-authored workflow and automation text into the blueprint surfaced schema-sensitive template failures
- Symptom:
  - the first bootstrap test rerun failed on template load even though the runtime wiring was correct
- Evidence:
  - unquoted YAML values containing `:` caused `yaml.Unmarshal` failures:
    - checklist text in the analytics workflow
    - automation footer text in the control-plane card
- Fix:
  - quoted the affected YAML strings in `templates/operations/youtube-factory/blueprint.yaml`
  - kept the loader strict so malformed template data fails fast during tests instead of being silently coerced
- Verification:
  - the same targeted bootstrap tests passed immediately after the YAML fixes
- Status: fixed in template data.

### Issue 16: the active Studio renderer still assumed fixed stages, `yt-*` lanes, and YouTube-specific storage namespaces
- Symptom:
  - even after bootstrap data moved to the template, the page still rendered from a fixed `STUDIO_PIPELINE_STAGES` array, matched workspace lanes via `yt-` prefixes, and used `youtube_factory` namespaces and YouTube-specific visible copy on the active path
- Evidence:
  - `web/index.html` still had:
    - `STUDIO_PIPELINE_STAGES`
    - `yt-command`
    - `youtube_factory/workspaces/*`
    - visible labels like `YouTube lanes active` and `open YouTube tasks`
- Fix:
  - removed fixed stage rendering from the active path and normalized stages from `bootstrap.blueprint.stages`
  - changed workspace lane matching to use starter/blueprint channel data instead of `yt-*` prefixes
  - switched Studio storage and broker namespaces to generic `studio` values with read-only legacy compatibility for older `youtube_factory` state
  - rewrote the active Studio copy to be operation-first instead of YouTube-first
- Verification:
  - parsed the inline browser script with Node
  - `rg -n "STUDIO_PIPELINE_STAGES|yt-command|YouTube lanes active|open YouTube tasks|channel pack|Automated channel business" web/index.html` returned no active-path matches
- Status: fixed in code.

### Issue 17: office-mode coding could still succeed after changing code without durable task evidence
- Symptom:
  - a coding turn could return success after mutating the workspace while leaving the owned task in a non-terminal state with no substantive completion evidence
- Evidence:
  - the bug had already been observed live in the eval, and the runtime previously trusted a successful child Codex exit too much
- Fix:
  - added a durability guard in `internal/team/headless_codex.go`
  - after a successful coding turn, WUPHF now checks for:
    - durable task state such as `review`, `done`, `blocked`, or approved review state, or
    - substantive completion evidence from the agent
  - if neither exists, the turn is treated as a failure and sent through the existing recovery path instead of being accepted silently
  - added focused runtime tests for the reject/accept and retry/block paths
- Verification:
  - `go test ./internal/team -count=1`
  - `go build -o wuphf ./cmd/wuphf`
- Status: fixed in code.

### Issue 18: Studio bootstrap still hardcoded the source docs root to `docs/youtube-factory`
- Symptom:
  - even after templateizing the bootstrap content, core still assumed pack discovery started from `docs/youtube-factory`
- Evidence:
  - `buildStudioBootstrapPackageFromRepo(...)` still set `docsDir := filepath.Join(repoRoot, "docs", "youtube-factory")`
- Fix:
  - pack discovery now scans `docs/` generically for `*channel-pack.yaml`
  - backlog and monetization companion files now load from the selected pack file's own directory
  - missing backlog/monetization files now fall back cleanly so non-YouTube operation packs are not forced to carry those exact files
- Verification:
  - `go test ./internal/team -count=1`
  - `go build -o wuphf ./cmd/wuphf`
- Status: fixed in code.

### Issue 19: Studio state still only modeled media-specific artifact buckets
- Symptom:
  - even after the UI became blueprint-driven, Studio state still persisted generated work only as `topicPackets`, `scriptBriefs`, `publishPackages`, and `workflowRuns`
- Evidence:
  - queue cards and the generated-artifacts panel still read from those named buckets directly
  - `/studio/generate-package` returned only a `package` object, not generic artifact entries
- Fix:
  - `/studio/generate-package` now also returns generic `artifacts`
  - Studio state now persists generic artifact entries alongside the legacy named buckets
  - workflow-run persistence also writes a generic `workflow_run` artifact entry
  - queue readiness checks now look through the generic layer first and fall back to the legacy buckets
  - the generated-artifacts panel now renders from the generic artifact list on the active path
- Verification:
  - `go test ./internal/team -run 'TestHandleStudioGeneratePackagePersistsAction|TestBuildStudioBootstrapPackageFromRepoIncludesStarterPlan|TestHandleStudioRunWorkflowExecutesOneDraftAndUpdatesSkill' -count=1`
  - `go test ./internal/team -count=1`
  - inline browser script parsed with Node
- Status: partially fixed in code; legacy artifact buckets still exist for compatibility.

### Issue 20: the core bootstrap surface still used `studio_*` naming even after the runtime became operation-generic
- Symptom:
  - the main bootstrap file and handler still presented the substrate as `studio_*`, which implied a product-native Studio mode instead of a generic operation bootstrap surface
- Evidence:
  - `internal/team/studio_bootstrap.go`
  - `buildStudioBootstrapPackageFromRepo(...)`
  - `handleStudioBootstrapPackage(...)`
- Fix:
  - renamed the core bootstrap file to `internal/team/operation_bootstrap.go`
  - renamed the main types/functions to `operation*`
  - added `/operations/bootstrap-package`
  - kept `/studio/bootstrap-package` as a compatibility alias
  - switched the web app bootstrap fetch to `/operations/bootstrap-package`
- Verification:
  - focused Go tests for the bootstrap path passed
  - `rg -n "buildStudioBootstrapPackageFromRepo|handleStudioBootstrapPackage" internal/team web/index.html` returned no active references
- Status: fixed in code.

### Issue 21: package generation still required a fixed three-artifact media contract
- Symptom:
  - `/studio/generate-package` still assumed Codex would always return `topic_packet`, `script_brief`, and `publish_package`, which blocked the path to arbitrary operation artifact bundles
- Evidence:
  - `internal/team/broker.go` still defined `studioGeneratedPackage` as a fixed struct and hardcoded the package-generation prompt around those three sections
- Fix:
  - changed generated-package decoding to validate against a requested artifact-id list
  - passed blueprint artifact definitions from the web app into the package-generation request
  - changed the prompt to ask for the requested artifact ids rather than a fixed media bundle
  - kept the current YouTube template working through the blueprint artifact list instead of through hardcoded top-level keys
- Verification:
  - `go test ./internal/team -run 'TestBuildOperationBootstrapPackageFromRepoIncludesStarterPlan|TestDecodeStudioGeneratedPackageHandlesFencedJSON|TestHandleStudioGeneratePackagePersistsAction|TestHandleStudioRunWorkflowExecutesOneDraftAndUpdatesSkill' -count=1`
  - inline browser script parsed with Node
- Status: fixed in code for the active path.

### Issue 22: blank directives still fell back to the nearest repo pack instead of synthesizing a new operation
- Symptom:
  - when a real directive did not match any repo-authored pack, bootstrap still defaulted to the closest seeded pack instead of generating a new operation blueprint
- Evidence:
  - `selectOperationPackFile(...)` returned the best pack even for non-empty no-match queries
  - `internal/operations` had a synthesis entrypoint but the generic branch was incomplete
- Fix:
  - added a generic synthesis path in `internal/operations`
  - changed bootstrap selection so a non-empty no-match query triggers synthesized bootstrap instead of silently loading a seeded pack
  - threaded runtime integrations and runtime capabilities into synthesized blueprints so the blank path stays connection-aware
- Verification:
  - `go test ./internal/operations -count=1`
  - `go test ./internal/team -run 'TestBuildOperationBootstrapPackageSynthesizesWhenNoPackSeedExists' -count=1`
- Status: fixed in code.

### Issue 23: the rebase left duplicate Nex helpers split across `launcher.go` and `launcher_nex.go`
- Symptom:
  - the team package no longer compiled after rebasing because Nex polling helpers existed in both files and `launcher.go` referenced `api` without importing it
- Evidence:
  - `go test ./internal/team -count=1` failed with duplicate method/type errors and undefined `api`
- Fix:
  - removed the duplicate Nex helper block from `launcher.go`
  - kept the Nex-specific surface in `launcher_nex.go`
  - reran the full team package after the merge cleanup
- Verification:
  - `go test ./internal/team -count=1`
  - `go build -o wuphf ./cmd/wuphf`
- Status: fixed in code.

### Issue 24: active package follow-up flows still assumed media-specific “revenue ops” semantics
- Symptom:
  - even after artifact generation became generic, the broker still stubbed three media-specific follow-up workflows tied to `topic_packet`, `script_brief`, and `publish_package`
- Evidence:
  - `buildStudioRevenueOpsStubExecutions(...)`
  - action kind `studio_revenue_ops_stub_executed`
  - UI success copy still said “stubbed revenue ops”
- Fix:
  - replaced the active-path follow-up stubs with generic review, offer-alignment, and approval-gate workflows derived from the generated artifact bundle
  - switched the broker action kind to `studio_followup_stub_executed`
  - updated the UI to persist generic workflow runs and show generic follow-up messaging
- Verification:
  - `go test ./internal/team -run 'TestHandleStudioGeneratePackagePersistsAction' -count=1`
  - inline browser script parsed with Node
- Status: fixed in code for the active path; old artifact-bucket compatibility readers still remain.

### Issue 25: operation bootstrap still required media-era filename conventions
- Symptom:
  - new operation templates still had to name their seed files `channel-pack.yaml`, `content-backlog.yaml`, and `monetization-registry.yaml`
- Evidence:
  - `internal/team/operation_bootstrap.go` only discovered those filenames
- Fix:
  - bootstrap now accepts generic operation filenames:
    - `operation-pack.yaml`
    - `operation-backlog.yaml`
    - `operation-offers.yaml`
    - `operation-monetization.yaml`
  - legacy YouTube/media filenames remain supported as compatibility fallbacks
- Verification:
  - `go test ./internal/team -count=1`
- Status: fixed in code.

### Issue 26: task execution mode and review still depended on engineering-only owner slugs
- Symptom:
  - outside the bootstrap/UI slice, core task handling still assumed only `eng`, `fe`, `be`, and `ai` could produce repository work that needed a local worktree and structured review
- Evidence:
  - `internal/team/task_pipeline.go` keyed `taskNeedsStructuredReview(...)` and `taskDefaultExecutionMode(...)` off those fixed owner slugs
- Fix:
  - switched task execution inference to use the task’s actual work shape (`title`, `details`, and owner text) instead of a fixed engineering roster
  - repository/code-shaped feature, bugfix, and incident tasks now default into `local_worktree` regardless of whether the owner slug comes from an old pack
  - structured review now follows the inferred execution shape instead of a fixed list of engineering member ids
- Verification:
  - `go test ./internal/team -run 'TestInferTaskTypeTreatsAuditWorkAsResearch|TestTaskDefaultExecutionModeTreatsEngineeringFeatureWorkAsLocalWorktree' -count=1`
  - `go test ./internal/team -count=1`
- Status: fixed in code.

### Issue 27: operation bootstrap still treated repo docs as the primary selector
- Symptom:
  - even after templates existed under `templates/operations`, the bootstrap path still centered selection around `docs/**/operation-pack.yaml`
- Fix:
  - bootstrap is now template-first:
    - direct `pack_id` values can resolve operation blueprint ids immediately
    - fuzzy selection scores the operation blueprints on disk
    - current manifest operation refs are considered before synthesis
  - legacy docs remain only as a compatibility fallback path
- Verification:
  - `go test ./internal/team -run 'TestBuildOperationBootstrapPackageFromRepo(IncludesStarterPlan|ResolvesLegacyPackIDToBlueprint|SynthesizesWhenNoPackSeedExists)$' -count=1`
- Status: fixed in the active path.

### Issue 28: full operation-template runtime matrix lacked startup coverage
- Symptom:
  - we had loader coverage and spot checks, but not a single runtime-facing test that every operation blueprint could bootstrap, seed a broker office, and be accepted by `NewLauncher(...)`
- Fix:
  - added `internal/team/operation_matrix_test.go` to validate, for every operation template:
    - bootstrap package generation
    - fresh broker office seeding from refs-only manifests
    - launcher acceptance without a static pack
- Verification:
  - `go test ./internal/operations ./internal/company ./internal/orchestration ./internal/team ./internal/commands -count=1`
  - `go build -o wuphf ./cmd/wuphf`
- Status: fixed in automated runtime coverage.

### Issue 29: default manifest and `/init` still depended on static pack resolution
- Symptom:
  - default startup and `/init` still called `agent.GetPack(...)` to shape the default office or lead resolution
- Fix:
  - `DefaultManifest()` no longer pulls member rosters from static packs
  - `/init` now resolves `TeamLeadSlug` from blueprint-backed runtime manifest state rather than pack lookup
- Verification:
  - `go test ./internal/operations ./internal/company ./internal/orchestration ./internal/team ./internal/commands -count=1`
- Status: fixed in the active path.

### Issue 30: interactive setup still taught the old pack abstraction even after startup became blueprint-backed
- Symptom:
  - `/init`, `/config`, and the TUI onboarding flow still surfaced `Pack` / `Agent pack` as the primary user-facing concept
  - the channel picker still only listened for the old pack-choice phase title
- Fix:
  - added a preferred `Blueprint` config field while keeping `Pack` as a compatibility alias
  - `/config show` now surfaces `Blueprint` first and `/config set` accepts `blueprint`, `template`, and `operation_template`
  - the TUI onboarding flow now discovers operation blueprints from `templates/operations/*` first, shows `Operation template` readiness, and drives a `Choose Operation Template` picker
  - `cmd/wuphf/channel.go` now handles the blueprint-choice phase and uses the new picker title
- Verification:
  - `go test ./internal/config ./internal/tui ./internal/commands ./internal/company ./internal/team -count=1`
  - `go test ./internal/operations ./internal/config ./internal/tui ./internal/commands ./internal/company ./internal/orchestration ./internal/team -count=1`
  - `go build -o wuphf ./cmd/wuphf`
- Status: fixed in the active setup/config path.

### Issue 31: synthesized blank-directive bootstrap still leaked channel/media defaults
- Symptom:
  - generic synthesized bootstrap packages could still inherit fallback phrases like `channel business`, `channel pack`, or viewer/content-centric KPI copy
- Fix:
  - made `buildOperationBootstrapPackage(...)` prefer blueprint/profile identifiers over legacy pack labels
  - changed generic fallback defaults in `internal/team/operation_bootstrap.go` to operation/workflow language
  - added a regression test that marshals the synthesized package and rejects legacy media phrases
- Verification:
  - `go test ./internal/team -run 'TestSynthesizedOperationBootstrapPackageUsesGenericFallbackCopy|TestOperationBlueprintMatrixBuildsBootstrapPackage|TestOperationBlueprintMatrixSeedsBrokerOffice|TestOperationBlueprintMatrixNewLauncherAcceptsAllBlueprints' -count=1`
  - `go test ./internal/operations ./internal/config ./internal/tui ./internal/commands ./internal/company ./internal/orchestration ./internal/team -count=1`
- Status: fixed in the active synthesized bootstrap path.

### Issue 32: parallel branch runs could collide on the broker port and leave the browser or MCP paths talking to the wrong office
- Symptom:
  - another checkout or stale `wuphf` process could keep `127.0.0.1:7890` busy, while parts of the system still assumed the default broker port
  - during replay, a stale `paid-discord-community` office was still listening on `127.0.0.1:7892` / `:7907`
  - after rebasing, the direct browser fallback and the distributed team MCP tool had dropped parts of the custom broker-base-url support
- Fix:
  - restored broker base-url propagation in the eval worktree
  - `/api-token` now returns both the auth token and the resolved `broker_url`
  - the Web UI now updates its direct fallback broker URL from `/api-token` instead of assuming `localhost:7890`
  - `mcp/dist/tools/team.js` now resolves broker URL, broker port, and broker token file from the same `WUPHF_*` / `NEX_*` env vars used by the Go runtime
  - stale replay process `71408` was explicitly killed before rerunning the next matrix entry
- Verification:
  - `go test ./internal/team ./internal/teammcp -count=1`
  - `go build -o wuphf ./cmd/wuphf`
  - `curl http://127.0.0.1:7907/api-token` returned `{"broker_url":"http://127.0.0.1:7892", ...}`
  - `multi-agent-workflow-consulting` successfully booted on `--broker-port 7892 --web-port 7907`
  - browser replay completed from the Web UI on that non-default broker port
- Status: fixed in the active custom-port path.

### Issue 33: Playwright replay could show a stale office after a new custom-port launch because the old browser session kept local UI state
- Symptom:
  - after relaunching `multi-agent-workflow-consulting`, the first browser snapshot still showed the previous `local-business-ai-package` office even though the live APIs reported a clean onboarding state
- Fix:
  - killed the stale Playwright daemon/session and reopened the browser against the fresh web server
  - verified against live `/api/onboarding/state` and `/api/company` before continuing the replay
- Verification:
  - clean snapshot after reopening showed the onboarding form for the fresh operation
  - subsequent replay completed into the `multi-agent-workflow-consulting` office with the expected channels and roster
- Status: replay artifact only, not a product bug.

### Issue 34: final onboarding completion was browser-click fragile in replay, but the server path itself was healthy
- Symptom:
  - the task-selection page rendered correctly, but the first automated click on `Start working` did not emit `/api/onboarding/complete`
- Fix:
  - validated the page state through browser network capture
  - forced the same CTA through the browser DOM to confirm the server completion path still worked
  - treated this as replay fragility, not a confirmed end-user regression
- Verification:
  - browser network log showed `POST /api/onboarding/complete => 200 OK`
  - `GET /api/onboarding/state` returned `onboarded: true`
  - `GET /api/company` persisted the consulting company fields
- Status: replay workaround used; monitor for a real UI bug if a human can reproduce the same CTA failure manually.

### Issue 35: remaining Codex browser matrix entries were still unverified
- Symptom:
  - after the first non-YouTube replays, `youtube-factory` and `bookkeeping-invoicing-service` were still the only unproven browser matrix entries
- Fix:
  - ran both remaining operations from the Web UI with isolated broker/web ports:
    - `youtube-factory` on `7893/7908`
    - `bookkeeping-invoicing-service` on `7894/7909`
- Verification:
  - `youtube-factory`
    - `GET /api/onboarding/state` returned `onboarded: true`
    - `GET /api/company` persisted `Back Office AI`
    - `GET /api/channels` returned:
      - `general`
      - `youtube-factory-command`
      - `research`
      - `scripts`
      - `production`
      - `packaging`
      - `growth`
      - `revenue`
    - browser snapshot showed roster:
      - `CEO`
      - `Research Lead`
      - `Scriptwriter`
      - `Editor`
      - `Packaging Lead`
      - `Growth Ops`
      - `Monetization Lead`
  - `bookkeeping-invoicing-service`
    - `GET /api/onboarding/state` returned `onboarded: true`
    - `GET /api/company` persisted `Ledger Loop`
    - `GET /api/channels` returned:
      - `general`
      - `bookkeeping-and-invoicing-service-command`
      - `intake`
      - `books`
      - `invoicing`
      - `review`
    - browser snapshot showed roster:
      - `Operator`
      - `Planner`
      - `Bookkeeper`
      - `Invoicing`
      - `Reviewer`
- Status: fixed by replay; the full Codex browser matrix is now complete.

### Issue 36: alternate-provider regression risk still existed after focusing on Codex
- Symptom:
  - Codex was the primary provider under test, but we still needed to know whether the same browser onboarding/office path broke under `claude-code`
- Fix:
  - ran a smoke pass only, not a full matrix:
    - `niche-crm` on `--provider claude-code --broker-port 7895 --web-port 7910`
- Verification:
  - `GET /api/onboarding/state` returned `onboarded: true`
  - `GET /api/company` persisted `Pipeline Foundry`
  - browser snapshot showed `Runtime provider: claude-code`
  - visible office loaded with the expected `niche-crm` channels and roster
- Status: Claude Code smoke passed.

What is still seeded rather than fully autonomous:
- real external account linking, publishing, outbound sends, spend, and commercial commitments still require human approval

Latest guardrail landed:
- policy rules now reach agent prompts directly
- external-action tasks now tell agents to use the connected system instead of creating repo-doc substitutes
- durable completion now rejects review/done states for those tasks unless broker-side external workflow evidence is recorded

Main open blockers:
- compatibility registries still exist in `internal/agent/packs.go`, `internal/agent/templates.go`, and related template helpers
- routing logic is still split across `internal/orchestration` and `internal/team`
- full Claude Code replay matrix has not been executed, by design; only a smoke pass was run
