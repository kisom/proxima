package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"git.sr.ht/~kisom/proxima/database"
	"git.sr.ht/~kisom/proxima/physics"
	"git.sr.ht/~kisom/proxima/rat"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Upstream struct {
	ctx context.Context
	db  *pgxpool.Pool
}

func NewUpstream() (*Upstream, error) {
	ctx := context.Background()
	db, err := database.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &Upstream{
		ctx: ctx,
		db:  db,
	}, nil
}

func (srv *Upstream) postUpdate(w http.ResponseWriter, r *http.Request) {
	log.Println("POST update")
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if !CheckBasicAuth(r) {
		log.Println("unauthenticated POST")
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			log.Println("BODY:")
			log.Println(string(body))
		}
		http.Error(w, "not allowed", http.StatusUnauthorized)
		return
	}

	update, err := UpdateFromReader(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = update.Store(srv.ctx, srv.db)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dur := rat.DurationSeconds(rat.Float(update.Elapsed))
	log.Println("update received:", physics.TimeString(dur))
	w.Write([]byte("OK"))
}

func (srv *Upstream) serveIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("serve index")
	created, update, err := FetchLastUpdate(srv.ctx, srv.db)
	if err != nil {
		if err != pgx.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	page := PageFromUpdate(created, update)
	buf := &bytes.Buffer{}
	err = indexTemplate.ExecuteTemplate(buf, "index", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.Copy(w, buf)
}

func (srv *Upstream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		srv.postUpdate(w, r)
	case http.MethodGet:
		srv.serveIndex(w, r)
	default:
		http.Error(w, fmt.Sprintf("method %s not implemented", r.Method),
			http.StatusBadRequest)
	}

}
