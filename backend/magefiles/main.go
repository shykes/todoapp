//go:build mage

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dagger/cloak/engine"
)

func Deploy(ctx context.Context) {
	if err := engine.Start(ctx, &engine.Config{}, func(ctx engine.Context, c *TodoAppClient) error {
		// User can configure netlify site name with $NETLIFY_SITE_NAME
		siteName, ok := os.LookupEnv("NETLIFY_SITE_NAME")
		if !ok {
			user, _ := os.LookupEnv("USER")
			siteName = fmt.Sprintf("%s-dagger-todoapp", user)
		}
		fmt.Printf("Using Netlify site name %q\n", siteName)

		// Deploy using the todoapp deploy extension
		resp, err := c.Deploy(ctx, siteName)
		if err != nil {
			return err
		}

		// Print deployment info to the user
		fmt.Println("URL:", resp.Todoapp.Deploy)

		return nil
	}); err != nil {
		panic(err)
	}
}
