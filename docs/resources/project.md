---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sentry_project Resource - terraform-provider-sentry"
subcategory: ""
description: |-
  Sentry Project resource.
---

# sentry_project (Resource)

Sentry Project resource.

## Example Usage

```terraform
# Create a project
resource "sentry_project" "default" {
  organization = "my-organization"

  teams = ["my-first-team", "my-second-team"]
  name  = "Web App"
  slug  = "web-app"

  platform    = "javascript"
  resolve_age = 720

  default_rules = false

  filters = {
    blacklisted_ips = ["127.0.0.1", "0.0.0.0/8"]
    releases        = ["1.*", "[!3].[0-9].*"]
    error_messages  = ["TypeError*", "*: integer division or modulo by zero"]
  }

  fingerprinting_rules  = <<-EOT
    # force all errors of the same type to have the same fingerprint
    error.type:DatabaseUnavailable -> system-down
    # force all memory allocation errors to be grouped together
    stack.function:malloc -> memory-allocation-error
  EOT
  grouping_enhancements = <<-EOT
    # remove all frames above a certain function from grouping
    stack.function:panic_handler ^-group
    # mark all functions following a prefix in-app
    stack.function:mylibrary_* +app
  EOT
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name for the project.
- `organization` (String) The organization of this resource.
- `teams` (Set of String) The slugs of the teams to create the project for.

### Optional

- `client_security` (Attributes) Configure origin URLs which Sentry should accept events from. This is used for communication with clients like [sentry-javascript](https://github.com/getsentry/sentry-javascript). (see [below for nested schema](#nestedatt--client_security))
- `default_key` (Boolean) Whether to create a default key on project creation. By default, Sentry will create a key for you. If you wish to manage keys manually, set this to false and create keys using the `sentry_key` resource. Note that this only takes effect on project creation, not on project update.
- `default_rules` (Boolean) Whether to create a default issue alert. Defaults to true where the behavior is to alert the user on every new issue.
- `digests_max_delay` (Number) The maximum amount of time (in seconds) to wait between scheduling digests for delivery.
- `digests_min_delay` (Number) The minimum amount of time (in seconds) to wait between scheduling digests for delivery after the initial scheduling.
- `filters` (Attributes) Custom filters for this project. (see [below for nested schema](#nestedatt--filters))
- `fingerprinting_rules` (String) This can be used to modify the fingerprint rules on the server with custom rules. Rules follow the pattern `matcher:glob -> fingerprint, values`. To learn more about fingerprint rules, [read the docs](https://docs.sentry.io/concepts/data-management/event-grouping/fingerprint-rules/).
- `grouping_enhancements` (String) This can be used to enhance the grouping algorithm with custom rules. Rules follow the pattern `matcher:glob [v^]?[+-]flag`. To learn more about stack trace rules, [read the docs](https://docs.sentry.io/concepts/data-management/event-grouping/stack-trace-rules/).
- `platform` (String) The platform for this project. Use `other` for platforms not listed. Valid values are: `other`, `android`, `apple`, `apple-ios`, `apple-macos`, `bun`, `capacitor`, `cordova`, `dart`, `deno`, `dotnet`, `dotnet-aspnet`, `dotnet-aspnetcore`, `dotnet-awslambda`, `dotnet-gcpfunctions`, `dotnet-maui`, `dotnet-uwp`, `dotnet-winforms`, `dotnet-wpf`, `dotnet-xamarin`, `electron`, `elixir`, `flutter`, `go`, `go-echo`, `go-fasthttp`, `go-fiber`, `go-gin`, `go-http`, `go-iris`, `go-martini`, `go-negroni`, `godot`, `ionic`, `java`, `java-log4j2`, `java-logback`, `java-spring`, `java-spring-boot`, `javascript`, `javascript-angular`, `javascript-astro`, `javascript-ember`, `javascript-gatsby`, `javascript-nextjs`, `javascript-nuxt`, `javascript-react`, `javascript-react-router`, `javascript-remix`, `javascript-solid`, `javascript-solidstart`, `javascript-svelte`, `javascript-sveltekit`, `javascript-tanstackstart-react`, `javascript-vue`, `kotlin`, `minidump`, `native`, `native-qt`, `nintendo-switch`, `nintendo-switch-2`, `node`, `node-awslambda`, `node-azurefunctions`, `node-cloudflare-pages`, `node-cloudflare-workers`, `node-connect`, `node-express`, `node-fastify`, `node-gcpfunctions`, `node-hapi`, `node-koa`, `node-nestjs`, `php`, `php-laravel`, `php-symfony`, `playstation`, `powershell`, `python`, `python-aiohttp`, `python-asgi`, `python-awslambda`, `python-bottle`, `python-celery`, `python-chalice`, `python-django`, `python-falcon`, `python-fastapi`, `python-flask`, `python-gcpfunctions`, `python-pylons`, `python-pymongo`, `python-pyramid`, `python-quart`, `python-rq`, `python-sanic`, `python-serverless`, `python-starlette`, `python-tornado`, `python-tryton`, `python-wsgi`, `react-native`, `ruby`, `ruby-rack`, `ruby-rails`, `rust`, `unity`, `unreal`, and `xbox`.
- `resolve_age` (Number) Hours in which an issue is automatically resolve if not seen after this amount of time.
- `slug` (String) The optional slug for this project.

### Read-Only

- `features` (Set of String)
- `id` (String) The ID of this resource.
- `internal_id` (String) The internal ID for this project.

<a id="nestedatt--client_security"></a>
### Nested Schema for `client_security`

Optional:

- `allowed_domains` (Set of String) A list of allowed domains. Examples: https://example.com, *, *.example.com, *:80.
- `scrape_javascript` (Boolean) Enable JavaScript source fetching. Allow Sentry to scrape missing JavaScript source context when possible.
- `security_token` (String) Security Token. Outbound requests matching Allowed Domains will have the header "{security_token_header}: {security_token}" appended.
- `security_token_header` (String) Security Token Header. Outbound requests matching Allowed Domains will have the header "{security_token_header}: {security_token}" appended.
- `verify_tls_ssl` (Boolean) Verify TLS/SSL. Outbound requests will verify TLS (sometimes known as SSL) connections.


<a id="nestedatt--filters"></a>
### Nested Schema for `filters`

Optional:

- `blacklisted_ips` (Set of String) Filter events from these IP addresses. (e.g. 127.0.0.1 or 10.0.0.0/8)
- `error_messages` (Set of String) Filter events by error messages. Allows [glob pattern matching](https://en.wikipedia.org/wiki/Glob_(programming)). (e.g. TypeError* or *: integer division or modulo by zero)
- `releases` (Set of String) Filter events from these releases. Allows [glob pattern matching](https://en.wikipedia.org/wiki/Glob_(programming)). (e.g. 1.* or [!3].[0-9].*)

## Import

Import is supported using the following syntax:

```shell
# import using the organization and team slugs from the URL:
# https://sentry.io/settings/[org-slug]/projects/[project-slug]/
terraform import sentry_project.default org-slug/project-slug
```
