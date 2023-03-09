import React from "react";

function Post(props: { postHTML: string }) {
  const postHTML = props.postHTML;
  return (
    <html>
      <head>
        <meta charSet="utf-8" />
        <title>Post</title>
      </head>
      <body>
        <main dangerouslySetInnerHTML={{ __html: postHTML }} />
      </body>
    </html>
  );
}

export default Post;
