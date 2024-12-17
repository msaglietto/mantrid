package middleware

import (
	"context"
	"fmt"

	"github.com/msaglietto/mantrid/internal/logging"
)

func WithErrorHandling(ctx context.Context, fn func() error) error {
	logger := logging.FromContext(ctx)

	err := fn()
	if err != nil {
		logger.Error("operation failed", "error", err)
		return fmt.Errorf("operation failed: %w", err)
	}

	return nil
}
