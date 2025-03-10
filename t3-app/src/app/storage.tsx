"use client";

import { UseStorageStore } from "@/lib/store";
import { formatStorage } from "@/lib/utils";
import React from "react";

const Storage = () => {
  const { v1total, v2total, v1used, v2used , v1hosts, v2hosts} = UseStorageStore();
  return (
    <div className="flex w-72 flex-col text-sm">
      <div className="flex w-full justify-between gap-4 border-b border-green-500">
        <h1></h1>
        <h2>Total</h2>
        <h2>Used</h2>
      </div>
      <div className="flex w-full justify-between gap-4">
        <div>V1-({v1hosts})</div>
        <div>{formatStorage(false, v1total)}</div>
        <div>{formatStorage(false, v1total - v1used)}</div>
      </div>
      <div className="flex w-full justify-between gap-4">
        <div>V2-({v2hosts})</div>
        <div>{formatStorage(true, v2total)}</div>
        <div>{formatStorage(true, v2total - v2used)}</div>
      </div>
    </div>
  );
};

export default Storage;
