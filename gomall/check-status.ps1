# Check service status
$LOCALHOST = "127.0.0.1"

Write-Host "=== Gomall Service Status Check ===" -ForegroundColor Cyan
Write-Host ""

# Check local infrastructure services
Write-Host "Infrastructure Services:" -ForegroundColor Yellow
$infraServices = @(
    @{Name="MySQL"; Port=3306},
    @{Name="Redis"; Port=6379},
    @{Name="Consul"; Port=8500}
)

foreach ($svc in $infraServices) {
    Write-Host "  $($svc.Name) ($($svc.Port))..." -NoNewline
    $result = Test-NetConnection -ComputerName $LOCALHOST -Port $svc.Port -WarningAction SilentlyContinue -InformationLevel Quiet
    if ($result) {
        Write-Host " OK" -ForegroundColor Green
    } else {
        Write-Host " FAILED" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "Application Services:" -ForegroundColor Yellow

# Check local application services
$appServices = @(
    @{Name="user"; Port=8881},
    @{Name="product"; Port=8882},
    @{Name="cart"; Port=8883},
    @{Name="payment"; Port=8884},
    @{Name="checkout"; Port=8885},
    @{Name="order"; Port=8886},
    @{Name="email"; Port=8887},
    @{Name="frontend"; Port=8080}
)

foreach ($svc in $appServices) {
    Write-Host "  $($svc.Name) ($($svc.Port))..." -NoNewline
    $result = Test-NetConnection -ComputerName localhost -Port $svc.Port -WarningAction SilentlyContinue -InformationLevel Quiet
    if ($result) {
        Write-Host " OK" -ForegroundColor Green
    } else {
        Write-Host " NOT READY" -ForegroundColor Yellow
    }
}

Write-Host ""
Write-Host "Access URLs:" -ForegroundColor Cyan
Write-Host "  Frontend:  http://localhost:8080" -ForegroundColor White
Write-Host "  Consul UI: http://localhost:8500" -ForegroundColor White
Write-Host "  Jaeger UI: http://localhost:16686" -ForegroundColor White
Write-Host "  Grafana:   http://localhost:3000" -ForegroundColor White
Write-Host "  Prometheus: http://localhost:9090" -ForegroundColor White
