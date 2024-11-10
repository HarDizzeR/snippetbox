package main

import (
    "database/sql"
	"flag"
	"net/http"
	"log/slog"
    "html/template"
	"os"
    "time"
    "crypto/tls"

    "snippetbox.minaasaad.net/internal/models"

    "github.com/alexedwards/scs/v2"
    "github.com/alexedwards/scs/mysqlstore"
    "github.com/go-playground/form/v4"
    _ "github.com/go-sql-driver/mysql"
)

type application struct {
    logger *slog.Logger
    snippets *models.SnippetModel
    templateCache map[string]*template.Template
    formDecoder   *form.Decoder
    sessionManager *scs.SessionManager
    users *models.UserModel
}


func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")
    dsn := flag.String("dsn", "web:1234@/snippetbox?parseTime=true", "MySQL data source name")
    flag.Parse()

    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    db, err := openDB(*dsn)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

    defer db.Close()

    templateCache, err := newTemplateCache()
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

    formDecoder := form.NewDecoder()

    sessionManager := scs.New()
    sessionManager.Store = mysqlstore.New(db)
    sessionManager.Lifetime = 12 * time.Hour

    app := &application{
        logger: logger,
        snippets: &models.SnippetModel{DB: db},
        users: &models.UserModel{DB: db},
        templateCache: templateCache,
        formDecoder: formDecoder,
        sessionManager: sessionManager,
    }

    tlsConfig := &tls.Config{
        CurvePreferences: []tls.CurveID{
            tls.X25519,tls.CurveP256,
        },
    }

    srv := &http.Server{
        Addr: *addr,
        Handler: app.routes(),
        ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
        TLSConfig: tlsConfig,
        IdleTimeout: time.Minute,
        ReadTimeout: 5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    logger.Info("starting server", "addr", *addr)
    
    err = srv.ListenAndServeTLS("./security/cert.pem", "./security/key.pem")
    logger.Error(err.Error())
    os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        db.Close()
        return nil, err
    }

    return db, nil
}