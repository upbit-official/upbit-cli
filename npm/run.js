#!/usr/bin/env node

"use strict";

const path = require("path");
const fs = require("fs");
const { spawnSync } = require("child_process");
const { getPlatform } = require("./platform");

let platform;
try {
  platform = getPlatform();
} catch (err) {
  process.stderr.write(`${err.message}\n`);
  process.exit(1);
}

const binPath = path.join(__dirname, "bin", platform.binary);

if (!fs.existsSync(binPath)) {
  process.stderr.write(`upbit binary not found, installing...\n`);
  const install = spawnSync(process.execPath, [path.join(__dirname, "install.js")], {
    cwd: __dirname,
    stdio: "inherit",
  });
  if (install.error) {
    process.stderr.write(`Install failed: ${install.error.message}\n`);
    process.exit(1);
  }
  if ((install.status ?? 1) !== 0) process.exit(install.status ?? 1);
}

const result = spawnSync(binPath, process.argv.slice(2), {
  cwd: process.cwd(),
  stdio: "inherit",
});

if (result.error) {
  process.stderr.write(`Error running upbit: ${result.error.message}\n`);
  process.exit(1);
}

process.exit(result.status ?? 1);
