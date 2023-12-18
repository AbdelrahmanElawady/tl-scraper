package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/threefoldtech/tl-scrapper/endpoints"
)

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Scrap all test data from TestLodge APIs",
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
		projectsPath := filepath.Join(outputDir, projectsDir)
		suitesPath := filepath.Join(outputDir, suitesDir)
		suiteSectionsPath := filepath.Join(outputDir, suiteSectionsDir)
		casesPath := filepath.Join(outputDir, testCasesDir)
		requirementsPath := filepath.Join(outputDir, requirementsDir)
		requirementDocumentsPath := filepath.Join(outputDir, requirementDocumentsDir)
		plansPath := filepath.Join(outputDir, plansDir)
		planContentsPath := filepath.Join(outputDir, planContentsDir)
		err = os.MkdirAll(projectsPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		err = os.MkdirAll(suitesPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		err = os.MkdirAll(suiteSectionsPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		err = os.MkdirAll(casesPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		err = os.MkdirAll(requirementsPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		err = os.MkdirAll(requirementDocumentsPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		err = os.MkdirAll(plansPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		err = os.MkdirAll(planContentsPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		if len(projects) == 0 {
			projects, err = caller.GetProjects(cmd.Context(), &projectsPath)
			if err != nil {
				log.Error().Err(err).Send()
			}
		}
		suiteIDs, err := caller.GetSuites(cmd.Context(), projects, &suitesPath)
		if err != nil {
			log.Error().Err(err).Send()
		}
		suiteSectionIDs, err := caller.GetSuiteSections(cmd.Context(), suiteIDs, &suiteSectionsPath)
		if err != nil {
			log.Error().Err(err).Send()
		}
		err = caller.GetTestCases(context.Background(), suiteSectionIDs, &casesPath)
		if err != nil {
			log.Error().Err(err).Send()
		}
		rdIDs, err := caller.GetRequirementDocuments(cmd.Context(), projects, &requirementDocumentsPath)
		if err != nil {
			log.Error().Err(err).Send()
		}
		_, err = caller.GetRequirements(cmd.Context(), rdIDs, &requirementsPath)
		if err != nil {
			log.Error().Err(err).Send()
		}
		planIDs, err := caller.GetPlans(cmd.Context(), projects, &plansPath)
		if err != nil {
			log.Error().Err(err).Send()
		}
		_, err = caller.GetPlanContents(cmd.Context(), planIDs, &planContentsPath)
		if err != nil {
			log.Error().Err(err).Send()
		}
	},
}

func init() {
	scrapCmd.AddCommand(allCmd)

	allCmd.Flags().IntSliceP("projects", "p", []int{}, "specify the projects to get the test cases from")
}
