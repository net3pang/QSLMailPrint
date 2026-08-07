const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('electronAPI', {
  isElectron: true,
  getPrinters: () => ipcRenderer.invoke('printers:list'),
  saveRecord: (collection, record) => ipcRenderer.invoke('storage:save', collection, record),
  printEnvelope: options => ipcRenderer.invoke('print:envelope', options)
});
