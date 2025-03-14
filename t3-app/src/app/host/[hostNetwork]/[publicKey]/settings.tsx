"use client";
import type {
  HostScan,
  Rhpv2Settings,
  Rhpv3Settings,
  Rhpv4Settings,
} from "@/lib/types";
import {
  convertPrice,
  convertSectorsToBytes,
  formatStorage,
} from "@/lib/utils";
import { useQuery } from "@tanstack/react-query";
import React from "react";

type Props = { network: string; publicKey: string };

const V1Settings = ({
  settings,
}: {
  settings: Rhpv2Settings;
  pricetable: Rhpv3Settings;
}) => {
  return (
    <>
      <div>V2: No</div>
      <div>
        Accepting Contracts: {settings.acceptingcontracts ? "Yes" : "No"}
      </div>
      <div>Max Collateral: {convertPrice(settings.maxcollateral)}</div>
      <div>Max Contract Duration: {settings.maxduration / 144} Days</div>
      <div>Remaining Storage: {formatStorage(settings.remainingstorage)}</div>
      <div>Total Storage: {formatStorage(settings.totalstorage)}</div>
    </>
  );
};

const V2Settings = ({ settings }: { settings: Rhpv4Settings }) => {
  return (
    <>
      <div>V2: Yes</div>
      <div>
        Accepting Contracts: {settings.acceptingContracts ? "Yes" : "No"}
      </div>
      <div>Max Collateral: {convertPrice(settings.maxCollateral)}</div>
      <div>
        Max Contract Duration: {settings.maxContractDuration / 144} Days
      </div>
      <div>
        Remaining Storage:{" "}
        {formatStorage(convertSectorsToBytes(settings.remainingStorage))}
      </div>
      <div>
        Total Storage:{" "}
        {formatStorage(convertSectorsToBytes(settings.totalStorage))}
      </div>
    </>
  );
};

const Settings = ({ network, publicKey }: Props) => {
  const data = useQuery({
    queryKey: ["hostscan", network, publicKey],
    queryFn: async () => {
      return (await fetch(
        `/api/v1/hostscan/?network=${network === "mainnet" ? "main" : "zen"}&publicKey=${publicKey}`,
      ).then((res) => res.json())) as HostScan;
    },
  });

  if (data.isError) {
    return <div>Error: {data.error.message}</div>;
  }

  if (data.isLoading) {
    return <div className="animate-pulse">Scanning host...</div>;
  }

  if (!data.data?.settings) {
    return <div>Failed scanning host</div>;
  }

  return (
    <div className="flex flex-col gap-2 rounded-lg p-2 shadow-md shadow-green-500">
      <h1 className="pb-6 text-2xl font-bold">Settings</h1>
      {data.data?.v2 ? (
        <V2Settings settings={data.data.rhpV4Settings} />
      ) : (
        <V1Settings
          settings={data.data?.settings}
          pricetable={data.data?.priceTable}
        />
      )}
    </div>
  );
};

export default Settings;
