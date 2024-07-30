package cmd

import (
	"bytes"
	"context"
	"fmt"
	"net/url"

	"github.com/VladislavLisovenko/antibruteforce/internal/httpclient"
	"github.com/VladislavLisovenko/antibruteforce/internal/server"
	"github.com/spf13/cobra"
)

var resetCommand = &cobra.Command{
	Use:   "reset block",
	Short: "Reset block by login and ip",
	RunE: func(_ *cobra.Command, _ []string) error {
		var err error
		var b []byte
		hc := httpclient.New(cfg.HostInfo.Host)
		vs := url.Values{}
		vs.Set(server.IPField, network)
		vs.Set(server.LoginField, login)

		b, err = hc.Get(context.Background(), "reset", vs)
		if err != nil {
			return err
		}

		if !bytes.Equal(b, successResponse) {
			return fmt.Errorf("wrong response. Body: %s", b)
		}

		fmt.Println("Block successfully reset")
		return nil
	},
}

func init() {
	resetCommand.Flags().StringVar(&network, "n", "", "IP address")
	resetCommand.Flags().StringVar(&login, "l", "", "Login")
}
