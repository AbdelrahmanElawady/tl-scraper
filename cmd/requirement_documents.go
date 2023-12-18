package cmd

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/threefoldtech/tl-scrapper/endpoints"
)

const requirementDocumentsDir = "requirement_documents"

// requirementDocumentsCmd represents the requirement-documents command
var requirementDocumentsCmd = &cobra.Command{
	Use:   "requirement-documents",
	Short: "Scrap requirement-documents from TestLodge APIs",
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
		requirementDocumentsPath := filepath.Join(outputDir, requirementDocumentsDir)
		err = os.MkdirAll(requirementDocumentsPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		_, err = caller.GetRequirementDocuments(cmd.Context(), projects, &requirementDocumentsPath)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	},
}

func init() {
	scrapCmd.AddCommand(requirementDocumentsCmd)

	requirementDocumentsCmd.Flags().IntSliceP("projects", "p", []int{}, "specify the projects to get the test requirement documents from")
}
