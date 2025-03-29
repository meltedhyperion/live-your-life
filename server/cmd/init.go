package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/meltedhyperion/globetrotter/server/db/pg_db"
	"github.com/meltedhyperion/globetrotter/server/logger"
	"github.com/meltedhyperion/globetrotter/server/util"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		logger.Log.Error(err)
	}
}
func InitServer(app *App) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(loggerMiddleware(logger.Log))

	// setup cors
	r.Use(cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		ExposedHeaders:   []string{"Authorization"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	}).Handler)

	initHandler(app, r)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "3000"
	}
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	srv := http.Server{
		Addr:    addr,
		Handler: r,
	}
	app.Srv = &srv

	walkFunc := func(method string, route string, handler http.Handler, middleware ...func(http.Handler) http.Handler) error {
		fmt.Printf("\t\t%s %s\n", util.PadStringTo(method, 7), route)
		return nil
	}

	fmt.Print("\t\tRegistered Routes: \n\n")
	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Error logging routes. Err: %s\n", err.Error())
	}
}

func InitDB(app *App) {
	logger.Log.Info("connecting to db")
	tDB, err := sql.Open("postgres", os.Getenv("DB_URI"))
	if err != nil {
		logger.Log.Fatal("CANNOT INIT DB", err)
	}
	err = tDB.Ping()
	if err != nil {
		logger.Log.Fatal("CANNOT PING DB", err)
	}
	q, err := pg_db.Prepare(context.Background(), tDB)
	if err != nil {
		logger.Log.Fatal("CANNOT PREPARE DB", err)
	}
	logger.Log.Info("connected to pg db")
	app.store = pg_db.NewStore(tDB, q)
}

func loggerMiddleware(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestBody []byte
			if r.Body != nil {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					logger.Errorw("Error reading request body", "error", err)
				}
				requestBody = body
				r.Body = io.NopCloser(bytes.NewReader(body))
			}

			queryParams := make(map[string]string)
			queryValues, _ := url.ParseQuery(r.URL.RawQuery)
			for key, values := range queryValues {
				queryParams[key] = strings.Join(values, ", ")
			}
			excludedHeaders := map[string]bool{
				"Authorization": true,
			}

			requestHeaders := make(http.Header)
			for key, values := range r.Header {
				if !excludedHeaders[key] {
					requestHeaders[key] = values
				}
			}
			next.ServeHTTP(w, r)
			logger.Infow("HTTP Request",
				"request_headers", requestHeaders,
				"request_method", r.Method,
				"request_url", r.URL.String(),
				"request_query", queryParams,
				"request_payload", string(requestBody),
				"response_status", w.Header().Get("Status"),
				"response_statuscode", w.Header().Get("StatusCode"),
			)
		})
	}
}
