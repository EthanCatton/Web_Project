Write-Host "Running"
Write-Host "Docker and Node js are not included in this install - please install seperately!" 
Start-Sleep -Seconds 2

pip install nltk
pip install fastapi
pip install uvicorn 

Start-Process -FilePath "python" -ArgumentList ".\Backend\Setup\Setup.py"   
Start-Process -FilePath "npm.cmd" -WorkingDirectory ".\Frontend\React\Web-React" -ArgumentList "install" 
 
