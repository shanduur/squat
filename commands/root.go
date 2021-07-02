package commands

import (
	"log"

	"github.com/shanduur/squat/server"
	"github.com/spf13/cobra"
)

var (
	port string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "squat",
	Short: "Server that provides simple SQL data generation functionality",
	Long:  `This is server for SQL data generation app.`,
	Run: func(cmd *cobra.Command, args []string) {
		srv := server.New(port)

		log.Printf("Server is listening on %s", port)
		if err := srv.Run(); err != nil {
			log.Fatalf("server exited unexpectedly: %s\n", err.Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringVarP(&port, "port", "p", ":8080", "port on which REST API listenes")
}
