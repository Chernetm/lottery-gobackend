@echo off
set GO_EXE="C:\Program Files\Go\bin\go.exe"
%GO_EXE% mod init lottery-backend
if %errorlevel% neq 0 (
    echo [ERROR] Failed to init go module
    exit /b %errorlevel%
)
%GO_EXE% get github.com/gin-gonic/gin
%GO_EXE% get gorm.io/gorm
%GO_EXE% get gorm.io/driver/mysql
%GO_EXE% get github.com/joho/godotenv
%GO_EXE% get github.com/golang-jwt/jwt/v5
%GO_EXE% get golang.org/x/crypto/bcrypt
echo [SUCCESS] Setup complete
