package main

import "github.com/kimxuanhong/go-logger/logger"

func main() {
	err := logger.Init()
	if err != nil {
		panic(err)
	}

	logger.Log.WithField("logger", "main").Info("Ứng dụng đã khởi động")

	someFunc()
}

func someFunc() {
	logger.Log.WithField("logger", "main.someFunc").Warn("Gọi hàm someFunc")
}
