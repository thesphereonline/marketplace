"use client"

import { useWallet } from "@/lib/wallet"
import { useEffect, useState } from "react"
import { Wallet } from "lucide-react" // optional: clean icon

export default function WalletConnect() {
  const { address, connect } = useWallet()
  const [displayAddress, setDisplayAddress] = useState<string | null>(null)

  useEffect(() => {
    if (address) {
      setDisplayAddress(`${address.slice(0, 6)}...${address.slice(-4)}`)
    } else {
      setDisplayAddress(null)
    }
  }, [address])

  return (
    <div className="flex items-center gap-3">
      {address ? (
        <div className="flex items-center gap-2 text-sm text-green-400 font-medium">
          <Wallet size={16} className="text-green-400" />
          <span>{displayAddress}</span>
        </div>
      ) : (
        <button
          onClick={connect}
          className="px-4 py-2 rounded-xl text-sm font-medium bg-blue-600 hover:bg-blue-700 text-white transition"
        >
          Connect Wallet
        </button>
      )}
    </div>
  )
}
