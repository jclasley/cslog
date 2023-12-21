package cslog

import (
	"context"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtx(t *testing.T) {
	opt := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	opt = WithoutTime(opt)

	l := slog.New(slog.NewJSONHandler(os.Stdout, opt))
	ctx := FromBackground(l)

	Debug(ctx, "hello", slog.String("from", "outer"))

	ctx = WithAttrs(ctx, slog.String("from", "outer"))

	f := func(ctx context.Context) {
		ctx = WithAttrs(ctx, slog.String("from", "inner"))

		Debug(ctx, "hello")
		ctx = WithAttrs(ctx, slog.Int("count", 1))
		Debug(ctx, "goodbye")
	}

	f(WithGroup(ctx, "f"))
	Debug(ctx, "goodbye")
}

func TestReplaceFn(t *testing.T) {
	opt := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if slices.Contains(groups, "secret") {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue(strings.Repeat("*", len(a.Value.String()))),
				}
			}
			return a
		},
	}

	f := new(Fixture)
	tt := []struct {
		name       string
		w          io.Writer
		initOpts   *slog.HandlerOptions
		forDisplay bool
	}{
		{
			"stdout -- inspection",
			os.Stdout,
			opt,
			true,
		},
		{
			name:     "fixture",
			w:        f,
			initOpts: opt,
		},
		{
			name:     "fixture with nil start",
			w:        f,
			initOpts: nil,
		},
	}
	// print it
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			l := slog.New(slog.NewJSONHandler(tc.w, WithoutTime(opt)))
			l.Debug("hello", slog.Group("secret", slog.String("value", "foo")))

			if tc.forDisplay {
				return
			}
			require.NotContains(t, f.Body, "time")
			require.NotContains(t, f.Body, "foo")
			require.Contains(t, f.Body, "value")
			require.Contains(t, f.Body, "secret")

			f.Reset()
		})
	}
}

type Fixture struct {
	Body string
}

func (f *Fixture) Write(b []byte) (int, error) {
	f.Body = string(b)
	return len(b), nil
}

func (f *Fixture) Reset() {
	f.Body = ""
}
