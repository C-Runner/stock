<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { CandlestickChart, BarChart, LineChart } from 'echarts/charts'
import {
  TooltipComponent,
  LegendComponent,
  GridComponent,
  DataZoomComponent,
  MarkLineComponent
} from 'echarts/components'
import type { EChartsOption } from 'echarts'
import type { PricePoint, MAData, RSIData, KDJData, MACDData, BOLLData } from '../api'

use([
  CanvasRenderer,
  CandlestickChart,
  BarChart,
  LineChart,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  DataZoomComponent,
  MarkLineComponent
])

interface Props {
  priceData: PricePoint[]
  maData?: MAData[]
  emaData?: MAData[]
  rsiData?: RSIData[]
  kdjData?: KDJData
  macdData?: MACDData
  bollData?: BOLLData
  height?: string
}

const props = withDefaults(defineProps<Props>(), {
  height: '700px'
})

const windowWidth = ref(window.innerWidth)
const isMobile = computed(() => windowWidth.value < 768)

const handleResize = () => {
  windowWidth.value = window.innerWidth
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})

// Use backend-provided MA data directly
const ma5Data = computed(() => props.maData?.find(m => m.period === 5)?.values.slice(-ma20BaseLength.value) ?? [])
const ma10Data = computed(() => props.maData?.find(m => m.period === 10)?.values.slice(-ma20BaseLength.value) ?? [])
const ma20BaseLength = computed(() => props.maData?.find(m => m.period === 20)?.values.length ?? 0)
const ma20Data = computed(() => props.maData?.find(m => m.period === 20)?.values ?? [])

// Use backend-provided EMA data directly
const ema12Data = computed(() => props.emaData?.find(e => e.period === 12)?.values.slice(-ma20BaseLength.value) ?? [])
const ema26Data = computed(() => props.emaData?.find(e => e.period === 26)?.values.slice(-ma20BaseLength.value) ?? [])

// Process data for candlestick chart - using aligned data
const candlestickData = computed(() =>
  alignedPriceData.value.map(p => [p.open, p.close, p.low, p.high])
)

const dates = computed(() =>
  props.priceData.map(p => {
    if (isMobile.value) {
      const parts = p.date.split('-')
      return `${parts[1]}/${parts[2]}`
    }
    return p.date
  })
)

// Aligned price data for all charts (from index 19 to match indicator start)
const alignedPriceData = computed(() => props.priceData.slice(19))

// Align dates with indicator data (MA20 starts at index 19)
const alignedDates = computed(() => {
  const startIdx = 19 // MA20 period - 1
  return dates.value.slice(startIdx)
})

// Use backend-provided BOLL data directly
const bollDataComputed = computed(() => ({
  upper: props.bollData?.upper ?? [],
  mid: props.bollData?.mid ?? [],
  lower: props.bollData?.lower ?? []
}))

// RSI data from backend (time series) - slice to match aligned data length
const rsi6Data = computed(() => props.rsiData?.find(r => r.period === 6)?.values.slice(-ma20Data.value.length) ?? [])
const rsi12Data = computed(() => props.rsiData?.find(r => r.period === 12)?.values.slice(-ma20Data.value.length) ?? [])
const rsi24Data = computed(() => props.rsiData?.find(r => r.period === 24)?.values.slice(-ma20Data.value.length) ?? [])

// KDJ data from backend (time series) - slice to match aligned data length
const kData = computed(() => props.kdjData?.k.slice(-ma20Data.value.length) ?? [])
const dData = computed(() => props.kdjData?.d.slice(-ma20Data.value.length) ?? [])
const jData = computed(() => props.kdjData?.j.slice(-ma20Data.value.length) ?? [])

// Volume data - using aligned data
const volumeData = computed(() =>
  alignedPriceData.value.slice(-ma20Data.value.length).map(p => ({
    value: p.volume,
    itemStyle: {
      color: p.close >= p.open ? '#ef5350' : '#26a69a'
    }
  }))
)

