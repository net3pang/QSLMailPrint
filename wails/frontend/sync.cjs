const fs = require('node:fs');
const path = require('node:path');

const frontend = __dirname;
const projectRoot = path.resolve(frontend, '../..');
for (const file of ['app.js', 'styles.css']) {
  fs.copyFileSync(path.join(projectRoot, file), path.join(frontend, file));
}

let index = fs.readFileSync(path.join(projectRoot, 'index.html'), 'utf8');
if (!index.includes('<script src="app.js"></script>')) {
  throw new Error('根目录 index.html 缺少 app.js 脚本入口');
}
index = index.replace('<script src="app.js"></script>', '<script type="module" src="wails-bridge.js"></script>\n    <script type="module" src="app.js"></script>');
fs.writeFileSync(path.join(frontend, 'index.html'), index);
