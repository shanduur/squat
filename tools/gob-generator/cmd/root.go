package cmd

import (
	"log"

	"github.com/shanduur/squat/generator"
	"github.com/spf13/cobra"
)

var (
	in  string
	out string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gob-generator",
	Short: "Tool for generating gob file from JSON",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := generator.ReadDump(in, out); err != nil {
			log.Fatalf("failed to read %s and dump %s: %s", in, out, err.Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&in, "input-json", "i", "data.json", "input file")
	rootCmd.Flags().StringVarP(&out, "output-file", "o", "data.gob", "output file")
}
