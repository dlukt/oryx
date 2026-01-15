//
// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
//
import React from "react";

export default function IconButton({onClick, title, children}) {
  return (
    <button
      type="button"
      className="btn btn-link p-0 text-decoration-none text-reset align-baseline border-0 bg-transparent"
      style={{display: 'inline-block'}}
      title={title}
      aria-label={title}
      onClick={onClick}
    >
      {children}
    </button>
  );
}
