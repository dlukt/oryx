//
// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT
//
import React from "react";

// A reusable icon button component that replaces non-semantic div buttons.
// Ensures keyboard accessibility and proper semantics.
export default function IconButton({onClick, title, children, className}) {
  return (
    <button
      type="button"
      className={`btn btn-link p-0 text-decoration-none text-reset align-baseline border-0 bg-transparent ${className || ''}`}
      style={{display: 'inline-block'}}
      title={title}
      aria-label={title}
      onClick={onClick}
    >
      {children}
    </button>
  );
}
