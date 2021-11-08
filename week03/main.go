package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return srv.ListenAndServe()
	})

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can"t be cached, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	g.Go(func() error {
		<-quit
		log.Println("Shutdown Server ...")
		return srv.Shutdown(ctx)
	})

	err := g.Wait()
	fmt.Println(err)
	fmt.Println(ctx.Err())
}
