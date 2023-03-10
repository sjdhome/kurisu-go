import React from "react";
import { title } from "../constants.js";
import MyProfile from "../lib/MyProfile.js";

function Body() {
  return (
    <>
      <header>
        <MyProfile className="limit-content-width" />
      </header>
      <hr />
      <main></main>
      <footer></footer>
    </>
  );
}

function Index() {
  return (
    <html>
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <title>{title}</title>
        <link rel="stylesheet" href="css/index.css" />
        <link rel="stylesheet" href="css/global.css" />
      </head>
      <body>
        <Body />
      </body>
    </html>
  );
}

export default Index;
