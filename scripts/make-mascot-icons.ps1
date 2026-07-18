# Renders assets/mascot.png into every raster the app needs:
#   frontend/src/assets/mascot-512.png   UI (header, loading screen)
#   build/appicon.png                    Wails packaging master
#   build/windows/icon.ico               multi-size app icon (PNG entries)
# Run after replacing the mascot art:  powershell -File scripts\make-mascot-icons.ps1

$ErrorActionPreference = 'Stop'
Add-Type -AssemblyName System.Drawing

$root = Split-Path $PSScriptRoot -Parent
$img = [System.Drawing.Image]::FromFile((Join-Path $root 'assets\mascot.png'))

# Slight inward crop so the subject fills icon frames.
$inset = [int]($img.Width * 0.04)
$crop = New-Object System.Drawing.Rectangle($inset, $inset, ($img.Width - 2 * $inset), ($img.Height - 2 * $inset))

function New-Resized {
    param($Source, $Rect, [int]$Size)
    $bmp = New-Object System.Drawing.Bitmap($Size, $Size)
    $g = [System.Drawing.Graphics]::FromImage($bmp)
    $g.InterpolationMode = [System.Drawing.Drawing2D.InterpolationMode]::HighQualityBicubic
    $g.PixelOffsetMode = [System.Drawing.Drawing2D.PixelOffsetMode]::HighQuality
    $g.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::HighQuality
    $dest = New-Object System.Drawing.Rectangle(0, 0, $Size, $Size)
    $g.DrawImage($Source, $dest, $Rect, [System.Drawing.GraphicsUnit]::Pixel)
    $g.Dispose()
    return $bmp
}

(New-Resized $img $crop 512).Save((Join-Path $root 'frontend\src\assets\mascot-512.png'), [System.Drawing.Imaging.ImageFormat]::Png)
(New-Resized $img $crop 1024).Save((Join-Path $root 'build\appicon.png'), [System.Drawing.Imaging.ImageFormat]::Png)

# ICO container: header + directory entries + PNG payloads.
$sizes = 16, 20, 24, 32, 48, 64, 128, 256
$blobs = foreach ($s in $sizes) {
    $bmp = New-Resized $img $crop $s
    $ms = New-Object System.IO.MemoryStream
    $bmp.Save($ms, [System.Drawing.Imaging.ImageFormat]::Png)
    $bmp.Dispose()
    , $ms.ToArray()
}

$out = New-Object System.IO.MemoryStream
$bw = New-Object System.IO.BinaryWriter($out)
$bw.Write([uint16]0)              # reserved
$bw.Write([uint16]1)              # type: icon
$bw.Write([uint16]$sizes.Count)
$offset = 6 + 16 * $sizes.Count
for ($i = 0; $i -lt $sizes.Count; $i++) {
    $s = $sizes[$i]; $b = $blobs[$i]
    $dim = if ($s -ge 256) { 0 } else { $s }   # 0 means 256
    $bw.Write([byte]$dim); $bw.Write([byte]$dim)
    $bw.Write([byte]0); $bw.Write([byte]0)     # palette, reserved
    $bw.Write([uint16]1); $bw.Write([uint16]32) # planes, bpp
    $bw.Write([uint32]$b.Length); $bw.Write([uint32]$offset)
    $offset += $b.Length
}
foreach ($b in $blobs) { $bw.Write($b) }
$bw.Flush()
[System.IO.File]::WriteAllBytes((Join-Path $root 'build\windows\icon.ico'), $out.ToArray())
$img.Dispose()

Write-Host "Wrote mascot-512.png, appicon.png, and icon.ico ($($sizes -join ', ') px)."
