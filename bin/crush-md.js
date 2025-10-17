#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');

const binPath = path.join(__dirname, 'crush-md');

// Forward all arguments and spawn the binary
const child = spawn(binPath, process.argv.slice(2), { stdio: 'inherit' });

child.on('exit', (code) => {
  process.exit(code);
});