const option = computed<EChartsOption>(() => {
  const mobile = isMobile.value
  const leftMargin = mobile ? '3%' : '10%'
  const rightMargin = mobile ? '3%' : '8%'
  const fontSize = mobile ? 10 : 12
  const labelHeight = mobile ? 20 : 30

  return {
    backgroundColor: 'rgba(10, 10, 15, 0.8)',
    animation: false,
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'cross' },
      backgroundColor: 'rgba(30, 27, 75, 0.95)',
      borderColor: 'rgba(99, 102, 241, 0.3)',
      textStyle: { color: 'rgba(255, 255, 255, 0.9)', fontSize: mobile ? 11 : 12 }
    },
    legend: {
      data: mobile ? [] : ['Candlestick', 'MA5', 'MA10', 'MA20', 'BOLL-U', 'BOLL-L', 'Volume', 'RSI', 'KDJ'],
      top: 0,
      textStyle: { color: 'rgba(255, 255, 255, 0.6)', fontSize: mobile ? 10 : 12 },
      inactiveColor: 'rgba(255, 255, 255, 0.2)'
    },
    grid: mobile ? [
      // Grid 0: Candlestick + MA/EMA/BOLL (main chart) - mobile optimized
      { left: leftMargin, right: rightMargin, top: '6%', height: '36%' },
      // Grid 1: Volume
      { left: leftMargin, right: rightMargin, top: '48%', height: '10%' },
      // Grid 2: RSI
      { left: leftMargin, right: rightMargin, top: '62%', height: '14%' },
      // Grid 3: KDJ
      { left: leftMargin, right: rightMargin, top: '80%', height: '14%' }
    ] : [
      // Grid 0: Candlestick + MA/EMA/BOLL (main chart)
      { left: leftMargin, right: rightMargin, top: '8%', height: '38%' },
      // Grid 1: Volume
      { left: leftMargin, right: rightMargin, top: '52%', height: '10%' },
      // Grid 2: RSI
      { left: leftMargin, right: rightMargin, top: '66%', height: '12%' },
      // Grid 3: KDJ
      { left: leftMargin, right: rightMargin, top: '82%', height: '12%' }
    ],
    xAxis: (mobile ? [
      { type: 'category' as const, data: alignedDates.value, gridIndex: 0, axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } }, axisLabel: { color: 'rgba(255, 255, 255, 0.5)', fontSize, interval: 4, height: labelHeight } },
      { type: 'category' as const, data: alignedDates.value, gridIndex: 1, axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } }, axisLabel: { color: 'rgba(255, 255, 255, 0.5)', fontSize, interval: 4, height: labelHeight } },
      { type: 'category' as const, data: alignedDates.value, gridIndex: 2, axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } }, axisLabel: { color: 'rgba(255, 255, 255, 0.5)', fontSize, interval: 4, height: labelHeight } },
      { type: 'category' as const, data: alignedDates.value, gridIndex: 3, axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } }, axisLabel: { color: 'rgba(255, 255, 255, 0.5)', fontSize, interval: 4, height: labelHeight } }
    ] : [
      { type: 'category' as const, data: alignedDates.value, gridIndex: 0, axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } }, axisLabel: { color: 'rgba(255, 255, 255, 0.5)', fontSize, interval: 0, height: labelHeight } },
      { type: 'category' as const, data: alignedDates.value, gridIndex: 1, axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } }, axisLabel: { color: 'rgba(255, 255, 255, 0.5)', fontSize, interval: 0, height: labelHeight } },
      { type: 'category' as const, data: alignedDates.value, gridIndex: 2, axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } }, axisLabel: { color: 'rgba(255, 255, 255, 0.5)', fontSize, interval: 0, height: labelHeight } },
      { type: 'category' as const, data: alignedDates.value, gridIndex: 3, axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } }, axisLabel: { color: 'rgba(255, 255, 255, 0.5)', fontSize, interval: 0, height: labelHeight } }
    ]),
    yAxis: [
      { scale: true, gridIndex: 0, splitLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.05)' } }, axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } }, axisLabel: { color: 'rgba(255, 255, 255, 0.5)', fontSize, margin: mobile ? 4 : 8 } },
      { scale: true, gridIndex: 1, splitNumber: 2, splitLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.05)' } }, axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } }, axisLabel: { color: 'rgba(255, 255, 255, 0.5)', fontSize: mobile ? 8 : 10, margin: mobile ? 2 : 4, formatter: (val: number) => {
        if (val >= 100000000) return (val / 100000000).toFixed(1) + 'E'
        if (val >= 10000000) return (val / 10000000).toFixed(1) + 'C'
        if (val >= 1000000) return (val / 1000000).toFixed(1) + 'M'
        if (val >= 1000) return (val / 1000).toFixed(0) + 'K'
        return String(val)
      } } },
      { scale: true, gridIndex: 2, min: 0, max: 100, splitNumber: mobile ? 2 : 4, splitLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.05)' } }, axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } }, axisLabel: { color: 'rgba(255, 255, 255, 0.5)', fontSize: mobile ? 9 : 10, margin: mobile ? 2 : 6, formatter: (val: number) => val === 0 || val === 50 || val === 100 ? String(val) : '' } },
      { scale: true, gridIndex: 3, splitNumber: mobile ? 2 : 4, splitLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.05)' } }, axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } }, axisLabel: { color: 'rgba(255, 255, 255, 0.5)', fontSize: mobile ? 9 : 10, margin: mobile ? 2 : 6 } }
    ],
    dataZoom: mobile ? [
      { type: 'inside', xAxisIndex: [0, 1, 2, 3], start: 50, end: 100 },
      { type: 'slider', xAxisIndex: [0, 1, 2, 3], start: 50, end: 100, top: '96%', height: '4%' }
    ] : [
      { type: 'inside', xAxisIndex: [0, 1, 2, 3], start: 70, end: 100 },
      { type: 'slider', xAxisIndex: [0, 1, 2, 3], start: 70, end: 100, top: '96%', height: '3%' }
    ],
    series: [
      // Grid 0: Candlestick + MA/EMA/BOLL
      { name: 'Candlestick', type: 'candlestick', data: candlestickData.value, xAxisIndex: 0, yAxisIndex: 0, itemStyle: { color: '#ef5350', color0: '#26a69a', borderColor: '#ef5350', borderColor0: '#26a69a' }, barWidth: mobile ? '40%' : '60%' },
      { name: 'MA5', type: 'line', data: ma5Data.value, xAxisIndex: 0, yAxisIndex: 0, smooth: true, showSymbol: false, lineStyle: { width: mobile ? 1 : 1.5, color: '#ff6b6b' } },
      { name: 'MA10', type: 'line', data: ma10Data.value, xAxisIndex: 0, yAxisIndex: 0, smooth: true, showSymbol: false, lineStyle: { width: mobile ? 1 : 1.5, color: '#feca57' } },
      { name: 'MA20', type: 'line', data: ma20Data.value, xAxisIndex: 0, yAxisIndex: 0, smooth: true, showSymbol: false, lineStyle: { width: mobile ? 1 : 1.5, color: '#54a0ff' } },
      { name: 'EMA12', type: 'line', data: ema12Data.value, xAxisIndex: 0, yAxisIndex: 0, smooth: true, showSymbol: false, lineStyle: { width: mobile ? 1 : 1.5, color: '#9b59b6', type: 'dashed' } },
      { name: 'EMA26', type: 'line', data: ema26Data.value, xAxisIndex: 0, yAxisIndex: 0, smooth: true, showSymbol: false, lineStyle: { width: mobile ? 1 : 1.5, color: '#e91e63', type: 'dashed' } },
      { name: 'BOLL-U', type: 'line', data: bollDataComputed.value.upper, xAxisIndex: 0, yAxisIndex: 0, smooth: true, showSymbol: false, lineStyle: { width: mobile ? 1 : 1, color: 'rgba(255, 255, 255, 0.4)', type: 'dotted' } },
      { name: 'BOLL-L', type: 'line', data: bollDataComputed.value.lower, xAxisIndex: 0, yAxisIndex: 0, smooth: true, showSymbol: false, lineStyle: { width: mobile ? 1 : 1, color: 'rgba(255, 255, 255, 0.4)', type: 'dotted' } },

      // Grid 1: Volume
      { name: 'Volume', type: 'bar', data: volumeData.value, xAxisIndex: 1, yAxisIndex: 1, barWidth: mobile ? '40%' : '60%' },

      // Grid 2: RSI (simplify on mobile - only show RSI6)
      ...(mobile ? [
        { name: 'RSI6' as const, type: 'line' as const, data: rsi6Data.value, xAxisIndex: 2, yAxisIndex: 2, smooth: true, showSymbol: false, lineStyle: { width: 1, color: '#ff6b6b' } },
        { name: 'RSI_70' as const, type: 'line' as const, data: Array(ma20Data.value.length).fill(70), xAxisIndex: 2, yAxisIndex: 2, showSymbol: false, lineStyle: { color: 'rgba(255,107,107,0.5)', type: 'dashed' as const, width: 1 } },
        { name: 'RSI_30' as const, type: 'line' as const, data: Array(ma20Data.value.length).fill(30), xAxisIndex: 2, yAxisIndex: 2, showSymbol: false, lineStyle: { color: 'rgba(56,239,125,0.5)', type: 'dashed' as const, width: 1 } }
      ] : [
        { name: 'RSI6' as const, type: 'line' as const, data: rsi6Data.value, xAxisIndex: 2, yAxisIndex: 2, smooth: true, showSymbol: false, lineStyle: { width: 1.5, color: '#ff6b6b' } },
        { name: 'RSI12' as const, type: 'line' as const, data: rsi12Data.value, xAxisIndex: 2, yAxisIndex: 2, smooth: true, showSymbol: false, lineStyle: { width: 1.5, color: '#feca57' } },
        { name: 'RSI24' as const, type: 'line' as const, data: rsi24Data.value, xAxisIndex: 2, yAxisIndex: 2, smooth: true, showSymbol: false, lineStyle: { width: 1.5, color: '#54a0ff' } },
        { name: 'RSI_70' as const, type: 'line' as const, data: Array(ma20Data.value.length).fill(70), xAxisIndex: 2, yAxisIndex: 2, showSymbol: false, lineStyle: { color: 'rgba(255,107,107,0.4)', type: 'dashed' as const, width: 1 }, markLine: { silent: true, symbol: 'none', label: { show: true, formatter: '70', position: 'end' as const, color: 'rgba(255,107,107,0.8)' }, data: [{ yAxis: 70 }] } },
        { name: 'RSI_30' as const, type: 'line' as const, data: Array(ma20Data.value.length).fill(30), xAxisIndex: 2, yAxisIndex: 2, showSymbol: false, lineStyle: { color: 'rgba(56,239,125,0.4)', type: 'dashed' as const, width: 1 }, markLine: { silent: true, symbol: 'none', label: { show: true, formatter: '30', position: 'end' as const, color: 'rgba(56,239,125,0.8)' }, data: [{ yAxis: 30 }] } }
      ]),

      // Grid 3: KDJ
      { name: 'K', type: 'line', data: kData.value, xAxisIndex: 3, yAxisIndex: 3, smooth: true, showSymbol: false, lineStyle: { width: mobile ? 1 : 1.5, color: '#ff6b6b' } },
      { name: 'D', type: 'line', data: dData.value, xAxisIndex: 3, yAxisIndex: 3, smooth: true, showSymbol: false, lineStyle: { width: mobile ? 1 : 1.5, color: '#feca57' } },
      { name: 'J', type: 'line', data: jData.value, xAxisIndex: 3, yAxisIndex: 3, smooth: true, showSymbol: false, lineStyle: { width: mobile ? 1 : 1.5, color: '#9b59b6' } }
    ]
  }
})
</script>

<template>
  <v-chart
    :option="option"
    :autoresize="true"
    :style="{ height: isMobile ? '500px' : props.height }"
  />
</template>

<style scoped>
</style>
