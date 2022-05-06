package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/sethvargo/go-retry"
)

func WaitForServer(url string) (err error) {

	var b retry.Backoff

	b, err = retry.NewFibonacci(500 * time.Millisecond)
	if err == nil {

		// Stop after 4 retries, when the 5th attempt has failed. In this example, the worst case elapsed
		// time would be 1s + 1s + 2s + 3s = 7s.
		b = retry.WithMaxRetries(3, b)

		ctx := context.Background()
		retry.Do(ctx, b, func(ctx context.Context) error {

			var r *http.Response
			r, err = http.Get(url)
			if r != nil {
				if r.StatusCode == http.StatusOK {
					r.Body.Close()
					return nil
				}
				r.Body.Close()
			}

			return retry.RetryableError(err)

		})
	}

	return
}
