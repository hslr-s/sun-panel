<script setup lang="ts">
import { VueDraggable } from 'vue-draggable-plus'
import { ref } from 'vue'

const list = ref<[any]>([])
for (let i = 1; i < 10; i++) {
  list.value.push({
    id: 1,
    name: `哈哈哈哈${i}`,
  })
}
</script>

<template>
  <div>
    <!-- <div class="mt-4 grid lg:grid-cols-12 gap-5"> -->
    <VueDraggable v-model="list" item-key="id" class="mt-4 grid grid-cols-3 sm:grid-cols-6 md:grid-cols-8 lg:grid-cols-10 xl:12 gap-5">
      <div v-for="item, index in list" :key="index" class="bg-[pink] m-[10px]">
        {{ item.name }}
      </div>
    </VueDraggable>
    <!-- </div> -->
  </div>
</template>
