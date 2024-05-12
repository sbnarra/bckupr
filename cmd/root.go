package cmd

import (
	"github.com/sbnarra/bckupr/cmd/backup"
	"github.com/sbnarra/bckupr/cmd/daemon"
	"github.com/sbnarra/bckupr/cmd/delete"
	"github.com/sbnarra/bckupr/cmd/list"
	"github.com/sbnarra/bckupr/cmd/restore"
	"github.com/sbnarra/bckupr/cmd/rotate"
	"github.com/sbnarra/bckupr/cmd/version"
	"github.com/sbnarra/bckupr/internal/cmd/flags"
	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "bckupr",
	Short: "Docker volume/filesystem backup manager.",
	Long: `Bckupr is a tool to manage backups  your docker volumes and filesystem.
This application automates creating new backups and restoring data.`,
}

func init() {
	flags.Register(keys.Debug, Cmd.PersistentFlags())
	flags.Register(keys.ThreadLimit, Cmd.PersistentFlags())

	addGroup("daemon", "Server Commands:", Cmd,
		daemon.Cmd,
	)
	addGroup("backups", "Backup Commands:", Cmd,
		backup.Cmd,
		restore.Cmd,
		delete.Cmd,

		list.Cmd,
		rotate.Cmd,
	)

	Cmd.AddCommand(version.Cmd)
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
