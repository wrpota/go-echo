package cmd

import (
	"log"

	"github.com/spf13/cobra"
	_ "github.com/wrpota/go-echo/init"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cmd",
		Short: "cmd is a tool for framework",
	}

	Env string
)

// Execute ..
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err.Error())
	}
}
