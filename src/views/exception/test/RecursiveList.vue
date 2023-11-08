<template>
    <div v-for="(item, index) in items" :key="index" :class="isChilden ? 'ml-[20px]' : ''">
        <!-- <span class="items-center flex">
            <NAvatar round :size="30" :src="item.headImage ? item.headImage : defaultAvatar" />
        </span> -->


        <div
            class="flex items-center cursor-pointer my-[15px] p-[5px] hover:bg-slate-100 dark:hover:bg-slate-800 rounded-md">

            <div class="w-[25px]">
                <span v-if="item.children" class="cursor-pointer">
                    <!-- 收起展开：'rotate-90' :'rotate-0' -->
                    <SvgIcon width="18" height="18" class="list-expand" :class="item.extand?'rotate-90':'rotate-0'" icon="ic:round-play-arrow" />
                </span>
            </div>

            <span class="flex items-center mr-[5px] text-slate-400 hover:text-black dark:hover:text-white">
                {{ item.name }}
                <!-- <NDropdown trigger="click" :options="item.isTop === 1 ? unSetTopOptions : setTopOptions" size="small"
                @select="(key: string) => handleClick(key, item.id)">
                <SvgIcon icon="ri:more-fill" />
            </NDropdown> -->

            </span>


        </div>
        <div v-show="item.extand">
            <!-- 如果有子级，递归渲染 -->
            <recursive-list v-if="item.children" :is-childen="true" :items="item.children"></recursive-list>
        </div>

    </div>
</template>
  
<script setup lang="ts">
import { SvgIcon } from '@/components/common'
import { defineProps } from "vue"
defineProps<{
    items: Array<any>,
    isChilden?: boolean
}>()

</script>
  
<style>
.list-expand {
    transition-property: transform;
    transition-duration: .25s;
}
</style>