package main

import (
	"github.com/sbnarra/bckupr/cmd"
	cobraConfig "github.com/sbnarra/bckupr/internal/config/cobra"
	"github.com/spf13/cobra"
)

var bckupr = &cobra.Command{
	Use:   "bckupr",
	Short: "Docker volume/filesystem backup manager.",
	Long: `Bckupr is a tool to manage backups  your docker volumes and filesystem.
This application automates creating new backups and restoring data.`,
}

func main() {
	bckupr.Execute()
}

func init() {
	cobraConfig.InitGlobal(bckupr)
	addGroup("backups", "Backup Commands:", bckupr,
		cmd.Backup,
		cmd.Restore,
		cmd.List,
		cmd.Delete,
		cmd.Rotate,
	)
	addGroup("daemons", "Daemons Commands:", bckupr,
		cmd.Daemon,
		cmd.Cron)

	bckupr.AddCommand(cmd.Debug)
	bckupr.AddCommand(cmd.Version)
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
