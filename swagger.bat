@echo off
setlocal enabledelayedexpansion

echo.
echo ========================================
echo   Blog API Swagger Documentation
echo ========================================
echo.

REM 检查Go是否安装
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ 错误: 未找到Go环境，请先安装Go
    pause
    exit /b 1
)

REM 设置变量
set "PROJECT_ROOT=%~dp0"
set "DOCS_DIR=%PROJECT_ROOT%docs"

echo 📁 项目根目录: %PROJECT_ROOT%
echo 📄 文档目录: %DOCS_DIR%
echo.

REM 显示菜单
:menu
echo 请选择操作:
echo 1. 生成Swagger文档
echo 2. 启动服务查看文档
echo 3. 清理文档文件
echo 4. 安装Swagger工具
echo 5. 重新生成并启动服务
echo 6. 打开Swagger文档页面
echo 7. 打开Apifox导入页面
echo 8. 查看文档文件
echo 9. 退出
echo.
set /p choice="请输入选择 (1-9): "

if "%choice%"=="1" goto generate_docs
if "%choice%"=="2" goto start_server
if "%choice%"=="3" goto clean_docs
if "%choice%"=="4" goto install_swagger
if "%choice%"=="5" goto regenerate_and_start
if "%choice%"=="6" goto open_browser
if "%choice%"=="7" goto open_apifox
if "%choice%"=="8" goto view_docs
if "%choice%"=="9" goto exit
echo ❌ 无效选择，请重新输入
echo.
goto menu

:generate_docs
echo.
echo 🔨 生成Swagger文档...
swag init -g main.go -o ./docs --parseDependency --parseInternal
if errorlevel 1 (
    echo ❌ 生成文档失败，请检查swag工具是否已安装
    echo 💡 提示: 运行选项4安装Swagger工具
    pause
    goto menu
)
echo ✅ Swagger文档生成完成!
echo 📁 文档位置: %DOCS_DIR%
pause
goto menu

:start_server
echo.
echo 🚀 启动服务...
echo 📖 Swagger文档地址: http://localhost:8868/swagger/index.html
echo 🛑 按 Ctrl+C 停止服务
echo.
go run main.go
pause
goto menu

:clean_docs
echo.
echo 🧹 清理文档文件...
if exist "%DOCS_DIR%" (
    rmdir /s /q "%DOCS_DIR%"
    echo ✅ 文档文件清理完成
) else (
    echo ℹ️  文档目录不存在，无需清理
)
pause
goto menu

:install_swagger
echo.
echo 📦 安装Swagger工具...
go install github.com/swaggo/swag/cmd/swag@latest
if errorlevel 1 (
    echo ❌ 安装失败
) else (
    echo ✅ Swagger工具安装完成
    echo 💡 现在可以使用 swag 命令了
)
pause
goto menu

:regenerate_and_start
echo.
echo 🔄 重新生成文档并启动服务...
echo.
echo 1. 生成Swagger文档...
swag init -g main.go -o ./docs --parseDependency --parseInternal
if errorlevel 1 (
    echo ❌ 生成文档失败
    pause
    goto menu
)
echo ✅ 文档生成完成
echo.
echo 2. 启动服务...
echo 📖 Swagger文档地址: http://localhost:8868/swagger/index.html
echo 🛑 按 Ctrl+C 停止服务
echo.
go run main.go
pause
goto menu

:open_browser
echo.
echo 🌐 打开Swagger文档页面...
start http://localhost:8868/swagger/index.html
echo ✅ 已尝试打开浏览器
echo 💡 如果页面无法访问，请先启动服务 (选项2)
pause
goto menu

:open_apifox
echo.
echo 🎯 打开Apifox导入页面...
start http://localhost:8868/apifox
echo ✅ 已尝试打开Apifox导入页面
echo.
echo 💡 使用说明:
echo   • 页面提供一键导入链接和手动导入URL
echo   • 一键导入: 点击页面上的"一键导入Apifox"按钮
echo   • 手动导入: 复制URL到Apifox中手动导入
echo   • 导入URL: http://localhost:8868/swagger/doc.json
echo.
echo 如果页面无法访问，请先启动服务 (选项2)
pause
goto menu

:view_docs
echo.
echo 📄 查看文档文件...
if exist "%DOCS_DIR%" (
    echo 📁 文档目录内容:
    dir "%DOCS_DIR%" /b
    echo.
    echo 📊 文档统计:
    if exist "%DOCS_DIR%\swagger.json" (
        echo   ✅ swagger.json - JSON格式文档
    )
    if exist "%DOCS_DIR%\swagger.yaml" (
        echo   ✅ swagger.yaml - YAML格式文档
    )
    if exist "%DOCS_DIR%\docs.go" (
        echo   ✅ docs.go - Go代码文档
    )
    echo.
    echo 是否打开文档目录? (y/N)
    set /p open_dir=""
    if /i "%open_dir%"=="y" (
        explorer "%DOCS_DIR%"
    )
) else (
    echo ❌ 文档目录不存在
    echo 💡 请先生成文档 (选项1)
)
pause
goto menu

:exit
echo.
echo 👋 感谢使用 Blog API Swagger Documentation!
echo.
echo 💡 快速提醒:
echo   • 文档地址: http://localhost:8868/swagger/index.html
echo   • 生成命令: swag init -g main.go -o ./docs
echo   • 启动服务: go run main.go
echo.
pause
exit /b 0
