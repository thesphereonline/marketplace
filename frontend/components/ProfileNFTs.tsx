"use client"

import { useEffect, useState } from "react"
import NFTCard from "./NFTCard"

type Transaction = {
  id: string
  from: string
  to: string
  amount: number
  data: string
}

type Block = {
  transactions: Transaction[]
}

export default function ProfileNFTs() {
  const [nfts, setNfts] = useState<Transaction[]>([])
  const [address, setAddress] = useState<string | null>(null)

  useEffect(() => {
    const getWallet = async () => {
      if (typeof window.ethereum !== "undefined") {
        const accounts = await window.ethereum.request({ method: "eth_requestAccounts" })
        setAddress(accounts[0])
      }
    }

    getWallet()
  }, [])

  useEffect(() => {
    if (!address) return

    fetch(`${process.env.NEXT_PUBLIC_API_URL}/blocks`)
      .then((res) => res.json())
      .then((blocks: Block[]) => {
        const allTxs = blocks.flatMap((b) => b.transactions)
        const ownedNFTs = allTxs.filter(
          (tx) => tx.from === "system" && tx.to.toLowerCase() === address.toLowerCase()
        )
        setNfts(ownedNFTs.reverse())
      })
  }, [address])

  if (!address) return <p>🔌 Please connect your wallet.</p>
  if (nfts.length === 0) return <p>🖼️ No NFTs found for this wallet.</p>

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
      {nfts.map((nft, i) => (
        <NFTCard key={i} owner={nft.to} metadata={nft.data} txId={nft.id} />
      ))}
    </div>
  )
}
