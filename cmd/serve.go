package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"template.go/database"
	c "template.go/packages/common"
	"template.go/src/middleware"
	"template.go/src/routes"
)

func NewServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start http server",
		Run: func(cmd *cobra.Command, args []string) {
			host, _ := cmd.Flags().GetString("host")
			port, _ := cmd.Flags().GetString("port")
			dsn, _ := cmd.Flags().GetString("dsn")
			debug, _ := cmd.Flags().GetBool("debug")

			if debug {
				log.SetFlags(log.LstdFlags | log.Lshortfile)
			}

			gin.SetMode(gin.ReleaseMode)
			info := GetVersion()

			showVersion := flag.Bool("version", info.Show, "Show version information")
			flag.Parse()
			if *showVersion {
				slog.Info("System",
					slog.String("version", info.Version),
					slog.String("go", info.Go),
					slog.String("compiler", info.Compiler),
					slog.String("platform", info.Platform),
				)
			}

			if dsn != "" {
				database.Connect(dsn, debug)
			}
			server := NewServer(host, port)

			ctx, cancel := context.WithCancel(context.Background())
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

			go func() {
				<-sig
				slog.Info("Shutdown server...")
				cancel()
			}()

			go func() {
				server.Start()
			}()

			<-ctx.Done()

			if dsn != "" {
				if err := database.Disconnect(); err != nil {
					slog.Error("Close database connection", slog.Any("error", err))
					panic(err)
				}
				slog.Info("Database disconnected")
			}

			server.Stop(10 * time.Second)
		},
	}

	cmd.Flags().StringP("host", "H", c.Env("HOST", "localhost"), "Server host")
	cmd.Flags().StringP("port", "p", c.Env("PORT", "8000"), "Server port")

	return cmd
}

type Server struct {
	*http.Server
}

func NewServer(host, port string) *Server {
	db := database.DB()
	router := gin.Default()

	router.Static("/public", "storage/public")
	router.StaticFile("/favicon.ico", "public/favicon.ico")
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.Cors())
	routes.App(db, router.Group("/"))

	router.NoRoute(func(ctx *gin.Context) {
		switch ctx.NegotiateFormat(gin.MIMEJSON, gin.MIMEHTML) {
		case gin.MIMEJSON:
			ctx.JSON(404, gin.H{"error": "Not found"})
		default:
			ctx.String(404, "Not found")
		}
	})

	return &Server{
		&http.Server{
			Addr:    fmt.Sprintf("%s:%s", host, port),
			Handler: router,
		},
	}
}

func (server *Server) Start() {
	slog.Info("Listening on", slog.String("url", "http://"+server.Addr))
	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		slog.Error("Start Http server", slog.Any("error", err))
		panic("Start Http server")
	}
}

func (server *Server) Stop(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Stop Http server", slog.Any("error", err))
		panic("Stop Http server")
	}
	slog.Info("Http server stopped")
}
