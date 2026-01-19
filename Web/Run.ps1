Write-Host "Running"


Start-Process "docker" -WorkingDirectory ".\Backend" -ArgumentList "compose up --build" 
$pythonswitch=Start-Process -FilePath "python" -ArgumentList ".\Backend\Interface\Interface.py" -PassThru -NoNewWindow
$viteswitch=Start-Process -FilePath "npm.cmd" -WorkingDirectory ".\Frontend\React\Web-React" -ArgumentList "run","dev"  -PassThru -NoNewWindow

Write-Host "All main components coming online, please give a few seconds for Docker containers to fully start"
Write-Host "http://localhost:5173 via Vite"
Start-Sleep -Seconds 3 


Write-Host "Press any key to begin shutdown"

#taken from web
$x = $host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

Start-Process "docker" -WorkingDirectory ".\Backend" -ArgumentList "compose down" -NoNewWindow
Stop-Process -Id $pythonswitch.Id
Write-Host "First shutdown will need to download a small package!"
npx kill-port 5173 
Write-Host "Shutdown"
 
 