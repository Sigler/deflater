// Renders assets/icon.svg into the icon files Wails expects:
//   build/appicon.png        1024px master
//   build/windows/icon.ico   multi-size Windows icon (PNG-compressed)
// Run from frontend/: npm run icons

import { mkdir, readFile, writeFile } from "node:fs/promises";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { Resvg } from "@resvg/resvg-js";

const root = join(dirname(fileURLToPath(import.meta.url)), "..", "..");
const svg = await readFile(join(root, "assets", "icon.svg"), "utf8");

function renderPng(size) {
  const resvg = new Resvg(svg, { fitTo: { mode: "width", value: size } });
  return resvg.render().asPng();
}

// Master PNG for Wails.
await mkdir(join(root, "build", "windows"), { recursive: true });
await writeFile(join(root, "build", "appicon.png"), renderPng(1024));

// ICO container: header + directory entries + PNG payloads.
const sizes = [16, 20, 24, 32, 48, 64, 128, 256];
const pngs = sizes.map((s) => ({ size: s, data: renderPng(s) }));

const header = Buffer.alloc(6);
header.writeUInt16LE(0, 0); // reserved
header.writeUInt16LE(1, 2); // type: icon
header.writeUInt16LE(pngs.length, 4);

const entries = [];
let offset = 6 + pngs.length * 16;
for (const { size, data } of pngs) {
  const entry = Buffer.alloc(16);
  entry.writeUInt8(size >= 256 ? 0 : size, 0); // width, 0 means 256
  entry.writeUInt8(size >= 256 ? 0 : size, 1); // height
  entry.writeUInt8(0, 2); // palette
  entry.writeUInt8(0, 3); // reserved
  entry.writeUInt16LE(1, 4); // planes
  entry.writeUInt16LE(32, 6); // bits per pixel
  entry.writeUInt32LE(data.length, 8);
  entry.writeUInt32LE(offset, 12);
  entries.push(entry);
  offset += data.length;
}

const ico = Buffer.concat([header, ...entries, ...pngs.map((p) => p.data)]);
await writeFile(join(root, "build", "windows", "icon.ico"), ico);

console.log(`Wrote build/appicon.png and build/windows/icon.ico (${sizes.join(", ")} px).`);
