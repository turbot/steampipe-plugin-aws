package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go/middleware"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

// hydrateRateLimitMiddlewareID is the ID of the middleware that waits for a
// hydrate call's rate limiters before each AWS API request.
const hydrateRateLimitMiddlewareID = "SteampipeHydrateRateLimit"

// addHydrateRateLimitMiddleware makes every service client built from cfg wait
// for the rate limiters of the hydrate call that sends a request, as the plugin
// sets Plugin.ClientWaitsForHydrateRateLimit. A hydrate call then takes its
// token at its first request, and a call that sends no request does not wait,
// e.g., create_volume_permissions in aws_ebs_snapshot for a snapshot the
// account does not own.
//
// The middleware runs once per operation, before the retry loop, and returns
// at once in list and get calls, which the SDK rate limits itself.
func addHydrateRateLimitMiddleware(cfg *aws.Config) {
	cfg.APIOptions = append(cfg.APIOptions, func(stack *middleware.Stack) error {
		return stack.Initialize.Add(middleware.InitializeMiddlewareFunc(hydrateRateLimitMiddlewareID,
			func(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (middleware.InitializeOutput, middleware.Metadata, error) {
				plugin.WaitForHydrateRateLimit(ctx)
				return next.HandleInitialize(ctx, in)
			}), middleware.Before)
	})
}
