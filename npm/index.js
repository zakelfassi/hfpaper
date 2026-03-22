#!/usr/bin/env node
const { execFileSync } = require('child_process');
const path = require('path');

const ext = process.platform === 'win32' ? '.exe' : '';
const bin = path.join(__dirname, 'bin', `hfpaper${ext}`);

try {
  // Pass arguments through to the binary
  execFileSync(bin, process.argv.slice(2), { stdio: 'inherit' });
} catch (e) {
  // If the binary returns a non-zero exit code, exit with that code
  if (e.status !== undefined) {
    process.exit(e.status);
  }
  // If the binary isn't found or fails to run
  console.error('hfpaper binary not found or failed to execute.');
  console.error('Try reinstalling: npm install -g hfpaper');
  process.exit(1);
}
