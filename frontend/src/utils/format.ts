// Format volume with B (billion) and W (ten thousand) suffixes
export const formatVolume = (vol: number): string => {
  if (vol >= 100000000) return (vol / 100000000).toFixed(2) + 'B'
  if (vol >= 10000) return (vol / 10000).toFixed(2) + 'W'
  return vol.toLocaleString()
}

// Format amount with B/W suffixes
export const formatAmount = (amt: number): string => {
  if (amt >= 100000000) return (amt / 100000000).toFixed(2) + 'B'
  if (amt >= 10000) return (amt / 10000).toFixed(2) + 'W'
  return amt.toFixed(2)
}
