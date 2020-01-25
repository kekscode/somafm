package cmd

import (
	"github.com/kekscode/somafm/channels"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play <channel-id>",
	Short: "Play from soma.fm",
	Long:  `This starts playing music based on a given channel id`,
	Run: func(cmd *cobra.Command, args []string) {
		ch := channels.NewChannelList()
		chID, err := cmd.Flags().GetString("channel-id")
		if err != nil {
			log.Errorf("Error: %s", err)
		}
		ch.PlayChannel(chID)
	},
}

func init() {
	rootCmd.AddCommand(playCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	playCmd.Flags().StringP("channel-id", "c", "groovesalad", "The channel you want to tune into")
}
