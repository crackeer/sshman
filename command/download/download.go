package download

import (
	"github.com/spf13/cobra"
)

var (
	targetDir string
)

// NewDownloadCommand
//
//	@return *cobra.Command
func NewDownloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "download ",
		Run:   download,
		Args:  cobra.MinimumNArgs(1),
	}

	return cmd
}

func download(cmd *cobra.Command, args []string) {

}
