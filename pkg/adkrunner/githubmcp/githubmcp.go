package githubmcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/oauth2"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/mcptoolset"
)

const (
	githubMCPEndpoint = "https://api.githubcopilot.com/mcp/"
)

/*
Tools:
"add_comment_to_pending_review"
"add_issue_comment"
"assign_copilot_to_issue"
"create_branch"
"create_or_update_file"
"create_pull_request"
"create_repository"
"delete_file"
"fork_repository"
"get_commit"
"get_file_contents"
"get_label"
"get_latest_release"
"get_me"
"get_release_by_tag"
"get_tag"
"get_team_members"
"get_teams"
"issue_read"
"issue_write"
"list_branches"
"list_commits"
"list_issue_types"
"list_issues"
"list_pull_requests"
"list_releases"
"list_tags"
"merge_pull_request"
"pull_request_read"
"pull_request_review_write"
"push_files"
"request_copilot_review"
"search_code"
"search_issues"
"search_pull_requests"
"search_repositories"
"search_users"
"sub_issue_write"
"update_pull_request_branch"
"update_pull_request"
*/

func New(ctx context.Context, githubToken string, included, excluded []string) (tool.Toolset, error) {
	transport := &mcp.StreamableClientTransport{
		Endpoint:   githubMCPEndpoint,
		HTTPClient: oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})),
	}

	m := make(map[string]bool, len(included)+len(excluded))
	for _, i := range included {
		m[i] = true
	}
	for _, e := range excluded {
		m[e] = false
	}
	filter := func(ctx agent.ReadonlyContext, tool tool.Tool) bool {
		t, present := m[tool.Name()]
		if !present {
			// If there are included tools specified, then tools not explicitly
			// included are excluded.
			if len(included) > 0 {
				return false
			}
			return true
		}
		return t
	}

	return mcptoolset.New(mcptoolset.Config{
		Transport:  transport,
		ToolFilter: filter,
	})
}
