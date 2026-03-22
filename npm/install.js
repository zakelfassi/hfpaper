const fs = require('fs');
const path = require('path');
const https = require('https');

const pkg = require('./package.json');
const VERSION = `v${pkg.version}`;

const PLATFORM_MAPPING = {
  'darwin': 'darwin',
  'linux': 'linux',
  'win32': 'windows'
};

const ARCH_MAPPING = {
  'x64': 'amd64',
  'arm64': 'arm64'
};

function fail(msg) {
  console.warn(`\n[hfpaper] Warning: ${msg}`);
  console.warn('[hfpaper] The CLI might not work. You may need to install it manually from GitHub releases.\n');
  // Exit successfully so we don't break the entire npm install process for the user
  process.exit(0);
}

const platform = PLATFORM_MAPPING[process.platform];
const arch = ARCH_MAPPING[process.arch];

if (!platform || !arch) {
  fail(`Unsupported platform/arch: ${process.platform}/${process.arch}`);
}

let binaryName = `hfpaper-${platform}-${arch}`;
if (platform === 'windows') binaryName += '.exe';

// Use the tag from package.json version
const url = `https://github.com/zakelfassi/hfpaper/releases/download/${VERSION}/${binaryName}`;
const binDir = path.join(__dirname, 'bin');
const destName = platform === 'windows' ? 'hfpaper.exe' : 'hfpaper';
const destPath = path.join(binDir, destName);

if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

console.log(`[hfpaper] Downloading ${VERSION} for ${platform}/${arch}...`);

function download(downloadUrl, dest) {
  https.get(downloadUrl, (response) => {
    // Handle redirects (GitHub releases usually redirect to S3/Azure)
    if (response.statusCode === 301 || response.statusCode === 302) {
      if (response.headers.location) {
        download(response.headers.location, dest);
      } else {
        fail('Redirect with no location header');
      }
      return;
    }

    if (response.statusCode !== 200) {
      fail(`Failed to download: HTTP ${response.statusCode} (${downloadUrl})`);
      return;
    }

    const file = fs.createWriteStream(dest);
    response.pipe(file);

    file.on('finish', () => {
      file.close(() => {
        console.log('[hfpaper] Download complete.');
        if (platform !== 'windows') {
          try {
            fs.chmodSync(dest, 0o755);
          } catch (e) {
            console.warn('[hfpaper] Failed to chmod binary: ' + e.message);
          }
        }
      });
    });

    file.on('error', (err) => {
      fs.unlink(dest, () => {}); // Delete failed file
      fail(`File write error: ${err.message}`);
    });

  }).on('error', (err) => {
    fail(`Network error: ${err.message}`);
  });
}

download(url, destPath);
