@echo off
setlocal

set PROTO_FILE=%1
set DATA=%2
set SERVER=%3
set SERVICE=%4

C:\Users\Kirill\go\bin\grpcurl.exe -proto "%PROTO_FILE%" -d "%DATA%" -plaintext %SERVER% %SERVICE%

endlocal