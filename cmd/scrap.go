package cmd

import (
	"errors"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// scrapCmd represents the scrap command
var scrapCmd = &cobra.Command{
	Use:   "scrap",
	Short: "A brief description of your command",
}

func init() {
	rootCmd.AddCommand(scrapCmd)

	scrapCmd.PersistentFlags().StringP("env", "e", ".env", "key-value file with the user credentials")

	scrapCmd.PersistentFlags().StringP("output-dir", "o", "", "directory to dump api results in")
}

func getUserData(cmd *cobra.Command) (key, email, accountID string, err error) {
	var env string
	env, err = cmd.Flags().GetString("env")
	if err != nil {
		return
	}

	config, err := godotenv.Read(env)
	if err != nil {
		return
	}
	if value, ok := config["key"]; !ok || value == "" {
		err = errors.New("missing 'key' field")
		return
	} else {
		key = value
	}
	if value, ok := config["email"]; !ok || value == "" {
		err = errors.New("missing 'email' field")
		return
	} else {
		email = value
	}
	if value, ok := config["accountID"]; !ok || value == "" {
		err = errors.New("missing 'accountID' field")
		return
	} else {
		accountID = value
	}
	return
}
