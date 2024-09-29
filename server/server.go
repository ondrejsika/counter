package server

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ondrejsika/counter/backend_inmemory"
	"github.com/ondrejsika/counter/backend_mongodb"
	"github.com/ondrejsika/counter/backend_postgres"
	"github.com/ondrejsika/counter/backend_redis"
	"github.com/ondrejsika/counter/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

var RunTimestamp time.Time
var Logger zerolog.Logger

var promRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "counter",
		Name:      "requests_total",
		Help:      "Total number of requests per endpoint",
	}, []string{"path"})

var promCounter = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Namespace: "counter",
		Name:      "counter",
		Help:      "Current counter value",
	})

type StatusResponse struct {
	Hostname             string `json:"hostname"`
	RunTimestampUnix     int    `json:"run_timestamp_unix"`
	RunTimestampRFC3339  string `json:"run_timestamp_rfc3339"`
	RunTimestampUnixDate string `json:"run_timestamp_unixdate"`
}

func indexPlainText(w http.ResponseWriter, hostname string, count int, extraText string) {
	if extraText != "" {
		fmt.Fprintf(w, "%s %s %d", extraText, hostname, count)
	} else {
		fmt.Fprintf(w, "%s %d", hostname, count)
	}
	fmt.Fprint(w, "\n")
}

