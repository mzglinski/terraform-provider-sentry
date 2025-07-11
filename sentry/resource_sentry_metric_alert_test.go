package sentry

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mzglinski/go-sentry/v2/sentry"
	"github.com/mzglinski/terraform-provider-sentry/internal/acctest"
	"github.com/mzglinski/terraform-provider-sentry/internal/tfutils"
)

func TestAccSentryMetricAlert_basic(t *testing.T) {
	teamName := acctest.RandomWithPrefix("tf-team")
	projectName := acctest.RandomWithPrefix("tf-project")
	alertName := acctest.RandomWithPrefix("tf-issue-alert")
	rn := "sentry_metric_alert.test"

	var alertID string

	check := func(alertName string) resource.TestCheckFunc {
		return resource.ComposeTestCheckFunc(
			testAccCheckSentryMetricAlertExists(rn, &alertID),
			resource.TestCheckResourceAttr(rn, "organization", acctest.TestOrganization),
			resource.TestCheckResourceAttr(rn, "project", projectName),
			resource.TestCheckResourceAttr(rn, "name", alertName),
			resource.TestCheckResourceAttr(rn, "environment", ""),
			resource.TestCheckResourceAttr(rn, "dataset", "generic_metrics"),
			resource.TestCheckResourceAttr(rn, "event_types.#", "1"),
			resource.TestCheckResourceAttr(rn, "event_types.0", "transaction"),
			resource.TestCheckResourceAttr(rn, "query", "http.url:http://testservice.com/stats"),
			resource.TestCheckResourceAttr(rn, "aggregate", "p50(transaction.duration)"),
			resource.TestCheckResourceAttr(rn, "time_window", "50"),
			resource.TestCheckResourceAttr(rn, "threshold_type", "0"),
			resource.TestCheckResourceAttr(rn, "resolve_threshold", "100"),
			resource.TestCheckResourceAttr(rn, "comparison_delta", "100"),
			resource.TestCheckResourceAttrPtr(rn, "internal_id", &alertID),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckSentryMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryMetricAlertConfig(teamName, projectName, alertName),
				Check:  check(alertName),
			},
			{
				Config: testAccSentryMetricAlertConfig(teamName, projectName, alertName+"-renamed"),
				Check:  check(alertName + "-renamed"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSentryMetricAlertDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sentry_metric_alert" {
			continue
		}

		org, project, id, err := tfutils.SplitThreePartId(rs.Primary.ID, "organization-slug", "project-slug", "alert-id")
		if err != nil {
			return err
		}

		ctx := context.Background()
		alert, resp, err := acctest.SharedClient.MetricAlerts.Get(ctx, org, project, id)
		if err == nil {
			if alert != nil {
				return errors.New("metric alert still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}

	return nil
}

func testAccCheckSentryMetricAlertExists(n string, gotAlertID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		org, project, alertID, err := tfutils.SplitThreePartId(rs.Primary.ID, "organization-slug", "project-slug", "alert-id")
		if err != nil {
			return err
		}
		ctx := context.Background()
		gotAlert, _, err := acctest.SharedClient.MetricAlerts.Get(ctx, org, project, alertID)
		if err != nil {
			return err
		}
		*gotAlertID = sentry.StringValue(gotAlert.ID)
		return nil
	}
}

func testAccSentryMetricAlertConfig(teamName, projectName, alertName string) string {
	return testAccSentryProjectConfig_team(teamName, projectName) + fmt.Sprintf(`
resource "sentry_metric_alert" "test" {
	organization      = sentry_project.test.organization
	project           = sentry_project.test.id
	name              = "%[1]s"
	dataset           = "generic_metrics"
	event_types       = ["transaction"]
	query             = "http.url:http://testservice.com/stats"
	aggregate         = "p50(transaction.duration)"
	time_window       = 50.0
	threshold_type    = 0
	resolve_threshold = 100.0
	comparison_delta  = 100.0

	trigger {
		action {
			type              = "email"
			target_type       = "team"
			target_identifier = sentry_team.test.internal_id
		}

		alert_threshold   = 1000
		label             = "critical"
		resolve_threshold = 100.0
		threshold_type    = 0
	}

	trigger {
		action {
			type              = "email"
			target_type       = "team"
			target_identifier = sentry_team.test.internal_id
		}
	
		alert_threshold   = 500
		label             = "warning"
		resolve_threshold = 100.0
		threshold_type    = 0
	}
}
	`, alertName)
}
