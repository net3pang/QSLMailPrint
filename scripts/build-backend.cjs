const { existsSync, mkdirSync } = require('node:fs');
const { spawnSync } = require('node:child_process');
const path = require('node:path');

const root = path.resolve(__dirname, '..');
const backendRoot = path.join(root, 'backend');
const goos = process.env.GOOS || process.platform;
const goarch = process.env.GOARCH || process.arch;
const extension = goos === 'windows' ? '.exe' : '';
const outputDir = process.env.BACKEND_OUTPUT_DIR
  ? path.resolve(root, process.env.BACKEND_OUTPUT_DIR)
  : path.join(root, 'backend', 'bin');
const output = path.join(outputDir, `qsl-mail-backend${extension}`);

mkdirSync(path.dirname(output), {recursive:true});
const result = spawnSync('go', ['build', '-trimpath', '-o', output, '.'], {
  cwd: backendRoot,
  env: {...process.env, GOOS:goos, GOARCH:goarch, CGO_ENABLED:'0', GOCACHE:process.env.GOCACHE || path.join(require('node:os').tmpdir(), 'qsl-mail-go-cache')},
  stdio:'inherit'
});

if (result.error) throw result.error;
if (result.status !== 0 || !existsSync(output)) process.exit(result.status || 1);
console.log(`Built ${goos}/${goarch}: ${path.relative(root, output)}`);
