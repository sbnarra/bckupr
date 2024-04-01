go run . -h > docs/snippets/cli-help.txt
go run . version -h > docs/snippets/cli-version-help.txt

go run . backup -h > docs/snippets/cli-backup-help.txt
go run . restore -h > docs/snippets/cli-restore-help.txt

go run . daemon -h > docs/snippets/cli-daemon-help.txt
go run . cron -h > docs/snippets/cli-cron-help.txt
# go run . web -h > docs/snippets/cli-web-help.txt