package log

import (
	"context"
	"log"
	"qq/pkg/qqcontext"
	"reflect"
	"time"
)

type Args map[string]interface{}

func Critical(ctx context.Context, message string, args ...Args) {
	print("CRITICAL", ctx, message, args)
}

func Error(ctx context.Context, message string, args ...Args) {
	print("ERROR", ctx, message, args)
}

func Warning(ctx context.Context, message string, args ...Args) {
	print("WARNING", ctx, message, args)
}

func Info(ctx context.Context, message string, args ...Args) {
	print("INFO", ctx, message, args)
}

func Debug(ctx context.Context, message string, args ...Args) {
	print("DEBUG", ctx, message, args)
}

func print(severity string, ctx context.Context, message string, argsArr []Args) {
	currentTime := time.Now()
	userId := qqcontext.GetUserIdValue(ctx)

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	log.Printf("%v [%s] %s: %s\n",
		currentTime.Format("2006.01.02 15:04:05.000"), severity, userId, message)

	for _, args := range argsArr {
		for k, v := range args {
			if reflect.ValueOf(v).Kind() == reflect.Struct {
				log.Printf("%s: %+v\n", k, v)
				continue
			}

			log.Printf("%s: %v\n", k, v)
		}
	}
}
