"use client";

import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts";

import {
  type ChartConfig,
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
} from "@/components/ui/chart";
import { useQuery } from "@tanstack/react-query";
import type { Scan } from "@/lib/types";
import { useMemo } from "react";

// const chartData = [
//   { month: "January", desktop: 186, mobile: 80 },
//   { month: "February", desktop: 305, mobile: 200 },
//   { month: "March", desktop: 237, mobile: 120 },
//   { month: "April", desktop: 73, mobile: 190 },
//   { month: "May", desktop: 209, mobile: 130 },
//   { month: "June", desktop: 214, mobile: 140 },
// ]

const chartConfig = {
  RHP2: {
    label: "RHP2",
    color: "#2563eb",
  },
  RHP3: {
    label: "RHP3",
    color: "#60a5fa",
  },
  RHP4: {
    label: "RHP4",
    color: "#38b245",
  },
} satisfies ChartConfig;

export function Component({
  network,
  host,
  v4,
}: {
  network: string;
  host: string;
  v4: boolean;
}) {
  const data = useQuery({
    queryKey: ["scan", network, host],
    queryFn: async () => {
      // return await sdk.scanHost(network, host.$id);
      return (await fetch(
        `/api/v1/scan/?network=${network === "mainnet" ? "main" : "zen"}&publicKey=${host}`,
      ).then((res) => res.json())) as Scan[];
    },
  });

  const chartDataV4 = useMemo(() => {
    if (!data.data) {
      return [];
    }

    return data.data.reverse().map((scan) => {
      return {
        date: scan.createdAt,
        RHP2: scan.rhp2v4delay,
        RHP3: scan.rhp3v4delay,
        RHP4: scan.rhp4v4delay,
      };
    });
  }, [data.data]);

  const chartDataV6 = useMemo(() => {
    if (!data.data) {
      console.log("no data");
      return [];
    }

    return data.data.reverse().map((scan) => {
      return {
        date: scan.createdAt,
        RHP2: scan.rhp2v6delay,
        RHP3: scan.rhp3v6delay,
        RHP4: scan.rhp4v6delay,
      };
    });
  }, [data.data]);

  return (
    <ChartContainer
      config={chartConfig}
      className="max-h-[350px] min-h-[200px] w-full"
    >
      <BarChart accessibilityLayer data={v4 ? chartDataV4 : chartDataV6}>
        <CartesianGrid vertical={false} />
        <XAxis
          dataKey="date"
          tickLine={false}
          tickMargin={10}
          axisLine={false}
          tickFormatter={(value) =>
            new Date(value as string).toLocaleString("en-US", {
              month: "short",
              day: "numeric",
              hour: "numeric",
              minute: "numeric",
            })
          }
        />
        <YAxis tickLine={false} tickMargin={10} axisLine={false} />
        <Bar dataKey="RHP2" fill="var(--color-RHP2)" radius={4} />
        <Bar dataKey="RHP3" fill="var(--color-RHP3)" radius={4} />
        <Bar dataKey="RHP4" fill="var(--color-RHP4)" radius={4} />
        <ChartLegend content={<ChartLegendContent />} />
      </BarChart>
    </ChartContainer>
  );
}
