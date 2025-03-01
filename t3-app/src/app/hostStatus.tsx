"use client";
import sdk from "@/lib/sdk";
import type { HostDoc, Network } from "@/lib/types";
import React, { useState } from "react";
import {
  Globe,
  HeartPulse,
  OctagonMinus,
  Shell,
  Mail,
  Handshake,
  Shield,
  ReceiptText,
  // Vote,
} from "lucide-react";
import TooltipWrapper from "@/components/TooltipWrapper";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { useQuery } from "@tanstack/react-query";
// import { UseVersionStore } from "@/lib/store";

export const Scan = ({
  network,
  host,
}: {
  network: Network;
  host: HostDoc;
}) => {
  // const { version } = UseVersionStore();
  const data = useQuery({
    queryKey: ["scan", network, host.publicKey],
    queryFn: async () => {
      return await sdk.scanHost(network, host.$id);
    },
  });

  if (data.isLoading) {
    return (
      <div className="flex animate-pulse items-center text-xs">
        <div className="flex items-center gap-2">
          <Shield className="h-4 w-4" />
          Scanning...
        </div>
      </div>
    );
  }
  if (data.isError) {
    return (
      <div className="flex items-center text-xs">
        <div className="flex items-center gap-2">
          <Shield className="h-4 w-4" />
          {data.error?.message}
        </div>
      </div>
    );
  }
  if (data.data?.documents.length === 0) {
    return (
      <div className="flex items-center text-xs">
        <div className="flex items-center gap-2">
          <Shield className="h-4 w-4" />
          No scan data
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-1 items-center justify-start gap-4 rounded-full border-b border-green-500 px-2 py-1 text-xs shadow">
      {/* Accepting Contracts */}
      <div className="flex items-center gap-2 pl-2">
        <div>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger className="flex cursor-default items-center">
                <ReceiptText
                  className={`h-4 w-4 ${data.data?.documents[0]?.acceptingContracts ? "text-green-500" : "text-red-500"}`}
                />
              </TooltipTrigger>
              <TooltipContent>
                <p>
                  {data.data?.documents[0]?.acceptingContracts
                    ? "Accepting Contracts"
                    : "Not Accepting Contracts"}
                </p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
        {/* Release */}
        {/* <div>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger className="flex cursor-default items-center">
                <Vote
                  className={`h-4 w-4 ${data.data?.documents[0]?.release === "hostd " +version ? "text-green-500" : "text-red-500"}`}
                />
              </TooltipTrigger>
              <TooltipContent>
                <p>
                  {data.data?.documents[0]?.release === version
                    ? "OK"
                    : "Current: " + data.data?.documents[0]?.release + ", Latest: " + version }
                </p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div> */}
      </div>
      {/* IP */}
      <div>
        {/* V4 */}
        <div
          className={
            data.data?.documents[0]?.v4Addr !== ""
              ? "w-36 truncate text-nowrap text-green-500"
              : "w-36 truncate text-nowrap text-red-500"
          }
        >
          IPv4 {data.data?.documents[0]?.v4Addr}
        </div>
        {/* V6 */}
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger
              className={
                data.data?.documents[0]?.v6Addr !== ""
                  ? "w-36 truncate text-nowrap text-green-500"
                  : "w-36 truncate text-nowrap text-left text-red-500"
              }
            >
              IPv6 {data.data?.documents[0]?.v6Addr}
            </TooltipTrigger>
            <TooltipContent>
              <p>{data.data?.documents[0]?.v6Addr}</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </div>
      <div className="flex flex-1 items-center justify-between gap-2">
        {/* RHP2 */}
        <div>
          <div
            className={`flex items-center gap-2 ${data.data?.documents[0]?.rhp2V4Delay && data.data?.documents[0]?.rhp2V4Delay > 250 ? "text-orange-500" : ""}`}
          >
            <Shield
              className={`h-4 w-4 ${data.data?.documents[0]?.rhp2V4 ? "text-green-500" : "text-red-500"}`}
            />
            <div className="flex items-center gap-2">
              <div>{data.data?.documents[0]?.rhp2Port}</div>
              {data.data?.documents[0]?.rhp2V4Delay !== 0 &&
                (data.data?.documents[0]?.rhp2V4Delay ?? 0) < 3000 && (
                  <div
                    className={
                      !data.data?.documents[0]?.rhp2V4 ? "hidden" : "flex"
                    }
                  >
                    {data.data?.documents[0]?.rhp2V4Delay}ms
                  </div>
                )}
            </div>
          </div>
          <div
            className={`flex items-center gap-2 ${data.data?.documents[0]?.rhp2V6Delay && data.data?.documents[0]?.rhp2V6Delay > 250 ? "text-orange-500" : ""}`}
          >
            <Shield
              className={`h-4 w-4 ${data.data?.documents[0]?.rhp2V6 ? "text-green-500" : "text-red-500"}`}
            />
            <div className="flex items-center gap-2">
              <div>{data.data?.documents[0]?.rhp2Port}</div>
              {data.data?.documents[0]?.rhp2V6Delay !== 0 &&
                (data.data?.documents[0]?.rhp2V6Delay ?? 0) < 3000 && (
                  <div
                    className={
                      !data.data?.documents[0]?.rhp2V6 ? "hidden" : "flex"
                    }
                  >
                    {data.data?.documents[0]?.rhp2V6Delay}ms
                  </div>
                )}
            </div>
          </div>
        </div>
        {/* RHP3 */}
        <div>
          <div
            className={`flex items-center gap-2 ${data.data?.documents[0]?.rhp3V4Delay && data.data?.documents[0]?.rhp3V4Delay > 250 ? "text-orange-500" : ""}`}
          >
            <Shield
              className={`h-4 w-4 ${data.data?.documents[0]?.rhp3V4 ? "text-green-500" : "text-red-500"}`}
            />
            <div className="flex items-center gap-2">
              <div>{data.data?.documents[0]?.rhp3Port}</div>
              {data.data?.documents[0]?.rhp3V4Delay !== 0 &&
                (data.data?.documents[0]?.rhp3V4Delay ?? 0) < 3000 && (
                  <div
                    className={
                      !data.data?.documents[0]?.rhp3V4 ? "hidden" : "flex"
                    }
                  >
                    {data.data?.documents[0]?.rhp3V4Delay}ms
                  </div>
                )}
            </div>
          </div>
          <div
            className={`flex items-center gap-2 ${data.data?.documents[0]?.rhp3V6Delay && data.data?.documents[0]?.rhp3V6Delay > 250 ? "text-orange-500" : ""}`}
          >
            <Shield
              className={`h-4 w-4 ${data.data?.documents[0]?.rhp3V6 ? "text-green-500" : "text-red-500"}`}
            />
            <div className="flex items-center gap-2">
              <div>{data.data?.documents[0]?.rhp3Port}</div>
              {data.data?.documents[0]?.rhp3V6Delay !== 0 &&
                (data.data?.documents[0]?.rhp3V6Delay ?? 0) < 3000 && (
                  <div
                    className={
                      !data.data?.documents[0]?.rhp2V6 ? "hidden" : "flex"
                    }
                  >
                    {data.data?.documents[0]?.rhp3V6Delay}ms
                  </div>
                )}
            </div>
          </div>
        </div>
        {/* RHP4 */}
        <div />
        {/* <div>
          <div
            className={`flex items-center gap-2 ${data.data?.documents[0]?.rhp4V4Delay && data.data?.documents[0]?.rhp4V4Delay > 250 ? "text-orange-500" : ""}`}
          >
            <Shield
              className={`h-4 w-4 ${data.data?.documents[0]?.rhp4V4 ? "text-green-500" : "text-red-500"}`}
            />
            <div>
              {data.data?.documents[0]?.rhp4Port}-
              {data.data?.documents[0]?.rhp4V4Delay}ms
            </div>
          </div>
          <div
            className={`flex items-center gap-2 ${data.data?.documents[0]?.rhp4V6Delay && data.data?.documents[0]?.rhp4V6Delay > 250 ? "text-orange-500" : ""}`}
          >
            <Shield
              className={`h-4 w-4 ${data.data?.documents[0]?.rhp4V6 ? "text-green-500" : "text-red-500"}`}
            />
            <div>
              {data.data?.documents[0]?.rhp4Port}-
              {data.data?.documents[0]?.rhp4V6Delay}ms
            </div>
          </div>
        </div> */}
      </div>
    </div>
  );
};

export const Email = ({
  network,
  host,
}: {
  network: Network;
  host: HostDoc;
}) => {
  const [open, setOpen] = useState(false);
  const [email, setEmail] = useState("");
  const [status, setStatus] = useState("");
  const [disabled, setDisabled] = useState(false);

  async function handleEmail() {
    setStatus("processing request");
    setDisabled(true);
    try {
      const response = await sdk.sendHostEmail(network, host.publicKey, email);
      // console.log(response);
      if (response) {
        setStatus("success");
        setOpen(false);
      } else {
        setStatus("Something went wrong, please try again");
      }
    } catch (error) {
      console.error(error);
      setStatus("Something went wrong, please try again");
    }
    setDisabled(false);
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <div className="flex items-center justify-center rounded-full border border-red-500 px-2">
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger>
                <Mail />
              </TooltipTrigger>
              <TooltipContent>
                <p>Configure Email Alerts</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </div>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[1425px]">
        <DialogHeader>
          <DialogTitle className="flex items-center justify-center gap-2">
            {host.$id} - {host.netAddress} - {host.publicKey}
          </DialogTitle>
          <DialogDescription className="flex flex-col items-center gap-2">
            Leave your email address to (un)subscribe for notifications.
            <Input
              value={email}
              type="email"
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Email"
              className="w-96"
            />
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <div className="flex w-full flex-col items-center justify-evenly gap-2 text-center">
            <Button onClick={handleEmail} disabled={disabled}>
              Send
            </Button>
            <div>{status}</div>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export const RenderError = ({ error }: { error: string }) => {
  // Split the error message into parts based on ":"
  const parts = error.split(":");
  return (
    <div className="flex items-center gap-1 whitespace-nowrap px-2">
      {parts.map((part: string, index: number) => {
        if (part.includes("failed to get host settings")) {
          return (
            <span key={index} className="flex items-center gap-1">
              <TooltipWrapper content={part}>
                <Globe />
                {part.split("failed to get host settings")[1]}
              </TooltipWrapper>
            </span>
          );
        } else if (part.includes("failed to connect to host")) {
          return (
            <span key={index} className="flex items-center gap-1">
              <TooltipWrapper content={part}>
                <Shell />
              </TooltipWrapper>
              {part.split("failed to connect to host")[1]}
            </span>
          );
        } else if (part.includes("failed to establish v2 transport")) {
          return (
            <span key={index} className="flex items-center gap-1">
              <TooltipWrapper content={part}>
                <OctagonMinus />
              </TooltipWrapper>
              {part.split("failed to establish v2 transport")[1]}
            </span>
          );
        } else if (part.includes("failed to parse net address")) {
          return (
            <span key={index} className="flex items-center gap-1">
              <TooltipWrapper content={part}>
                <HeartPulse />
              </TooltipWrapper>
              {part.split("failed to parse net address")[1]}
            </span>
          );
        } else if (part.includes("handshake signature was invalid")) {
          return (
            <span key={index} className="flex items-center gap-1">
              <TooltipWrapper content={part}>
                <Handshake />
              </TooltipWrapper>
              {part.split("handshake signature was invalid")[1]}
            </span>
          );
          // } else if (part.includes('Post "http')) {
          //   return <span key={index}>{part.split('Post "http')[1]}</span>;
          // } else if (part.includes("//bench")) {
          //   return <span key={index}>{part.split("//bench")[1]}</span>;
          // } else if (part.includes('8484/scan"')) {
          //   return <span key={index}>{part.split('8484/scan"')[1]}</span>;
          // } else if (part.includes("dial tcp")) {
          //   return <span key={index}>{part.split("dial tcp")[1]}</span>;
        } else {
          // return <span key={index}>{part}</span>;
          return;
        }
      })}
    </div>
  );
};

export const HostStatus = ({ host }: { host: HostDoc }) => {
  return (
    <div className="flex items-center justify-between gap-2 text-nowrap">
      <div className="flex items-center gap-2">
        <span
          className={`text-lg font-bold ${host.online ? "text-green-500" : "text-red-500"}`}
        >
          {host.online ? "Online" : "Offline"}
        </span>
        <span className="-mr-3 text-xs">
          {host.onlineSince && new Date(host.onlineSince).toLocaleString()}
        </span>
        <span className="text-xs">
          {host.offlineSince && new Date(host.offlineSince).toLocaleString()}
        </span>
        {/* <RenderError error={host.error ?? ""} /> */}
      </div>

      {/* <Email network={network} host={host} /> */}
    </div>
  );
};
