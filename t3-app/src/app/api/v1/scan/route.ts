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
      try {
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
        if (main.length === 0) {
          return Response.json([], {
            status: 500,
          });
        }
        // sort desc on cratedAt
        main.sort((a, b) => {
          return (
            new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
          );
        });
        return Response.json(main);
      } catch {
        return Response.json([], {
          status: 500,
        });
      }
    case "zen":
      try {
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
        if (zen.length === 0) {
          return Response.json([], {
            status: 500,
          });
        }
        return Response.json(zen);
      } catch {
        return Response.json([], {
          status: 500,
        });
      }
    default:
      return Response.json([], {
        status: 500,
      });
  }
}
