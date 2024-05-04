<script lang="ts" setup>
import { ref } from 'vue'
import { Container } from './components'
import { useDockerStore } from './store/docker'
// import Header from './layout/Header/index.vue'
const inverted = ref(false)

const dockerStore = useDockerStore()

const menuOptions = [
  {
    label: '容器',
    key: 'container',
  },
  //   {
  //     label: "",
  //     key: 'compose',
  //   },
  {
    label: '镜像',
    key: 'image',
  },
  {
    label: '网络',
    key: 'network',
  },

]
</script>

<template>
  <div v-if="dockerStore.actionMenuBar === 'container'" class="p-[20px]">
    <Container />
  </div>
</template>

<style scoped>
.docker-menu /deep/  .n-menu-item-content::before {
  border-radius: 8px;
}
</style>
