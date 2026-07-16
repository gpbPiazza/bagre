// Package log doc here
package log

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

type closeFunc func() error

func InitLogger() (*slog.Logger, closeFunc, error) {
	var (
		handlers []slog.Handler
		closers  []closeFunc
	)

	// handlers = append(handlers, slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
	// 	ReplaceAttr: replaceAttr,
	// }))

	handlers = append(handlers, tint.NewTextHandler(os.Stderr, &tint.Options{
		ReplaceAttr: replaceAttr,
	}))

	close := func() error {
		var errs []error
		for _, closer := range closers {
			errs = append(errs, closer())
		}
		return errors.Join(errs...)
	}

	logger := slog.New(slog.NewMultiHandler(handlers...))

	return logger, close, nil
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == errorKey {
		err, ok := a.Value.Any().(error)
		if !ok {
			return a
		}

		if multiErr, ok := errors.AsType[multiError](err); ok {
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
