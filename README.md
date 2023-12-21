# cslog

A simple wrapper to put the logger in a context, so that you can simply pass the context around.
Examples include attaching relevant information _contextually_ to an HTTP request, such that information is preserved
downstream, without having to also pass around a logger. You're probably going to pass the context anyway.

## Example

```go
import (
	"context"
	"log/slog"
	"os"
	"time"

    "github.com/jclasley/cslog"
)

func main() {
	var opt *slog.HandlerOptions
	// uncomment below to have the debug show up!
	// opt = &slog.HandlerOptions{Level: slog.LevelDebug}

	ctx := cslog.FromBackground(slog.New(slog.NewJSONHandler(os.Stdout, cslog.WithoutTime(opt))))

	cslog.Info(ctx, "hello!")
	y, _, _ := time.Now().Local().Date()
	ctx = cslog.WithAttrs(ctx, slog.Int("year", y))

	got := getFromDB(cslog.WithGroup(ctx, "getFromDB"), "1234567890")
	if want := 43; got != want {
		cslog.Error(ctx, "wrong result!", slog.Int("wanted", want), slog.Int("got", got))
	}
}

func getFromDB(ctx context.Context, userID string) int {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    ctx = cslog.WithAttrs(ctx, slog.String("userID", userID))

    // call DB
    const fakeResult = 42
    cslog.Debug(ctx, "got from DB", slog.Int("result", 42))
    return fakeResult
}
```

## Option funcs

The `ReplaceAttr` function used by `slog` is awesome, and I can forsee adding more option functions in the future. For now, I've only implemented what I've found to be the most useful, which is to drop `time` from the attributes.

You can check out the tests to see a nifty one I created that hides secret values. Basically, if you have some group `secret`, all attribute **values** added to that will be replaced with as many `*` as the value has. So `slog.Group("secret", slog.Group("db", slog.String("password", "12345")))` will result in an attribute (assuming a JSON handler) that looks like `"secret":{"db":{"password":"*****"}}`. You could probably also configure this so secrets are shown at the `DEBUG` level, and so on. Pretty cool.
