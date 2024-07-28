package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/VladislavLisovenko/antibruteforce/internal/httpclient"
	"github.com/VladislavLisovenko/antibruteforce/internal/server"
	"github.com/spf13/cobra"
)

var blackListCommand = &cobra.Command{
	Use:   "blacklist",
	Short: "Add/remove from blacklist",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		method := args[0]
		if _, ok := allowedMethods[method]; !ok {
			return errNotFoundMethod
		}

		var err error
		var b []byte
		hc := httpclient.New(cfg.Host)
		vs := url.Values{}
		vs.Set(server.IPField, network)

		if method == "add" {
			b, err = hc.Post(context.Background(), "blackList", vs)
		} else {
			b, err = hc.Delete(context.Background(), "blackList", vs)
		}

		if err := checkResponse(b, err); err != nil {
			return err
		}

		fmt.Println("Address is successfully added to blacklist")
		return nil
	},
}
