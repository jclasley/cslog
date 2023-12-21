package cslog

import (
	"context"
	"log/slog"
	"os"
	"time"
)

func ExampleCtx() {
	getFromDB := func(ctx context.Context, userID string) int {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		ctx = WithAttrs(ctx, slog.String("userID", userID))

		// call DB
		const fakeResult = 42
		Debug(ctx, "got from DB", slog.Int("result", 42))
		return fakeResult
	}

	var opt *slog.HandlerOptions
	// uncomment below to have the debug show up!
	// opt = &slog.HandlerOptions{Level: slog.LevelDebug}

	ctx := FromBackground(slog.New(slog.NewJSONHandler(os.Stdout, WithoutTime(opt))))

	Info(ctx, "hello!")
	y, _, _ := time.Now().Local().Date()
	ctx = WithAttrs(ctx, slog.Int("year", y))

	got := getFromDB(WithGroup(ctx, "getFromDB"), "foobar!!!")
	if want := 43; got != want {
		Error(ctx, "wrong result!", slog.Int("wanted", want), slog.Int("got", got))
	}
}
