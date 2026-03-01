# Fix Go module dependencies for all services
Write-Host "=== Fixing Go Module Dependencies ===" -ForegroundColor Cyan
Write-Host ""

$services = @("user", "product", "cart", "payment", "checkout", "order", "email", "frontend")
$rootDir = $PSScriptRoot

foreach ($svc in $services) {
    Write-Host "Processing $svc..." -ForegroundColor Yellow
    $svcPath = Join-Path $rootDir "app\$svc"

    if (Test-Path $svcPath) {
        Push-Location $svcPath

        Write-Host "  Running go mod tidy..." -NoNewline
        $output = go mod tidy 2>&1

        if ($LASTEXITCODE -eq 0) {
            Write-Host " OK" -ForegroundColor Green
        } else {
            Write-Host " FAILED" -ForegroundColor Red
            Write-Host "  Error: $output" -ForegroundColor Red
        }

        Pop-Location
    } else {
        Write-Host "  Service directory not found!" -ForegroundColor Red
    }

    Write-Host ""
}

Write-Host "=== Done ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Check for any errors above"
Write-Host "2. Run .\check-status.ps1 to verify services"
Write-Host "3. Start services using the manual startup guide"
