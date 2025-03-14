"use client";
import type { AuthOutput } from "@/lib/types";
// import sdk from "@/lib/sdk";
import { ArrowLeft } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import React, { Suspense, useEffect, useState } from "react";

const Auth = () => {
  const router = useRouter();
  const [output, setOutput] = useState<AuthOutput | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<{ message: string } | null>(null);
  const query = useSearchParams();
  const secret = query.get("otp");
  const publicKey = query.get("publicKey");
  const email = query.get("email");
  const network = query.get("network");

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setOutput(null);
      try {
        // const data = await sdk.verifyOtp(
        //   network as Network,
        //   secret!,
        //   publicKey!,
        //   email!,
        //   expire!,
        // );
        const data = await fetch(
          `/api/auth/otp?publicKey=${publicKey}&email=${email}&network=${network}&secret=${secret}`,
          {
            method: "PUT",
          },
        );
        if (data.ok) {
          const resp = (await data.json()) as AuthOutput;
          setOutput(resp);
          // console.log(data.text());
        } else {
          setError({ message: "Invalid OTP" });
        }
      } catch (error) {
        console.log(error);
        setError(error as { message: string });
      }
      setLoading(false);
    };
    if (secret && publicKey && email) {
      void fetchData();
    }
  }, [network, email, publicKey, secret]);

  return (
    <div className="flex min-h-screen flex-col items-center gap-12 text-emerald-500">
      <div className="flex w-full items-center justify-between gap-2">
        <div className="flex gap-2">
          <ArrowLeft
            className="h12 w-12 cursor-pointer text-red-500"
            onClick={() => router.push("/")}
          />
        </div>
        <div className="text-4xl">Change Subscription</div>
        <div />
      </div>

      {loading && (
        <p className="animate-pulse text-2xl text-green-800">
          Please wait... verifying data.
        </p>
      )}
      {!loading && output?.address && (
        <div>
          <div className="text-lg text-green-500">
            Updated {output.address} - {output.publicKey}
          </div>
          <div
            className={`${output.message === "enabled" ? "text-green-500" : "text-red-500"}`}
          >
            Notifications for {output.email} are {output.message}.
          </div>
        </div>
      )}
      {!loading && error && <div className="text-red-500">{error.message}</div>}
    </div>
  );
};

const Page = () => {
  return (
    <Suspense>
      <Auth />
    </Suspense>
  );
};

export default Page;
