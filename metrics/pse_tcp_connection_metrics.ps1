function Get-PidToNameMap($processes) {
	$hash = @{}
	$processes | ForEach-Object {
		$hash[$_.Id.ToString()] = ($_.UserName -split '\\')[1]
	}
	return $hash
}
$pidToName = Get-PidToNameMap(Get-Process -IncludeUserName `
								| Where-Object { $_.ProcessName -eq "w3wp" } `
								| Select-Object -Property Id, UserName)

Get-NetTCPConnection `
| Group-Object -Property OwningProcess, State `
| Select-Object -Property Count, `
							@{Name = "OwningProcess"; Expression = { $_.Group[0].OwningProcess.ToString() } }, `
							@{Name = "State"; Expression = { $_.Group[0].State.ToString() } } `
| Sort-Object -Property OwningProcess, State `
| ForEach-Object {
	if ($pidToName.ContainsKey($_.OwningProcess)) {
		@{
			# metric value
			'value'  = $_.Count
			# metric labels
			'labels' = @{
				'state'   = $_.State
				'service' = $pidToName[$_.OwningProcess]
			}
		}
	}
} | ConvertTo-Json -Compress