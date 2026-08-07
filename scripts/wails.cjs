const {spawnSync} = require('node:child_process');
const path = require('node:path');
const fs = require('node:fs');

const command = process.argv[2] || 'build';
if (!['dev', 'build'].includes(command)) {
  console.error('Usage: npm run wails:dev | npm run wails:build');
  process.exit(2);
}

function findWails() {
  if (process.env.WAILS_BIN) return process.env.WAILS_BIN;
  const lookup = spawnSync('wails', ['version'], {stdio: 'ignore'});
  if (!lookup.error) return 'wails';
  const goPath = spawnSync('go', ['env', 'GOPATH'], {encoding: 'utf8'}).stdout?.trim();
  if (!goPath) return 'wails';
  const binary = process.platform === 'win32' ? 'wails.exe' : 'wails';
  const fallback = path.join(goPath, 'bin', binary);
  return fs.existsSync(fallback) ? fallback : 'wails';
}

const result = spawnSync(findWails(), [command, ...process.argv.slice(3)], {
  cwd: path.resolve(__dirname, '../wails'),
  stdio: 'inherit',
  env: process.env
});

if (result.error) {
  console.error('找不到 Wails CLI。安装：go install github.com/wailsapp/wails/v2/cmd/wails@v2.10.2');
  process.exit(1);
}
process.exit(result.status ?? 1);
