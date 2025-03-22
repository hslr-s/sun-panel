import { PanelPanelConfigStyleEnum } from '@/enums'

export function correctionNumberByCardStyle(v: number, cardStyle: PanelPanelConfigStyleEnum): number {
  let keepNum = 0
  if (cardStyle === PanelPanelConfigStyleEnum.small)
    keepNum = 1
  else if (cardStyle === PanelPanelConfigStyleEnum.info)
    keepNum = 2

  return correctionNumber(v, keepNum)
}

export function correctionNumber(v: number, keepNum = 2): number {
  return v === 0 ? 0 : Number(v.toFixed(keepNum))
}

export function allCPUUsages(cpuInfo: SystemMonitor.CPUInfo | null) {
  if (!cpuInfo)
    return 0
  let countPercentage = 0
  for (let i = 0; i < cpuInfo.usages.length; i++) {
    const element = cpuInfo.usages[i]
    countPercentage += element
  }
  return countPercentage / cpuInfo.usages.length
}
