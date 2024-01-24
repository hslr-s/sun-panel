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
