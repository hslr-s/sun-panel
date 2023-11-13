<script lang="ts" setup>
import { ref } from 'vue'
import { NLayout, NLayoutSider, NMenu, NSpace } from 'naive-ui'
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
  <NSpace vertical>
    <NLayout>
      <Header />
      <NLayout has-sider>
        <NLayoutSider
          bordered
          collapse-mode="width"
          :collapsed-width="0"
          :width="150"
          :native-scrollbar="false"
          show-trigger
          :inverted="inverted"
          class="h-[100vh]"
        >
          <NMenu
            :value="dockerStore.actionMenuBar"
            :inverted="inverted"
            :collapsed-width="0"
            :default-value="dockerStore.actionMenuBar"
            :collapsed-icon-size="0"
            :options="menuOptions"
            class="docker-menu"
            @update:value="(key: string) => dockerStore.actionMenuBar = key"
          />
        </NLayoutSider>
        <NLayout class="h-[calc(100vh-60px)] dark:bg-slate-900 bg-slate-100 ">
          <div v-if="dockerStore.actionMenuBar === 'container'" class="p-[20px]">
            <Container />
          </div>
        </NLayout>
      </NLayout>
    </NLayout>
  </NSpace>
</template>

<style scoped>
.docker-menu /deep/  .n-menu-item-content::before {
  border-radius: 8px;
}
</style>
