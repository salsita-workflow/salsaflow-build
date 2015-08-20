package modules

import (
	"jira"

	"github.com/salsaflow/salsaflow/modules/issue_trackers/pivotaltracker"
)

var issueTrackerFactories = map[string]IssueTrackerFactory{
	jira.Id:           jira.Factory,
	pivotaltracker.Id: pivotaltracker.Factory,
}
