"use client";
import type {
  HostScan,
  Rhpv2Settings,
  Rhpv3Settings,
  Rhpv4Settings,
} from "@/lib/types";
import {
  convertPrice,
  convertPricePerBlock,
  convertSectorsToBytes,
  formatStorage,
  toSia,
} from "@/lib/utils";
import { useQuery } from "@tanstack/react-query";
import React from "react";

type Props = { network: string; publicKey: string };

const V1Settings = ({
  netAddress,
  settings,
  pricetable,
  publicKey,
}: {
  netAddress: string;
  settings: Rhpv2Settings;
  pricetable: Rhpv3Settings;
  publicKey: string;
}) => {
  return (
    <>
      <div className="mb-4 flex w-full justify-center border-b border-green-500 text-2xl font-bold">
        {netAddress}
      </div>
      <div
        className="flex w-full cursor-pointer justify-center"
        onClick={() => navigator.clipboard.writeText(publicKey)}
      >
        {publicKey}
      </div>
      <div className="flex max-w-96 flex-col gap-2 p-10">
        <div className="flex justify-between gap-2">
          <div>V2:</div>
          <div>No</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Accepting Contracts:</div>
          <div>{settings.acceptingcontracts ? "Yes" : "No"}</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Max Collateral:</div>
          <div>{convertPrice(settings.maxcollateral)}</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Max Contract Duration:</div>
          <div>{settings.maxduration / 144} Days</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Remaining Storage:</div>
          <div>{formatStorage(settings.remainingstorage)}</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Total Storage:</div>
          <div>{formatStorage(settings.totalstorage)}</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Release:</div>
          <div>{settings.release}</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Storage Price:</div>
          <div>{convertPricePerBlock(settings.storageprice)}</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Upload Price:</div>
          <div>{toSia(pricetable.uploadbandwidthcost, false)}</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Download Price:</div>
          <div>{toSia(pricetable.downloadbandwidthcost, false)}</div>
        </div>
      </div>
    </>
  );
};

const V2Settings = ({
  netAddress,
  settings,
  publicKey,
}: {
  netAddress: string;
  settings: Rhpv4Settings;
  publicKey: string;
}) => {
  return (
    <>
      <div className="mb-4 flex w-full justify-center border-b border-green-500 text-2xl font-bold">
        {netAddress}
      </div>
      <div className="flex w-full justify-center">{publicKey}</div>
      <div className="flex max-w-96 flex-col gap-2 p-10">
        <div className="flex justify-between gap-2">
          <div>V2:</div>
          <div>Yes</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Accepting Contracts:</div>
          <div>{settings.acceptingContracts ? "Yes" : "No"}</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Max Collateral:</div>
          <div>{convertPrice(settings.maxCollateral)}</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Max Contract Duration:</div>
          <div>{settings.maxContractDuration / 144} Days</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Remaining Storage:</div>
          <div>
            {formatStorage(convertSectorsToBytes(settings.remainingStorage))}
          </div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Total Storage:</div>
          <div>
            {formatStorage(convertSectorsToBytes(settings.totalStorage))}
          </div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Release:</div>
          <div>{settings.release}</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Storage Price:</div>
          <div>{convertPricePerBlock(settings.prices.storagePrice)}</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Upload Price:</div>
          <div>{toSia(settings.prices.ingressPrice, false)}</div>
        </div>
        <div className="flex justify-between gap-2">
          <div>Download Price:</div>
          <div>{toSia(settings.prices.egressPrice, false)}</div>
        </div>
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
      {/* <h1 className="pb-6 text-2xl font-bold">Settings</h1> */}
      {data.data?.v2 ? (
        <V2Settings
          netAddress={data.data.netAddress}
          settings={data.data.rhpV4Settings}
          publicKey={publicKey}
        />
      ) : (
        <V1Settings
          netAddress={data.data.netAddress}
          settings={data.data?.settings}
          pricetable={data.data?.priceTable}
          publicKey={publicKey}
        />
      )}
    </div>
  );
};

export default Settings;
