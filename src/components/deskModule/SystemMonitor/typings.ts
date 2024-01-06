export enum MonitorType {
  'cpu' = 'cpu', // 图标风格
  'memory' = 'memory', // 详情风格
  'disk' = 'disk',
}

export interface CardStyle {
  background: string
}

export interface MonitorData {
  monitorType: MonitorType
  extendParam?: { [key: string]: [value:any] } | any
  description?: string
  // cardStyle: CardStyle
}

export interface ProgressStyle {
  color: string
  railColor: string
  height: number
}

export interface GenericProgressStyleExtendParam {
  progressColor: string
  progressRailColor: string
  color: string
  backgroundColor: string
}

export interface DiskExtendParam extends GenericProgressStyleExtendParam {
  path: string
}
