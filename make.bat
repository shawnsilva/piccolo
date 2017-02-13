@echo off

set EXITCODE=0

if [%1]==[] (
    echo :: Error: Incorrect usage
    set EXITCODE=1
    goto usage
    )

for %%a in (build,clean) do (
    if x%1==x%%a (
        goto %%a
        )
)
echo.
echo :: Unknown make target: %1
echo :: Expected one of the following: "all", "clean".
set EXITCODE=1
goto end


:build
md build 2>NUL
call .\scripts\windows\build.bat %CD%
if not errorlevel 1 goto end
echo.
echo :: BUILD FAILED
set EXITCODE=%ERRORLEVEL%
goto end

:clean
go clean -x .\...
goto end

:usage
echo :: USAGE: make.bat (build,clean)
goto end

:end
exit /B %EXITCODE%
