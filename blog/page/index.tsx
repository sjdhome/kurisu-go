import React from "react";
import { title } from "../constants.js";

function Body() {
  return (
    <>
      <h1>Hello world</h1>
    </>
  );
}

function Index() {
  return (
    <html>
      <head>
        <meta charSet="utf-8" />
        <title>{title}</title>
      </head>
      <body>
        <Body />
      </body>
    </html>
  );
}

export default Index;
