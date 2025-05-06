package main

import (
	"context"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v42/github"
	"github.com/sethvargo/go-githubactions"
	"golang.org/x/oauth2"
)

var (
	markdownCommentRegex = regexp.MustCompile(`\<\!\-\-\-.*\-\-\>`)
)

type config struct {
	githubToken              string
	templatePath             string
	exemptLabels             []string
	comment                  bool
	commentEmptyDescription  string
	commentTemplateNotFilled string
	commentGithubToken       string

	repoOwner string
	repoName  string
	prNumber  int
}

func generateConfig() *config {
	cfg := config{}

	cfg.githubToken = githubactions.GetInput("repo-token")
	cfg.templatePath = githubactions.GetInput("template-path")
	cfg.exemptLabels = strings.Split(githubactions.GetInput("exempt-labels"), ",")
	cfg.comment, _ = strconv.ParseBool(githubactions.GetInput("comment"))
	cfg.commentEmptyDescription = githubactions.GetInput("comment-empty-description")
	cfg.commentTemplateNotFilled = githubactions.GetInput("comment-template-not-filled")
	cfg.commentGithubToken = githubactions.GetInput("comment-github-token")

	cfg.prNumber, _ = strconv.Atoi(githubactions.GetInput("pr-number"))
	cfg.repoOwner = githubactions.GetInput("repo-owner")
	cfg.repoName = githubactions.GetInput("repo-name")

	return &cfg
}

func fetchTemplate() (string, error) {
	data, err := os.ReadFile(cfg.templatePath)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func newGithubClient(token string) *github.Client {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func normalizeDescription(description string) string {
	description = strings.Replace(description, "\r\n", "\n", -1)
	description = markdownCommentRegex.ReplaceAllString(description, "")
	description = strings.TrimSpace(description)

	return description
}

var cfg *config

func main() {
	cfg = generateConfig()

	template, err := fetchTemplate()
	template = normalizeDescription(template)

	if err != nil {
		githubactions.Infof("Failed to fetch template: %s, will continue without template", err)
	}

	githubClient := newGithubClient(cfg.githubToken)

	pr, _, err := githubClient.PullRequests.Get(context.Background(), cfg.repoOwner, cfg.repoName, cfg.prNumber)

	if err != nil {
		githubactions.Fatalf("Failed to get PR: %s", err)
	}

	skipCheck := false
	for _, label := range pr.Labels {
		for _, exemptLabel := range cfg.exemptLabels {
			if label.GetName() == strings.Trim(exemptLabel, " ") {
				skipCheck = true
				break
			}
		}
	}

	if skipCheck {
		githubactions.Infof("Skipping check because of exempt label")
		os.Exit(0)
	}

	description := normalizeDescription(pr.GetBody())
	var errorMsg string
	if len(description) == 0 {
		errorMsg = cfg.commentEmptyDescription
	} else if len(description) <= len(template) {
		errorMsg = cfg.commentTemplateNotFilled
	}

	if errorMsg != "" {
		if cfg.comment {
			if cfg.commentGithubToken != "" {
				githubClient = newGithubClient(cfg.commentGithubToken)
			}

			_, _, err := githubClient.Issues.CreateComment(context.Background(), cfg.repoOwner, cfg.repoName, cfg.prNumber, &github.IssueComment{
				Body: &errorMsg,
			})

			if err != nil {
				githubactions.Fatalf("Failed to create comment: %s", err)
			}
		}

		githubactions.Fatalf(errorMsg)
	}

	githubactions.Infof("Description is valid")
}
