package cmd

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/threefoldtech/tl-scrapper/endpoints"
)

const requirementsDir = "requirements"

// requirementsCmd represents the requirements command
var requirementsCmd = &cobra.Command{
	Use:   "requirements",
	Short: "Scrap requirements from TestLodge APIs",
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
		rdIDs, err := caller.GetRequirementDocuments(cmd.Context(), projects, nil)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		requirementsPath := filepath.Join(outputDir, requirementsDir)
		err = os.MkdirAll(requirementsPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		_, err = caller.GetRequirements(cmd.Context(), rdIDs, &requirementsPath)
		if err != nil {
			log.Fatal().Err(err).Send()
		}

	},
}

func init() {
	scrapCmd.AddCommand(requirementsCmd)

	requirementsCmd.Flags().IntSliceP("projects", "p", []int{}, "specify the projects to get the test requirements from")
}
