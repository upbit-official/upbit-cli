#!/usr/bin/env node

"use strict";

const crypto = require("crypto");
const fs = require("fs");
const path = require("path");
const os = require("os");
const { pipeline } = require("stream/promises");
const { createWriteStream, mkdirSync, rmSync } = require("fs");
const { spawnSync } = require("child_process");
const { getPlatform } = require("./platform");

const { version } = require("./package.json");
const REPO = process.env.UPBIT_REPO || "upbit-official/upbit-cli";
const INSTALL_DIR = path.join(__dirname, "bin");
const VERSION_FILE = path.join(INSTALL_DIR, ".version");
const TOKEN = process.env.GITHUB_TOKEN || process.env.GH_TOKEN || null;

function releaseUrl(filename) {
  return `https://github.com/${REPO}/releases/download/v${version}/${filename}`;
}

// GitHub's /releases/download/ URL does not accept Authorization headers.
// For private repos we must use the GitHub API (/releases/assets/{id}) instead,
// which returns a 302 to a signed CDN URL (auth stripped automatically on redirect).
let _assetIndex = null;
async function getAssetIndex() {
  if (_assetIndex) return _assetIndex;
  const url = `https://api.github.com/repos/${REPO}/releases/tags/v${version}`;
  const res = await fetch(url, {
    headers: {
      Authorization: `Bearer ${TOKEN}`,
      Accept: "application/vnd.github+json",
      "X-GitHub-Api-Version": "2022-11-28",
    },
  });
  if (!res.ok) throw new Error(`GitHub API: ${res.status} ${res.statusText} — ${url}`);
  const data = await res.json();
  _assetIndex = Object.fromEntries(data.assets.map((a) => [a.name, a.id]));
  return _assetIndex;
}

async function resolveUrl(filename) {
  if (!TOKEN) return releaseUrl(filename);
  const index = await getAssetIndex();
  const id = index[filename];
  if (!id) throw new Error(`Asset not found in release: ${filename}`);
  return `https://api.github.com/repos/${REPO}/releases/assets/${id}`;
}

async function apiFetch(url) {
  const headers = { "User-Agent": "upbit-cli-installer" };
  if (TOKEN) {
    headers.Authorization = `Bearer ${TOKEN}`;
    headers.Accept = "application/octet-stream";
  }
  return fetch(url, { headers, redirect: "follow" });
}

async function fetchText(filename) {
  const url = await resolveUrl(filename);
  const res = await apiFetch(url);
  if (!res.ok) throw new Error(`GET ${url}: ${res.status} ${res.statusText}`);
  return res.text();
}

async function downloadFile(filename, dest) {
  const url = await resolveUrl(filename);
  const res = await apiFetch(url);
  if (!res.ok) throw new Error(`GET ${url}: ${res.status} ${res.statusText}`);
  if (!res.body) throw new Error(`Empty response body: ${url}`);
  const { Readable } = require("stream");
  await pipeline(Readable.fromWeb(res.body), createWriteStream(dest));
}

function sha256File(filePath) {
  return crypto.createHash("sha256").update(fs.readFileSync(filePath)).digest("hex").toLowerCase();
}

function extract(archivePath, destDir) {
  let result;
  if (archivePath.endsWith(".zip")) {
    if (process.platform === "win32") {
      const ps = archivePath.replace(/'/g, "''");
      const pd = destDir.replace(/'/g, "''");
      result = spawnSync("powershell.exe", [
        "-NoProfile", "-NonInteractive", "-Command",
        `Expand-Archive -LiteralPath '${ps}' -DestinationPath '${pd}' -Force`,
      ], { stdio: "pipe" });
    } else {
      result = spawnSync("unzip", ["-q", "-o", archivePath, "-d", destDir], { stdio: "pipe" });
    }
  } else {
    result = spawnSync("tar", ["xf", archivePath, "-C", destDir], { stdio: "pipe" });
  }
  if (result.error) throw new Error(`Extract failed: ${result.error.message}`);
  if ((result.status ?? 1) !== 0) {
    throw new Error(`Extract failed (exit ${result.status}): ${(result.stderr ?? "").toString()}`);
  }
}

async function install() {
  if ((process.env.CI || process.env.GITHUB_ACTIONS) && !TOKEN) {
    process.stderr.write(`Warning: CI environment detected without token — skipping binary installation.\n`);
    return;
  }

  const platform = getPlatform();
  const artifactName = platform.artifact(version);
  const binPath = path.join(INSTALL_DIR, platform.binary);

  if (fs.existsSync(binPath) && fs.existsSync(VERSION_FILE)) {
    if (fs.readFileSync(VERSION_FILE, "utf8").trim() === version) {
      process.stderr.write(`upbit v${version} is already installed.\n`);
      return;
    }
  }

  if (fs.existsSync(INSTALL_DIR)) rmSync(INSTALL_DIR, { recursive: true, force: true });
  mkdirSync(INSTALL_DIR, { recursive: true });

  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), "upbit-install-"));
  const archivePath = path.join(tmpDir, artifactName);

  try {
    const checksumsName = `upbit_${version}_checksums.txt`;
    process.stderr.write(`Downloading checksums...\n`);
    const checksumsText = await fetchText(checksumsName);

    process.stderr.write(`Downloading ${artifactName}...\n`);
    await downloadFile(artifactName, archivePath);

    const expectedHash = checksumsText
      .split("\n")
      .map((line) => line.trim())
      .filter((line) => line.endsWith(artifactName))
      .map((line) => line.split(/\s+/)[0].toLowerCase())[0];

    if (!expectedHash) {
      throw new Error(`No checksum entry found for ${artifactName}`);
    }

    const actualHash = sha256File(archivePath);
    if (actualHash !== expectedHash) {
      throw new Error(`SHA256 mismatch!\n  Expected: ${expectedHash}\n  Actual:   ${actualHash}`);
    }
    process.stderr.write(`Checksum verified.\n`);

    process.stderr.write(`Extracting...\n`);
    extract(archivePath, INSTALL_DIR);

    if (process.platform !== "win32") fs.chmodSync(binPath, 0o755);

    fs.writeFileSync(VERSION_FILE, version);
    process.stderr.write(`upbit v${version} installed successfully.\n`);
  } finally {
    rmSync(tmpDir, { recursive: true, force: true });
  }
}

install().catch((err) => {
  process.stderr.write(`Error installing upbit: ${err.message}\n`);
  process.exit(1);
});
