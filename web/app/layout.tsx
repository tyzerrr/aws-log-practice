import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "AWS Log Practice",
  description: "Product operations dashboard",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja" className="h-full antialiased">
      <body className="flex min-h-full flex-col">{children}</body>
    </html>
  );
}
