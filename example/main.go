package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	gel "github.com/geldata/gel-go"
	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/gelerr"
	"github.com/geldata/gel-go/geltypes"
)

type Movie struct {
	ID          geltypes.UUID          `gel:"id"`
	Title       geltypes.OptionalStr   `gel:"title"`
	Year        geltypes.OptionalInt64 `gel:"year"`
	Description geltypes.OptionalStr   `gel:"description"`
}

type app struct {
	gelClient *gel.Client
}

func NewApp(opts gelcfg.Options) *app {
	client, err := gel.CreateClient(opts)
	if err != nil {
		log.Fatalf("Error creating Gel client: %v", err)
	}
	return &app{gelClient: client}
}

func (a *app) createMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	m := make(map[string]any)
	m["title"] = movie.Title
	m["year"] = movie.Year
	m["description"] = movie.Description

	query := `
		insert Movie {
			title := <str>$title,
			year := <int64>$year,
			description := <str>$description,
		};
	`
	if err := a.gelClient.QuerySingle(context.Background(), query, &movie, m); err != nil {
		http.Error(w, fmt.Sprintf("Error creating movie: %v", err), http.StatusInternalServerError)
		return
	}
	writeJSON(w, movie)
}

func (a *app) getMovies(w http.ResponseWriter, r *http.Request) {
	var movies []Movie
	err := a.gelClient.Query(context.Background(), "select Movie { ** } order by .year desc", &movies)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching movies: %v", err), http.StatusInternalServerError)
		return
	}
	if len(movies) == 0 {
		err = json.NewEncoder(w).Encode([]string{})
		if err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, movies)
}

func (a *app) updateMovie(w http.ResponseWriter, r *http.Request) {
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if movie.ID == (geltypes.UUID{}) {
		http.Error(w, "Movie ID is required", http.StatusBadRequest)
		return
	}

	query := `
		select (
			update Movie
			filter .id = <uuid>$id
			set {
				title := <optional str>$title ?? .title,
				year := <optional int64>$year ?? .year,
				description := <optional str>$description ?? .description,
			}
		) {id, title, year, description};
	`
	args := map[string]any{
		"id":          movie.ID,
		"year":        movie.Year,
		"title":       movie.Title,
		"description": movie.Description,
	}
	if err := a.gelClient.QuerySingle(context.Background(), query, &movie, args); err != nil {
		var ge gelerr.Error
		if errors.As(err, &ge) && ge.Category(gelerr.NoDataError) {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error updating movie: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, movie)
}

func (a *app) getMovie(w http.ResponseWriter, r *http.Request) {
	id, err := idFromReq(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	var movie Movie
	query := `
		select Movie { ** }
		filter .id = <uuid>$id;
	`
	args := map[string]any{"id": id}
	if err := a.gelClient.QuerySingle(context.Background(), query, &movie, args); err != nil {
		var ge gelerr.Error
		if errors.As(err, &ge) && ge.Category(gelerr.NoDataError) {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error fetching movie: %v", err), http.StatusInternalServerError)
		return
	}
	writeJSON(w, movie)
}

func (a *app) deleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := idFromReq(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	query := `
		delete Movie
		filter .id = <uuid>$id;
	`
	args := map[string]any{"id": id}
	if err := a.gelClient.QuerySingle(context.Background(), query, &Movie{}, args); err != nil {
		var ge gelerr.Error
		if errors.As(err, &ge) && ge.Category(gelerr.NoDataError) {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error deleting movie: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func idFromReq(r *http.Request) (geltypes.UUID, error) {
	id := r.URL.Query().Get("id")
	if id == "" {
		return geltypes.UUID{}, fmt.Errorf("missing 'id' parameter")
	}
	uuid, err := geltypes.ParseUUID(id)
	if err != nil {
		return geltypes.UUID{}, fmt.Errorf("invalid UUID: %v", err)
	}
	return uuid, nil
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

func (a *app) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /movie", a.createMovie)
	mux.HandleFunc("GET /movies", a.getMovies)
	mux.HandleFunc("GET /movie", a.getMovie)
	mux.HandleFunc("PUT /movie", a.updateMovie)
	mux.HandleFunc("DELETE /movie", a.deleteMovie)
	return mux
}

func main() {
	app := NewApp(gelcfg.Options{})
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", app.Routes())
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
