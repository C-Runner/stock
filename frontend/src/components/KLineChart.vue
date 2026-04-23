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

// Process data for candlestick chart
const candlestickData = computed(() =>
  props.priceData.map(p => [p.open, p.close, p.low, p.high])
)

const dates = computed(() =>
  props.priceData.map(p => {
    // On mobile, show fewer date labels
    if (isMobile.value) {
      const parts = p.date.split('-')
      return `${parts[1]}/${parts[2]}`
    }
    return p.date
  })
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
      data: mobile ? [] : ['Candlestick', 'MA5', 'MA10', 'MA20', 'Volume'],
      top: 0,
      textStyle: { color: 'rgba(255, 255, 255, 0.6)', fontSize: mobile ? 10 : 12 }
    },
    grid: [
      { left: leftMargin, right: rightMargin, top: '8%', height: mobile ? '45%' : '50%' },
      { left: leftMargin, right: rightMargin, top: mobile ? '60%' : '68%', height: mobile ? '18%' : '20%' }
    ],
    xAxis: [
      {
        type: 'category',
        data: dates.value,
        gridIndex: 0,
        axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } },
        axisLabel: {
          color: 'rgba(255, 255, 255, 0.5)',
          fontSize,
          interval: mobile ? 4 : 0,
          height: labelHeight
        }
      },
      {
        type: 'category',
        data: dates.value,
        gridIndex: 1,
        axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } },
        axisLabel: {
          color: 'rgba(255, 255, 255, 0.5)',
          fontSize,
          interval: mobile ? 4 : 0,
          height: labelHeight
        }
      }
    ],
    yAxis: [
      {
        scale: true,
        gridIndex: 0,
        splitLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.05)' } },
        axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } },
        axisLabel: {
          color: 'rgba(255, 255, 255, 0.5)',
          fontSize,
          margin: mobile ? 4 : 8
        }
      },
      {
        scale: true,
        gridIndex: 1,
        splitNumber: 2,
        splitLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.05)' } },
        axisLine: { lineStyle: { color: 'rgba(255, 255, 255, 0.1)' } },
        axisLabel: {
          color: 'rgba(255, 255, 255, 0.5)',
          fontSize,
          margin: mobile ? 4 : 8
        }
      }
    ],
    dataZoom: mobile ? [
      { type: 'inside', xAxisIndex: [0, 1], start: 60, end: 100 },
      { type: 'slider', xAxisIndex: [0, 1], start: 60, end: 100, top: '88%', height: '4%' }
    ] : [
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
        },
        barWidth: mobile ? '40%' : '60%'
      },
      {
        name: 'MA5',
        type: 'line',
        data: getMAData(5),
        xAxisIndex: 0,
        yAxisIndex: 0,
        smooth: true,
        showSymbol: false,
        lineStyle: { width: mobile ? 1 : 1.5, color: '#ff6b6b' }
      },
      {
        name: 'MA10',
        type: 'line',
        data: getMAData(10),
        xAxisIndex: 0,
        yAxisIndex: 0,
        smooth: true,
        showSymbol: false,
        lineStyle: { width: mobile ? 1 : 1.5, color: '#feca57' }
      },
      {
        name: 'MA20',
        type: 'line',
        data: getMAData(20),
        xAxisIndex: 0,
        yAxisIndex: 0,
        smooth: true,
        showSymbol: false,
        lineStyle: { width: mobile ? 1 : 1.5, color: '#54a0ff' }
      },
      {
        name: 'Volume',
        type: 'bar',
        data: volumeData.value,
        xAxisIndex: 1,
        yAxisIndex: 1,
        barWidth: mobile ? '40%' : '60%'
      }
    ]
  }
})
</script>

<template>
  <v-chart
    :option="option"
    :autoresize="true"
    :style="{ height: isMobile ? '300px' : props.height }"
  />
</template>

<style scoped>
</style>
