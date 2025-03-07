import { env } from "@/env";
import type { Scan } from "@/lib/types";
import { type NextRequest } from "next/server";

export async function GET(request: NextRequest) {
  // Do whatever you want
  const searchParams = request.nextUrl.searchParams;
  const network = searchParams.get("network");
  const publicKey = searchParams.get("publicKey");

  // console.log(network, search, online);

  switch (network) {
    case "main":
      const resMain = await fetch(
        `${env.NEXT_PUBLIC_NETWORK_MAIN_URL}/v1/scan?publicKey=${publicKey}`,
      );
      if (!resMain.ok) {
        return Response.json([], {
          status: 500,
        });
      }
      const main = (await resMain.json()) as Scan[];
      // console.log(main);
      return Response.json(main[0]);
    case "zen":
      const resZen = await fetch(
        `${env.NEXT_PUBLIC_NETWORK_ZEN_URL}/v1/scan?publicKey=${publicKey}`,
      );
      if (!resZen.ok) {
        return Response.json([], {
          status: 500,
        });
      }
      const zen = (await resZen.json()) as Scan[];
      // console.log(zen);
      if (zen.length === 1) {
        return Response.json(zen[0]);
      }
    default:
      return Response.json([], {
        status: 500,
      });
  }
}
