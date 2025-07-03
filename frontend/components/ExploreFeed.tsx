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
  index: number
  hash: string
  transactions: Transaction[]
}

export default function ExploreFeed() {
  const [nfts, setNfts] = useState<Transaction[]>([])

  useEffect(() => {
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/blocks`)
      .then(res => res.json())
      .then((blocks: Block[]) => {
        const allTxs = blocks.flatMap(b => b.transactions)
        const nftTxs = allTxs.filter(tx => tx.from === "system" && tx.data.length > 0)
        setNfts(nftTxs.reverse()) // show newest first
      })
  }, [])

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
      {nfts.length === 0 && <p>No NFTs yet. Mint one!</p>}
      {nfts.map((nft, i) => (
        <NFTCard key={i} owner={nft.to} metadata={nft.data} txId={nft.id} />
      ))}
    </div>
  )
}
