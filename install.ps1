# Claude PP installation script for Windows

$ErrorActionPreference = "Stop"

# Configuration
$Repo = "DandaAkhilReddy/claude_pp"
$InstallDir = if ($env:CLAUDE_PP_INSTALL_DIR) { $env:CLAUDE_PP_INSTALL_DIR } else { "$env:USERPROFILE\.local\bin" }

function Get-LatestVersion {
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        return $response.tag_name
    } catch {
        return "v0.1.0"
    }
}

function Main {
    Write-Host "Claude PP Installer for Windows"
    Write-Host ""

    # Detect architecture
    $arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
    Write-Host "Detected architecture: windows/$arch"

    # Get latest version
    $version = Get-LatestVersion
    Write-Host "Version: $version"

    # Construct download URL
    $binary = "claude_pp-windows-$arch.exe"
    $downloadUrl = "https://github.com/$Repo/releases/download/$version/$binary"

    # Create install directory
    if (!(Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }

    $installPath = Join-Path $InstallDir "claude_pp.exe"

    # Download binary
    Write-Host "Downloading $binary..."
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $installPath
    } catch {
        Write-Host "Error: Failed to download binary"
        Write-Host "URL: $downloadUrl"
        exit 1
    }

    Write-Host ""
    Write-Host "Claude PP installed to $installPath"

    # Check if install dir is in PATH
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -notlike "*$InstallDir*") {
        Write-Host ""
        Write-Host "Add $InstallDir to your PATH:"
        Write-Host ""
        Write-Host "  `$env:PATH = `"`$env:PATH;$InstallDir`""
        Write-Host ""
        Write-Host "Or add it permanently in System Properties > Environment Variables"
    }

    Write-Host ""
    Write-Host "Run 'claude_pp setup' to configure for your AI coding tool."
}

Main
