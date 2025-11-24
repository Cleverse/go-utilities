package logger

import (
	"context"
	"fmt"
	"log/slog"
)

func durationToMsAttrReplacer(groups []string, attr slog.Attr) slog.Attr {
	if attr.Value.Kind() == slog.KindDuration {
		return slog.Attr{
			Key:   attr.Key,
			Value: slog.Int64Value(attr.Value.Duration().Milliseconds()),
		}
	}
	return attr
}

// middlewareDurationHuman is a middleware that adds a new string attribute for each duration attribute to render formatted duration
func middlewareDurationHuman() middleware {
	return func(next handleFunc) handleFunc {
		return func(ctx context.Context, rec slog.Record) error {
			rec.Attrs(func(attr slog.Attr) bool {
				if attr.Value.Kind() == slog.KindDuration {
					duration := attr.Value.Duration()
					rec.AddAttrs(
						slog.String(
							fmt.Sprintf("%s_human", attr.Key),
							duration.String(),
						),
					)
				}
				return true
			})
			return next(ctx, rec)
		}
	}
}
