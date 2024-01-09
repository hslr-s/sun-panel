import type { MonitorData } from './typings'
import { useModuleConfig } from '@/store/modules'

const modelName = 'systemMonitor'
const moduleConfig = useModuleConfig()

export async function saveAll(value: MonitorData[]) {
  return await moduleConfig.saveToCloud(modelName, { list: value })
}

export async function getAll(): Promise< MonitorData[]> {
  const res = await moduleConfig.getValueByNameFromCloud<{ list: MonitorData[] }>(modelName)
  if (res.code === 0 && res.data && res.data.list)
    return res.data.list
  return []
}

export async function add(value: MonitorData): Promise<boolean> {
  let success = true
  let newData: MonitorData[] = []
  try {
    const data = await getAll()
    if (data)
      newData = data

    newData.push(value)
    const res = await saveAll(newData)
    if (res.code !== 0)
      console.log('save failed', res)
  }
  catch (error) {
    console.error(error)
    success = false
  }
  return success
}

export async function saveByIndex(index: number | undefined, value: MonitorData): Promise<boolean> {
  if (!index)
    index = 0

  let success = true
  let newData: MonitorData[] = []
  try {
    const data = await getAll()
    if (data)
      newData = data

    newData[index] = value
    const res = await saveAll(newData)
    if (res.code !== 0)
      console.log('save failed', res)
  }
  catch (error) {
    console.error(error)
    success = false
  }
  return success
}

export async function getByIndex(index: number): Promise<MonitorData | null> {
  try {
    const data = await getAll()
    if (data[index])
      return data[index]
  }
  catch (error) {

  }

  return null
}

export async function deleteByIndex(index: number): Promise<boolean> {
  let success = true
  try {
    const data = await getAll()
    if (data[index])
      data.splice(index, 1)
    await saveAll(data)
  }
  catch (error) {
    success = false
    console.error(error)
  }

  return success
}
