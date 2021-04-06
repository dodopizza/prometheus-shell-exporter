$services = Get-Process -IncludeUserName | Where-Object {$_.ProcessName -eq "w3wp"} 
$connections = Get-NetTCPConnection
@(
	ForEach ($service in $services) 
	{
		$connections `
			| Where-Object {$service.Id -eq $_.OwningProcess} `
			| Group-Object -Property OwningProcess, State `
			| Select -Property Count, Name, @{Name="Service";Expression={ $service.UserName }} `
			| ForEach-Object {
				@{
					# metric value
					'value' = $_.Count
					# metric labels
					'labels' = @{
						'state' = ($_.Name -split ', ')[1]
						'service' = ($_.Service -split '\\')[1]
					}
				}
			}
	}
) | ConvertTo-Json -Compress