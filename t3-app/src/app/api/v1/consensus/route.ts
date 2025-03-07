import { env } from "@/env";
import type { Consensus } from "@/lib/types";
import { type NextRequest } from "next/server";

export async function GET(request: NextRequest) {
  // Do whatever you want
  const searchParams = request.nextUrl.searchParams;
  const network = searchParams.get("network");

  // console.log(network, search, online);

  switch (network) {
    case "main":
      const resMain = await fetch(
        `${env.NEXT_PUBLIC_NETWORK_MAIN_URL}/v1/consensus`,
      );
      if (!resMain.ok) {
        return Response.json(
          {},
          {
            status: 500,
          },
        );
      }
      const main = (await resMain.json()) as Consensus;
      // console.log(main);
      return Response.json(main);
    case "zen":
      const resZen = await fetch(
        `${env.NEXT_PUBLIC_NETWORK_ZEN_URL}/v1/consensus`,
      );
      if (!resZen.ok) {
        return Response.json(
          {},
          {
            status: 500,
          },
        );
      }
      const zen = (await resZen.json()) as Consensus;
      console.log(zen);
      return Response.json(zen);
    default:
      return Response.json(
        {},
        {
          status: 500,
        },
      );
  }
}
