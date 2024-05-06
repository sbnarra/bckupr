package pkg

import (
	"context"
	"fmt"

	"github.com/sbnarra/bckupr/pkg/client"
)

func djk() {
	c := client.NewAPIClient(client.NewConfiguration())
	backup, res, err := c.BackupAPI.
		CreateBackup(context.TODO()).
		CreateBackup(client.CreateBackup{}).
		Execute()
	if err != nil {

	}
	fmt.Println(backup.Id)
	fmt.Println(res.Header.Get("dryRun"))
}
