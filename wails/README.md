# QSL Mail Lite

这是 QSL Mail 的 Wails 轻量桌面壳。正式界面源代码仍在仓库根目录的 `index.html`、`app.js` 和 `styles.css`，构建前由 `frontend/sync.cjs` 同步，因此 Electron 与 Wails 的界面保持一致。

## 开发

在仓库根目录执行：

```bash
npm run wails:dev
```

也可以进入本目录执行 `wails dev`。

## 构建

```bash
npm run wails:build
```

构建产物位于 `wails/build/bin/`。Wails 使用系统 WebView，不会把 Chromium 一起打包。Wails 壳通过 `wails-bridge.js` 提供打印机查询和 Go 本地数据保存接口；打印仍沿用浏览器系统打印流程，以保持现有预览和打印设置行为。

如果构建环境无法访问 Go 模块代理，先执行 `go mod tidy`，再运行构建命令。发布跨平台安装包应在对应的 macOS、Windows、Linux runner 上分别构建。
