$ver = Get-Content .\docs\version
$label = Get-Variable GITHUB_LABEL -valueOnly
$runId = Get-Variable GITHUB_RUN_ID -valueOnly
$runNumber = Get-Variable GITHUB_RUN_NUMBER -valueOnly
$ver = -join($ver, "-",$label,".",$runId,".",$runId);

New-Item -Path "docs" -Name "version" -ItemType "file" -Value "${ver}" -Force
