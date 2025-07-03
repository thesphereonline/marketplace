"use client"

import { useEffect, useState } from "react"
import NFTCard from "@/components/NFTCard"

type Transaction = {
  id: string
  from: string
  to: string
  amount: number
  data: string
}

export default function MarketPage() {
  const [listings, setListings] = useState<Transaction[]>([])

  useEffect(() => {
  fetch(`${process.env.NEXT_PUBLIC_API_URL}/listings`)
    .then((res) => res.json())
    .then((listings: any[]) => {
      setListings(
        listings.map((l) => ({
          id: l.id,
          from: l.owner,
          to: "marketplace",
          amount: l.price,
          data: "LIST:" + l.metadata,
        }))
      )
    })
}, [])

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">🛒 Marketplace</h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
        {listings.map((tx, i) => (
          <NFTCard
            key={i}
            owner={tx.from}
            metadata={tx.data.replace("LIST:", "")}
            txId={tx.id}
            isListed
          />
        ))}
      </div>
    </div>
  )
}
