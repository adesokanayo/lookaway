import type { Metadata } from "next";
import { Fraunces, Source_Sans_3 } from "next/font/google";
import "./globals.css";

const display = Fraunces({
  variable: "--font-display-loaded",
  subsets: ["latin"],
});

const sans = Source_Sans_3({
  variable: "--font-sans-loaded",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Lookaway",
  description:
    "Every 20 minutes the screen takes over so you look about 20 feet away for 20 seconds.",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="en" className={`${display.variable} ${sans.variable} h-full antialiased`}>
      <body className="min-h-full bg-background text-foreground">{children}</body>
    </html>
  );
}
