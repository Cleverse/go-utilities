package logger

import (
	"context"
	"log/slog"

	"github.com/Cleverse/go-utilities/logger/slogx"
	"github.com/Cleverse/go-utilities/logger/stacktrace"
	"github.com/lmittmann/tint"
)

// middlewareErrorStackTrace is a middleware that adds the error stack trace when ErrorKey or "err" is found.
func middlewareErrorStackTrace() middleware {
	return func(next handleFunc) handleFunc {
		return func(ctx context.Context, rec slog.Record) error {
			rec.Attrs(func(attr slog.Attr) bool {
				if attr.Key == slogx.ErrorKey || attr.Key == "err" {
					err := attr.Value.Any()
					if err, ok := err.(error); ok && err != nil {
						rec.AddAttrs(
							slogx.Stringer(
								slogx.ErrorStackTraceKey,
								stacktrace.ExtractErrorStackTraces(err),
							),
						)
					}
				}
				return true
			})
			return next(ctx, rec)
		}
	}
}

func errorAttrReplacer(groups []string, attr slog.Attr) slog.Attr {
	if len(groups) == 0 {
		switch attr.Key {
		case slogx.ErrorKey, "err":
			if err, ok := attr.Value.Any().(error); ok {
				aErr := tint.Attr(9, slogx.String(slogx.ErrorKey, err.Error()))
				aErr.Key = slogx.ErrorKey
				return aErr
			}
		case slogx.ErrorStackTraceKey:
			type stackDetails struct {
				Error  string   `json:"error"`
				Stacks []string `json:"stacks"`
			}
			if st, ok := attr.Value.Any().(stacktrace.ErrorStackTraces); ok {
				errsStacks := make([]stackDetails, 0)
				for _, errStack := range st {
					errsStacks = append(errsStacks, stackDetails{
						Error:  errStack.Error(),
						Stacks: errStack.StackTrace.FramesStrings(),
					})
				}
				return slog.Attr{Key: slogx.ErrorStackTraceKey, Value: slog.AnyValue(errsStacks)}
			}
		}
	}
	return attr
}
