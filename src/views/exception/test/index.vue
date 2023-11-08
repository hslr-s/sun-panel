<script setup lang="ts">
import RecursiveList from './RecursiveList.vue'

const multiLevelData = [
  {
    name: '开发笔记',
    children: [
      {
        name: 'Level 2 Item 1',
        extand: true,
        children: [
          { name: 'SAI-Chat开发' },
          { name: '笔记本项目' },
        ],
      },
      { name: '测试项目' },
    ],
  },
  {
    name: '学习笔记',
    children: [
      { name: 'Blender' },
      { name: 'mongo' },
    ],
  },
]
</script>

<template>
  <div>
    <RecursiveList :items="multiLevelData" />
  </div>
</template>
