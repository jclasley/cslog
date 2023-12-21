package cslog

import "log/slog"

type replaceFn = func([]string, slog.Attr) slog.Attr

// WithoutTime adds a ReplaceAttr function to remove `time` from the attributes on each log.
//
// If the passed in handler options are nil, it uses the default.
func WithoutTime(o *slog.HandlerOptions) *slog.HandlerOptions {
	var f replaceFn
	replaceF := func(groups []string, a slog.Attr) slog.Attr {
		if groups == nil && a.Key == "time" {
			return slog.Attr{}
		}
		if f != nil {
			return f(groups, a)
		}
		return a
	}

	if o == nil {
		return &slog.HandlerOptions{
			ReplaceAttr: replaceF,
		}
	}

	f = o.ReplaceAttr
	o.ReplaceAttr = replaceF
	return o
}
