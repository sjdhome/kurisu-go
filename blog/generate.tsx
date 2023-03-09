import React from "react";
import * as ReactDOMServer from "react-dom/server";
import Showdown from "showdown";
import fs from "node:fs/promises";
import { cwd } from "node:process";
import Post from "./page/post.js";

const convertor = new Showdown.Converter();

await fs.rm(`${cwd()}/blog/dist`, { recursive: true, force: true });
await fs.mkdir(`${cwd()}/blog/dist/post`, { recursive: true });

async function renderPages() {
  const pages = await fs.readdir(`${cwd()}/blog/page`);
  for (const pageFilename of pages) {
    if (!pageFilename.endsWith(".js")) continue;
    if (pageFilename === "post.js") continue; /* We'll render this page later */
    const pageModule = await import(`${cwd()}/blog/page/${pageFilename}`);
    const Page = pageModule.default;
    const html = ReactDOMServer.renderToString(<Page />);
    await fs.writeFile(
      `${cwd()}/blog/dist/${pageFilename.replace(".js", ".html")}`,
      `<!DOCTYPE html>${html}`
    );
  }
  const posts = await fs.readdir(`${cwd()}/blog/post`);
  for (const postFilename of posts) {
    const postMarkdown = await fs.readFile(
      `${cwd()}/blog/post/${postFilename}`,
      "utf-8"
    );
    const postHTML = convertor.makeHtml(postMarkdown);
    const post = ReactDOMServer.renderToString(<Post postHTML={postHTML} />);
    await fs.writeFile(
      `${cwd()}/blog/dist/post/${postFilename.replace(".md", ".html")}`,
      `<!DOCTYPE html>${post}`
    );
  }
}

await renderPages();
