import matter from "gray-matter";
import fs from "node:fs/promises";
import { cwd } from "node:process";
import React from "react";
import * as ReactDOMServer from "react-dom/server";
import Showdown from "showdown";
import { title } from "./constants.js";
import Post from "./page/post.js";
import { PostData } from "./PostData.js";
import fse from "fs-extra";

const convertor = new Showdown.Converter();

function render() {
  const renderPage = async () => {
    const pages = await fs.readdir(`${cwd()}/blog/page`);
    for (const pageFilename of pages) {
      if (!pageFilename.endsWith(".js")) continue;
      if (pageFilename === "post.js")
        continue; /* We'll render this page later */
      const pageModule = await import(`${cwd()}/blog/page/${pageFilename}`);
      const Page = pageModule.default;
      const html = ReactDOMServer.renderToString(<Page />);
      await fs.writeFile(
        `${cwd()}/blog/dist/${pageFilename.replace(".js", ".html")}`,
        `<!DOCTYPE html>${html}`
      );
    }
  };

  const renderPost = async () => {
    const posts = await fs.readdir(`${cwd()}/blog/post`);
    for (const postFilename of posts) {
      const rawPost = await fs.readFile(
        `${cwd()}/blog/post/${postFilename}`,
        "utf-8"
      );
      const postMatter = matter(rawPost);
      const postData = postMatter.data as PostData;
      const postHTML = convertor.makeHtml(postMatter.content);
      const post = ReactDOMServer.renderToString(
        <Post title={`${postData.title} | ${title}`} postHTML={postHTML} />
      );
      await fs.writeFile(
        `${cwd()}/blog/dist/post/${postFilename.replace(".md", ".html")}`,
        `<!DOCTYPE html>${post}`
      );
    }
  };

  return Promise.all([renderPage(), renderPost()]);
}

await fs.rm(`${cwd()}/blog/dist`, { recursive: true, force: true });
await fs.mkdir(`${cwd()}/blog/dist/post`, { recursive: true });
await Promise.all([
  render(),
  fse.copy(`${cwd()}/blog/css`, `${cwd()}/blog/dist/css`),
  fse.copy(`${cwd()}/blog/img`, `${cwd()}/blog/dist/img`),
]);
