@echo off

if not exist %1 exit /B 1
cd %1

for /f %%i in ('git describe --always --dirty') do set GIT_VERSION=%%i
for /f %%i in ('git rev-parse --abbrev-ref HEAD') do set GIT_BRANCH=%%i

echo :: Downloading Dependencies
cd cmd\piccolo
go get -v .\...
cd %1

echo :: Building Piccolo...
go build^
    -v^
    -ldflags "-X github.com/jatgam/goutils/version.gitVersion=%GIT_VERSION% -X github.com/jatgam/goutils/version.gitBranch=%GIT_BRANCH%"^
    -gcflags "all=-trimpath=%GOPATH%"^
    -asmflags "all=-trimpath=%GOPATH%"^
    -o build\piccolo.exe cmd\piccolo\piccolo.go
if errorlevel 1 exit /B 1
