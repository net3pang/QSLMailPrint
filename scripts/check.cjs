const { spawnSync } = require('node:child_process');
const path = require('node:path');

const root = path.resolve(__dirname, '..');
const backendRoot = path.join(root, 'backend');
const goCache = process.env.GOCACHE || path.join(require('node:os').tmpdir(), 'qsl-mail-go-cache');

function run(command, args, options = {}) {
  const result = spawnSync(command, args, {stdio:'inherit', ...options});
  if (result.error) throw result.error;
  if (result.status !== 0) process.exit(result.status || 1);
}

run(process.execPath, ['--check', path.join(root, 'app.js')]);
run('go', ['test', './...'], {cwd:backendRoot, env:{...process.env, GOCACHE:goCache}});
