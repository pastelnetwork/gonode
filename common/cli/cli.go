package cli

import "context"

// ActionFn represents the application's entry point function.
type ActionFn func(ctx context.Context, args []string) error
