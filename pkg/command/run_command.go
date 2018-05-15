package command

import (
	"strings"

	"github.com/arquillian/ike-prow-plugins/pkg/github/client"
	"github.com/arquillian/ike-prow-plugins/pkg/log"
	gogh "github.com/google/go-github/github"
)

// RunCommentPrefix is used as a command to bypass test presence validation
const RunCommentPrefix = "/run"

// RunCmd represents a command that is triggered by "/run"
type RunCmd struct {
	UserPermissionService *PermissionService
	WhenAddedOrCreated    DoFunction
}

// Perform executes the set DoFunctions for the given IssueCommentEvent (when all conditions are fulfilled)
func (c *RunCmd) Perform(client ghclient.Client, log log.Logger, comment *gogh.IssueCommentEvent) error {
	user := c.UserPermissionService
	var RunCommand = &CmdExecutor{Command: RunCommentPrefix}

	RunCommand.
		When(Triggered).
		By(AnyOf(user.Admin, user.PRReviewer, user.PRApprover, user.PRCreator)).
		Then(c.WhenAddedOrCreated)

	return RunCommand.Execute(client, log, comment)
}

// Matches returns true when the given IssueCommentEvent content prefix is "/run"
func (c *RunCmd) Matches(comment *gogh.IssueCommentEvent) bool {
	body := strings.TrimSpace(*comment.Comment.Body)
	return strings.HasPrefix(body, RunCommentPrefix)
}

// ContainsRunCmdWithPluginNameOrAll returns true when the given IssueCommentEvent content contains "/run plugin-name" or "/run all"
func (c *RunCmd) ContainsRunCmdWithPluginNameOrAll(pluginName string, comment *gogh.IssueCommentEvent) bool {
	body := strings.TrimSpace(*comment.Comment.Body)
	return strings.HasPrefix(body, RunCommentPrefix) && (strings.Contains(body, pluginName) || strings.Contains(body, "all"))
}
