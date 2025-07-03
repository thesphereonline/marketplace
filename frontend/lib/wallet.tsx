"use client"

import { createContext, useContext, useEffect, useState } from "react"

type WalletContextType = {
  address: string | null
  connect: () => Promise<void>
}

const WalletContext = createContext<WalletContextType>({
  address: null,
  connect: async () => {},
})

export const WalletProvider = ({ children }: { children: React.ReactNode }) => {
  const [address, setAddress] = useState<string | null>(null)

  const connect = async () => {
    if (typeof window.ethereum !== "undefined") {
      const accounts = await window.ethereum.request({ method: "eth_requestAccounts" })
      setAddress(accounts[0])
    } else {
      alert("Please install MetaMask")
    }
  }

  useEffect(() => {
    if (typeof window.ethereum !== "undefined") {
      window.ethereum.request({ method: "eth_accounts" }).then((accs: string[]) => {
        if (accs[0]) setAddress(accs[0])
      })
    }
  }, [])

  return (
    <WalletContext.Provider value={{ address, connect }}>
      {children}
    </WalletContext.Provider>
  )
}

export const useWallet = () => useContext(WalletContext)
