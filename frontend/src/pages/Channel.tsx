import React from "react";
import { useParams } from "react-router-dom";

export default function Channel() {
  const { channelId } = useParams();
  return (
    <div>
      <h2>Channel: {channelId}</h2>
      <p>Channel info and posts (stub).</p>
    </div>
  );
}
