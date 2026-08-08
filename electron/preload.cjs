const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('electronAPI', {
  isElectron: true,
  getPrinters: () => ipcRenderer.invoke('printers:list'),
  saveRecord: (collection, record) => ipcRenderer.invoke('storage:save', collection, record),
  listRecords: (collection) => ipcRenderer.invoke('storage:list', collection).then(result => result?.success ? result.records : []),
  deleteRecord: (collection, id) => ipcRenderer.invoke('storage:delete', collection, id),
  printEnvelope: options => ipcRenderer.invoke('print:envelope', options)
});
