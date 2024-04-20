package cmd

import (
	cobraConfig "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/spf13/cobra"
)

var Bckupr = &cobra.Command{
	Use:   "bckupr",
	Short: "Docker volume/filesystem backup manager.",
	Long: `Bckupr is a tool to manage backups  your docker volumes and filesystem.
This application automates creating new backups and restoring data.`,
	RunE: runDaemon,
}

func init() {
	cobraConfig.InitGlobal(Bckupr)
	addGroup("backups", "Backup Commands:", Bckupr,
		Backup,
		Restore,
		List,
		Delete,
		Rotate,
	)

	Bckupr.AddCommand(Debug)
	Bckupr.AddCommand(Version)
}

func addGroup(id string, title string, root *cobra.Command, subs ...*cobra.Command) {
	root.AddGroup(&cobra.Group{
		ID:    id,
		Title: title,
	})
	for _, sub := range subs {
		sub.GroupID = id
		root.AddCommand(sub)
	}
}
