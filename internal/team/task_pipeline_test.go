package team

import "testing"

func TestInferTaskTypeTreatsAuditWorkAsResearch(t *testing.T) {
	got := inferTaskType("eng", "Audit repo and define the fastest path to a working web UI", "Working plan for the faceless YouTube business build.")
	if got != "research" {
		t.Fatalf("inferTaskType returned %q, want research", got)
	}
}

func TestTaskDefaultExecutionModeTreatsEngineeringFeatureWorkAsLocalWorktree(t *testing.T) {
	if got := taskDefaultExecutionMode("operator", "feature", "Implement the repository-backed approval engine", "Ship the code path and tests."); got != "local_worktree" {
		t.Fatalf("taskDefaultExecutionMode returned %q, want local_worktree", got)
	}
	if got := taskDefaultExecutionMode("generalist", "bugfix", "Fix the API regression in the worker", "Debug the code path and update the repo."); got != "local_worktree" {
		t.Fatalf("taskDefaultExecutionMode returned %q, want local_worktree", got)
	}
	if got := taskDefaultExecutionMode("gtm", "launch", "Launch the outbound sequence", "Coordinate the campaign and approvals."); got != "office" {
		t.Fatalf("taskDefaultExecutionMode returned %q, want office", got)
	}
}

func TestTaskRequiresRealExternalExecution(t *testing.T) {
	if !taskRequiresRealExternalExecution(&teamTask{
		Title:   "Create one new Notion proof artifact for the intake slice",
		Details: "Use the connected Notion workspace and leave evidence in channel.",
		Owner:   "builder",
		Status:  "in_progress",
	}) {
		t.Fatal("expected Notion create task to require real external execution")
	}
	if taskRequiresRealExternalExecution(&teamTask{
		Title:   "Generate a local preview packet for Slack handoff",
		Details: "Preview only. Do not post externally yet.",
		Owner:   "builder",
		Status:  "in_progress",
	}) {
		t.Fatal("expected preview-only task not to require real external execution")
	}
}
