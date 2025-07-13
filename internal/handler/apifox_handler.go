package handler

import (
	"blog/internal/utils"

	"github.com/gin-gonic/gin"
)

type ApifoxHandler struct{}

func NewApifoxHandler() *ApifoxHandler {
	return &ApifoxHandler{}
}

// GetApifoxImportInfo 获取Apifox导入信息
// @Summary 获取Apifox导入信息
// @Description 获取Apifox一键导入所需的信息和链接
// @Tags Apifox导入
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "导入信息"
// @Router /api/v1/apifox/import [get]
func (h *ApifoxHandler) GetApifoxImportInfo(c *gin.Context) {
	host := c.Request.Host
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	baseURL := scheme + "://" + host

	importInfo := map[string]interface{}{
		"title":       "Blog API Apifox导入指南",
		"description": "一键导入Blog API到Apifox的完整指南",
		"api_info": map[string]interface{}{
			"title":       "Blog API",
			"version":     "1.0.0",
			"description": "博客系统API文档",
		},
		"import_urls": map[string]string{
			"openapi_json": baseURL + "/swagger/doc.json",
			"openapi_yaml": baseURL + "/docs/swagger.yaml",
			"swagger_ui":   baseURL + "/swagger/index.html",
		},
		"apifox_import_steps": []map[string]interface{}{
			{
				"step":        1,
				"title":       "打开Apifox应用",
				"description": "启动Apifox桌面应用或访问网页版",
			},
			{
				"step":        2,
				"title":       "创建或选择项目",
				"description": "在Apifox中创建新项目或选择现有项目",
			},
			{
				"step":        3,
				"title":       "导入API",
				"description": "点击「导入」按钮，选择「URL导入」或「OpenAPI」",
			},
			{
				"step":        4,
				"title":       "输入导入URL",
				"description": "使用以下URL进行导入：" + baseURL + "/swagger/doc.json",
				"url":         baseURL + "/swagger/doc.json",
			},
			{
				"step":        5,
				"title":       "确认导入",
				"description": "检查导入设置，点击确认完成导入",
			},
		},
		"quick_import_link": "apifox://import?url=" + baseURL + "/swagger/doc.json",
		"manual_import_url": baseURL + "/swagger/doc.json",
		"statistics": map[string]interface{}{
			"total_apis": 14,
			"groups": []map[string]interface{}{
				{"name": "用户管理", "count": 3},
				{"name": "文章管理", "count": 5},
				{"name": "文件管理", "count": 1},
				{"name": "测试管理", "count": 5},
			},
		},
		"tips": []string{
			"建议使用JSON格式的URL进行导入，兼容性更好",
			"导入后可以在Apifox中直接测试所有接口",
			"支持JWT认证，导入后需要配置认证信息",
			"如果导入失败，请检查网络连接和URL是否正确",
		},
	}

	utils.Success(c, importInfo)
}

