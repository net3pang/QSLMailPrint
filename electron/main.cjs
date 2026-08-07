const { app, BrowserWindow, ipcMain } = require('electron');
const { existsSync } = require('node:fs');
const { spawn } = require('node:child_process');
const http = require('node:http');
const path = require('node:path');

let mainWindow;
let backendProcess;

function mmToMicrons(value) {
  return Math.max(1, Math.round(Number(value || 1) * 1000));
}

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1440,
    height: 960,
    minWidth: 980,
    minHeight: 720,
    backgroundColor: '#f3f6f6',
    webPreferences: {
      preload: path.join(__dirname, 'preload.cjs'),
      contextIsolation: true,
      nodeIntegration: false
    }
  });
  mainWindow.loadFile(path.join(__dirname, '..', 'index.html'));
}

function startGoBackend() {
  const filename = process.platform === 'win32' ? 'qsl-mail-backend.exe' : 'qsl-mail-backend';
  const candidates = [
    path.join(process.resourcesPath || '', 'backend', filename),
    path.join(__dirname, '..', 'backend', 'bin', filename),
    path.join(__dirname, '..', 'backend', filename)
  ];
  const binary = candidates.find(existsSync);
  if (binary) backendProcess = spawn(binary, [], {stdio:'ignore'});
}

function saveRecord(collection, record) {
  return new Promise((resolve, reject) => {
    const payload = Buffer.from(JSON.stringify(record));
    const request = http.request({
      hostname:'127.0.0.1', port:38765, path:`/api/${collection}`, method:'POST',
      headers:{'Content-Type':'application/json','Content-Length':payload.length}
    }, response => {
      let body = '';
      response.on('data', chunk => { body += chunk; });
      response.on('end', () => {
        if (response.statusCode >= 200 && response.statusCode < 300) resolve(JSON.parse(body));
        else reject(new Error(body));
      });
    });
    request.on('error', reject);
    request.write(payload);
    request.end();
  });
}

ipcMain.handle('printers:list', async event => {
  const window = BrowserWindow.fromWebContents(event.sender);
  if (!window) return [];
  return window.webContents.getPrintersAsync();
});
ipcMain.handle('storage:save', async (_event, collection, record) => { try { return {success:true, record:await saveRecord(collection, record)}; } catch (error) { return {success:false, reason:error.message}; } });

ipcMain.handle('print:envelope', async (event, options = {}) => {
  const window = BrowserWindow.fromWebContents(event.sender);
  if (!window) return {success:false, reason:'打印窗口不存在'};
  return new Promise(resolve => {
    const landscape = Boolean(options.landscape);
    const width = landscape ? options.width : options.height;
    const height = landscape ? options.height : options.width;
    window.webContents.print({
      silent: options.silent === true,
      deviceName: options.deviceName || undefined,
      printBackground: false,
      landscape: false,
      scaleFactor: Math.max(10, Math.min(200, Number(options.scale) || 100)),
      pageSize: {width: mmToMicrons(width), height: mmToMicrons(height)},
      margins: {marginType:'custom', top:0, bottom:0, left:0, right:0}
    }, (success, failureReason) => resolve({success, reason:failureReason || ''}));
  });
});

app.whenReady().then(() => {
  startGoBackend();
  createWindow();
  app.on('activate', () => { if (BrowserWindow.getAllWindows().length === 0) createWindow(); });
});

app.on('before-quit', () => { if (backendProcess) backendProcess.kill(); });
app.on('window-all-closed', () => { if (process.platform !== 'darwin') app.quit(); });
