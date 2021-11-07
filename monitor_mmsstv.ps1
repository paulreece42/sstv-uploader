##
# Originally from StackOverflow: https://superuser.com/questions/226828/how-to-monitor-a-folder-and-trigger-a-command-line-action-when-a-file-is-created
#
### SET FOLDER TO WATCH + FILES TO WATCH + SUBFOLDERS YES/NO
    $watcher = New-Object System.IO.FileSystemWatcher
    $watcher.Path = "C:\Ham\MMSSTV\History"
    $watcher.Filter = "*.bmp*"
    $watcher.IncludeSubdirectories = $true
    $watcher.EnableRaisingEvents = $true  

### DEFINE ACTIONS AFTER AN EVENT IS DETECTED
#    $action = { $path = $Event.SourceEventArgs.FullPath
#                $changeType = $Event.SourceEventArgs.ChangeType
#                $logline = "$(Get-Date), $changeType, $path"
#                Add-content "C:\log.txt" -value $logline
#              }    

$action = { 

sleep 2
$path = $Event.SourceEventArgs.FullPath
                $changeType = $Event.SourceEventArgs.ChangeType
                $logline = "$(Get-Date), $changeType, $path"
                Add-content "C:\log.txt" -value $logline
C:\curl.exe -XPOST -H "Bearer: foobarbaz" -H "Content-Type: multipart/form-data" --form "file=@$path" http://hackdetroit.city:14230/sstv/

}
### DECIDE WHICH EVENTS SHOULD BE WATCHED 
#    Register-ObjectEvent $watcher "Created" -Action $action
    Register-ObjectEvent $watcher "Changed" -Action $action
#    Register-ObjectEvent $watcher "Deleted" -Action $action
#    Register-ObjectEvent $watcher "Renamed" -Action $action
    while ($true) {sleep 5}
