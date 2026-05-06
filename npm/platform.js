#!/usr/bin/env node

"use strict";

const PLATFORM_MAP = {
  "darwin-x64":   { artifact: (v) => `upbit_${v}_macos_amd64.zip`,    binary: "upbit" },
  "darwin-arm64": { artifact: (v) => `upbit_${v}_macos_arm64.zip`,    binary: "upbit" },
  "linux-x64":    { artifact: (v) => `upbit_${v}_linux_amd64.tar.gz`, binary: "upbit" },
  "linux-arm64":  { artifact: (v) => `upbit_${v}_linux_arm64.tar.gz`, binary: "upbit" },
  "linux-ia32":   { artifact: (v) => `upbit_${v}_linux_386.tar.gz`,   binary: "upbit" },
  "linux-arm":    { artifact: (v) => `upbit_${v}_linux_armv6.tar.gz`, binary: "upbit" },
  "win32-x64":    { artifact: (v) => `upbit_${v}_windows_amd64.zip`,  binary: "upbit.exe" },
  "win32-arm64":  { artifact: (v) => `upbit_${v}_windows_arm64.zip`,  binary: "upbit.exe" },
  "win32-ia32":   { artifact: (v) => `upbit_${v}_windows_386.zip`,    binary: "upbit.exe" },
};

function getPlatform() {
  const key = `${process.platform}-${process.arch}`;
  const entry = PLATFORM_MAP[key];
  if (!entry) {
    throw new Error(
      `Unsupported platform: ${process.platform} ${process.arch}\n` +
      `Supported: ${Object.keys(PLATFORM_MAP).join(", ")}`
    );
  }
  return entry;
}

module.exports = { getPlatform };
