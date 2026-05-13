#!/usr/bin/env node
import { spawnSync } from 'node:child_process';
import { mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const scriptDir = path.dirname(fileURLToPath(import.meta.url));
const siteRoot = path.resolve(scriptDir, '..');
const repoRoot = path.resolve(siteRoot, '..', '..');
const reportsDir = path.join(siteRoot, 'public', 'reports');
const frameworkSourcePaths = ['backend', 'desktop', 'frontend', 'internal/componenttmpl'];
const clocExcludeArgs = [
  '--fullpath',
  '--exclude-dir=testdata',
  '--not-match-f=(^|/)([^/]+_test\\.go|[^/]+\\.test\\.[^/]+)$',
];

mkdirSync(reportsDir, { recursive: true });

function run(command, args, options = {}) {
  console.log(`$ ${[command, ...args].join(' ')}`);
  const result = spawnSync(command, args, {
    cwd: options.cwd ?? repoRoot,
    encoding: 'utf-8',
    stdio: options.capture ? ['ignore', 'pipe', 'pipe'] : 'inherit',
  });

  if (result.status !== 0) {
    if (options.capture) {
      process.stdout.write(result.stdout ?? '');
      process.stderr.write(result.stderr ?? '');
    }
    throw new Error(`${command} exited with status ${result.status}`);
  }

  return options.capture ? result.stdout.trim() : '';
}

const generatedAt = new Date().toISOString();
const commit = run('git', ['rev-parse', '--short', 'HEAD'], { capture: true });
const clocJsonPath = path.join(reportsDir, 'cloc.json');
const clocTextPath = path.join(reportsDir, 'cloc.txt');
const coverageProfilePath = path.join(reportsDir, 'coverage.out');
const coverageHtmlPath = path.join(reportsDir, 'coverage.html');
const coverageTextPath = path.join(reportsDir, 'coverage.txt');

run('cloc', [
  ...frameworkSourcePaths,
  ...clocExcludeArgs,
  '--json',
  `--out=${clocJsonPath}`,
]);
run('cloc', [
  ...frameworkSourcePaths,
  ...clocExcludeArgs,
  `--out=${clocTextPath}`,
]);

const modulePath = run('go', ['list', '-m'], { capture: true });
const frameworkPackagePrefixes = frameworkSourcePaths.map((sourcePath) =>
  `${modulePath}/${sourcePath.replaceAll(path.sep, '/')}`,
);
const isFrameworkPackage = (pkg) =>
  frameworkPackagePrefixes.some((prefix) => pkg === prefix || pkg.startsWith(`${prefix}/`));
const packages = run('go', ['list', './...'], { capture: true })
  .split('\n')
  .filter((pkg) => pkg && isFrameworkPackage(pkg));

if (packages.length === 0) {
  throw new Error('No framework packages found for coverage report');
}

const coverageCommand = 'go test <framework packages> -coverprofile=coverage.out -covermode=atomic';

run('go', ['test', ...packages, `-coverprofile=${coverageProfilePath}`, '-covermode=atomic']);
run('go', ['tool', 'cover', `-html=${coverageProfilePath}`, `-o=${coverageHtmlPath}`]);
run('go', ['tool', 'cover', `-func=${coverageProfilePath}`, `-o=${coverageTextPath}`]);

const cloc = JSON.parse(readFileSync(clocJsonPath, 'utf-8'));
const coverageText = readFileSync(coverageTextPath, 'utf-8');
const coverageTotalLine = coverageText
  .split('\n')
  .map((line) => line.trim())
  .find((line) => line.startsWith('total:'));
const coverageTotal = coverageTotalLine?.match(/([0-9]+(?:\.[0-9]+)?)%/)?.[1] ?? null;
const languages = Object.entries(cloc)
  .filter(([language]) => !['header', 'SUM'].includes(language))
  .map(([language, stats]) => ({
    language,
    files: stats.nFiles ?? 0,
    blank: stats.blank ?? 0,
    comment: stats.comment ?? 0,
    code: stats.code ?? 0,
  }))
  .sort((a, b) => b.code - a.code);

const summary = {
  generatedAt,
  commit,
  cloc: {
    files: cloc.SUM?.nFiles ?? 0,
    blank: cloc.SUM?.blank ?? 0,
    comment: cloc.SUM?.comment ?? 0,
    code: cloc.SUM?.code ?? 0,
    languages,
  },
  coverage: {
    total: coverageTotal,
    totalLine: coverageTotalLine ?? null,
    command: coverageCommand,
    packageCount: packages.length,
  },
  links: {
    clocText: '/marionette/reports/cloc.txt',
    clocJson: '/marionette/reports/cloc.json',
    coverageHtml: '/marionette/reports/coverage.html',
    coverageText: '/marionette/reports/coverage.txt',
    coverageProfile: '/marionette/reports/coverage.out',
  },
};

writeFileSync(path.join(reportsDir, 'summary.json'), `${JSON.stringify(summary, null, 2)}\n`);
console.log(`Wrote reports to ${reportsDir}`);
