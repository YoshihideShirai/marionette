import fs from 'node:fs/promises';
import path from 'node:path';

export type ReportSummary = {
  generatedAt: string | null;
  commit: string | null;
  cloc: {
    files: number;
    blank: number;
    comment: number;
    code: number;
    languages: Array<{
      language: string;
      files: number;
      blank: number;
      comment: number;
      code: number;
    }>;
  };
  coverage: {
    total: string | null;
    totalLine: string | null;
  };
  links: {
    clocText: string;
    clocJson: string;
    coverageHtml: string;
    coverageText: string;
    coverageProfile: string;
  };
};

const fallbackSummary: ReportSummary = {
  generatedAt: null,
  commit: null,
  cloc: {
    files: 0,
    blank: 0,
    comment: 0,
    code: 0,
    languages: [],
  },
  coverage: {
    total: null,
    totalLine: null,
  },
  links: {
    clocText: '/marionette/reports/cloc.txt',
    clocJson: '/marionette/reports/cloc.json',
    coverageHtml: '/marionette/reports/coverage.html',
    coverageText: '/marionette/reports/coverage.txt',
    coverageProfile: '/marionette/reports/coverage.out',
  },
};

export async function getReportSummary(): Promise<ReportSummary> {
  const summaryPath = path.resolve(process.cwd(), 'public', 'reports', 'summary.json');

  try {
    const raw = await fs.readFile(summaryPath, 'utf-8');
    return { ...fallbackSummary, ...JSON.parse(raw) };
  } catch (error) {
    const err = error as Error & { code?: string };
    if (err.code !== 'ENOENT') {
      console.warn(`Unable to read report summary at ${summaryPath}: ${err.message}`);
    }
    return fallbackSummary;
  }
}

export function formatNumber(value: number): string {
  return new Intl.NumberFormat('en-US').format(value);
}

export function formatGeneratedAt(value: string | null, locale: string): string {
  if (!value) return locale === 'ja-JP' ? '未生成' : 'Not generated';
  return new Intl.DateTimeFormat(locale, {
    dateStyle: 'medium',
    timeStyle: 'short',
    timeZone: 'UTC',
  }).format(new Date(value));
}
