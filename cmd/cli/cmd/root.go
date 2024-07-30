package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/VladislavLisovenko/antibruteforce/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfg            config.Config
	allowedMethods = map[string]struct{}{
		"add":    {},
		"remove": {},
	}
	network           string
	login             string
	errNotFoundMethod = errors.New("method not found")
	successResponse   = []byte("{\"ok\":true}")
)

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "Antibruteforce client",
}

func init() {
	var err error
	cfg, err = config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd.AddCommand(blackListCommand)
	rootCmd.AddCommand(whiteListCommand)
	rootCmd.AddCommand(resetCommand)

	whiteListCommand.Flags().StringVar(&network, "n", "", "IP address with mask")
	blackListCommand.Flags().StringVar(&network, "n", "", "IP address with mask")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkResponse(b []byte, err error) error {
	if err != nil {
		return err
	}

	if !bytes.Equal(b, successResponse) {
		return fmt.Errorf("wrong response. Body: %s", b)
	}

	return nil
}
