package main

import (
	"context"
	"fmt"
	echo "go.buf.build/loopholelabs/frisbee/loopholelabs/echo-example"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c, err := echo.NewClient("127.0.0.1:8080", nil, nil)
	if err != nil {
		panic(err)
	}

	err = c.Connect()
	if err != nil {
		panic(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	var t time.Time

	var res *echo.Response
	req := echo.NewRequest()
	req.Message = "Test Message"
	i := 0
	for {
		select {
		case <-stop:
			err = c.Close()
			if err != nil {
				panic(err)
			}
			return
		default:
			req.Message = fmt.Sprintf("%d", i)
			log.Printf("Sending Request %s\n", req.Message)
			t = time.Now()
			res, err = c.Echo(context.Background(), req)
			if err != nil {
				panic(err)
			}
			log.Printf("Received Response %s in %s\n", res.Message, time.Since(t))
			time.Sleep(time.Second)
			i++
		}
	}
}
