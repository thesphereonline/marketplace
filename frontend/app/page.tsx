"use client"

import ExploreFeed from "@/components/ExploreFeed"
import { motion } from "framer-motion"

export default function Home() {
  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
    >
      <h1 className="text-3xl font-bold mb-6">🪐 Explore NFTs on Sphere</h1>
      <ExploreFeed />
    </motion.div>
  )
}
