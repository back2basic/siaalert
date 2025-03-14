import React from "react";
import { Component as Chart } from "./chart";
import Settings from "./settings";

const Page = async ({
  params,
}: {
  params: Promise<{ hostNetwork: string; publicKey: string }>;
}) => {
  const publicKey = decodeURIComponent((await params).publicKey);
  const network = decodeURIComponent((await params).hostNetwork);
  return (
    <div className="flex w-full flex-wrap gap-2 p-6">
      <div className="flex flex-1 flex-col gap-10">
        <Settings network={network} publicKey={publicKey} />
      </div>
      <div className="flex flex-1 flex-col gap-12">
        <div className="flex flex-col items-center text-2xl">
          <div>IPv4 Delay (ms)</div>
          <Chart host={publicKey} network={network} v4={true} />
        </div>
        <div className="flex flex-col items-center text-2xl">
          <div>IPv6 Delay (ms)</div>
          <Chart host={publicKey} network={network} v4={false} />
        </div>
      </div>
    </div>
  );
};

export default Page;
