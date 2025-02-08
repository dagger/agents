package main

import (
	"context"
	"dagger/progress/internal/dagger"
	"fmt"
	"strings"
	"time"
)


// NewProgress creates a new progress report for tracking tasks on a GitHub issue
func New(
	// A unique identifier for this progress report in the given issue.
	// Using the same key on the same issue will overwrite the same comment in the issue
	key string,
	// GitHub authentication token
	token *dagger.Secret,
	// Github repository to send updates
	repo string,
	// Issue number to report progress on
	issue int,
) Progress {
	return Progress{
		Token: token,
		Repo:  repo,
		Issue: issue,
		Key:   key,
	}
}

// A progress report via a github issue comment
type Progress struct{
	Token   *dagger.Secret
	Repo    string // +private
	Issue   int    // +private
	Key     string // +private
	Title   string
	Summary string
	Tasks   []Task
}

type Task struct {
	Key         string
	Description string
	Status      string
}

// Write a new summary for the progress report.
// Any previous summary is overwritten.
// This function only stages the change. Call publish to actually apply it.
func (r Progress) WriteSummary(
	ctx context.Context,
	// The text of the summary, markdown-formatted
	// It will be formatted as-is in the comment, after the title and before the task list
	summary string,
) (Progress, error) {
	r.Summary = summary
	return r, nil
}

// Append new text to the summary, without overwriting it
// This function only stages the change. Call publish to actually apply it.
func (r Progress) AppendSummary(
	ctx context.Context,
	// The text of the summary, markdown-formatted
	// It will be formatted as-is in the comment, after the title and before the task list
	summary string,
) (Progress, error) {
	if r.Summary == "" {
		r.Summary = summary
		return r, nil
	}
	sep := "\n"
	// If the current summary already ends with a newline,
	// don't add another one to avoid double newlines
	if strings.HasSuffix(r.Summary, "\n") {
		sep = ""
	}
	// Trim whitespace from current summary, add separator and new summary
	r.Summary = strings.TrimSpace(r.Summary) + sep + summary
	return r, nil
}

// Report the starting of a new task
// This function only stages the change. Call publish to actually apply it.
func (r Progress) StartTask(
	ctx context.Context,
	// A unique key for the task. Not sent in the comment. Use to update the task status later.
	key string,
	// The task description. It will be formatted as a cell in the first column of a markdown table
	description string,
	// The task status. It will be formatted as a cell in the second column of a markdown table
	status string,
) (Progress, error) {
	r.Tasks = append(r.Tasks, Task{
		Key:         key,
		Description: description,
		Status:      status,
	})
	return r, nil
}

// Write a new title for the progress report.
// Any previous title is overwritten.
// This function only stages the change. Call publish to actually apply it.
func (r Progress) WriteTitle(
	ctx context.Context,
	// The summary. It should be a single line of unformatted text.
	// It will be formatted as a H2 title in the markdown body of the comment
	title string,
) (Progress, error) {
	r.Title = strings.ToTitle(title)
	return r, nil
}

// Update the status of a task
// This function only stages the change. Call publish to actually apply it.
func (r Progress) UpdateTask(
	ctx context.Context,
	// A unique key for the task. Use to update the task status later.
	key string,
	// The task status. It will be formatted as a cell in the second column of a markdown table
	status string,
) (Progress, error) {
	for i := range r.Tasks {
		if r.Tasks[i].Key == key {
			r.Tasks[i].Status = status
			return r, nil
		}
	}
	return r, fmt.Errorf("no task at key %s", key)
}

// Publish all staged changes to the status update.
// This will cause a single comment on the target issue to be either
// created, or updated in-place.
func (r Progress) Publish(ctx context.Context) error {
	var contents string
	if r.Title != "" {
		contents = "## " + r.Title + "\n\n"
	}
	if r.Summary != "" {
		contents += r.Summary + "\n\n"
	}
	if len(r.Tasks) > 0 {
		contents += "### Tasks\n\n"
		contents += "<table>\n<tr><th>Description</th><th>Status</th></tr>\n"
		for _, task := range r.Tasks {
			contents += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>\n", task.Description, task.Status)
		}
		contents += "</table>\n"
	}
	contents += fmt.Sprintf("\n<sub>*Last update: %s*<sub>\n", time.Now().Local().Format("2006-01-02 15:04:05 MST"))
	comment := dag.GithubComment(r.Token, r.Repo, dagger.GithubCommentOpts{Issue: r.Issue, MessageID: r.Key})
	_, err := comment.Create(ctx, contents)
	return err
}
