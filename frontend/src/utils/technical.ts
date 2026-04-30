export const getLatestValue = (values: number[]): number => {
  if (!values || values.length === 0) return 0
  return values[values.length - 1] ?? 0
}

export const getRSIClass = (rsi: number): string => {
  if (rsi > 70) return 'overbought'
  if (rsi < 30) return 'oversold'
  return ''
}

export const getKDJClass = (value: number): string => {
  if (value > 80) return 'overbought'
  if (value < 20) return 'oversold'
  return ''
}

export const getWRClass = (wr: number): string => {
  if (wr > -20) return 'overbought'
  if (wr < -80) return 'oversold'
  return ''
}
