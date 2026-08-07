const path = require('node:path');
const { notarize } = require('@electron/notarize');

module.exports = async context => {
  if (context.electronPlatformName !== 'darwin') return;

  const appleId = process.env.APPLE_ID;
  const appleIdPassword = process.env.APPLE_APP_SPECIFIC_PASSWORD;
  const teamId = process.env.APPLE_TEAM_ID;
  if (!appleId || !appleIdPassword || !teamId) {
    console.warn('macOS notarization skipped: APPLE_ID, APPLE_APP_SPECIFIC_PASSWORD, and APPLE_TEAM_ID are required.');
    return;
  }

  const appName = context.packager.appInfo.productFilename;
  await notarize({
    appBundleId: context.packager.appInfo.id,
    appPath: path.join(context.appOutDir, `${appName}.app`),
    appleId,
    appleIdPassword,
    teamId
  });
};
