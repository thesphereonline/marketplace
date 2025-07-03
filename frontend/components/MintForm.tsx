"use client"

import { useState } from "react"

export default function MintForm() {
  const [owner, setOwner] = useState("")
  const [meta, setMeta] = useState("")
  const [result, setResult] = useState("")

  const handleSubmit = async () => {
    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/mint`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ owner, meta }),
    })

    if (res.ok) {
      const json = await res.json()
      setResult(`✅ Minted! Block: ${json.hash.slice(0, 8)}...`)
    } else {
      setResult("❌ Mint failed")
    }
  }

  return (
    <div className="space-y-4">
      <input
        className="border px-3 py-2 w-full rounded"
        placeholder="Wallet address"
        value={owner}
        onChange={(e) => setOwner(e.target.value)}
      />
      <input
        className="border px-3 py-2 w-full rounded"
        placeholder="NFT metadata"
        value={meta}
        onChange={(e) => setMeta(e.target.value)}
      />
      <button className="bg-black text-white px-4 py-2 rounded" onClick={handleSubmit}>
        Mint NFT
      </button>
      {result && <p className="text-sm mt-2">{result}</p>}
    </div>
  )
}
