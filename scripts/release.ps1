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

# Stamp the version everywhere it lives. WriteAllText with a BOM-less
# encoding: PowerShell 5.1's Set-Content -Encoding UTF8 writes a BOM,
# which breaks strict JSON parsers.
$utf8 = New-Object System.Text.UTF8Encoding($false)

$main = (Get-Content "$root\main.go" -Raw) -replace 'const appVersion = "[^"]+"', "const appVersion = `"$Version`""
[System.IO.File]::WriteAllText("$root\main.go", $main, $utf8)

$wails = (Get-Content "$root\wails.json" -Raw) -replace '"productVersion": "[^"]+"', "`"productVersion`": `"$Version`""
[System.IO.File]::WriteAllText("$root\wails.json", $wails, $utf8)

git add main.go wails.json
git commit -m "Release v$Version"
git tag "v$Version"
git push
git push origin "v$Version"

Write-Host "Tagged v$Version. GitHub Actions is building the release."
