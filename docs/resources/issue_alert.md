---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sentry_issue_alert Resource - terraform-provider-sentry"
subcategory: ""
description: |-
  Create an Issue Alert Rule for a Project. See the Sentry Documentation https://docs.sentry.io/api/alerts/create-an-issue-alert-rule-for-a-project/ for more information.
  NOTE: Since v0.15.0, the conditions, filters, and actions attributes which are JSON strings have been deprecated in favor of conditions_v2, filters_v2, and actions_v2 which are lists of objects.
---

# sentry_issue_alert (Resource)

Create an Issue Alert Rule for a Project. See the [Sentry Documentation](https://docs.sentry.io/api/alerts/create-an-issue-alert-rule-for-a-project/) for more information.

**NOTE:** Since v0.15.0, the `conditions`, `filters`, and `actions` attributes which are JSON strings have been deprecated in favor of `conditions_v2`, `filters_v2`, and `actions_v2` which are lists of objects.

## Example Usage

```terraform
resource "sentry_issue_alert" "main" {
  organization = sentry_project.main.organization
  project      = sentry_project.main.id
  name         = "My issue alert"

  action_match = "any"
  filter_match = "any"
  frequency    = 30

  conditions_v2 = [
    { first_seen_event = {} },
    { regression_event = {} },
    { reappeared_event = {} },
    { new_high_priority_issue = {} },
    { existing_high_priority_issue = {} },
    {
      event_frequency = {
        comparison_type = "count"
        value           = 100
        interval        = "1h"
      }
    },
    {
      event_frequency = {
        comparison_type     = "percent"
        comparison_interval = "1w"
        value               = 100
        interval            = "1h"
      }
    },
    {
      event_unique_user_frequency = {
        comparison_type = "count"
        value           = 100
        interval        = "1h"
      }
    },
    {
      event_unique_user_frequency = {
        comparison_type     = "percent"
        comparison_interval = "1w"
        value               = 100
        interval            = "1h"
      }
    },
    {
      event_frequency_percent = {
        comparison_type = "count"
        value           = 100
        interval        = "1h"
      }
    },
    {
      event_frequency_percent = {
        comparison_type     = "percent"
        comparison_interval = "1w"
        value               = 100
        interval            = "1h"
      }
    },
  ]

  filters_v2 = [
    {
      age_comparison = {
        comparison_type = "older"
        value           = 10
        time            = "minute"
      }
    },
    {
      issue_occurrences = {
        value = 10
      }
    },
    {
      assigned_to = {
        target_type = "Unassigned"
      }
    },
    {
      assigned_to = {
        target_type       = "Team"
        target_identifier = sentry_team.test.internal_id // Note: This is the internal ID of the team rather than the slug
      }
    },
    {
      latest_adopted_release = {
        oldest_or_newest = "oldest"
        older_or_newer   = "older"
        environment      = "test"
      }
    },
    { latest_release = {} },
    {
      issue_category = {
        value = "Error"
      }
    },
    {
      event_attribute = {
        attribute = "message"
        match     = "CONTAINS"
        value     = "test"
      }
    },
    {
      event_attribute = {
        attribute = "message"
        match     = "IS_SET"
      }
    },
    {
      tagged_event = {
        key   = "key"
        match = "CONTAINS"
        value = "value"
      }
    },
    {
      tagged_event = {
        key   = "key"
        match = "NOT_SET"
      }
    },
    {
      level = {
        match = "EQUAL"
        level = "error"
      }
    },
  ]

  actions_v2 = [/* Please see below for examples */]

}

#
# Send a notification to Suggested Assignees
#

resource "sentry_issue_alert" "member_alert" {
  actions_v2 = [
    {
      notify_email = {
        target_type      = "IssueOwners"
        fallthrough_type = "ActiveMembers"
      }
    },
  ]
  // ...
}

#
# Send a notification to a Member
#

data "sentry_organization_member" "member" {
  organization = data.sentry_organization.test.id
  email        = "test@example.com"
}

resource "sentry_issue_alert" "member_alert" {
  actions_v2 = [
    {
      notify_email = {
        target_type       = "Member"
        target_identifier = data.sentry_organization_member.member.internal_id
        fallthrough_type  = "AllMembers"
      }
    },
  ]
  // ...
}

#
# Send a notification to a Team
#

data "sentry_team" "team" {
  organization = sentry_project.test.organization
  slug         = "my-team"
}

resource "sentry_issue_alert" "team_alert" {
  actions_v2 = [
    {
      notify_email = {
        target_type       = "Team"
        target_identifier = data.sentry_team.team.internal_id
        fallthrough_type  = "AllMembers"
      }
    },
  ]
  // ...
}

#
# Send a Slack notification
#

# Retrieve a Slack integration
data "sentry_organization_integration" "slack" {
  organization = sentry_project.test.organization

  provider_key = "slack"
  name         = "Slack Workspace" # Name of your Slack workspace
}

resource "sentry_issue_alert" "slack_alert" {
  actions_v2 = [
    {
      slack_notify_service = {
        workspace = data.sentry_organization_integration.slack.id
        channel   = "#warning"
        tags      = ["environment", "level"]
        notes     = "Please <http://example.com|click here> for triage information"
      }
    },
  ]
  // ...
}

#
# Send a Microsoft Teams notification
#

# Retrieve a MS Teams integration
data "sentry_organization_integration" "msteams" {
  organization = sentry_project.test.organization

  provider_key = "msteams"
  name         = "My Team" # Name of your Microsoft Teams team
}

resource "sentry_issue_alert" "slack_alert" {
  actions_v2 = [
    {
      msteams_notify_service = {
        team    = data.sentry_organization_integration.msteams.id
        channel = "General"
      }
    },
  ]
  // ...
}

#
# Send a Discord notification
#

data "sentry_organization_integration" "discord" {
  organization = sentry_project.test.organization

  provider_key = "discord"
  name         = "Discord Server" # Name of your Discord server
}

resource "sentry_issue_alert" "discord_alert" {
  actions_v2 = [
    {
      discord_notify_service = {
        server     = data.sentry_organization_integration.discord.id
        channel_id = "94732897"
        tags       = ["browser", "user"]
      }
    },
  ]
  // ...
}

#
# Create a Jira Ticket
#

data "sentry_organization_integration" "jira" {
  organization = sentry_project.test.organization

  provider_key = "jira"
  name         = "JIRA" # Name of your Jira server
}

resource "sentry_issue_alert" "jira_alert" {
  actions_v2 = [
    {
      jira_create_ticket = {
        integration = data.sentry_organization_integration.jira.id
        project     = "349719"
        issue_type  = "1"
      }
    },
  ]
  // ...
}

#
# Create a Jira Server Ticket
#

data "sentry_organization_integration" "jira_server" {
  organization = sentry_project.test.organization

  provider_key = "jira_server"
  name         = "JIRA" # Name of your Jira server
}

# TODO
resource "sentry_issue_alert" "jira_server_alert" {
  actions_v2 = [
    {
      jira_server_create_ticket = {
        integration = data.sentry_organization_integration.jira_server.id
        project     = "349719"
        issue_type  = "1"
      }
    },
  ]
  // ...
}

#
# Create a GitHub Issue
#

data "sentry_organization_integration" "github" {
  organization = sentry_project.test.organization

  provider_key = "github"
  name         = "GitHub"
}

resource "sentry_issue_alert" "github_alert" {
  actions_v2 = [
    {
      github_create_ticket = {
        integration = data.sentry_organization_integration.github.id
        repo        = "default"
        assignee    = "Baxter the Hacker"
        labels      = ["bug", "p1"]
      }
    },
  ]
  // ...
}

#
# Create an Azure DevOps work item
#

data "sentry_organization_integration" "vsts" {
  organization = sentry_project.test.organization

  provider_key = "vsts"
  name         = "Azure DevOps"
}

resource "sentry_issue_alert" "vsts_alert" {
  actions_v2 = [
    {
      azure_devops_create_ticket = {
        integration    = data.sentry_organization_integration.vsts.id
        project        = "0389485"
        work_item_type = "Microsoft.VSTS.WorkItemTypes.Task"
      }
    },
  ]
  // ...
}

#
# Send a PagerDuty notification
#

data "sentry_organization_integration" "pagerduty" {
  organization = sentry_project.test.organization
  provider_key = "pagerduty"
  name         = "PagerDuty"
}

resource "sentry_integration_pagerduty" "pagerduty" {
  organization    = data.sentry_organization_integration.pagerduty.organization
  integration_id  = data.sentry_organization_integration.pagerduty.id
  service         = "issue-alert-service"
  integration_key = "issue-alert-integration-key"
}

resource "sentry_issue_alert" "pagerduty_alert" {
  actions_v2 = [
    {
      pagerduty_notify_service = {
        account  = sentry_integration_pagerduty.pagerduty.integration_id
        service  = sentry_integration_pagerduty.pagerduty.id
        severity = "default"
      }
    },
  ]
  // ...
}

#
# Send an Opsgenie notification
#

data "sentry_organization_integration" "opsgenie" {
  organization = sentry_project.test.organization
  provider_key = "opsgenie"
  name         = "Opsgenie"
}

resource "sentry_integration_opsgenie" "opsgenie" {
  organization    = data.sentry_organization_integration.opsgenie.organization
  integration_id  = data.sentry_organization_integration.opsgenie.id
  team            = "issue-alert-team"
  integration_key = "my-integration-key"
}

resource "sentry_issue_alert" "opsgenie_alert" {
  actions_v2 = [
    {
      opsgenie_notify_team = {
        account  = sentry_integration_opsgenie.opsgenie.integration_id
        team     = sentry_integration_opsgenie.opsgenie.id
        priority = "P1"
      }
    },
  ]
  // ...
}

#
# Send a notification via an integration
#

resource "sentry_issue_alert" "notification_alert" {
  actions_v2 = [
    {
      notify_event_service = {
        # Sourced from: https://terraform-provider-sentry.sentry.io/settings/developer-settings/<service>/
        service = "my-service"
      }
    },
  ]
  // ...
}

#
# Send a notification to a Sentry app
#

resource "sentry_issue_alert" "sentry_app" {
  actions_v2 = [
    {
      notify_event_sentry_app = {
        sentry_app_installation_uuid = "my-sentry-app-installation-uuid"
        settings = {
          key1 = "value1"
          key2 = "value2"
          key3 = "value3"
        }
      }
    },
  ]
  // ...
}

#
# Send a notification (for all legacy integrations)
#

resource "sentry_issue_alert" "notification_alert" {
  actions_v2 = [
    { notify_event = {} },
  ]
  // ...
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action_match` (String) Trigger actions when an event is captured by Sentry and `any` or `all` of the specified conditions happen. Valid values are: `all`, and `any`.
- `frequency` (Number) Perform actions at most once every `X` minutes for this issue.
- `name` (String) The issue alert name.
- `organization` (String) The organization of this resource.
- `project` (String) The project of this resource.

### Optional

- `actions` (String, Deprecated) **Deprecated** in favor of `actions_v2`. A list of actions that take place when all required conditions and filters for the rule are met. In JSON string format.
- `actions_v2` (Attributes List) A list of actions that take place when all required conditions and filters for the rule are met. (see [below for nested schema](#nestedatt--actions_v2))
- `conditions` (String, Deprecated) **Deprecated** in favor of `conditions_v2`. A list of triggers that determine when the rule fires. In JSON string format.
- `conditions_v2` (Attributes List) A list of triggers that determine when the rule fires. (see [below for nested schema](#nestedatt--conditions_v2))
- `environment` (String) Perform issue alert in a specific environment.
- `filter_match` (String) A string determining which filters need to be true before any actions take place. Required when a value is provided for `filters`. Valid values are: `all`, `any`, and `none`.
- `filters` (String, Deprecated) **Deprecated** in favor of `filters_v2`. A list of filters that determine if a rule fires after the necessary conditions have been met. In JSON string format.
- `filters_v2` (Attributes List) A list of filters that determine if a rule fires after the necessary conditions have been met. (see [below for nested schema](#nestedatt--filters_v2))
- `owner` (String) The ID of the team or user that owns the rule.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--actions_v2"></a>
### Nested Schema for `actions_v2`

Optional:

- `azure_devops_create_ticket` (Attributes) Create an Azure DevOps work item in `integration`. (see [below for nested schema](#nestedatt--actions_v2--azure_devops_create_ticket))
- `discord_notify_service` (Attributes) Send a notification to the `server` Discord server in the channel with ID or URL: `channel_id` and show tags `tags` in the notification. (see [below for nested schema](#nestedatt--actions_v2--discord_notify_service))
- `github_create_ticket` (Attributes) Create a GitHub issue in `integration`. (see [below for nested schema](#nestedatt--actions_v2--github_create_ticket))
- `github_enterprise_create_ticket` (Attributes) Create a GitHub Enterprise issue in `integration`. (see [below for nested schema](#nestedatt--actions_v2--github_enterprise_create_ticket))
- `jira_create_ticket` (Attributes) Create a Jira issue in `integration`. (see [below for nested schema](#nestedatt--actions_v2--jira_create_ticket))
- `jira_server_create_ticket` (Attributes) Create a Jira Server issue in `integration`. (see [below for nested schema](#nestedatt--actions_v2--jira_server_create_ticket))
- `msteams_notify_service` (Attributes) Send a notification to the `team` Team to `channel`. (see [below for nested schema](#nestedatt--actions_v2--msteams_notify_service))
- `notify_email` (Attributes) Send a notification to `target_type` and if none can be found then send a notification to `fallthrough_type`. (see [below for nested schema](#nestedatt--actions_v2--notify_email))
- `notify_event` (Attributes) Send a notification to all legacy integrations. (see [below for nested schema](#nestedatt--actions_v2--notify_event))
- `notify_event_sentry_app` (Attributes) Send a notification to a Sentry app. (see [below for nested schema](#nestedatt--actions_v2--notify_event_sentry_app))
- `notify_event_service` (Attributes) Send a notification via an integration. (see [below for nested schema](#nestedatt--actions_v2--notify_event_service))
- `opsgenie_notify_team` (Attributes) Send a notification to Opsgenie account `account` and team `team` with `priority` priority. (see [below for nested schema](#nestedatt--actions_v2--opsgenie_notify_team))
- `pagerduty_notify_service` (Attributes) Send a notification to PagerDuty account `account` and service `service` with `severity` severity. (see [below for nested schema](#nestedatt--actions_v2--pagerduty_notify_service))
- `slack_notify_service` (Attributes) Send a notification to the `workspace` Slack workspace to `channel` (optionally, an ID: `channel_id`) and show tags `tags` and notes `notes` in notification. (see [below for nested schema](#nestedatt--actions_v2--slack_notify_service))

<a id="nestedatt--actions_v2--azure_devops_create_ticket"></a>
### Nested Schema for `actions_v2.azure_devops_create_ticket`

Required:

- `integration` (String) The integration ID.
- `project` (String) The ID of the Azure DevOps project.
- `work_item_type` (String) The type of work item to create.

Read-Only:

- `name` (String)


<a id="nestedatt--actions_v2--discord_notify_service"></a>
### Nested Schema for `actions_v2.discord_notify_service`

Required:

- `channel_id` (String) The ID of the channel to send the notification to. You must enter either a channel ID or a channel URL, not a channel name
- `server` (String) The integration ID associated with the Discord server.

Optional:

- `tags` (Set of String) A string of tags to show in the notification.

Read-Only:

- `name` (String)


<a id="nestedatt--actions_v2--github_create_ticket"></a>
### Nested Schema for `actions_v2.github_create_ticket`

Required:

- `integration` (String) The integration ID associated with GitHub.
- `repo` (String) The name of the repository to create the issue in.

Optional:

- `assignee` (String) The GitHub user to assign the issue to.
- `labels` (Set of String) A list of labels to assign to the issue.

Read-Only:

- `name` (String)


<a id="nestedatt--actions_v2--github_enterprise_create_ticket"></a>
### Nested Schema for `actions_v2.github_enterprise_create_ticket`

Required:

- `integration` (String) The integration ID associated with GitHub Enterprise.
- `repo` (String) The name of the repository to create the issue in.

Optional:

- `assignee` (String) The GitHub user to assign the issue to.
- `labels` (Set of String) A list of labels to assign to the issue.

Read-Only:

- `name` (String)


<a id="nestedatt--actions_v2--jira_create_ticket"></a>
### Nested Schema for `actions_v2.jira_create_ticket`

Required:

- `integration` (String) The integration ID associated with Jira.
- `issue_type` (String) The ID of the type of issue that the ticket should be created as.
- `project` (String) The ID of the Jira project.

Read-Only:

- `name` (String)


<a id="nestedatt--actions_v2--jira_server_create_ticket"></a>
### Nested Schema for `actions_v2.jira_server_create_ticket`

Required:

- `integration` (String) The integration ID associated with Jira Server.
- `issue_type` (String) The ID of the type of issue that the ticket should be created as.
- `project` (String) The ID of the Jira Server project.

Read-Only:

- `name` (String)


<a id="nestedatt--actions_v2--msteams_notify_service"></a>
### Nested Schema for `actions_v2.msteams_notify_service`

Required:

- `channel` (String) The name of the channel to send the notification to.
- `team` (String) The integration ID associated with the Microsoft Teams team.

Read-Only:

- `channel_id` (String)
- `name` (String)


<a id="nestedatt--actions_v2--notify_email"></a>
### Nested Schema for `actions_v2.notify_email`

Required:

- `target_type` (String) Valid values are: `IssueOwners`, `Team`, and `Member`.

Optional:

- `fallthrough_type` (String) Who the notification should be sent to if there are no suggested assignees. Valid values are: `AllMembers`, `ActiveMembers`, and `NoOne`.
- `target_identifier` (String) The ID of the Member or Team the notification should be sent to. Only required when `target_type` is `Team` or `Member`.

Read-Only:

- `name` (String)


<a id="nestedatt--actions_v2--notify_event"></a>
### Nested Schema for `actions_v2.notify_event`

Read-Only:

- `name` (String)


<a id="nestedatt--actions_v2--notify_event_sentry_app"></a>
### Nested Schema for `actions_v2.notify_event_sentry_app`

Required:

- `sentry_app_installation_uuid` (String)

Optional:

- `settings` (Map of String)

Read-Only:

- `name` (String)


<a id="nestedatt--actions_v2--notify_event_service"></a>
### Nested Schema for `actions_v2.notify_event_service`

Required:

- `service` (String) The slug of the integration service. Sourced from `https://terraform-provider-sentry.sentry.io/settings/developer-settings/<service>/`.

Read-Only:

- `name` (String)


<a id="nestedatt--actions_v2--opsgenie_notify_team"></a>
### Nested Schema for `actions_v2.opsgenie_notify_team`

Required:

- `account` (String)
- `priority` (String)
- `team` (String)

Read-Only:

- `name` (String)


<a id="nestedatt--actions_v2--pagerduty_notify_service"></a>
### Nested Schema for `actions_v2.pagerduty_notify_service`

Required:

- `account` (String)
- `service` (String)
- `severity` (String)

Read-Only:

- `name` (String)


<a id="nestedatt--actions_v2--slack_notify_service"></a>
### Nested Schema for `actions_v2.slack_notify_service`

Required:

- `channel` (String) The name of the channel to send the notification to (e.g., #critical, Jane Schmidt).
- `workspace` (String) The integration ID associated with the Slack workspace.

Optional:

- `notes` (String) Text to show alongside the notification. To @ a user, include their user id like `@<USER_ID>`. To include a clickable link, format the link and title like `<http://example.com|Click Here>`.
- `tags` (Set of String) A string of tags to show in the notification.

Read-Only:

- `channel_id` (String) The ID of the channel to send the notification to.
- `name` (String)



<a id="nestedatt--conditions_v2"></a>
### Nested Schema for `conditions_v2`

Optional:

- `event_frequency` (Attributes) When the `comparison_type` is `count`, the number of events in an issue is more than `value` in `interval`. When the `comparison_type` is `percent`, the number of events in an issue is `value` % higher in `interval` compared to `comparison_interval` ago. (see [below for nested schema](#nestedatt--conditions_v2--event_frequency))
- `event_frequency_percent` (Attributes) When the `comparison_type` is `count`, the percent of sessions affected by an issue is more than `value` in `interval`. When the `comparison_type` is `percent`, the percent of sessions affected by an issue is `value` % higher in `interval` compared to `comparison_interval` ago. (see [below for nested schema](#nestedatt--conditions_v2--event_frequency_percent))
- `event_unique_user_frequency` (Attributes) When the `comparison_type` is `count`, the number of users affected by an issue is more than `value` in `interval`. When the `comparison_type` is `percent`, the number of users affected by an issue is `value` % higher in `interval` compared to `comparison_interval` ago. (see [below for nested schema](#nestedatt--conditions_v2--event_unique_user_frequency))
- `existing_high_priority_issue` (Attributes) Sentry marks an existing issue as high priority. (see [below for nested schema](#nestedatt--conditions_v2--existing_high_priority_issue))
- `first_seen_event` (Attributes) A new issue is created. (see [below for nested schema](#nestedatt--conditions_v2--first_seen_event))
- `new_high_priority_issue` (Attributes) Sentry marks a new issue as high priority. (see [below for nested schema](#nestedatt--conditions_v2--new_high_priority_issue))
- `reappeared_event` (Attributes) The issue changes state from ignored to unresolved. (see [below for nested schema](#nestedatt--conditions_v2--reappeared_event))
- `regression_event` (Attributes) The issue changes state from resolved to unresolved. (see [below for nested schema](#nestedatt--conditions_v2--regression_event))

<a id="nestedatt--conditions_v2--event_frequency"></a>
### Nested Schema for `conditions_v2.event_frequency`

Required:

- `comparison_type` (String) Valid values are: `count`, and `percent`.
- `value` (Number)

Optional:

- `comparison_interval` (String) `m` for minutes, `h` for hours, `d` for days, and `w` for weeks. Valid values are: `5m`, `15m`, `1h`, `1d`, `1w`, and `30d`.
- `interval` (String) `m` for minutes, `h` for hours, `d` for days, and `w` for weeks. Valid values are: `1m`, `5m`, `15m`, `1h`, `1d`, `1w`, and `30d`.

Read-Only:

- `name` (String)


<a id="nestedatt--conditions_v2--event_frequency_percent"></a>
### Nested Schema for `conditions_v2.event_frequency_percent`

Required:

- `comparison_type` (String) Valid values are: `count`, and `percent`.
- `interval` (String) `m` for minutes, `h` for hours. Valid values are: `5m`, `10m`, `30m`, and `1h`.
- `value` (Number)

Optional:

- `comparison_interval` (String) `m` for minutes, `h` for hours, `d` for days, and `w` for weeks. Valid values are: `5m`, `15m`, `1h`, `1d`, `1w`, and `30d`.

Read-Only:

- `name` (String)


<a id="nestedatt--conditions_v2--event_unique_user_frequency"></a>
### Nested Schema for `conditions_v2.event_unique_user_frequency`

Required:

- `comparison_type` (String) Valid values are: `count`, and `percent`.
- `value` (Number)

Optional:

- `comparison_interval` (String) `m` for minutes, `h` for hours, `d` for days, and `w` for weeks. Valid values are: `5m`, `15m`, `1h`, `1d`, `1w`, and `30d`.
- `interval` (String) `m` for minutes, `h` for hours, `d` for days, and `w` for weeks. Valid values are: `1m`, `5m`, `15m`, `1h`, `1d`, `1w`, and `30d`.

Read-Only:

- `name` (String)


<a id="nestedatt--conditions_v2--existing_high_priority_issue"></a>
### Nested Schema for `conditions_v2.existing_high_priority_issue`

Read-Only:

- `name` (String)


<a id="nestedatt--conditions_v2--first_seen_event"></a>
### Nested Schema for `conditions_v2.first_seen_event`

Read-Only:

- `name` (String)


<a id="nestedatt--conditions_v2--new_high_priority_issue"></a>
### Nested Schema for `conditions_v2.new_high_priority_issue`

Read-Only:

- `name` (String)


<a id="nestedatt--conditions_v2--reappeared_event"></a>
### Nested Schema for `conditions_v2.reappeared_event`

Read-Only:

- `name` (String)


<a id="nestedatt--conditions_v2--regression_event"></a>
### Nested Schema for `conditions_v2.regression_event`

Read-Only:

- `name` (String)



<a id="nestedatt--filters_v2"></a>
### Nested Schema for `filters_v2`

Optional:

- `age_comparison` (Attributes) The issue is older or newer than `value` `time`. (see [below for nested schema](#nestedatt--filters_v2--age_comparison))
- `assigned_to` (Attributes) The issue is assigned to no one, team, or member. (see [below for nested schema](#nestedatt--filters_v2--assigned_to))
- `event_attribute` (Attributes) The event's `attribute` value `match` `value`. (see [below for nested schema](#nestedatt--filters_v2--event_attribute))
- `issue_category` (Attributes) The issue's category is equal to `value`. (see [below for nested schema](#nestedatt--filters_v2--issue_category))
- `issue_occurrences` (Attributes) The issue has happened at least `value` times (Note: this is approximate). (see [below for nested schema](#nestedatt--filters_v2--issue_occurrences))
- `latest_adopted_release` (Attributes) The {oldest_or_newest} adopted release associated with the event's issue is {older_or_newer} than the latest adopted release in {environment}. (see [below for nested schema](#nestedatt--filters_v2--latest_adopted_release))
- `latest_release` (Attributes) The event is from the latest release. (see [below for nested schema](#nestedatt--filters_v2--latest_release))
- `level` (Attributes) The event's level is `match` `level`. (see [below for nested schema](#nestedatt--filters_v2--level))
- `tagged_event` (Attributes) The event's tags match `key` `match` `value`. (see [below for nested schema](#nestedatt--filters_v2--tagged_event))

<a id="nestedatt--filters_v2--age_comparison"></a>
### Nested Schema for `filters_v2.age_comparison`

Required:

- `comparison_type` (String) Valid values are: `older`, and `newer`.
- `time` (String) Valid values are: `minute`, `hour`, `day`, and `week`.
- `value` (Number)

Read-Only:

- `name` (String)


<a id="nestedatt--filters_v2--assigned_to"></a>
### Nested Schema for `filters_v2.assigned_to`

Required:

- `target_type` (String) Valid values are: `Unassigned`, `Team`, and `Member`.

Optional:

- `target_identifier` (String) The target's ID. Only required when `target_type` is `Team` or `Member`.

Read-Only:

- `name` (String)


<a id="nestedatt--filters_v2--event_attribute"></a>
### Nested Schema for `filters_v2.event_attribute`

Required:

- `attribute` (String) Valid values are: `message`, `platform`, `environment`, `type`, `error.handled`, `error.unhandled`, `error.main_thread`, `exception.type`, `exception.value`, `user.id`, `user.email`, `user.username`, `user.ip_address`, `http.method`, `http.url`, `http.status_code`, `sdk.name`, `stacktrace.code`, `stacktrace.module`, `stacktrace.filename`, `stacktrace.abs_path`, `stacktrace.package`, `unreal.crash_type`, `app.in_foreground`, `os.distribution_name`, `os.distribution_version`, `symbolicated_in_app`, `ota_updates.channel`, `ota_updates.runtime_version`, and `ota_updates.update_id`.
- `match` (String) The comparison operator. Valid values are: `CONTAINS`, `ENDS_WITH`, `EQUAL`, `GREATER_OR_EQUAL`, `GREATER`, `IS_SET`, `IS_IN`, `LESS_OR_EQUAL`, `LESS`, `NOT_CONTAINS`, `NOT_ENDS_WITH`, `NOT_EQUAL`, `NOT_SET`, `NOT_STARTS_WITH`, `NOT_IN`, and `STARTS_WITH`.

Optional:

- `value` (String)

Read-Only:

- `name` (String)


<a id="nestedatt--filters_v2--issue_category"></a>
### Nested Schema for `filters_v2.issue_category`

Required:

- `value` (String) Valid values are: `Error`, `Performance`, `Profile`, `Cron`, `Replay`, `Feedback`, `Uptime`, `Metric_Alert`, `Test_Notification`, `Outage`, `Metric`, `Db_Query`, `Http_Client`, `Frontend`, and `Mobile`.

Read-Only:

- `name` (String)


<a id="nestedatt--filters_v2--issue_occurrences"></a>
### Nested Schema for `filters_v2.issue_occurrences`

Required:

- `value` (Number)

Read-Only:

- `name` (String)


<a id="nestedatt--filters_v2--latest_adopted_release"></a>
### Nested Schema for `filters_v2.latest_adopted_release`

Required:

- `environment` (String)
- `older_or_newer` (String) Valid values are: `older`, and `newer`.
- `oldest_or_newest` (String) Valid values are: `oldest`, and `newest`.

Read-Only:

- `name` (String)


<a id="nestedatt--filters_v2--latest_release"></a>
### Nested Schema for `filters_v2.latest_release`

Read-Only:

- `name` (String)


<a id="nestedatt--filters_v2--level"></a>
### Nested Schema for `filters_v2.level`

Required:

- `level` (String) Valid values are: `sample`, `debug`, `info`, `warning`, `error`, and `fatal`.
- `match` (String) The comparison operator. Valid values are: `EQUAL`, `GREATER_OR_EQUAL`, and `LESS_OR_EQUAL`.

Read-Only:

- `name` (String)


<a id="nestedatt--filters_v2--tagged_event"></a>
### Nested Schema for `filters_v2.tagged_event`

Required:

- `key` (String) The tag.
- `match` (String) The comparison operator. Valid values are: `CONTAINS`, `ENDS_WITH`, `EQUAL`, `GREATER_OR_EQUAL`, `GREATER`, `IS_SET`, `IS_IN`, `LESS_OR_EQUAL`, `LESS`, `NOT_CONTAINS`, `NOT_ENDS_WITH`, `NOT_EQUAL`, `NOT_SET`, `NOT_STARTS_WITH`, `NOT_IN`, and `STARTS_WITH`.

Optional:

- `value` (String)

Read-Only:

- `name` (String)

## Import

Import is supported using the following syntax:

```shell
# import using the organization, project slugs and rule id from the URL:
# https://sentry.io/organizations/[org-slug]/alerts/rules/[project-slug]/[rule-id]/details/
terraform import sentry_issue_alert.default org-slug/project-slug/rule-id
```
