const { existsSync } = require("fs");
const { join, resolve } = require("path");
const { spawnSync } = require("child_process");

if (process.env.CI || process.env.HUSKY === "0") {
  console.log("Skipping husky install in CI or when HUSKY=0");
  process.exit(0);
}

const webDir = process.cwd();
const repoDir = resolve(webDir, "..");
const localGitDir = join(webDir, ".git");
const repoGitDir = join(repoDir, ".git");

let cwd = webDir;
let hooksDir = ".husky";

if (existsSync(localGitDir)) {
  cwd = webDir;
  hooksDir = ".husky";
} else if (existsSync(repoGitDir)) {
  cwd = repoDir;
  hooksDir = "web/.husky";
} else {
  console.log("Skipping husky install because no .git directory was found");
  process.exit(0);
}

const localHusky = join(
  webDir,
  "node_modules",
  ".bin",
  process.platform === "win32" ? "husky.cmd" : "husky",
);

const command = existsSync(localHusky)
  ? localHusky
  : process.platform === "win32"
    ? "npx.cmd"
    : "npx";
const args = existsSync(localHusky)
  ? ["install", hooksDir]
  : ["husky", "install", hooksDir];

const result = spawnSync(command, args, {
  cwd,
  stdio: "inherit",
});

if (result.error) {
  throw result.error;
}

process.exit(result.status || 0);
