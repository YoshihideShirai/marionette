import fs from 'node:fs/promises';
import path from 'node:path';
import MarkdownIt from 'markdown-it';

export type ApiLang = 'en' | 'ja';

export type ApiDocEntry = {
  slug: string;
  fileName: string;
  title: string;
  sourcePath: string;
  markdown: string;
};

export type ApiNavItem = {
  slug: string;
  title: string;
  href: string;
};

const markdownIt = new MarkdownIt({ html: false, linkify: true, typographer: true });

const repoRoot = path.resolve(process.cwd(), '..', '..');
const apiDocsBaseDir = path.join(repoRoot, 'docs', 'api');
const apiDocsJaDir = path.join(apiDocsBaseDir, 'ja');

export async function getApiDocs(lang: ApiLang): Promise<ApiDocEntry[]> {
  const [baseFileNames, localizedFileNames] = await Promise.all([
    markdownFileNames(apiDocsBaseDir),
    lang === 'ja' ? markdownFileNames(apiDocsJaDir) : Promise.resolve([]),
  ]);
  const localizedFiles = new Set(localizedFileNames);

  return Promise.all(
    baseFileNames.map(async (fileName) => {
      const sourcePath = lang === 'ja' && localizedFiles.has(fileName)
        ? path.join(apiDocsJaDir, fileName)
        : path.join(apiDocsBaseDir, fileName);
      const markdown = await fs.readFile(sourcePath, 'utf-8');

      return {
        slug: slugFromFileName(fileName),
        fileName,
        title: titleFromMarkdown(markdown, fileName),
        sourcePath,
        markdown,
      };
    }),
  );
}

export async function getApiDoc(lang: ApiLang, slug: string): Promise<ApiDocEntry | undefined> {
  const docs = await getApiDocs(lang);
  return docs.find((doc) => doc.slug === slug);
}

export async function getApiNavItems(lang: ApiLang, langRoot: string): Promise<ApiNavItem[]> {
  const docs = await getApiDocs(lang);
  return docs.map((doc) => ({
    slug: doc.slug,
    title: doc.title,
    href: `${langRoot}/api/${doc.slug}/`,
  }));
}

export function renderApiMarkdown(markdown: string): string {
  return markdownIt.render(rewriteMarkdownLinks(markdown));
}

function rewriteMarkdownLinks(markdown: string): string {
  return markdown.replace(/\]\(\.\/([^)#]+?)\.md(#[^)]+)?\)/g, (_match, slug, hash = '') => `](../${slug}/${hash})`);
}

async function markdownFileNames(dir: string): Promise<string[]> {
  const fileNames = await fs.readdir(dir);
  return fileNames.filter((name) => name.endsWith('.md')).sort();
}

function slugFromFileName(fileName: string): string {
  return fileName.replace(/\.md$/, '');
}

function titleFromMarkdown(markdown: string, fileName: string): string {
  const heading = markdown.match(/^#\s+(.+)$/m) ?? markdown.match(/^##\s+(.+)$/m);
  return heading?.[1]?.trim() || slugFromFileName(fileName);
}
