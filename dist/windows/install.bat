@echo off
echo Installing gnode...

REM Criar diretório de instalação
set "INSTALL_DIR=%USERPROFILE%\bin"
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

REM Copiar executável
echo Copying gnode.exe...
copy "gnode.exe" "%INSTALL_DIR%\" >nul
if errorlevel 1 (
    echo Error: Could not copy gnode.exe
    pause
    exit /b 1
)

REM Verificar se já está no PATH
echo %PATH% | find "%INSTALL_DIR%" >nul
if not errorlevel 1 (
    echo gnode already in PATH
    goto :test
)

REM Adicionar ao PATH do usuário
echo Adding to PATH...
for /f "tokens=2*" %%a in ('reg query "HKCU\Environment" /v PATH 2^>nul') do set "USER_PATH=%%b"
if not defined USER_PATH set "USER_PATH="

if defined USER_PATH (
    reg add "HKCU\Environment" /v PATH /t REG_EXPAND_SZ /d "%USER_PATH%;%INSTALL_DIR%" /f >nul
) else (
    reg add "HKCU\Environment" /v PATH /t REG_EXPAND_SZ /d "%INSTALL_DIR%" /f >nul
)

if errorlevel 1 (
    echo Error: Could not update PATH
    echo Please add %INSTALL_DIR% to your PATH manually
    pause
    exit /b 1
)

echo PATH updated successfully!

:test
echo.
echo Testing installation...
"%INSTALL_DIR%\gnode.exe" --version >nul 2>&1
if errorlevel 1 (
    echo Warning: gnode may not be working correctly
) else (
    echo gnode installed successfully!
)

echo.
echo Installation completed!
echo.
echo Usage:
echo   gnode install v18.17.0
echo   gnode use v18.17.0
echo.
echo Note: You may need to restart your terminal for PATH changes to take effect.
echo.
pause