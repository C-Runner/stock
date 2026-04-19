<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { CandlestickChart, BarChart, LineChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  DataZoomComponent
} from 'echarts/components'
import type { EChartsOption } from 'echarts'
import type { PricePoint, MAData } from '../api'

// Register ECharts components
use([
  CanvasRenderer,
  CandlestickChart,
  BarChart,
  LineChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  DataZoomComponent
])

interface Props {
  priceData: PricePoint[]
  maData?: MAData[]
  height?: string
}

const props = withDefaults(defineProps<Props>(), {
  height: '450px'
})

// Process data for candlestick chart
const candlestickData = computed(() =>
  props.priceData.map(p => [p.open, p.close, p.low, p.high])
)

const dates = computed(() =>
  props.priceData.map(p => p.date)
)

const volumeData = computed(() =>
  props.priceData.map(p => ({
    value: p.volume,
    itemStyle: {
      color: p.close >= p.open ? '#ef5350' : '#26a69a'
    }
  }))
)

// Get MA value for a period
const getMAValue = (period: number): number => {
  const ma = props.maData?.find(m => m.period === period)
  return ma?.value ?? 0
}

// MA line data - replicate the value across all dates for display
const getMAData = (period: number) => {
  const value = getMAValue(period)
  if (!value) return []
  return Array(dates.value.length).fill(value)
}

const option = computed<EChartsOption>(() => ({
  backgroundColor: '#1a1a1a',
  animation: false,
  tooltip: {
    trigger: 'axis',
    axisPointer: { type: 'cross' },
    backgroundColor: '#333',
    borderColor: '#555',
    textStyle: { color: '#fff' }
  },
  legend: {
    data: ['Candlestick', 'MA5', 'MA10', 'MA20', 'Volume'],
    top: 0,
    textStyle: { color: '#999' }
  },
  grid: [
    { left: '10%', right: '8%', top: '10%', height: '50%' },
    { left: '10%', right: '8%', top: '68%', height: '20%' }
  ],
  xAxis: [
    {
      type: 'category',
      data: dates.value,
      gridIndex: 0,
      axisLine: { lineStyle: { color: '#333' } },
      axisLabel: { color: '#666' }
    },
    {
      type: 'category',
      data: dates.value,
      gridIndex: 1,
      axisLine: { lineStyle: { color: '#333' } },
      axisLabel: { color: '#666' }
    }
  ],
  yAxis: [
    {
      scale: true,
      gridIndex: 0,
      splitLine: { lineStyle: { color: '#222' } },
      axisLine: { lineStyle: { color: '#333' } },
      axisLabel: { color: '#666' }
    },
    {
      scale: true,
      gridIndex: 1,
      splitNumber: 2,
      splitLine: { lineStyle: { color: '#222' } },
      axisLine: { lineStyle: { color: '#333' } },
      axisLabel: { color: '#666' }
    }
  ],
  dataZoom: [
    { type: 'inside', xAxisIndex: [0, 1], start: 70, end: 100 },
    { type: 'slider', xAxisIndex: [0, 1], start: 70, end: 100, top: '92%' }
  ],
  series: [
    {
      name: 'Candlestick',
      type: 'candlestick',
      data: candlestickData.value,
      xAxisIndex: 0,
      yAxisIndex: 0,
      itemStyle: {
        color: '#ef5350',
        color0: '#26a69a',
        borderColor: '#ef5350',
        borderColor0: '#26a69a'
      }
    },
    {
      name: 'MA5',
      type: 'line',
      data: getMAData(5),
      xAxisIndex: 0,
      yAxisIndex: 0,
      smooth: true,
      showSymbol: false,
      lineStyle: { width: 1, color: '#ff6b6b' }
    },
    {
      name: 'MA10',
      type: 'line',
      data: getMAData(10),
      xAxisIndex: 0,
      yAxisIndex: 0,
      smooth: true,
      showSymbol: false,
      lineStyle: { width: 1, color: '#feca57' }
    },
    {
      name: 'MA20',
      type: 'line',
      data: getMAData(20),
      xAxisIndex: 0,
      yAxisIndex: 0,
      smooth: true,
      showSymbol: false,
      lineStyle: { width: 1, color: '#54a0ff' }
    },
    {
      name: 'Volume',
      type: 'bar',
      data: volumeData.value,
      xAxisIndex: 1,
      yAxisIndex: 1
    }
  ]
}))
</script>

<template>
  <v-chart
    ref="chartRef"
    :option="option"
    :autoresize="true"
    :style="{ height: props.height }"
  />
</template>

<style scoped>
</style>
