package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rolandwarburton/playwright-server/controllers"
	"github.com/rolandwarburton/playwright-server/routes"
)

type Server struct {
	port           int
	mode           string
	trustedProxies []string
	routes         []routes.Route
}

func main() {
	gin.SetMode("debug")

	// create a router for the whole framework
	router := gin.New()
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	v1 := router.Group("/agent")
	root := router.Group("")

	wsController := controllers.NewWSController()
	accountController := controllers.NewAgentController()

	// ACCOUNT ROUTES
	agentMiddleware := &routes.Middleware{
		// GET:    []func() gin.HandlerFunc{},
	}
	route, _ := routes.GetRoute("create", *agentMiddleware)
	route.Register("GET", v1, accountController.CreateAgent)
	route, _ = routes.GetRoute("example", *agentMiddleware)
	route.Register("GET", v1, accountController.ExampleAction)

	route, _ = routes.GetRoute("ws/:id", routes.Middleware{})
	route.Register("GET", root, wsController.WS)

	// create a server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// run the server in a go routine
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
	// create a channel to wait for an interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)

	// when SIGINT or SIGTERM is sent to the process, then notify the "quit" channel
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// once quit is sent a signal by above signal.Notify we can start shutting down the server
	<-quit

	// the timeout returns ctx which has a Done() channel, we can wait for the channel to be called
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutdown server...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown:", err)
	} else {
		// the server has closed, cancel the timer
		fmt.Println("server shutdown successful")
		cancel()
	}

	// catch ctx.Done (when N seconds has elapsed, or cancel has been called)
	<-ctx.Done()
	log.Println("server exiting: ", ctx.Err())
}
