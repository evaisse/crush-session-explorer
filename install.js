#!/usr/bin/env node

const https = require('https');
const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

// Get the package version from package.json
const packageJson = require('./package.json');
const version = packageJson.version;

// Determine platform and architecture
function getPlatform() {
  const platform = process.platform;
  const arch = process.arch;

  let platformName;
  let archName;

  // Map Node.js platform names to our binary naming
  switch (platform) {
    case 'darwin':
      platformName = 'darwin';
      break;
    case 'linux':
      platformName = 'linux';
      break;
    case 'win32':
      platformName = 'windows';
      break;
    default:
      throw new Error(`Unsupported platform: ${platform}`);
  }

  // Map Node.js arch names to our binary naming
  switch (arch) {
    case 'x64':
      archName = 'amd64';
      break;
    case 'arm64':
      archName = 'arm64';
      break;
    default:
      throw new Error(`Unsupported architecture: ${arch}`);
  }

  return { platformName, archName };
}

function getBinaryName() {
  const { platformName, archName } = getPlatform();
  const ext = platformName === 'windows' ? '.exe' : '';
  return `crush-md-${platformName}-${archName}${ext}`;
}

function getDownloadUrl() {
  const binaryName = getBinaryName();
  return `https://github.com/evaisse/crush-session-explorer/releases/download/v${version}/${binaryName}`;
}

function download(url, destination) {
  return new Promise((resolve, reject) => {
    console.log(`Downloading ${url}...`);
    
    const file = fs.createWriteStream(destination);
    
    https.get(url, (response) => {
      // Handle redirects
      if (response.statusCode === 302 || response.statusCode === 301) {
        download(response.headers.location, destination).then(resolve).catch(reject);
        return;
      }

      if (response.statusCode !== 200) {
        reject(new Error(`Failed to download: ${response.statusCode} ${response.statusMessage}`));
        return;
      }

      response.pipe(file);
      
      file.on('finish', () => {
        file.close();
        resolve();
      });
    }).on('error', (err) => {
      fs.unlink(destination, () => {}); // Delete the file on error
      reject(err);
    });
  });
}

async function install() {
  try {
    const binaryName = getBinaryName();
    const downloadUrl = getDownloadUrl();
    const binDir = path.join(__dirname, 'bin');
    
    // Ensure bin directory exists
    if (!fs.existsSync(binDir)) {
      fs.mkdirSync(binDir, { recursive: true });
    }

    // Determine the final binary name (without platform suffix)
    const finalBinaryName = process.platform === 'win32' ? 'crush-md.exe' : 'crush-md';
    const binaryPath = path.join(binDir, finalBinaryName);

    // Download the binary
    await download(downloadUrl, binaryPath);

    // Make it executable on Unix-like systems
    if (process.platform !== 'win32') {
      fs.chmodSync(binaryPath, 0o755);
    }

    console.log(`✅ Successfully installed crush-md v${version}`);
    console.log(`Binary location: ${binaryPath}`);
  } catch (error) {
    console.error('❌ Installation failed:', error.message);
    console.error('\nPlease try installing manually from:');
    console.error(`https://github.com/evaisse/crush-session-explorer/releases/tag/v${version}`);
    process.exit(1);
  }
}

install();
