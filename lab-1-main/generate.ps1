New-Item -ItemType Directory -Force -Path data | Out-Null

for ($i = 1; $i -le 20; $i++) {
    $size = [math]::Round(512KB * $i)
    $outputFile = "data/input$i.dat"
    Write-Host "Generating $outputFile of size $size bytes..."
    bin\gensort-amd64.exe "$size" $outputFile
}
