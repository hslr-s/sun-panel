import { inArray } from '@/utils/cmn'
import type { SystemApp } from '@/utils/cmn/builtInApps'
export function getAppList(systemApps: SystemApp[], userRole: number): SystemApp[] {
  const applist: SystemApp[] = []
  systemApps.forEach((systemApp) => {
    if (!systemApp.roles || inArray(systemApp.roles, userRole)) {
      applist.push(systemApp)
    }
  })

  return applist
}
