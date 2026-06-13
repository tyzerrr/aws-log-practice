import type { Metadata } from "next";
import "./globals.css";
import Footer from "@/components/Footer";
import Header from "@/components/Header";

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
    <html lang="ja" className="h-screen antialiased">
      <body className="flex min-h-screen flex-col">
        <Header />
        <main className={"flex-1"}>{children}</main>
        <Footer />
      </body>
    </html>
  );
}
