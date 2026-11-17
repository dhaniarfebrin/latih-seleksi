// Command api serves the drill API.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"latihseleksi/api/internal/auth"
	"latihseleksi/api/internal/bank"
	"latihseleksi/api/internal/drill"
	"latihseleksi/api/internal/exam"
	"latihseleksi/api/internal/export"
	"latihseleksi/api/internal/gen/catalog"
	"latihseleksi/api/internal/httpx"
	"latihseleksi/api/internal/insight"
	"latihseleksi/api/internal/review"
	"latihseleksi/api/internal/spec"
	"latihseleksi/api/internal/store"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return errors.New("DATABASE_URL is required")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	st, err := store.Open(ctx, dbURL)
	if err != nil {
		return err
	}
	defer st.Close()
	if err := st.Migrate(ctx); err != nil {
		return err
	}

	specs, err := spec.Load()
	if err != nil {
		return err
	}
	items, err := bank.Load()
	if err != nil {
		return err
	}
	examSvc := exam.New(st, specs, exam.NewGeneratedSource(catalog.All()), items)

	reviews := review.New(st)

	// FIREBASE_PROJECT_ID is not a secret. Empty leaves /v1/exams* closed (503).
	var verifier *auth.Verifier
	if pid := os.Getenv("FIREBASE_PROJECT_ID"); pid != "" {
		verifier = auth.NewFirebase(pid)
	} else {
		slog.Warn("FIREBASE_PROJECT_ID not set: /v1/exams* answers 503")
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           httpx.NewHandler(drill.New(st, catalog.All(), items, reviews), examSvc, insight.New(st, catalog.All(), reviews), export.New(st), reviews, st, httpx.Config{CORSOrigin: os.Getenv("CORS_ORIGIN"), Auth: verifier}),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	errc := make(chan error, 1)
	go func() { errc <- srv.ListenAndServe() }()
	slog.Info("listening", "addr", srv.Addr)

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}
	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
