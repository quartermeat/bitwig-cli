param(
    [int]$GraceSeconds = 10,
    [switch]$NoForce,
    [string]$ProjectPath,
    [string]$BitwigExe = "C:\Program Files\Bitwig Studio\Bitwig Studio.exe"
)

$ErrorActionPreference = "Stop"

if (-not (Test-Path -LiteralPath $BitwigExe)) {
    throw "Bitwig executable not found: $BitwigExe"
}

if ($ProjectPath) {
    $ProjectPath = (Resolve-Path -LiteralPath $ProjectPath).Path
    if (-not $ProjectPath.EndsWith(".bwproject", [System.StringComparison]::OrdinalIgnoreCase)) {
        throw "ProjectPath must point to a .bwproject file: $ProjectPath"
    }
}

function Get-BitwigProcess {
    Get-Process |
        Where-Object {
            if ($_.ProcessName -like "Bitwig*") {
                return $true
            }

            try {
                return $_.Path -eq $BitwigExe
            } catch {
                return $false
            }
        } |
        Sort-Object Id -Unique
}

$processes = Get-BitwigProcess

foreach ($process in $processes) {
    if ($process.MainWindowHandle -ne 0) {
        Write-Host "Closing Bitwig window for process $($process.Id)..."
        [void]$process.CloseMainWindow()
    }
}

$deadline = (Get-Date).AddSeconds($GraceSeconds)
while ((Get-Date) -lt $deadline) {
    $remaining = Get-BitwigProcess
    if (-not $remaining) {
        break
    }
    Start-Sleep -Milliseconds 250
}

$remaining = Get-BitwigProcess

if ($remaining) {
    if ($NoForce) {
        $ids = ($remaining | ForEach-Object { $_.Id }) -join ", "
        throw "Bitwig did not exit within $GraceSeconds seconds. Still running: $ids"
    }

    foreach ($process in $remaining) {
        Write-Host "Force-stopping Bitwig process $($process.Id)..."
        Stop-Process -Id $process.Id -Force
    }
}

if ($ProjectPath) {
    Write-Host "Starting Bitwig with project: $ProjectPath"
    Start-Process -FilePath $BitwigExe -ArgumentList @("`"$ProjectPath`"") -WindowStyle Normal
} else {
    Write-Host "Starting Bitwig..."
    Start-Process -FilePath $BitwigExe -WindowStyle Normal
}
