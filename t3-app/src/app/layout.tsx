import "@/styles/globals.css";

import { GeistSans } from "geist/font/sans";
import { type Metadata } from "next";
import Providers from "./providers";
import Footer from "./footer";
import Image from "next/image";
import LatestRelease from "./hostdVersion";
import Link from "next/link";
import Storage from "./storage";

export const metadata: Metadata = {
  title: "Sia Alert",
  description: "Control Alerts for your Sia Host",
  icons: [{ rel: "icon", url: "/favicon.ico" }],
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en" className={`${GeistSans.variable}`}>
      <body className="flex min-h-screen w-full flex-col justify-between bg-black bg-gradient-to-b from-[#2e026d] to-[#15162c] font-mono text-emerald-500">
        <Providers>
          <header className="flex w-full items-center justify-between gap-2 p-2 pb-6">
            <div className="flex items-center gap-6">
              <Link href="/" className="flex items-center gap-2">
                <Image src={"/logo.png"} alt="logo" width={64} height={64} />
                <h1 className="text-4xl font-bold underline">Sia Host Alert</h1>
              </Link>
              {/* <Link href="/api-doc" className="flex items-center gap-2 underline">API</Link> */}
            </div>
            <Storage />
            <LatestRelease />
          </header>
          {children}
          <Footer />
        </Providers>
      </body>
    </html>
  );
}
