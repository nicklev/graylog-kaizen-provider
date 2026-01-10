@echo off
REM Provider configuration
set PROVIDER_NAME=kaizen
set VERSION=0.0.1
set BINARY_NAME=terraform-provider-%PROVIDER_NAME%_v%VERSION%.exe

echo Building %BINARY_NAME%...
go build -o "%GOPATH%\bin\%BINARY_NAME%" .

if %errorlevel% equ 0 (
    echo Successfully built %BINARY_NAME%
    echo Location: %GOPATH%\bin\%BINARY_NAME%
) else (
    echo Build failed
    exit /b 1
)
