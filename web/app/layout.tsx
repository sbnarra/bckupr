import type { Metadata } from "next";
import { Inter } from "next/font/google";
import {Providers} from "./providers";

import "./globals.css";
import Navbar from "./navbar";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Bckupr",
  description: "Bckupr Web Interface",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
  
}>) {
  return (
    <html lang="en" className='light'>
      <body>
        <Providers>
          <Navbar />
            {children}
        </Providers>
      </body>
    </html>
  );
}
