$ErrorActionPreference = "Stop";

Write-Output "DB_PROTOCOL=$env:DB_PROTOCOL"
Write-Output "DB_HOST=$env:DB_HOST"
Write-Output "DB_PORT=$env:DB_PORT"
Write-Output "DB_DEFAULT_CHARACTER_SET=$env:DB_DEFAULT_CHARACTER_SET"
Write-Output "DB_EXPORT_GZIP=$env:DB_EXPORT_GZIP"
Write-Output "DB_EXPORT_FILE_PATH=$env:DB_EXPORT_FILE_PATH"
Write-Output "DB_NAME=$env:DB_NAME"
Write-Output "DB_USERNAME=$env:DB_USERNAME"
Write-Output "DB_PASSWORD=$env:DB_PASSWORD"
Write-Output "DB_ARGS=$env:DB_ARGS"

if ($LastExitCode -gt 0) { exit $LastExitCode }
