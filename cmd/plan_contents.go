package cmd

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/threefoldtech/tl-scrapper/endpoints"
)

const planContentsDir = "plan_contents"

// planContentsCmd represents the plan-contents command
var planContentsCmd = &cobra.Command{
	Use:   "plan-contents",
	Short: "Scrap test plan contents from TestLodge APIs",
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
		planIDs, err := caller.GetPlans(cmd.Context(), projects, nil)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		planContentsPath := filepath.Join(outputDir, planContentsDir)
		err = os.MkdirAll(planContentsPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		_, err = caller.GetPlanContents(cmd.Context(), planIDs, &planContentsPath)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	},
}

func init() {
	scrapCmd.AddCommand(planContentsCmd)

	planContentsCmd.Flags().IntSliceP("projects", "p", []int{}, "specify the projects to get the test suite sections from")
}
