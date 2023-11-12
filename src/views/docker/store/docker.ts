import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useDockerStore = defineStore('docker', () => {
  const state = ref<Docker.State>()

  return { state }
})