func indexHTML(w http.ResponseWriter, hostname string, count int, extraText string) {
	extraTextBlock := ""
	if extraText != "" {
		extraTextBlock = `<h2>` + extraText + `</h2>`
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<!DOCTYPE html>
	<html lang="en">
	<head>
  <meta charset="UTF-8">
	<style>
	html, body {
		height: 100%;
	}
	.center-parent {
		width: 100%;
		height: 100%;
		display: table;
		text-align: center;
	}
	.center-parent > .center-child {
		display: table-cell;
		vertical-align: middle;
	}
	</style>
	<style>
	h1 {
		font-family: Arial;
		font-size: 5em;
	}
	h2 {
		font-family: Arial;
		font-size: 2em;
	}
	</style>
	<link rel="icon" href="/favicon.ico">
	</head>
	<body>
	<section class="center-parent">
		<div class="center-child">
			`+extraTextBlock+`
			<h1>
				`+strconv.Itoa(count)+`
			</h1>
			<h2>`+hostname+`</h2>
		</div>
	</section>
	</body></html>
	`)
}

func versionAPI(w http.ResponseWriter, r *http.Request) {
	promRequestsTotal.With(prometheus.Labels{"path": r.URL.Path}).Inc()
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(map[string]string{
		"version": version.Version,
	})
	fmt.Fprint(w, string(data))
}

func livez(w http.ResponseWriter, r *http.Request) {
	promRequestsTotal.With(prometheus.Labels{"path": r.URL.Path}).Inc()
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(map[string]bool{
		"live": true,
	})
	fmt.Fprint(w, string(data))
}

func readyz(w http.ResponseWriter, r *http.Request) {
	promRequestsTotal.With(prometheus.Labels{"path": r.URL.Path}).Inc()
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(map[string]bool{
		"ready": true,
	})
	fmt.Fprint(w, string(data))
}

func status(w http.ResponseWriter, r *http.Request) {
	promRequestsTotal.With(prometheus.Labels{"path": r.URL.Path}).Inc()
	hostname, _ := os.Hostname()
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(StatusResponse{
		Hostname:             hostname,
		RunTimestampUnix:     int(RunTimestamp.Unix()),
		RunTimestampRFC3339:  RunTimestamp.Format(time.RFC3339),
		RunTimestampUnixDate: RunTimestamp.Format(time.UnixDate),
	})
	fmt.Fprint(w, string(data))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	promRequestsTotal.With(prometheus.Labels{"path": r.URL.Path}).Inc()
	w.WriteHeader(http.StatusNotFound)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	promRequestsTotal.With(prometheus.Labels{"path": r.URL.Path}).Inc()
	promhttp.Handler().ServeHTTP(w, r)
}

func BaseServer(
	doCountFunc func() (int, error),
	getCountFunc func() (int, error),
) {
	var err error

	prometheus.MustRegister(promRequestsTotal)
	prometheus.MustRegister(promCounter)

	Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	RunTimestamp = time.Now()

	hostname, _ := os.Hostname()

	port := "80"
	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	}

	slowStart := 0
	envSlowStart := os.Getenv("SLOW_START")
	if envSlowStart != "" {
		slowStart, err = strconv.Atoi(envSlowStart)
		if err != nil {
			Logger.Fatal().Str("hostname", hostname).Msg(`cannot parse integer form SLOW_START value "` + envSlowStart + `", original Go error: ` + err.Error())
		}
	}

	extraText := os.Getenv("EXTRA_TEXT")

	go func() {
		for {
			count, _ := getCountFunc()
			promCounter.Set(float64(count))
			time.Sleep(1 * time.Second)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		promRequestsTotal.With(prometheus.Labels{"path": r.URL.Path}).Inc()
		hostname, _ := os.Hostname()

		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404 Not Found\n")

			Logger.Info().
				Str("hostname", hostname).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Msg(r.Method + " " + r.URL.Path + " 404 Not Found")
			return
		}

		counter, _ := doCountFunc()
		// Check if User-Agent header exists
		if userAgentList, ok := r.Header["User-Agent"]; ok {
			// Check if User-Agent header has some data
			if len(userAgentList) > 0 {
				// If User-Agent starts with curl, use plain text
				if strings.HasPrefix(userAgentList[0], "curl") {
					indexPlainText(w, hostname, counter, extraText)
				} else {
					// If User-Agent header presents and not starts with curl
					// use HTML (Chrome, Safari, Firefox, ...)
					indexHTML(w, hostname, counter, extraText)
				}
			}
		} else {
			// If User-Agent header doesn't exists, use plain text
			indexPlainText(w, hostname, counter, extraText)
		}
		Logger.Info().
			Str("hostname", hostname).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("counter", counter).
			Msg(r.Method + " " + r.URL.Path)
	})
	http.HandleFunc("/api/counter", func(w http.ResponseWriter, r *http.Request) {
		promRequestsTotal.With(prometheus.Labels{"path": r.URL.Path}).Inc()
		hostname, _ := os.Hostname()
		counter, _ := doCountFunc()
		w.Header().Set("Content-Type", "application/json")
		type Response struct {
			Counter   int    `json:"counter"`
			Hostname  string `json:"hostname"`
			Version   string `json:"version"`
			ExtraText string `json:"extra_text"`
		}
		data, _ := json.Marshal(Response{
			Counter:   counter,
			Hostname:  hostname,
			Version:   version.Version,
			ExtraText: extraText,
		})
		fmt.Fprint(w, string(data))
		fmt.Fprint(w, "\n")
		Logger.Info().
			Str("hostname", hostname).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("counter", counter).
			Msg(r.Method + " " + r.URL.Path)
	})
	http.HandleFunc("/api/version", versionAPI)
	http.HandleFunc("/version", versionAPI)
	http.HandleFunc("/api/livez", livez)
	http.HandleFunc("/livez", livez)
	http.HandleFunc("/api/readyz", readyz)
	http.HandleFunc("/readyz", readyz)
	http.HandleFunc("/api/status", status)
	http.HandleFunc("/status", status)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/metrics", metricsHandler)

	Logger.Info().Str("hostname", hostname).Msg("Starting server counter " + version.Version + " ...")

	for i := 0; i < slowStart; i++ {
		Logger.Info().Str("hostname", hostname).Msgf("Starting in %d seconds ...", slowStart-i)
		time.Sleep(1 * time.Second)
	}

	Logger.Info().Str("hostname", hostname).Msg("Server counter " + version.Version + " started on 0.0.0.0:" + port + ", see http://127.0.0.1:" + port)
	err = http.ListenAndServe("0.0.0.0:"+port, nil)
	if err != nil {
		Logger.Fatal().Str("hostname", hostname).Msg(err.Error())
	}
}

func Server() {
	hostname, _ := os.Hostname()

	backend := "redis"
	envBackend := os.Getenv("BACKEND")
	if envBackend != "" {
		backend = envBackend
	}

	var doCountFunc func() (int, error)
	var getCountFunc func() (int, error)

	if backend == "redis" {
		redisHost := "127.0.0.1"
		envRedisHost := os.Getenv("REDIS")
		if envRedisHost != "" {
			redisHost = envRedisHost
		}
		doCountFunc = func() (int, error) { return backend_redis.DoCountRedis(redisHost, hostname) }
		getCountFunc = func() (int, error) { return backend_redis.GetCountRedis(redisHost, hostname) }
	} else if backend == "inmemory" {
		doCountFunc = func() (int, error) { return backend_inmemory.DoCountInMemory() }
		getCountFunc = func() (int, error) { return backend_inmemory.GetCountInMemory() }
	} else if backend == "postgres" {
		postgresHost := "127.0.0.1"
		envPostgresHost := os.Getenv("POSTGRES_HOST")
		if envPostgresHost != "" {
			postgresHost = envPostgresHost
		}

		postgresUser := "postgres"
		envPostgresUser := os.Getenv("POSTGRES_USER")
		if envPostgresUser != "" {
			postgresUser = envPostgresUser
		}

		postgresPassword := "pg"
		envPostgresPassword := os.Getenv("POSTGRES_PASSWORD")
		if envPostgresPassword != "" {
			postgresPassword = envPostgresPassword
		}

		postgresDatabase := "postgres"
		envPostgresDatabase := os.Getenv("POSTGRES_DATABASE")
		if envPostgresDatabase != "" {
			postgresDatabase = envPostgresDatabase
		}

		doCountFunc = func() (int, error) {
			return backend_postgres.DoCountPostgres(
				postgresHost, 5432, postgresUser, postgresPassword, postgresDatabase, hostname,
			)
		}
		getCountFunc = func() (int, error) {
			return backend_postgres.GetCountPostgres(
				postgresHost, 5432, postgresUser, postgresPassword, postgresDatabase, hostname,
			)
		}
	} else if backend == "mongodb" {
		mongodbURI := "mongodb://127.0.0.1:27017"
		envMongodbURI := os.Getenv("MONGODB_URI")
		if envMongodbURI != "" {
			mongodbURI = envMongodbURI
		}
		doCountFunc = func() (int, error) { return backend_mongodb.DoCountMongoDB(mongodbURI, hostname) }
		getCountFunc = func() (int, error) { return backend_mongodb.GetCountMongoDB(mongodbURI, hostname) }
	} else {
		log.Fatalf(`no backend "%s" exists, you can use "redis" (default), "postgres", or "inmemory"\n`, backend)
	}

	BaseServer(
		doCountFunc,
		getCountFunc,
	)
}
