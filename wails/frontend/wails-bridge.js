import {GetPrinters, SaveRecord} from './wailsjs/go/main/App.js';

window.electronAPI = {
  isElectron: false,
  getPrinters: () => GetPrinters(),
  saveRecord: async (collection, record) => ({
    success: true,
    record: await SaveRecord(collection, record)
  })
};
