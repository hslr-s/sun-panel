<script setup lang="ts">
    import {ref} from "vue"
    const a = ref<any>()
</script>
<template>
    <div>
            <VueDraggable
              v-model="items"
              item-key="id"
              :animation="300"
              class="mx-auto mt-4 grid grid-cols-4 sm:grid-cols-6 md:grid-cols-8 lg:grid-cols-10 xl:12 gap-5"
              filter=".not-drag"
            >
              <div v-for="item, index in items" :key="index" @contextmenu="(e) => handleContextMenu(e, item)">
                <div
                  class="sunpanel w-[70px] h-[70px] mx-auto rounded-2xl cursor-pointer transition-all duration-200 hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)] bg-[#2a2a2a6b]"
                  @click="handleItemClick(item)"
                >
                  <ItemIcon :item-icon="item.icon" />
                </div>
                <div class="text-center app-icon-text-shadow cursor-pointer mt-[2px]" :style="{ color: panelState.panelConfig.iconTextColor }" @click="handleItemClick(item)">
                  <span>{{ item.title }}</span>
                </div>
              </div>

              <div class="not-drag">
                <div class="w-[70px] h-[70px] mx-auto rounded-2xl cursor-pointer transition-all duration-200 hover:shadow-[0_0_20px_10px_rgba(0,0,0,0.2)]" @click="handleAddAppClick">
                  <ItemIcon :item-icon="{ itemType: 3, text: 'subway:add', bgColor: '#343434' }" />
                </div>
                <div class="text-center app-icon-text-shadow cursor-pointer mt-[2px]" :style="{ color: panelState.panelConfig.iconTextColor }" @click="handleAddAppClick">
                  添加图标
                </div>
              </div>

            </vuedraggable>
          </div>
</template>