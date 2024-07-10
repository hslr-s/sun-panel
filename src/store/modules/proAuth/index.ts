import { defineStore } from 'pinia'
import { ref } from 'vue'
import { proIsExpired as proIsExpiredApi } from '@/api/openness'

interface ProIsExpired {
  isExpired: boolean
  hideProBadge: boolean
}

export const useProAuthorizeStore = defineStore('proAuthorize', () => {
  const proIsExpired = ref(true)
  const hideProBadge = ref(true)
  async function getIsExpired() {
    try {
      const { data } = await proIsExpiredApi<ProIsExpired>()
      proIsExpired.value = data.isExpired
      hideProBadge.value = !data.isExpired && data.hideProBadge
    }
    catch (error) {

    }
  }
  getIsExpired()
  return { isExpired: proIsExpired, hideProBadge, refreshProAuthStatus: getIsExpired }
})
