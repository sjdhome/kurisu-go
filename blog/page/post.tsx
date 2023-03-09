import React from "react";
import { title } from "../constants.js";

function Post(props: { title: string; postHTML: string }) {
  const { postHTML, title } = props;
  return (
    <html>
      <head>
        <meta charSet="utf-8" />
        <title>{title}</title>
      </head>
      <body>
        <main dangerouslySetInnerHTML={{ __html: postHTML }} />
      </body>
    </html>
  );
}

export default Post;
