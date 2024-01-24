declare namespace SystemMonitor {

    interface CPUInfo {
        coreCount: number
        cpuNum: number
        model: string
        usages: number[]
    }

    interface DiskInfo {
        mountpoint: string
        total: number
        used: number
        free: number
        usedPercent: number
    }

    interface NetIOCountersInfo {
        bytesSent: number
        bytesRecv: number
        name: string
    }

    interface MemoryInfo {
        total: number
        used: number
        free: number
        usedPercent: number
    }

    interface GetAllRes {
        cpuInfo: CPUInfo
        diskInfo: DiskInfo[]
        netIOCountersInfo: NetIOCountersInfo[]
        memoryInfo: MemoryInfo
    }

    interface Mountpoint{
        device:string
        mountpoint:string
        fstype:string
    }
}