import "./globals.css"
import type { Metadata } from "next"
import { Inter } from "next/font/google"
import WalletConnect from "@/components/WalletConnect"
import { WalletProvider } from "@/lib/wallet"

import Link from "next/link"

const inter = Inter({ subsets: ["latin"] })

export const metadata: Metadata = {
  title: "Sphere Market",
  description: "A sleek next-gen NFT marketplace on a custom Layer 1",
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className={`${inter.className} bg-gradient-to-br from-black via-zinc-900 to-zinc-800 text-white`}>
        <WalletProvider>
        <header className="sticky top-0 z-50 bg-white/5 backdrop-blur-xl shadow-sm border-b border-white/10">
          <nav className="max-w-6xl mx-auto px-6 py-4 flex items-center justify-between">
            <Link href="/" className="text-xl font-bold tracking-tight text-white hover:text-cyan-400 transition-colors">
              🌐 Sphere Market
            </Link>
            <div className="flex gap-6 items-center text-sm font-medium text-zinc-300">
              <Link href="/mint" className="hover:text-white transition-all duration-200">Mint</Link>
              <Link href="/market" className="hover:text-white transition-all duration-200">Market</Link>
              <Link href="/profile" className="hover:text-white transition-all duration-200">Profile</Link>
              <WalletConnect />
            </div>
          </nav>
        </header>
        </WalletProvider>
        <main className="max-w-5xl mx-auto px-6 py-10 animate-fadeIn">{children}</main>
      </body>
    </html>
  )
}
