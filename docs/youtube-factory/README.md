# YouTube Factory GTM Kit

Date: 2026-04-14
Owner: GTM

This folder turns the faceless-channel strategy into day-one operating assets the repo can use.

Replay and eval traces:
- `../evals/2026-04-14-youtube-business-e2e.md`
  - issue-by-issue Codex eval log
- `../evals/2026-04-14-youtube-user-journey.md`
  - replayable Web UI journey with roadblocks, fixes, and current expected outcomes for a future Claude Code rerun

## Included Assets

- `channel-operating-manual.md`
  - niche rationale
  - audience thesis
  - content pillars
  - packaging rubric
  - publishing cadence
  - launch KPIs
- `launch-collateral.md`
  - channel description
  - banner and about copy
  - pinned comment and description templates
  - lead magnet copy
  - audit offer copy
  - sponsor one-sheet starter
- `launch-week-operating-brief.md`
  - locked week-one wedge and audience promise
  - first 10 upload slate
  - CTA ladder and launch-week monetization order
  - top 3 factory KPIs
  - missing business inputs the generator still needs
- `single-job-factory-channel-thesis.md`
  - audience, promise, and format constraints for the narrow `wuphf video-job` cut line
  - why the current stubbed job pack already matches a workflow-first YouTube format
  - 10 starter topic angles that fit the current generated output
- `wuphf-operator-channel-package.md`
  - product-led channel recommendation for the live WUPHF GTM motion
  - wedge, positioning, monetization ladder, KPI targets, and first 30-video slate
- `wuphf-operator-channel-pack.yaml`
  - machine-readable WUPHF operator channel seed using the existing channel-pack schema
  - brand, audience, CTA, playlist, QA, and approval defaults for the product-led wedge
- `channel-strategy.yaml`
  - machine-readable brand, ICP, content pillar, cadence, SEO, and distribution defaults
- `default-channel-pack.yaml`
  - machine-readable default channel seed for the Studio and publish generator
  - brand, render, CTA, playlist, QA, and approval defaults in one pack
- `content-backlog.yaml`
  - first 30-video backlog with priorities, CTAs, search intent, and thumbnail angles
- `episode-launch-packets/vid_01-inbox-operator.yaml`
  - publish-ready dry-run launch packet for the flagship inbox-operator episode
  - title, thumbnail brief, CTA map, description, pinned comment, and QA checklist
- `episode-launch-packets/wuphf-vid_01-steer-mid-flight.yaml`
  - publish-ready dry-run launch packet for the WUPHF flagship see-and-steer episode
  - title, thumbnail brief, CTA map, description, pinned comment, chapters, and QA checklist
- `seo-distribution-playbook.md`
  - search-intent model
  - playlist and metadata rules
  - repurposing workflow
  - per-upload distribution checklist
- `monetization-registry.yaml`
  - structured ladder of offers
  - CTA routing rules
  - approved affiliate categories
  - digital products and service offers
  - UTM defaults and disclosure copy
- `partner-matrix.yaml`
  - workflow-to-affiliate and sponsor mapping
  - disclosure defaults
  - commercialization guardrails
- `automation-sops.yaml`
  - stage-by-stage operating procedures
  - pass gates, stop conditions, and handoff artifacts
- `factory-system-design.md`
  - end-to-end automation architecture
  - stack recommendation for this repo
  - QA boundaries, failure handling, and monetization instrumentation
- `youtube-data-dry-run-workflow.md`
  - reusable workflow skill spec for YouTube Data API dry runs
  - covers channel lookup, video metadata fetch, publish-package handoff, and upload/publish placeholder validation
- `youtube-data-dry-run-checklist.yaml`
  - machine-readable QA checklist for the YouTube Data dry-run lane
  - defines must-pass, warn, and block conditions for non-mutating connector validation

## Strategic Call

- Channel thesis: `AI Back Office for Small Teams`
- Launch wedge: founder-led service businesses and agencies with 2-20 people
- Brand call: `Back Office AI`
- Format bias: long-form first, cutdowns second
- Revenue model: lead magnet -> affiliate -> low-ticket templates -> audit/sprint -> sponsors -> ads

## Live GTM Note

The assets above package the generic faceless-channel business. The additive `wuphf-operator-channel-package.md` captures the sharper product-led recommendation for the current WUPHF-vs-Paperclip launch window: lead with the WUPHF operator channel first, keep `Back Office AI` as a later expansion lane.

That recommendation is now machine-readable too in `wuphf-operator-channel-pack.yaml` and `episode-launch-packets/wuphf-vid_01-steer-mid-flight.yaml`, so the same generator path can run against the WUPHF wedge without reworking the schema.

## Script Packet Generator

The current downstream cut line is `channel brief -> script packet`.

