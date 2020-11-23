package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"time"

	"github.com/liam-jones-lucout/golangtest/internal/pkg/db"
	"github.com/liam-jones-lucout/golangtest/internal/pkg/mylogger"
	"github.com/liam-jones-lucout/golangtest/internal/pkg/spaceshipmodels"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type key int

const (
	requestIDKey   key    = 0
	successPayload string = "{\"success\": true}"
)

var (
	port    string
	healthy int32
)

func main() {
	flag.StringVar(&port, "port", ":5000", "server listen address")
	flag.Parse()

	cfg := zap.NewProductionConfig()

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := mylogger.NewLogger()
	if err != nil {
		panic("Failed to create logger!")
	}

	dataAccess := db.NewDb(logger)
	if err := dataAccess.Initiate(); err != nil {
		logger.Fatal("Failed to initialise database", zap.Error(err))
	}

	logger.Info("Server is starting...")

	router := http.NewServeMux()
	router.Handle("/", index())
	router.Handle("/search", search(logger, dataAccess))
	router.Handle("/ship/", crud(logger, dataAccess))

	nextRequestID := func() string {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	server := &http.Server{
		Addr:         port,
		Handler:      tracing(nextRequestID)(logging(logger)(router)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Info("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatal("Could not gracefully shutdown the server:", zap.Error(err))
		}
		close(done)
	}()

	logger.Info("Server is ready to handle requests at", zap.String("listenAddr", port))
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Could not listen:", zap.String("listenAddr", port), zap.Error(err))
	}

	<-done
	logger.Info("Server stopped")
}

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello, World!")
	})
}

func search(logger *zap.Logger, dataAccess *db.Db) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		name := query["name"]
		class := query["class"]
		status := query["status"]

		results, err := dataAccess.Search(name, class, status)
		if err != nil {
			logger.Error("Failed to search database", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		payload, err := json.Marshal(results)
		if err != nil {
			logger.Error("Failed to marshal search results", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info("returning payload:", zap.String("payload", string(payload)))

		if _, err := w.Write(payload); err != nil {
			logger.Error("Failed to write results to response body", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func crud(logger *zap.Logger, dataAccess *db.Db) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := strings.TrimPrefix(r.URL.Path, "/ship/")

		switch r.Method {
		case http.MethodGet:
			result, err := dataAccess.Get(id)
			if err != nil {
				logger.Error("Failed to search database", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			payload, err := json.Marshal(result)
			if err != nil {
				logger.Error("Failed to marshal search results", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			logger.Info("returning payload:", zap.String("payload", string(payload)))

			if _, err := w.Write(payload); err != nil {
				logger.Error("Failed to write results to response body", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		case http.MethodDelete:
			err := dataAccess.Delete(id)
			if err != nil {
				logger.Error("Failed to search database", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if _, err := w.Write([]byte(successPayload)); err != nil {
				logger.Error("Failed to write results to response body", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

		case http.MethodPut:
			payload := spaceshipmodels.Spaceship{}
			payloadbytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logger.Error("Failed to search database", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err = json.Unmarshal(payloadbytes, &payload)
			if err != nil {
				logger.Error("Failed to unmarshal payload", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err = dataAccess.Update(id, payload)
			if err != nil {
				logger.Error("Failed to search database", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if _, err := w.Write([]byte(successPayload)); err != nil {
				logger.Error("Failed to write results to response body", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

	})
}

func logging(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				requestID, ok := r.Context().Value(requestIDKey).(string)
				if !ok {
					requestID = "unknown"
				}
				logger.Info(requestID, zap.String("Method", r.Method), zap.String("path", r.URL.Path), zap.String("remoteaddr", r.RemoteAddr), zap.String("Uagent", r.UserAgent()))
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func tracing(nextRequestID func() string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = nextRequestID()
			}
			ctx := context.WithValue(r.Context(), requestIDKey, requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
