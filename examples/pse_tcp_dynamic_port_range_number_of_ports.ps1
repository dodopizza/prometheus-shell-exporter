@(
	@{}
    @{
        # metric value
        'value' = (Get-NetTCPSetting -SettingName 'Datacenter' ).DynamicPortRangeNumberOfPorts
        # metric labels
        'labels' = @{}
    }
) | ConvertTo-Json -Compress