- Input contract: `generated/channel-brief-inbox-operator.json`
- Output artifact: `generated/script-packet-inbox-operator.json`
- Review bundle: `generated/script-packet-inbox-operator-review-bundle/`
- Live-ready packet fixture: `generated/live-client-pilot/script-packet-inbox-operator.json`
- Generator: `go run ./cmd/youtube-script-packet --in docs/youtube-factory/generated/channel-brief-inbox-operator.json --out docs/youtube-factory/generated/script-packet-inbox-operator.json --bundle-dir docs/youtube-factory/generated/script-packet-inbox-operator-review-bundle`

The normalized brief stays intentionally narrow:

- `metadata`: brief id, version, source trace
- `channel`: brand and narration defaults
- `render`: scene order, duration, and visual direction
- `episode`: audience, workflow, promise, and proof asset
- `packaging`: final title, hook promise, and thumbnail constraints
- `cta`, `publish`, `qa`: the routing and guardrails narration should inherit

The generator emits one deterministic script packet JSON with:

- narration direction for voiceover
- scene-by-scene story beats for downstream scripting
- packaging guardrails so the narration stays aligned with title/thumbnail
- production notes and QA gates carried forward from the brief

The same command now emits a consultant-readable review bundle alongside the JSON packet:

- `summary.md`
  - one-page review brief with title, promise, CTA, QA focus, and story beats
- `slack-payload.json`
  - review-post stub for the dry-run handoff channel
- `google-drive-payload.json`
  - upload stub for the review folder/doc payload
- `notion-payload.json`
  - review-queue stub for checklist and metadata capture

To map the shipped review bundle into the next dry-run orchestration slice, run:

```bash
GOCACHE=/tmp/gocache go run ./cmd/review-bundle-handoff \
  --bundle-dir docs/youtube-factory/generated/script-packet-inbox-operator-review-bundle
```

This writes a sibling `*-review-handoff/` folder with:

- `workflow-run.json`
  - approval-gate verdict plus schema validation for Slack, Drive, and Notion payload contracts
- `slack-dispatch.json`
- `google-drive-dispatch.json`
- `notion-dispatch.json`
  - dry-run consumer handoff previews that stay blocked until the approval gate clears
- `handoff-summary.md`
  - operator-facing summary of blockers, release conditions, and next routing step

When the review bundle includes `approval-packet.json` and `approval-status.json`, the handoff runner uses those as the approval source of truth and carries their paths into `workflow-run.json`, `slack-dispatch.json`, and `google-drive-dispatch.json`. If they are absent, the runner falls back to the Notion preview payload so older bundles still replay.

The checked-in dry-run proof for that command now lives at:

- `docs/youtube-factory/generated/script-packet-inbox-operator-review-handoff/`
- `docs/youtube-factory/generated/live-client-pilot/script-packet-inbox-operator-review-handoff/`
- `docs/youtube-factory/generated/live-client-pilot/script-packet-inbox-operator-approved-review-handoff/`

Use the `live-client-pilot/` bundle path instead when you want the named-client variant of the same handoff lane.

The approved rerun proof uses `docs/youtube-factory/generated/live-client-pilot/script-packet-inbox-operator-approved-review-bundle/` as its source bundle. That checked-in fixture keeps `preview_only=true` in dry-run mode, but flips the gate to `release_ready` so Slack and Google Drive dispatches leave `hold_for_approval` with concrete handoff targets.

Use `--force-approve` when you need a release-ready dry-run preview without editing the fixture payloads.

The brief now also carries an `approval` block so the same fixture can move from internal dry-run to a named pilot client without hand-editing the packet narrative:

- approval mode and overall status
- named approvers with per-person status
- the live-ready packet path the next client handoff should use

## Why This Folder Exists

The repo already had a sound business-plan direction. What it lacked was reusable GTM collateral that can be handed to the Studio, the publish payload generator, and future sales/sponsor workflows without rethinking the business each run.

## New Operating Assets

- The brand system now includes multiple naming paths, visual direction, and packaging patterns in `channel-strategy.yaml`.
- The default channel seed now lives in `default-channel-pack.yaml`, so Studio state and future automation can start from a concrete bundle instead of stitching together prose docs at runtime.
- The WUPHF operator wedge now has its own publish-ready seed in `wuphf-operator-channel-pack.yaml`, using the same schema as the default channel pack.
- The backlog is now structured for automation in `content-backlog.yaml`, not just listed in prose.
- The first flagship launch packet now lives in `episode-launch-packets/vid_01-inbox-operator.yaml`, giving the factory a dry-run publish payload to wire up immediately.
- The WUPHF flagship packet now lives in `episode-launch-packets/wuphf-vid_01-steer-mid-flight.yaml`, so the product-led launch path can run in parallel with the faceless channel path.
- Commercial routing is now explicit in `partner-matrix.yaml`, so sponsor and affiliate choices follow the workflow instead of guesswork.
- Automation-ready SOPs now live in `automation-sops.yaml`, giving the control plane stage gates instead of vague GTM advice.
- The YouTube connector dry-run lane now has a reusable operator spec in `youtube-data-dry-run-workflow.md` and a machine-readable validator checklist in `youtube-data-dry-run-checklist.yaml`, so metadata and publish-stub coverage can be tested before a real uploader exists.
