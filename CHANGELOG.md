# 更新日志

## 1.2.1 - 统一应用图标

- Electron 和 Wails 共用 QSL Mail 信封 / 无线电波应用图标。
- Electron 构建接入 macOS ICNS、Windows/Linux PNG 图标。
- Wails 使用同一 PNG 在各平台生成原生应用图标。

## 1.2.0 - 轻量版重构

- 增加 Wails 轻量桌面版，使用 Go + 系统 WebView 复用现有界面，减少 Electron/Chromium 运行时体积。
- Wails 与 Electron 共用 Go JSON 数据存储，任务和联系人可以跨桌面版本继续使用。
- 增加 `npm run wails:dev` 和 `npm run wails:build` 构建入口。

## 1.1.4 - 2026-08-08

- 增加 macOS hardened runtime、Developer ID 签名和 notarization 配置。
- 未配置 Apple 密钥时，构建日志会明确提示跳过公证，避免误认为发行包已通过 macOS 安全验证。

## 1.1.3 - 2026-08-08

- macOS x64 构建改用可用的 macOS 14 runner，避免因 macOS 13 runner 排队导致 Release 无法完成。

## 1.1.2 - 2026-08-08

- 修正 GitHub Actions 在标签构建时错误尝试直接发布的问题。
- 补充 Linux `.deb` 所需的作者邮箱信息。

## 1.1.1 - 2026-08-08

- Windows 增加免安装便携版 `.exe`，可直接双击运行。
- macOS 继续提供 `.dmg` / `.zip`，Linux 提供 `.AppImage`。

## 1.1.0 - 2026-08-08

- 增加 Electron 桌面版真实打印和 Go 本地数据服务。
- 增加任务历史、联系人地址簿、模板、CSV 导入和 CSV 模板下载。
- 增加国内 / 国际寄送模式，国内模式支持左上角六位邮编格的位置、字大小和字间距调整。
- 增加自定义信封尺寸、打印机选择、打印方向、缩放和 X/Y 校准。
- 增加打印前预览，打印时隐藏预览边框、信封纹样、折线、印章和装饰图形。
- 收发件人支持手机号、可调字号、位置和框宽，长地址自动换行并按内容自适应高度。
- 未设置呼号时统一显示 `NO CALL`。
- 增加 macOS、Windows、Linux 的 Electron 自动构建流程。
- Windows 增加免安装便携版 `.exe`，macOS 和 Linux 同时提供可直接启动的桌面发行格式。

## 1.0.0

- 初始版本。
