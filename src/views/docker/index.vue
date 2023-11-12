<script lang="ts" setup>
import { ref } from 'vue'
import { NLayout, NLayoutSider, NMenu, NSpace } from 'naive-ui'
import { useRoute, useRouter } from 'vue-router'
import { Container } from './components'
// import Header from './layout/Header/index.vue'
const inverted = ref(false)
const router = useRouter()
const route = useRoute()

const menuOptions = [
  {
    label: '容器',
    key: 'DockerContainer',
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

//   export default defineComponent({
//     setup () {
//       return {
//         inverted: ref(false),
//         menuOptions
//       }
//     }
//   })
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
          :width="240"
          :native-scrollbar="false"
          show-trigger
          :inverted="inverted"
          class="h-[calc(100vh-60px)]"
        >
          <NMenu
            :value="route.name as string"
            :inverted="inverted"
            :collapsed-width="0"
            :collapsed-icon-size="0"
            :options="menuOptions"
            :default-expanded-keys="['user-manage', 'system-manage']"
            @update:value="(key: string) => router.push({ name: key })"
          />
        </NLayoutSider>
        <NLayout class="h-[calc(100vh-60px)] dark:bg-slate-900 bg-slate-100 ">
          <div class="p-[20px]">
            <Container />
          </div>
          <!-- <router-view class=" h-[calc(100vh-60px)] p-[20px] bg-slate-100 dark:bg-slate-900" /> -->
          <!-- <RouterView v-slot="{ Component, route }">
            <component :is="Component" :key="route.fullPath" />
          </RouterView> -->
        </NLayout>
      </NLayout>
      <!-- <NLayoutFooter :inverted="inverted" bordered>
        Footer Footer Footer
      </NLayoutFooter> -->
    </NLayout>
  </NSpace>
</template>
