package cmd

import (
	"os"
	"path/filepath"

	"github.com/AbdelrahmanElawady/tl-scraper/endpoints"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const plansDir = "plans"

// plansCmd represents the plans command
var plansCmd = &cobra.Command{
	Use:   "plans",
	Short: "Scrap test plans from TestLodge APIs",
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
		plansPath := filepath.Join(outputDir, plansDir)
		err = os.MkdirAll(plansPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		_, err = caller.GetPlans(cmd.Context(), projects, &plansPath)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	},
}

func init() {
	scrapCmd.AddCommand(plansCmd)

	plansCmd.Flags().IntSliceP("projects", "p", []int{}, "specify the projects to get the test suites from")
}
