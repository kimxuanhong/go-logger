package main

import (
	"context"
	"github.com/kimxuanhong/go-logger/logger"
	"time"
)

func main() {
	err := logger.Init()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	log := logger.WithContext(ctx)
	log.Info("Ứng dụng đã khởi động")
	logger.Log.WithField("logger", "main").Info("Ứng dụng đã khởi động")

	someFunc()
}

func someFunc() {
	logger.Log.WithField("logger", "main.someFunc").Warn("Gọi hàm someFunc")
}
