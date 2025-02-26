"use client";
import { Switch } from "@/components/ui/switch";
import sdk from "@/lib/sdk";
import { useQuery } from "@tanstack/react-query";
import React, { useMemo, useState } from "react";
import { Email, HostStatus, RenderError, Scan } from "./hostStatus";
import { Input } from "@/components/ui/input";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  useReactTable,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  flexRender,
  type PaginationState,
} from "@tanstack/react-table";
import type { HostDoc, Network } from "@/lib/types";

const Hosts = () => {
  const [online, setOnline] = useState(true);
  const [search, setSearch] = useState("");
  const [network, setNetwork] = useState<Network>("mainnet");
  const [pagination, setPagination] = React.useState<PaginationState>({
    pageIndex: 0,
    pageSize: 15,
  });

  const consensusData = useQuery({
    queryKey: ["consensus", network],
    queryFn: async () => {
      return await sdk.getStatus(network);
    },
  });

  const data = useQuery({
    queryKey: ["hosts", network, online, search],
    queryFn: async () => {
      return await sdk.getHosts(network, search, online);
    },
    refetchInterval: 5 * 60 * 1000, // 5 min
  });

  const columns = useMemo(
    () => [
      {
        header: "NetAddress",
        accessorKey: "netAddress",
        cell: ({ getValue }: { getValue: () => string }) => (
          // <p className="max-w-96 text-lg font-bold truncate">{getValue()}</p>
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger className="flex items-center px-2">
                <p
                  className="max-w-64 truncate text-lg font-bold"
                  onClick={() => navigator.clipboard.writeText(getValue())}
                >
                  {getValue()}
                </p>
              </TooltipTrigger>
              <TooltipContent>
                <p>{getValue()}</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        ),
      },
      {
        header: "PublicKey",
        accessorKey: "publicKey",
        cell: ({ getValue }: { getValue: () => string }) => (
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger className="flex items-center px-2">
                <p
                  className="max-w-72 cursor-default truncate"
                  onClick={() => navigator.clipboard.writeText(getValue())}
                >
                  {getValue()}
                </p>
              </TooltipTrigger>
              <TooltipContent>
                <p>{getValue()}</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        ),
      },
      {
        header: "Status",
        accessorKey: "status",
        cell: ({ row }: { row: { original: HostDoc } }) => (
          <div className="flex items-center gap-6">
            <HostStatus host={row.original} />
            <Scan network={network} host={row.original} />
          </div>
        ),
      },
      {
        header: "",
        accessorKey: "error",
        cell: ({ row }: { row: { original: HostDoc } }) => (
          <RenderError error={row.original.error ?? ""} />
        ),
      },
      {
        header: "Email",
        accessorKey: "email",
        cell: ({ row }: { row: { original: HostDoc } }) => (
          <Email network={network} host={row.original} />
        ),
      },
    ],
    [network],
  );

  const table = useReactTable({
    data:
      data.data?.documents.sort((a, b) =>
        a.netAddress.localeCompare(b.netAddress),
      ) ?? [],
    columns,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    onPaginationChange: setPagination,
    state: {
      pagination,
    },
  });

  return (
    <div className="flex h-full w-full flex-col gap-6 p-2">
      <div className="flex items-center justify-around gap-6">
        <Input
          type="text"
          placeholder="Search NetAddress"
          className={`w-96 rounded-full ${network === "zen" ? "border-orange-500" : "border-green-500"}`}
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <div className="flex items-center gap-2">
          <div>
            Total Hosts {online ? " Online" : " Found"}: {data.data?.total ?? 0}
          </div>
          <Switch
            className={
              network === "zen" ? "border-orange-500" : "border-green-500"
            }
            checked={online}
            onCheckedChange={() => setOnline(!online)}
          />
        </div>
        <div className="flex flex-col">
          <div className="flex items-center gap-4">
            <div className="flex w-full items-center justify-between gap-2 text-2xl font-bold">
              <div>Network:</div>
              <div className={network === "zen" ? "text-orange-500" : ""}>
                {network.toLocaleUpperCase()}
              </div>
            </div>
            <Switch
              // disabled
              className={
                network === "zen" ? "border-orange-500" : "border-green-500"
              }
              checked={network === "mainnet"}
              onCheckedChange={() =>
                setNetwork(network === "mainnet" ? "zen" : "mainnet")
              }
            />
          </div>
          <div className="text-xs">
            Consensus height: {consensusData.data?.height} -{" "}
            {new Date(consensusData.data?.$updatedAt ?? "").toLocaleString()}
          </div>
        </div>
      </div>
      <div>
        {data.isLoading && <p className="animate-pulse">Loading Hosts...</p>}
        {data.isError && <p>Error: {data.error.message}</p>}
        {data.isSuccess && !data.data && (
          <p>Searching for all hosts is possible only via NetAddress.</p>
        )}
        {data.isSuccess && (data.data?.total ?? 0) > 0 && (
          <div>
            <table className="w-full">
              <thead>
                {table.getHeaderGroups().map((headerGroup) => (
                  <tr key={headerGroup.id}>
                    {headerGroup.headers.map((header) => (
                      <th key={header.id}>
                        {flexRender(
                          header.column.columnDef.header,
                          header.getContext(),
                        )}
                      </th>
                    ))}
                  </tr>
                ))}
              </thead>
              <tbody>
                {table.getRowModel().rows.map((row) => (
                  <tr key={row.id} className="text-sm">
                    {row.getVisibleCells().map((cell) => (
                      <td key={cell.id}>
                        {flexRender(
                          cell.column.columnDef.cell,
                          cell.getContext(),
                        )}
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>

            {/* Pagination */}
            <div className="h-6" />
            <div className="flex w-full items-end justify-center gap-2">
              <button
                className="rounded-full border border-red-500 px-2"
                onClick={() => table.firstPage()}
                disabled={!table.getCanPreviousPage()}
              >
                {"<<"}
              </button>
              <button
                className="rounded-full border border-red-500 px-2"
                onClick={() => table.previousPage()}
                disabled={!table.getCanPreviousPage()}
              >
                {"<"}
              </button>
              <button
                className="rounded-full border border-red-500 px-2"
                onClick={() => table.nextPage()}
                disabled={!table.getCanNextPage()}
              >
                {">"}
              </button>
              <button
                className="rounded-full border border-red-500 px-2"
                onClick={() => table.lastPage()}
                disabled={!table.getCanNextPage()}
              >
                {">>"}
              </button>
              <span className="flex items-center gap-1">
                <div>Page</div>
                <strong>
                  {table.getState().pagination.pageIndex + 1} of{" "}
                  {table.getPageCount().toLocaleString()}
                </strong>
              </span>
              <span className="flex items-center gap-1">
                | Go to page:
                <input
                  type="number"
                  min="1"
                  max={table.getPageCount()}
                  defaultValue={table.getState().pagination.pageIndex + 1}
                  onChange={(e) => {
                    const page = e.target.value
                      ? Number(e.target.value) - 1
                      : 0;
                    table.setPageIndex(page);
                  }}
                  className="w-16 rounded-full border border-red-500 bg-transparent px-2 text-center"
                />
              </span>
              <select
                value={table.getState().pagination.pageSize}
                onChange={(e) => {
                  table.setPageSize(Number(e.target.value));
                }}
                className="rounded-full border border-red-500 bg-transparent px-2 py-0.5"
              >
                {[15, 25, 50, 75, 150].map((pageSize) => (
                  <option key={pageSize} value={pageSize}>
                    Show {pageSize}
                  </option>
                ))}
              </select>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Hosts;
