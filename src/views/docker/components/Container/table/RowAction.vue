<script setup lang="ts">
import { NSwitch, useMessage } from 'naive-ui'
import { computed, ref } from 'vue'
import { useDockerStore } from '../../../store/docker'
import { DockerContainerInfoState } from '@/enums/docker'
import { execSwitchAction } from '@/api/docker/container'

const props = defineProps<{
  value: Docker.Container.Info
}>()

const dockerState = useDockerStore()
const data = ref(props.value)
const ms = useMessage()
const switchLoading = ref(false)

const switchShow = computed(() => {
  switch (data.value.state) {
    case DockerContainerInfoState.created:
    case DockerContainerInfoState.running:
    case DockerContainerInfoState.paused:
    case DockerContainerInfoState.dead:
    case DockerContainerInfoState.exited:
      return true

    case DockerContainerInfoState.restarting:
    case DockerContainerInfoState.removing:
      return true

    default:
      return false
  }
})

const switchValue = computed({
  get() {
    if (props.value.state === DockerContainerInfoState.running)
      return true
    return false
  },
  set() {
    // 更新状态
    switchLoading.value = true
    const acticonStatus = !switchValue.value
    execSwitchAction({ id: data.value.id, action: acticonStatus }).then(({ code, msg }) => {
      if (code !== 0)
        ms.warning(msg)
      else
        data.value.state = acticonStatus ? DockerContainerInfoState.running : DockerContainerInfoState.exited

      switchLoading.value = false
      dockerState.refreshContainerList()
    })
  },
})
</script>

<template>
  <div class="w-full">
    <div v-if="switchShow">
      <NSwitch v-model:value="switchValue" :loading="switchLoading" />
    </div>
    <div v-else>
      <NSwitch v-model:value="switchValue" :loading="switchLoading" disabled />
    </div>
  </div>
</template>
