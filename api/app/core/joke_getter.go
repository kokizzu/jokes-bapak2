package core

import (
	"context"
	"errors"
	"math/rand"

	"github.com/allegro/bigcache/v3"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pquerna/ffjson/ffjson"
)

// GetAllJSONJokes fetch the database for all the jokes then output it as a JSON []byte.
// Keep in mind, you will need to store it to memory yourself.
func GetAllJSONJokes(db *pgxpool.Pool, ctx *context.Context) ([]byte, error) {
	conn, err := db.Acquire(*ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var jokes []Joke
	results, err := conn.Query(*ctx, "SELECT \"id\",\"link\" FROM \"jokesbapak2\" ORDER BY \"id\"")
	if err != nil {
		return nil, err
	}
	defer results.Close()

	err = pgxscan.ScanAll(&jokes, results)
	if err != nil {
		return nil, err
	}

	data, err := ffjson.Marshal(jokes)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetRandomJokeFromCache returns a link string of a random joke from cache.
func GetRandomJokeFromCache(memory *bigcache.BigCache) (string, error) {
	jokes, err := memory.Get("jokes")
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}

	var data []Joke
	err = ffjson.Unmarshal(jokes, &data)
	if err != nil {
		return "", nil
	}

	// Return an error if the database is empty
	dataLength := len(data)
	if dataLength == 0 {
		return "", ErrEmpty
	}

	random := rand.Intn(dataLength)
	joke := data[random].Link

	return joke, nil
}

// CheckJokesCache checks if there is some value inside jokes cache.
func CheckJokesCache(memory *bigcache.BigCache) (bool, error) {
	_, err := memory.Get("jokes")
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// CheckTotalJokesCache literally does what the name is for
func CheckTotalJokesCache(memory *bigcache.BigCache) (bool, error) {
	_, err := memory.Get("total")
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetCachedJokeByID returns a link string of a certain ID from cache.
func GetCachedJokeByID(memory *bigcache.BigCache, id int) (string, error) {
	jokes, err := memory.Get("jokes")
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}

	var data []Joke
	err = ffjson.Unmarshal(jokes, &data)
	if err != nil {
		return "", nil
	}

	// This is a simple solution, might convert it to goroutines and channels sometime soon.
	for _, v := range data {
		if v.ID == id {
			return v.Link, nil
		}
	}

	return "", nil
}

// GetCachedTotalJokes
func GetCachedTotalJokes(memory *bigcache.BigCache) (int, error) {
	total, err := memory.Get("total")
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	return int(total[0]), nil
}