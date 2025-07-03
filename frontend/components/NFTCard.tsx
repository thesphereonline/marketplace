"use client"

import { useState } from "react"

type Props = {
  owner: string
  metadata: string
  txId: string
  isOwner?: boolean
  isListed?: boolean
}

export default function NFTCard({ owner, metadata, txId, isOwner, isListed }: Props) {
  const [price, setPrice] = useState("")
  const [msg, setMsg] = useState("")

  const listNFT = async () => {
    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/list`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        id: txId,
        from: owner,
        to: "marketplace",
        amount: parseInt(price),
        data: `LIST:${metadata}`,
      }),
    })
    if (res.ok) setMsg("✅ Listed!")
    else setMsg("❌ Failed to list")
  }

  const buyNFT = async () => {
  if (typeof window !== "undefined" && typeof window.ethereum !== "undefined") {
    const accounts = await window.ethereum.request({ method: "eth_requestAccounts" })
    const buyer = accounts[0]

    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/buy`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        id: txId,
        from: buyer,
        to: owner,
        amount: 250, // ideally: fetch actual listing price from backend
        data: `BUY:${metadata}`,
      }),
    })

    if (res.ok) setMsg("✅ Purchase complete!")
    else setMsg("❌ Purchase failed")
  } else {
    alert("Please install MetaMask or another Ethereum wallet.")
  }
}


  return (
    <div className="glass rounded-xl p-4 space-y-2">
      <p className="text-xs text-zinc-400">Owner: {owner.slice(0, 10)}...</p>
      <h3 className="text-lg font-semibold">{metadata}</h3>
      <p className="text-[10px] text-zinc-500">ID: {txId.slice(0, 12)}...</p>

      {isOwner && (
        <div className="flex items-center gap-2 mt-2">
          <input
            className="bg-zinc-900 border px-2 py-1 rounded text-sm w-24"
            placeholder="Price"
            value={price}
            onChange={(e) => setPrice(e.target.value)}
          />
          <button className="btn text-sm" onClick={listNFT}>List</button>
        </div>
      )}

      {isListed && !isOwner && (
        <button className="btn text-sm w-full mt-2" onClick={buyNFT}>
          💰 Buy Now
        </button>
      )}

      {msg && <p className="text-xs mt-1">{msg}</p>}
    </div>
  )
}
