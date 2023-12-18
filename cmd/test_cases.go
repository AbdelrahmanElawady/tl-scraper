package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/threefoldtech/tl-scrapper/endpoints"
)

const testCasesDir = "test_cases"

// testCasesCmd represents the test-cases command
var testCasesCmd = &cobra.Command{
	Use:   "test-cases",
	Short: "Scrap test cases from TestLodge APIs",
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
		suiteIDs, err := caller.GetSuites(cmd.Context(), projects, nil)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		suiteSectionIDs, err := caller.GetSuiteSections(cmd.Context(), suiteIDs, nil)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		casesPath := filepath.Join(outputDir, testCasesDir)

		err = os.MkdirAll(casesPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		err = caller.GetTestCases(context.Background(), suiteSectionIDs, &casesPath)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	},
}

func init() {
	scrapCmd.AddCommand(testCasesCmd)

	testCasesCmd.Flags().IntSliceP("projects", "p", []int{}, "specify the projects to get the test cases from")
}
