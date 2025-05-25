package main // ЭТО ДОЛЖНА БЫТЬ ПЕРВАЯ СТРОКА!

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"minityweb/backend/pkg/api"    // Убедитесь, что путь к модулю правильный
	"minityweb/backend/pkg/store"  // Убедитесь, что путь к модулю правильный

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type AppConfig struct {
	Port      string
	DBDriver  string
	DSN       string
	JWTSecret string
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		if fallback == "" && key != "DB_PASSWORD" {
			log.Printf("WARNING: Environment variable %s is not set and no fallback is provided.", key)
		} else if fallback == "" && key == "DB_PASSWORD" {
            log.Printf("INFO: DB_PASSWORD environment variable is not set. Using fallback if defined, or it might be empty.")
        }
		return fallback
	}
	return value
}

func LoadConfig() AppConfig {
	dbHost := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "minity_user")
	dbPassword := getEnv("DB_PASSWORD", "minity_password")
	dbName := getEnv("DB_NAME", "minity_db")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	jwtSecret := getEnv("JWT_SECRET", "your-very-secret-and-long-key-for-minity-app") // ЗАМЕНИТЕ ИЛИ УСТАНОВИТЕ ПЕРЕМЕННУЮ ОКРУЖЕНИЯ
	if jwtSecret == "your-very-secret-and-long-key-for-minity-app" {
		log.Println("WARNING: Using default JWT_SECRET. Please set a strong JWT_SECRET environment variable for production.")
	}

	return AppConfig{
		Port:      getEnv("PORT", "8080"),
		DBDriver:  getEnv("DB_DRIVER", "postgres"),
		DSN:       dsn,
		JWTSecret: jwtSecret,
	}
}

func main() { // ЭТО ТОЧКА ВХОДА
	cfg := LoadConfig()

	log.Println("Starting Minity server...")
	log.Printf("Port: %s, DB Driver: %s", cfg.Port, cfg.DBDriver)

	var db *sql.DB
	var err error

	maxRetries := 15
	retryInterval := 5 * time.Second

	log.Println("Attempting to connect to the database...")
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open(cfg.DBDriver, cfg.DSN)
		if err != nil {
			log.Printf("Error opening database connection (attempt %d/%d): %v. Retrying in %s...", i+1, maxRetries, err, retryInterval)
			time.Sleep(retryInterval)
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to the database.")
			break 
		}

		log.Printf("Failed to ping database (attempt %d/%d): %v. Retrying in %s...", i+1, maxRetries, err, retryInterval)
		if i < maxRetries-1 {
			db.Close() 
			time.Sleep(retryInterval)
		}
	}

	if err != nil {
		log.Fatalf("FATAL: Could not connect to database after %d attempts: %v", maxRetries, err)
	}
	defer db.Close()

	var appStore store.DataStore
	if cfg.DBDriver == "postgres" {
		appStore = store.NewPGStore(db)
	} else {
		log.Fatalf("Unsupported DB driver: %s", cfg.DBDriver)
	}

	apiHandler := api.NewAPIHandler(appStore, cfg.JWTSecret)

	router := mux.NewRouter()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("FATAL: Could not get working directory: %v", err)
	}
	log.Printf("Server's current working directory (os.Getwd()): %s", wd)

	staticDirRelative := "../frontend/" 
	absStaticDir, err := filepath.Abs(filepath.Join(wd, staticDirRelative))
	if err != nil {
		log.Fatalf("FATAL: Could not determine absolute path for static directory: %v", err)
	}
	log.Printf("Absolute static directory resolved to: %s", absStaticDir)

	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir(filepath.Join(absStaticDir, "css")))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir(filepath.Join(absStaticDir, "js")))))
	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir(filepath.Join(absStaticDir, "images")))))

	htmlFiles := []string{"index.html", "establishment.html", "admin.html"}
	for _, fileName := range htmlFiles {
		func(file string) {
			router.HandleFunc("/"+file, func(w http.ResponseWriter, r *http.Request) {
				filePath := filepath.Join(absStaticDir, file)
				if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
						http.NotFound(w, r)
						return
				} else if statErr != nil {
						http.Error(w, "Internal server error", http.StatusInternalServerError)
						return
				}
				http.ServeFile(w, r, filePath)
			}).Methods(http.MethodGet)
		}(fileName)
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join(absStaticDir, "index.html") 
		if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
			http.NotFound(w, r)
			return
		} else if statErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.ServeFile(w, r, filePath)
	}).Methods(http.MethodGet)

	apiHandler.RegisterRoutes(router)

	serverAddr := ":" + cfg.Port
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server listening on http://localhost%s", serverAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", serverAddr, err)
	}
}