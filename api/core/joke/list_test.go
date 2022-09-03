package joke_test

import (
	"context"
	"jokes-bapak2-api/core/joke"
	"testing"
	"time"
)

func TestListJokeFromBucket(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	jokes, err := joke.ListJokesFromBucket(ctx, bucket, cache)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(jokes) != 5 {
		t.Errorf("expected joke to have a length of 5, instead got %d", len(jokes))
	}
}