// GetApifoxQuickImport 获取Apifox快速导入页面
// @Summary Apifox快速导入页面
// @Description 显示Apifox快速导入的HTML页面
// @Tags Apifox导入
// @Produce html
// @Success 200 {string} string "HTML页面"
// @Router /apifox [get]
func (h *ApifoxHandler) GetApifoxQuickImport(c *gin.Context) {
	host := c.Request.Host
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	baseURL := scheme + "://" + host
	importURL := baseURL + "/swagger/doc.json"
	apifoxURL := "apifox://import?url=" + importURL

	html := `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Blog API - Apifox 一键导入</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 20px;
            padding: 40px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
            max-width: 600px;
            width: 100%;
            text-align: center;
        }
        .logo { font-size: 48px; margin-bottom: 20px; }
        h1 { color: #333; margin-bottom: 10px; font-size: 28px; }
        .subtitle { color: #666; margin-bottom: 30px; font-size: 16px; }
        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
            gap: 15px;
            margin: 30px 0;
        }
        .stat {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 10px;
            border-left: 4px solid #667eea;
        }
        .stat-number { font-size: 24px; font-weight: bold; color: #667eea; }
        .stat-label { font-size: 12px; color: #666; margin-top: 5px; }
        .buttons { margin: 30px 0; }
        .btn {
            display: inline-block;
            padding: 15px 30px;
            margin: 10px;
            border-radius: 10px;
            text-decoration: none;
            font-weight: bold;
            transition: all 0.3s ease;
            border: none;
            cursor: pointer;
            font-size: 16px;
        }
        .btn-primary {
            background: #667eea;
            color: white;
        }
        .btn-primary:hover {
            background: #5a6fd8;
            transform: translateY(-2px);
        }
        .btn-secondary {
            background: #f8f9fa;
            color: #333;
            border: 2px solid #dee2e6;
        }
        .btn-secondary:hover {
            background: #e9ecef;
            transform: translateY(-2px);
        }
        .url-box {
            background: #f8f9fa;
            border: 1px solid #dee2e6;
            border-radius: 8px;
            padding: 15px;
            margin: 20px 0;
            font-family: monospace;
            font-size: 14px;
            word-break: break-all;
            position: relative;
        }
        .copy-btn {
            position: absolute;
            top: 10px;
            right: 10px;
            background: #667eea;
            color: white;
            border: none;
            padding: 5px 10px;
            border-radius: 5px;
            cursor: pointer;
            font-size: 12px;
        }
        .steps {
            text-align: left;
            margin: 30px 0;
        }
        .step {
            display: flex;
            align-items: flex-start;
            margin: 15px 0;
        }
        .step-number {
            background: #667eea;
            color: white;
            width: 24px;
            height: 24px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 12px;
            font-weight: bold;
            margin-right: 15px;
            flex-shrink: 0;
        }
        .step-content h4 { color: #333; margin-bottom: 5px; }
        .step-content p { color: #666; font-size: 14px; }
        .tips {
            background: #fff3cd;
            border: 1px solid #ffeaa7;
            border-radius: 8px;
            padding: 15px;
            margin: 20px 0;
            text-align: left;
        }
        .tips h4 { color: #856404; margin-bottom: 10px; }
        .tips ul { margin-left: 20px; }
        .tips li { color: #856404; margin: 5px 0; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">🚀</div>
        <h1>Blog API</h1>
        <p class="subtitle">一键导入到 Apifox</p>
        
        <div class="stats">
            <div class="stat">
                <div class="stat-number">14</div>
                <div class="stat-label">API接口</div>
            </div>
            <div class="stat">
                <div class="stat-number">4</div>
                <div class="stat-label">功能分组</div>
            </div>
            <div class="stat">
                <div class="stat-number">100%</div>
                <div class="stat-label">文档覆盖</div>
            </div>
        </div>

        <div class="buttons">
            <a href="` + apifoxURL + `" class="btn btn-primary">🎯 一键导入 Apifox</a>
            <a href="` + baseURL + `/swagger/index.html" class="btn btn-secondary">📖 查看文档</a>
        </div>

        <div class="url-box">
            <button class="copy-btn" onclick="copyURL()">复制</button>
            <div id="import-url">` + importURL + `</div>
        </div>

        <div class="steps">
            <h3 style="margin-bottom: 20px; color: #333;">📋 手动导入步骤</h3>
            <div class="step">
                <div class="step-number">1</div>
                <div class="step-content">
                    <h4>打开 Apifox</h4>
                    <p>启动 Apifox 桌面应用或访问网页版</p>
                </div>
            </div>
            <div class="step">
                <div class="step-number">2</div>
                <div class="step-content">
                    <h4>选择导入</h4>
                    <p>点击「导入」按钮，选择「URL导入」或「OpenAPI」</p>
                </div>
            </div>
            <div class="step">
                <div class="step-number">3</div>
                <div class="step-content">
                    <h4>输入URL</h4>
                    <p>粘贴上方的导入URL，点击确认导入</p>
                </div>
            </div>
        </div>

        <div class="tips">
            <h4>💡 导入提示</h4>
            <ul>
                <li>建议使用 JSON 格式的 URL 进行导入，兼容性更好</li>
                <li>导入后可以在 Apifox 中直接测试所有接口</li>
                <li>支持 JWT 认证，导入后需要配置认证信息</li>
                <li>如果一键导入不工作，请使用手动导入方式</li>
            </ul>
        </div>
    </div>

    <script>
        function copyURL() {
            const url = document.getElementById('import-url').textContent;
            navigator.clipboard.writeText(url).then(() => {
                const btn = document.querySelector('.copy-btn');
                const originalText = btn.textContent;
                btn.textContent = '已复制';
                btn.style.background = '#28a745';
                setTimeout(() => {
                    btn.textContent = originalText;
                    btn.style.background = '#667eea';
                }, 2000);
            });
        }
    </script>
</body>
</html>`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}
