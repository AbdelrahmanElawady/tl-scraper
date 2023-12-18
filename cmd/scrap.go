package cmd

import (
	"github.com/spf13/cobra"
)

// scrapCmd represents the scrap command
var scrapCmd = &cobra.Command{
	Use:   "scrap",
	Short: "A brief description of your command",
}

func init() {
	rootCmd.AddCommand(scrapCmd)

	scrapCmd.PersistentFlags().StringP("key", "k", "", "key used in api authentication")
	scrapCmd.PersistentFlags().StringP("email", "e", "", "email used in api authentication")
	scrapCmd.PersistentFlags().StringP("account-id", "a", "", "account id used in api authentication")
	err := scrapCmd.MarkPersistentFlagRequired("key")
	if err != nil {
		panic(err)
	}
	err = scrapCmd.MarkPersistentFlagRequired("email")
	if err != nil {
		panic(err)
	}
	err = scrapCmd.MarkPersistentFlagRequired("account-id")
	if err != nil {
		panic(err)
	}
	scrapCmd.PersistentFlags().StringP("output-dir", "o", "", "directory to dump api results in")
}

func getUserData(cmd *cobra.Command) (key, email, accountID string, err error) {
	key, err = cmd.Flags().GetString("key")
	if err != nil {
		return
	}
	email, err = cmd.Flags().GetString("email")
	if err != nil {
		return
	}
	accountID, err = cmd.Flags().GetString("account-id")
	if err != nil {
		return
	}
	return
}
