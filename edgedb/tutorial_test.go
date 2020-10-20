// This source file is part of the EdgeDB open source project.
//
// Copyright 2020-present EdgeDB Inc. and the EdgeDB authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package edgedb

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Person struct {
	FirstName string `edgedb:"first_name"`
	LastName  string `edgedb:"last_name"`
}

type Movie struct {
	Title    string   `edgedb:"title"`
	Year     int64    `edgedb:"year"`
	Director Person   `edgedb:"director"`
	Actors   []Person `edgedb:"actors"`
}

func TestTutorial(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())
	dbName := fmt.Sprintf("test%v", rand.Intn(10_000))
	err := client.Execute(ctx, "CREATE DATABASE "+dbName)
	require.Nil(t, err)

	edb, err := Connect(
		ctx,
		Options{
			Host:     server.Host,
			Port:     server.Port,
			User:     server.User,
			Password: server.Password,
			Database: dbName,
			admin:    server.admin,
		},
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = edb.Close()
		require.Nil(t, err)
	}()

	err = edb.Execute(ctx, `
		START MIGRATION TO {
			module default {
				type Movie {
					required property title -> str;
					# the year of release
					property year -> int64;
					required link director -> Person;
					multi link actors -> Person;
				}
				type Person {
					required property first_name -> str;
					required property last_name -> str;
				}
			}
		};

		POPULATE MIGRATION;

		COMMIT MIGRATION;
	`)
	require.Nil(t, err)

	err = edb.Execute(ctx, `
		INSERT Movie {
			title := 'Blade Runner 2049',
			year := 2017,
			director := (
				INSERT Person {
					first_name := 'Denis',
					last_name := 'Villeneuve',
				}
			),
			actors := {
				(INSERT Person {
					first_name := 'Harrison',
					last_name := 'Ford',
				}),
				(INSERT Person {
					first_name := 'Ryan',
					last_name := 'Gosling',
				}),
				(INSERT Person {
					first_name := 'Ana',
					last_name := 'de Armas',
				}),
			}
		}`,
	)
	require.Nil(t, err)

	err = edb.Execute(ctx, `
		INSERT Movie {
				title := 'Dune',
				director := (
						SELECT Person
						FILTER
								# the last name is sufficient
								# to identify the right person
								.last_name = 'Villeneuve'
						# the LIMIT is needed to satisfy the single
						# link requirement validation
						LIMIT 1
				)
		};`,
	)
	require.Nil(t, err)

	var out []Movie
	err = edb.Query(ctx, `
		SELECT Movie {
				title,
				year,
				director: {
						first_name,
						last_name
				},
				actors: {
						first_name,
						last_name
				}
		}`,
		&out,
	)
	require.Nil(t, err)

	expected := []Movie{
		{
			Title: "Blade Runner 2049",
			Year:  int64(2017),
			Director: Person{
				FirstName: "Denis",
				LastName:  "Villeneuve",
			},
			Actors: []Person{
				{
					FirstName: "Harrison",
					LastName:  "Ford",
				},
				{
					FirstName: "Ryan",
					LastName:  "Gosling",
				},
				{
					FirstName: "Ana",
					LastName:  "de Armas",
				},
			},
		},
		{
			Title: "Dune",
			Director: Person{
				FirstName: "Denis",
				LastName:  "Villeneuve",
			},
			Actors: []Person{},
		},
	}

	assert.Equal(t, expected, out)
}
