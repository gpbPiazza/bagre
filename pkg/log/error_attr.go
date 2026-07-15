package log

import (
	"errors"
	"fmt"
	"log/slog"

	pkgerr "github.com/pkg/errors"
)

// Err is helper func to make a log Attr with erroKey
func Err(err error) slog.Attr {
	return slog.Any(errorKey, pkgerr.WithStack(err))
}

const errorKey = "error"

type stackTracer interface {
	error
	StackTrace() pkgerr.StackTrace // StackTrace is the method from stack type from the third party library
}

type multiError interface {
	error
	Unwrap() []error
}

type errWithAttrs struct {
	error
	attrs []slog.Attr
}

func ErrWithAttrs(err error, args ...any) error {
	return &errWithAttrs{
		error: err,
		attrs: argsToAttr(args),
	}
}

// argsToAttr turns a list of typed or untyped values into a slice of [slog.Attr].
// args[i] is treated as a key if it is a string or an [slog.Attr]; otherwise, it
// is treated as a value with key "!BADKEY".
func argsToAttr(args []any) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(args))
	for i := 0; i < len(args); {
		switch key := args[i].(type) {
		case slog.Attr:
			attrs = append(attrs, key)
			i++
		case string:
			if i+1 >= len(args) {
				attrs = append(attrs, slog.String("!BADKEY", key))
				i++
			} else {
				attrs = append(attrs, slog.Any(key, args[i+1]))
				i += 2
			}
		default:
			attrs = append(attrs, slog.Any("!BADKEY", args[i]))
			i++
		}
	}

	return attrs
}

func (e *errWithAttrs) Unwrap() error {
	return e.error
}

func (e *errWithAttrs) Attrs() []slog.Attr {
	return e.attrs
}

type attrError interface {
	error
	Attrs() []slog.Attr
}

// attrsFromErr recursively extracts all logging attributes from an error chain. In the
// case of duplicate keys, the outermost value takes precedence.
func attrsFromErr(err error) []slog.Attr {
	var out []slog.Attr

	for err != nil {
		attrErr, ok := errors.AsType[attrError](err)
		if ok {
			out = append(out, attrErr.Attrs()...)
		}
		err = errors.Unwrap(err)
	}

	return out
}

func errAttr(err error) []slog.Attr {
	var logAttr []slog.Attr

	logAttr = append(logAttr, slog.Attr{Key: "message", Value: slog.StringValue(err.Error())})

	attrsFromErr := attrsFromErr(err)
	if len(attrsFromErr) != 0 {
		logAttr = append(logAttr, attrsFromErr...)
	}

	stackErr, ok := errors.AsType[stackTracer](err)
	if ok {
		logAttr = append(logAttr,
			slog.Attr{
				Key:   "stack_trace",
				Value: slog.StringValue(fmt.Sprintf("%+v", stackErr.StackTrace())),
			})
	}

	return logAttr
}
