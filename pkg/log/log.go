// Package log doc here
package log

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

// CloseFunc flushes and releases every resource the logger owns.
type CloseFunc func() error

func InitLogger() (*slog.Logger, CloseFunc, error) {
	handlers := make([]slog.Handler, 0, 1)
	var closers []CloseFunc

	// handlers = append(handlers, slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
	// 	ReplaceAttr: replaceAttr,
	// }))

	handlers = append(handlers, tint.NewTextHandler(os.Stderr, &tint.Options{
		ReplaceAttr: replaceAttr,
	}))

	closeAll := func() error {
		errs := make([]error, 0, len(closers))
		for _, closer := range closers {
			errs = append(errs, closer())
		}
		return errors.Join(errs...)
	}

	logger := slog.New(slog.NewMultiHandler(handlers...))

	return logger, closeAll, nil
}

func replaceAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Key == errorKey {
		err, ok := a.Value.Any().(error)
		if !ok {
			return a
		}

		if multiErr, isMulti := errors.AsType[multiError](err); isMulti {
			var errAttrs []slog.Attr
			for i, e := range multiErr.Unwrap() {
				errAttrs = append(errAttrs,
					slog.GroupAttrs(fmt.Sprintf("error_%d", i+1),
						errAttr(e)...))
			}

			return slog.GroupAttrs("errors", errAttrs...)
		}

		return slog.GroupAttrs("error", errAttr(err)...)
	}

	return a
}
