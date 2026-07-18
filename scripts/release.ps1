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
# which breaks strict JSON parsers. Each replace is asserted to actually
# match, so a renamed field can never silently ship the old version.
$utf8 = New-Object System.Text.UTF8Encoding($false)

function Set-Version($file, $pattern, $replacement) {
    $before = Get-Content $file -Raw
    $after = $before -replace $pattern, $replacement
    if ($after -eq $before) { throw "version stamp did not match anything in $file" }
    [System.IO.File]::WriteAllText($file, $after, $utf8)
}

Set-Version "$root\main.go" 'const appVersion = "[^"]+"' "const appVersion = `"$Version`""
Set-Version "$root\wails.json" '"productVersion": "[^"]+"' "`"productVersion`": `"$Version`""

git add main.go wails.json
git commit -m "Release v$Version"
git tag "v$Version"
git push
git push origin "v$Version"

Write-Host "Tagged v$Version. GitHub Actions is building the release."
