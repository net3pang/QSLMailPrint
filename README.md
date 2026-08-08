# QSL Mail

面向业余无线电爱好者的跨平台信封制作与打印软件。支持收发件人姓名、呼号、地址、邮编和手机号，提供国内 / 国际寄送模式、信封尺寸自定义、可视化定位、打印预览和本地任务历史。

## 功能

- 收件人与发件人的姓名、呼号、地址、邮政编码和手机号
- CSV / TXT 导入，支持中文和常见英文列名，并可下载导入模板
- DL、C5、#10、C6 和自定义信封尺寸
- 横向 / 纵向预览，收发件人位置、字号和框宽调整
- 地址自动换行，框高度根据内容自动匹配
- 地址簿增删改、搜索和快速套用
- 模板保存、套用和删除
- 任务历史保存
- 打印机选择、实际信封纸张、缩放、X/Y 偏移和安全边界
- 打印前预览，可选择是否弹出系统打印框
- 打印时隐藏预览边框、信封折线、印章和装饰纹样
- 国内模式左上角六位邮编格，可调整位置、单字大小和字间距
- 未设置呼号时显示 `NO CALL`

## 直接运行浏览器版

不要直接双击 `index.html`。浏览器对 `file://` 页面限制较多，推荐在项目目录运行：

```bash
python3 -m http.server 4173
```

然后打开 <http://127.0.0.1:4173>。

浏览器版使用当前浏览器的 `localStorage` 保存草稿、任务、联系人、模板和设置。浏览器版打印始终会弹出系统打印窗口。

## 运行 Electron 桌面版

环境要求：Node.js 20+、Go 1.20+。

```bash
npm install
npm run build:backend
npm start
```

也可以直接运行：

```bash
npm run dev
```

如果没有编译 Go 后端，Electron 仍可启动，但任务和联系人只会保存在前端本地存储中。

## 运行轻量版（Wails）

轻量版使用 Go + 系统 WebView 打包，保留同一套界面和交互，不再捆绑 Chromium。运行已发布的 `.app`、`.exe` 或 Linux 可执行文件不需要安装 Go；Go 只在本地开发和构建时需要。

开发环境：Go 1.23+、Node.js 20+、Wails CLI 2.10.2。安装 CLI：

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@v2.10.2
```

启动开发版：

```bash
npm run wails:dev
```

构建当前平台的轻量版：

```bash
npm run wails:build
```

产物在 `wails/build/bin/`：macOS 为可双击打开的 `.app`，Windows 为 `.exe`，Linux 为可执行文件。轻量版的草稿、任务、联系人和模板仍使用前端本地存储；任务和联系人同时写入与 Electron 相同的 Go JSON 数据文件。轻量版打印使用系统打印流程，macOS 原生打印会按设置显示或跳过系统打印框，并使用当前信封尺寸。

## 保存位置

Electron 版使用两层本地保存：

- 草稿、位置、模式、模板和打印设置：Electron `localStorage`
- 任务和联系人：同时写入 Go 本地 JSON 数据文件

默认 Go 数据文件位置：

- macOS：`~/Library/Application Support/QSLMail/qsl-mail-data.json`
- Windows：`%AppData%\\QSLMail\\qsl-mail-data.json`
- Linux：`~/.config/QSLMail/qsl-mail-data.json`

也可以使用环境变量 `QSL_MAIL_DB` 指定数据文件位置。

## 开发检查

```bash
npm run check
```

该命令会检查前端 JavaScript 语法并运行 Go 测试。

## 构建安装包

本机架构构建：

```bash
npm run dist
```

产物会写入 `dist/`。Go 后端也可以单独交叉编译：

生成的桌面版可以直接双击使用：Windows 提供安装版 `.exe` 和免安装便携版 `.exe`；macOS 提供 `.dmg` 和 `.zip`，安装后运行 `.app`；Linux 提供 `.AppImage`。

### macOS 安全提示

macOS 要求下载的软件具备 Apple Developer ID 签名并完成 notarization。GitHub Actions 已经接入签名和公证流程，但需要在仓库的 `Settings → Secrets and variables → Actions` 中配置：

- `MAC_CERT_BASE64`：Developer ID Application 证书 `.p12` 的 base64 内容
- `MAC_CERT_PASSWORD`：`.p12` 导出密码
- `APPLE_ID`：Apple Developer 账号邮箱
- `APPLE_APP_SPECIFIC_PASSWORD`：Apple ID 专用密码
- `APPLE_TEAM_ID`：Apple Developer Team ID

未配置这些密钥时仍可构建，但 macOS 会把应用标记为未签名，下载后可能显示“已阻止恶意软件”。当前本地构建即处于未签名状态，首次打开时请在“系统设置 → 隐私与安全性”中点击“仍要打开”。

```bash
GOOS=darwin GOARCH=arm64 npm run build:backend
GOOS=darwin GOARCH=amd64 npm run build:backend
GOOS=linux GOARCH=amd64 npm run build:backend
GOOS=windows GOARCH=amd64 npm run build:backend
```

推送版本标签后，GitHub Actions 会分别构建 macOS、Windows 和 Linux 安装包：

```bash
git tag v1.3.0
git push origin v1.3.0
```

同一个 Release 还会上传 Wails 轻量版产物，名称以 `qsl-mail-lite-` 开头。轻量版 macOS 是 `.app`，Windows 是 `.exe`，Linux 是原生可执行文件；Electron 安装包仍用于需要静默打印、设备级打印参数或完整安装体验的场景。

## CSV 格式

第一行为表头，字段顺序如下：

```csv
姓名,呼号,地址,邮编,手机号
Takahashi Ken,JA7QXG,1-2-3 Chiyoda Tokyo,100-0001,090-1234-5678
```

## 版本

当前版本：`1.3.0`。详细变更见 [CHANGELOG.md](CHANGELOG.md)。

## 开源许可

本项目使用 MIT License，详见 [LICENSE](LICENSE)。
