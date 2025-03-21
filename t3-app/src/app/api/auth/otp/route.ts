import { env } from "@/env";
import type { AuthOutput } from "@/lib/types";
import { type NextRequest } from "next/server";

export async function POST(request: NextRequest) {
  // Do whatever you want

  const searchParams = request.nextUrl.searchParams;
  const network = searchParams.get("network");
  const publicKey = searchParams.get("publicKey");
  const email = searchParams.get("email");

  switch (network) {
    case "mainnet":
      try {
        const resMain = await fetch(
          `${env.NEXT_PUBLIC_NETWORK_MAIN_URL}/auth/otp?publicKey=${publicKey}&email=${email}`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            // body: JSON.stringify({ secret }),
          },
        );
        if (!resMain.ok) {
          return Response.json([], {
            status: 500,
          });
        }
        const main = (await resMain.json()) as { message: string };
        // console.log(main);
        return Response.json(main);
      } catch {
        return Response.json(
          {},
          {
            status: 500,
          },
        );
      }
    case "zen":
      try {
        const resZen = await fetch(
          `${env.NEXT_PUBLIC_NETWORK_ZEN_URL}/auth/otp?publicKey=${publicKey}&email=${email}`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            // body: JSON.stringify({ secret }),
          },
        );
        console.log(resZen.status);
        if (!resZen.ok) {
          return Response.json([], {
            status: 500,
          });
        }
        const zen = (await resZen.json()) as { message: string };
        // console.log(zen);
        return Response.json(zen);
      } catch {
        return Response.json(
          {},
          {
            status: 500,
          },
        );
      }
    default:
      return Response.json([], {
        status: 500,
      });
  }
}

export async function PUT(request: NextRequest) {
  // Do whatever you want

  const searchParams = request.nextUrl.searchParams;
  const network = searchParams.get("network");
  const publicKey = searchParams.get("publicKey");
  const secret = searchParams.get("secret");
  const email = searchParams.get("email");

  // console.log(network, publicKey, secret, email);
  switch (network) {
    case "main":
      try {
        const resMain = await fetch(
          `${env.NEXT_PUBLIC_NETWORK_MAIN_URL}/auth/otp?publicKey=${publicKey}&email=${email}&secret=${secret}`,
          {
            method: "PUT",
            headers: {
              "Content-Type": "application/json",
            },
            // body: JSON.stringify({ secret }),
          },
        );
        if (!resMain.ok) {
          return Response.json([], {
            status: 500,
          });
        }
        const main = (await resMain.json()) as AuthOutput;
        // console.log(main);
        return Response.json(main);
      } catch {
        return Response.json(
          {},
          {
            status: 500,
          },
        );
      }
    case "zen":
      try {
        const resZen = await fetch(
          `${env.NEXT_PUBLIC_NETWORK_ZEN_URL}/auth/otp?publicKey=${publicKey}&email=${email}&secret=${secret}`,
          {
            method: "PUT",
            headers: {
              "Content-Type": "application/json",
            },
            // body: JSON.stringify({ secret }),
          },
        );
        if (!resZen.ok) {
          return Response.json([], {
            status: 500,
          });
        }
        const zen = (await resZen.json()) as AuthOutput;
        // console.log(zen);
        return Response.json(zen);
      } catch {
        return Response.json(
          {},
          {
            status: 500,
          },
        );
      }
    default:
      return Response.json([], {
        status: 500,
      });
  }
}
