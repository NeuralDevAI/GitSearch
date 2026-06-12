# Installs gitsearch globally
Write-Host "Installing gitsearch..." -ForegroundColor Cyan

# Check if go is installed and compile
if (Get-Command go -ErrorAction SilentlyContinue) {
    Write-Host "Go compiler found. Compiling and installing..." -ForegroundColor Green
    go install .
    if ($LASTEXITCODE -eq 0) {
        Write-Host "gitsearch has been successfully installed via 'go install'!" -ForegroundColor Green
        Write-Host "Make sure your GOPATH bin directory ($env:USERPROFILE\go\bin) is in your PATH." -ForegroundColor Yellow
        exit 0
    } else {
        Write-Error "Failed to install gitsearch using 'go install'."
        exit 1
    }
} else {
    Write-Host "Go compiler not found. Trying to install precompiled binary..." -ForegroundColor Yellow
    if (Test-Path .\gitsearch.exe) {
        $installDir = "$env:USERPROFILE\go\bin"
        if (-not (Test-Path $installDir)) {
            New-Item -ItemType Directory -Force -Path $installDir | Out-Null
        }
        Copy-Item .\gitsearch.exe -Destination "$installDir\gitsearch.exe" -Force
        Write-Host "Copied gitsearch.exe to $installDir" -ForegroundColor Green
        
        # Check PATH
        $path = [Environment]::GetEnvironmentVariable("PATH", "User")
        if ($path -notlike "*$installDir*") {
            Write-Host "Adding $installDir to User PATH..." -ForegroundColor Yellow
            [Environment]::SetEnvironmentVariable("PATH", $path + ";" + $installDir, "User")
            $env:PATH += ";" + $installDir
        }
        Write-Host "gitsearch has been successfully installed!" -ForegroundColor Green
        exit 0
    } else {
        Write-Error "Neither Go compiler nor precompiled 'gitsearch.exe' was found in the current directory."
        exit 1
    }
}
