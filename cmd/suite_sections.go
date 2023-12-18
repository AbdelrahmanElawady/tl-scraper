package cmd

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/threefoldtech/tl-scrapper/endpoints"
)

const suiteSectionsDir = "suite_sections"

// suiteSectionsCmd represents the suite-sections command
var suiteSectionsCmd = &cobra.Command{
	Use:   "suite-sections",
	Short: "Scrap suite sections from TestLodge APIs",
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
		suiteSectionsPath := filepath.Join(outputDir, suiteSectionsDir)
		err = os.MkdirAll(suiteSectionsPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		_, err = caller.GetSuiteSections(cmd.Context(), suiteIDs, &suiteSectionsPath)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	},
}

func init() {
	scrapCmd.AddCommand(suiteSectionsCmd)

	suiteSectionsCmd.Flags().IntSliceP("projects", "p", []int{}, "specify the projects to get the test suite sections from")
}
