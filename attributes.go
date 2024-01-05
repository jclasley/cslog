package cslog

import "log/slog"

// Convenience attributes to avoid having to type slog.String("key", "value") all the time.

// Err is the equivalent of slog.String("err", err.Error())
func Err(err error) slog.Attr {
	return slog.String("err", err.Error())
}
