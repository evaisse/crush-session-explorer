// This file exists for require('crush-session-explorer') compatibility
// The main functionality is provided via the CLI binary

module.exports = {
  name: 'crush-session-explorer',
  description: 'A fast, lightweight CLI tool for exporting Crush chat sessions',
  binPath: require('path').join(__dirname, 'bin', process.platform === 'win32' ? 'crush-md.exe' : 'crush-md')
};
