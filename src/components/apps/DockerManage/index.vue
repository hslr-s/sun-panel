<script lang="ts" setup>
import { NCard } from 'naive-ui'
import { Container } from './components'
import { useDockerStore } from './store/docker'
// import Header from './layout/Header/index.vue'
// const inverted = ref(false)

const dockerStore = useDockerStore()

// const menuOptions = [
//   {
//     label: '容器',
//     key: 'container',
//   },
//   //   {
//   //     label: "",
//   //     key: 'compose',
//   //   },
//   {
//     label: '镜像',
//     key: 'image',
//   },
//   {
//     label: '网络',
//     key: 'network',
//   },

// ]

// console.log(inverted, menuOptions)
</script>

<template>
  <div v-if="dockerStore.actionMenuBar === 'container'" class="h-full overflow-auto p-2">
    <NCard style="border-radius:10px;height: 100%;" size="small" content-style="height: 100%">
      <Container />
    </NCard>
  </div>
</template>

<style scoped>
.docker-menu /deep/  .n-menu-item-content::before {
  border-radius: 8px;
}
</style>
