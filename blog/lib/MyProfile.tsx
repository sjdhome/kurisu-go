import React from "react";
import { motto } from "../constants.js";

function MyProfile(props: { className?: string }) {
  return (
    <div
      className={`${props.className} flex flex-column`}
      style={{
        paddingTop: "5rem",
        paddingBottom: "3rem",
        marginLeft: "1rem",
        marginRight: "1rem",
      }}
    >
      <div className="flex flex-align-center">
        <a href="/">
          <img
            src="img/earth.jpg"
            alt="avatar"
            width={64}
            height={64}
            className="round"
            style={{ marginRight: "1rem" }}
          />
        </a>
        <a href="/" className="black-font no-a-decoration">
          <h1 style={{ fontWeight: "500" }}>sjdhome blog</h1>
        </a>
      </div>
      <div style={{ marginLeft: "0.25rem", marginTop: "0.75rem" }}>
        <p className="no-margin small-font">{motto}</p>
      </div>
    </div>
  );
}

export default MyProfile;
