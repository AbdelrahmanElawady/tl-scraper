package cmd

import (
	"os"
	"path/filepath"

	"github.com/AbdelrahmanElawady/tl-scraper/endpoints"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const projectsDir = "projects"

// projectsCmd represents the projects command
var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Scrap projects from TestLodge APIs",
	Run: func(cmd *cobra.Command, args []string) {
		key, email, accountID, err := getUserData(cmd)
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		outputDir, err := cmd.Flags().GetString("output-dir")
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		caller, err := endpoints.NewCaller(key, email, accountID)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		projectsPath := filepath.Join(outputDir, projectsDir)
		err = os.MkdirAll(projectsPath, 0777)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		_, err = caller.GetProjects(cmd.Context(), &projectsPath)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	},
}

func init() {
	scrapCmd.AddCommand(projectsCmd)

}
