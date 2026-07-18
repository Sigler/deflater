# Cuts a release: stamps the version into the app, commits, tags, and
# pushes. GitHub Actions then builds the exe and publishes the release.
#
#   powershell -File scripts\release.ps1 0.2.0

param([Parameter(Mandatory)][string]$Version)

$ErrorActionPreference = 'Stop'
if ($Version -notmatch '^\d+\.\d+\.\d+$') {
    throw "Version must look like 1.2.3, got '$Version'"
}

$root = Split-Path $PSScriptRoot -Parent
Set-Location $root

if (git status --porcelain) {
    throw 'Working tree is not clean. Commit or stash first.'
}

# Stamp the version everywhere it lives.
$main = Get-Content "$root\main.go" -Raw
$main = $main -replace 'const appVersion = "[^"]+"', "const appVersion = `"$Version`""
Set-Content "$root\main.go" $main -NoNewline -Encoding UTF8

$wails = Get-Content "$root\wails.json" -Raw
$wails = $wails -replace '"productVersion": "[^"]+"', "`"productVersion`": `"$Version`""
Set-Content "$root\wails.json" $wails -NoNewline -Encoding UTF8

git add main.go wails.json
git commit -m "Release v$Version"
git tag "v$Version"
git push
git push origin "v$Version"

Write-Host "Tagged v$Version. GitHub Actions is building the release."
