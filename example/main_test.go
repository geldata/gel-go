package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/geldata/gel-go/gelcfg"
	"github.com/geldata/gel-go/geltypes"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func gelCLI(args []string, t *testing.T) string {
	cmd := exec.Command("gel", args...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("gelCLI failed: %v\nOutput:\n%s", err, output)
	}
	return strings.TrimSpace(string(output))
}

func TestExampleApp(t *testing.T) {
	t.Cleanup(func() {
		gelCLI([]string{"instance", "destroy", "-I", "example", "--force"}, t)
	})
	gelCLI([]string{"init", "--server-instance", "example", "--non-interactive"}, t)
	cfgDir := gelCLI([]string{"info", "--get", "config-dir"}, t)
	configFile := filepath.Join(cfgDir, "credentials", "example.json")
	creds, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	pw := gjson.Get(string(creds), "password").String()
	if pw == "" {
		t.Fatal("Password is empty, ensure the credentials file is set up correctly")
	}
	port := gjson.Get(string(creds), "port").Int()
	if port == 0 {
		t.Fatal("Port is not set, ensure the credentials file is set up correctly")
	}
	t.Setenv("GEL_DSN", fmt.Sprintf("gel://admin:%s@localhost:%d/main", pw, port))

	app := NewApp(gelcfg.Options{
		TLSOptions: gelcfg.TLSOptions{
			SecurityMode: gelcfg.TLSModeInsecure,
		},
	})
	server := httptest.NewServer(app.Routes())
	defer server.Close()

	var movieID geltypes.UUID

	t.Run("Create Movie", func(t *testing.T) {
		resp, err := http.Post(server.URL+"/movie", "application/json", strings.NewReader(`{"title":"Inception","year":2010,"description":"A mind-bending thriller"}`))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
		var movie Movie
		err = json.NewDecoder(resp.Body).Decode(&movie)
		assert.NoError(t, err)
		title, ok := movie.Title.Get()
		assert.True(t, ok)
		assert.Equal(t, "Inception", title)
		year, ok := movie.Year.Get()
		assert.True(t, ok)
		assert.Equal(t, int64(2010), year)
		desc, ok := movie.Description.Get()
		assert.True(t, ok)
		assert.Equal(t, "A mind-bending thriller", desc)
		movieID = movie.ID
	})
	t.Run("Get Movies", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/movies")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
		var movies []Movie
		err = json.NewDecoder(resp.Body).Decode(&movies)
		assert.NoError(t, err)
		assert.NotEmpty(t, movies)
		title, ok := movies[0].Title.Get()
		assert.True(t, ok)
		assert.Equal(t, "Inception", title)
		year, ok := movies[0].Year.Get()
		assert.True(t, ok)
		assert.Equal(t, int64(2010), year)
		desc, ok := movies[0].Description.Get()
		assert.True(t, ok)
		assert.Equal(t, "A mind-bending thriller", desc)
	})
	t.Run("Update Movie", func(t *testing.T) {
		req, err := http.NewRequest("PUT", server.URL+"/movie", strings.NewReader(fmt.Sprintf(`{"id":"%s","title":"Inception","year":2010,"description":"A mind-bending thriller with a twist"}`, movieID.String())))
		assert.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
		var movie Movie
		err = json.NewDecoder(resp.Body).Decode(&movie)
		assert.NoError(t, err)
		title, ok := movie.Title.Get()
		assert.True(t, ok)
		assert.Equal(t, "Inception", title)
		year, ok := movie.Year.Get()
		assert.True(t, ok)
		assert.Equal(t, int64(2010), year)
		desc, ok := movie.Description.Get()
		assert.True(t, ok)
		assert.Equal(t, "A mind-bending thriller with a twist", desc)
	})
	t.Run("Delete Movie", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/movie?id=%s", server.URL, movieID.String()), nil)
		assert.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.NoError(t, err)
	})
	t.Run("Get Movies After Deletion", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/movies")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()
		var movies []Movie
		err = json.NewDecoder(resp.Body).Decode(&movies)
		assert.NoError(t, err)
		assert.Empty(t, movies, "Expected no movies after deletion")
	})
}
