package providerdata

import (
	"github.com/mzglinski/go-sentry/v2/sentry"
	"github.com/mzglinski/terraform-provider-sentry/internal/apiclient"
)

type ProviderData struct {
	Client    *sentry.Client
	ApiClient *apiclient.ClientWithResponses
}
