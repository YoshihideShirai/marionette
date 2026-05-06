import fs from 'node:fs/promises';
import path from 'node:path';
import MarkdownIt from 'markdown-it';

const markdownIt = new MarkdownIt({ html: false, linkify: true, typographer: true });

const repoRoot = path.resolve(process.cwd(), '..', '..');
const apiDocsBaseDir = path.join(repoRoot, 'docs', 'api');
const apiDocsJaDir = path.join(apiDocsBaseDir, 'ja');

export async function getApiDocs(lang) {
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

export async function getApiDoc(lang, slug) {
  const docs = await getApiDocs(lang);
  return docs.find((doc) => doc.slug === slug);
}

export async function getApiNavItems(lang, langRoot) {
  const docs = await getApiDocs(lang);
  return docs.map((doc) => ({
    slug: doc.slug,
    title: doc.title,
    href: `${langRoot}/api/${doc.slug}/`,
  }));
}

export function renderApiMarkdown(markdown) {
  return markdownIt.render(rewriteMarkdownLinks(markdown));
}

function rewriteMarkdownLinks(markdown) {
  return markdown.replace(/\]\(\.\/([^)#]+?)\.md(#[^)]+)?\)/g, (_match, slug, hash = '') => `](../${slug}/${hash})`);
}

async function markdownFileNames(dir) {
  const fileNames = await fs.readdir(dir);
  return fileNames.filter((name) => name.endsWith('.md')).sort();
}

function slugFromFileName(fileName) {
  return fileName.replace(/\.md$/, '');
}

function titleFromMarkdown(markdown, fileName) {
  const heading = markdown.match(/^#\s+(.+)$/m) ?? markdown.match(/^##\s+(.+)$/m);
  return heading?.[1]?.trim() || slugFromFileName(fileName);
}
