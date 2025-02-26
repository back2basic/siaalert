"use client";

import sdk from "@/lib/sdk";
import { SiGithub } from "@icons-pack/react-simple-icons";
import { useQuery } from "@tanstack/react-query";
import { Donut } from "lucide-react";
import React from "react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";

const Footer = () => {
  const data = useQuery({
    queryKey: ["serverless"],
    queryFn: async () => {
      return await sdk.checkServerless();
    },
    refetchInterval: 60 * 1000,
  });

  return (
    <footer className="flex items-center justify-center gap-6 pb-2 pr-6">
      <div>Sia Host Alert © {new Date().getFullYear()} ©</div>
      <a href="https://github.com/back2basic/siaalert" target="_blank">
        <SiGithub className="h-6 w-6" />
      </a>
      <div className="flex items-center gap-2 text-xs">
        serverless:
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger>
              <Donut
                className={`h-4 w-4 ${data.isSuccess ? "text-green-500" : "text-red-500"} ${data.isLoading ? "animate-spin" : ""}`}
              />
            </TooltipTrigger>
            <TooltipContent>
              <p>{data.isSuccess ? "Online" : "Offline"}</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </div>
    </footer>
  );
};

export default Footer;
