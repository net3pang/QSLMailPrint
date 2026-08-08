(function installWailsBridge() {
  function invoke(method, args = []) {
    return new Promise((resolve, reject) => {
      const deadline = Date.now() + 5000;
      const attempt = () => {
        const app = window.go?.main?.App;
        if (typeof app?.[method] === 'function') {
          Promise.resolve(app[method](...args)).then(resolve, reject);
          return;
        }
        if (Date.now() >= deadline) {
          reject(new Error(`Wails method ${method} is unavailable`));
          return;
        }
        window.setTimeout(attempt, 25);
      };
      attempt();
    });
  }

  // Keep this a classic script: app.js must see electronAPI during its initialisation.
  window.electronAPI = {
    isElectron: false,
    getPrinters: () => invoke('GetPrinters'),
    saveTemplate: () => invoke('SaveTemplate'),
    saveRecord: (collection, record) => invoke('SaveRecord', [collection, record]).then(record => ({success: true, record})),
    listRecords: (collection) => invoke('ListRecords', [collection]),
    deleteRecord: (collection, id) => invoke('DeleteRecord', [collection, id]),
    printEnvelope: async options => {
      try {
        const result = await invoke('PrintEnvelope', [options]);
        if (result?.handled) return result;
      } catch (error) {
        console.warn('Wails native printing unavailable', error);
      }
      window.print();
      return {success: true};
    }
  };
})();
