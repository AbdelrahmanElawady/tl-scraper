package cmd

import (
	"os"
	"path/filepath"

	"github.com/AbdelrahmanElawady/tl-scraper/endpoints"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const suitesDir = "suites"

// suitesCmd represents the suites command
var suitesCmd = &cobra.Command{
	Use:   "suites",
	Short: "Scrap suites from TestLodge APIs",
	Run: func(cmd *cobra.Command, args []string) {

		key, email, accountID, err := getUserData(cmd)
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		outputDir, err := cmd.Flags().GetString("output-dir")
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		projects, err := cmd.Flags().GetIntSlice("projects")
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		caller, err := endpoints.NewCaller(key, email, accountID)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		suitesPath := filepath.Join(outputDir, suitesDir)
		err = os.MkdirAll(suitesPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		_, err = caller.GetSuites(cmd.Context(), projects, &suitesPath)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	},
}

func init() {
	scrapCmd.AddCommand(suitesCmd)

	suitesCmd.Flags().IntSliceP("projects", "p", []int{}, "specify the projects to get the test suites from")
}
