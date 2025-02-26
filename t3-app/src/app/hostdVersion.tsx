"use client";
import React, { useEffect } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { useQuery } from "@tanstack/react-query";
import type { GitHiubHostdRelease } from "@/lib/types";
import { Computer } from "lucide-react";
import { Textarea } from "@/components/ui/textarea";
import { UseVersionStore } from "@/lib/store";

const LatestRelease: React.FC = () => {
  const { setVersion } = UseVersionStore();
  const data = useQuery({
    queryKey: ["release"],
    queryFn: async () => {
      return (await fetch(
        "https://api.github.com/repos/SiaFoundation/hostd/releases/latest",
      ).then((res) => res.json())) as GitHiubHostdRelease;
    },
  });

  useEffect(() => {
    if (!data.data) return;
    setVersion(data.data?.tag_name);
  }, [data.data, setVersion]);

  if (data.isError) {
    return <div>Error: {data.error.message}</div>;
  }

  if (data.isLoading) {
    return (
      <div>
        <Computer className="animate-spin" />
      </div>
    );
  }

  return (
    <div>
      <Dialog>
        <DialogTrigger asChild>
          <div className="flex cursor-pointer items-end gap-2 rounded-full p-3 shadow-md shadow-green-500 hover:bg-green-500 hover:text-purple-500 hover:shadow-lg hover:shadow-slate-900">
            <div>Latest Hostd:</div>
            <div>{data.data?.tag_name}</div>
          </div>
        </DialogTrigger>
        <DialogContent className="bg-black sm:max-w-[825px]">
          <DialogHeader>
            <DialogTitle>Hostd {data.data?.tag_name}</DialogTitle>
            <DialogDescription className="p-3 text-lg">
              <Textarea
                readOnly
                value={data.data?.body}
                rows={15}
                className="border-0"
              />
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button>
              <a
                href={data.data?.html_url}
                target="_blank"
                rel="noopener noreferrer"
              >
                View on GitHub
              </a>
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default LatestRelease;
