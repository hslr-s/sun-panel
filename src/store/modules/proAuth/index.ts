import { defineStore } from 'pinia'
import { ref } from 'vue'
import { proIsExpired as proIsExpiredApi } from '@/api/openness'

interface ProIsExpired {
  isExpired: boolean
}

export const useProAuthorizeStore = defineStore('proAuthorize', () => {
  const proIsExpired = ref(true)
  async function getIsExpired() {
    try {
      const { data } = await proIsExpiredApi<ProIsExpired>()
      proIsExpired.value = data.isExpired
    }
    catch (error) {

    }
  }
  getIsExpired()
  return { isExpired: proIsExpired }
})
