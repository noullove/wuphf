package onboarding

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/nex-crm/wuphf/internal/operations"
)

// TaskTemplate describes a first-task suggestion shown during onboarding.
// Templates are scoped to a specific agent role via OwnerSlug.
type TaskTemplate struct {
	// ID is a stable, URL-safe identifier for the template.
	ID string `json:"id"`

	// Title is the short, human-readable task name.
	Title string `json:"title"`

	// Description is a single-sentence clarification shown below the title.
	Description string `json:"description"`

	// OwnerSlug is the agent slug that should receive this task.
	OwnerSlug string `json:"owner_slug"`
}

// DefaultTemplates returns the generic fallback starter tasks used when no
// blueprint-specific task list can be resolved.
func DefaultTemplates() []TaskTemplate {
	return []TaskTemplate{
		{ID: "landing", Title: "Draft the landing page", Description: "Hero, value props, one clear CTA. Not the WUPHF.com approach.", OwnerSlug: "executor"},
		{ID: "repo", Title: "Set up repo structure", Description: "Folders, README, CI scaffold. Dwight would document everything.", OwnerSlug: "executor"},
		{ID: "spec", Title: "Write the product spec", Description: "What we're building, why, and what done looks like. Michael would skip this step.", OwnerSlug: "planner"},
		{ID: "readme", Title: "Write the README", Description: "Installation, usage, one example. Short enough that someone actually reads it.", OwnerSlug: "planner"},
		{ID: "audit", Title: "Audit the competition", Description: "What they do, what they miss, where we win. No memos. Just findings.", OwnerSlug: "ceo"},
	}
}

// RevOpsTemplates preserves the existing legacy pack-specific starter set.
func RevOpsTemplates() []TaskTemplate {
	return []TaskTemplate{
		{ID: "pipeline_audit", Title: "Run a pipeline audit", Description: "CRM hygiene sweep — stale deals, missing fields, bad data. Find the leaks before forecast.", OwnerSlug: "analyst"},
		{ID: "meeting_prep", Title: "Prep me for my next call", Description: "One-page brief on the account, deal stage, stakeholders, and the ask. No fluff.", OwnerSlug: "ae"},
		{ID: "revive_closed_lost", Title: "Revive closed-lost leads", Description: "Surface deals lost 3–18 months ago with trigger events. Draft re-engagement outreach.", OwnerSlug: "sdr"},
		{ID: "score_inbound", Title: "Score new inbound", Description: "Rate unworked leads on fit and intent. Route Tier 1 to the AE within 24 hours.", OwnerSlug: "analyst"},
		{ID: "stalled_deals", Title: "Find stalled deals", Description: "Open pipeline with no activity in 10+ days. Diagnose the cause and recommend a next step.", OwnerSlug: "ops-lead"},
	}
}

// TemplatesForPack is a legacy alias retained for older callers that still
// talk about packs.
func TemplatesForPack(packSlug string) []TaskTemplate {
	return TemplatesForSelection("", packSlug)
}

func TemplatesForSelection(repoRoot, selection string) []TaskTemplate {
	repoRoot = resolveTemplatesRepoRoot(repoRoot)
	selection = strings.TrimSpace(selection)
	if repoRoot != "" && selection != "" {
		if blueprint, err := operations.LoadBlueprint(repoRoot, selection); err == nil {
			if templates := templatesFromBlueprint(blueprint); len(templates) > 0 {
				return templates
			}
		}
	}
	switch selection {
	case "revops":
		return RevOpsTemplates()
	default:
		return DefaultTemplates()
	}
}

func templatesFromBlueprint(blueprint operations.Blueprint) []TaskTemplate {
	out := make([]TaskTemplate, 0, len(blueprint.Starter.Tasks))
	for _, task := range blueprint.Starter.Tasks {
		title := strings.TrimSpace(task.Title)
		description := strings.TrimSpace(task.Details)
		owner := strings.TrimSpace(task.Owner)
		if title == "" || description == "" || owner == "" {
			continue
		}
		out = append(out, TaskTemplate{
			ID:          onboardingTemplateID(title),
			Title:       title,
			Description: description,
			OwnerSlug:   owner,
		})
		if len(out) == 5 {
			break
		}
	}
	return out
}

func onboardingTemplateID(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, " ", "-")
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastDash = false
		case r == '-' || r == '_' || r == '.':
			if !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		default:
			if !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

func resolveTemplatesRepoRoot(repoRoot string) string {
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return ""
		}
		repoRoot = cwd
	}
	for current := repoRoot; ; current = filepath.Dir(current) {
		if _, err := os.Stat(filepath.Join(current, "templates", "operations")); err == nil {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
	}
	return ""
}
