<script setup lang="ts">
import { NCard } from 'naive-ui'

defineProps<{
  nickName?: string
  proFlag?: string
  username?: string
  expired?: string
  cardStyle?: string | object
}>()
</script>

<template>
  <NCard :style="cardStyle" class="card">
    <div class="flex text-xl mb-5">
      <span>
        <slot name="nickName">
          {{ nickName }}
        </slot>
      </span>
      <span class="ml-auto">
        <slot name="proFlag">
          {{ proFlag }}
        </slot>
      </span>
    </div>
    <div>
      <slot name="username">
        {{ username }}
      </slot>
    </div>
    <div class="flex justify-end mt-10 text-xs">
      <slot name="expired">
        {{ expired }}
      </slot>
    </div>
    <div class="card-pattern" />
  </NCard>
</template>

<style scoped>
.card-pattern:before {
    position: absolute;
    width: 100px;
    display: block;
    height: 100px;
    content: "";
    border-radius: 50%;
    z-index: -1;
    right: -27px;
    top: -35px;
    background: #ffffff3b;
    box-shadow: -8px 21px 0px #ffffff1a;
}

.card-pattern:after {
    position: absolute;
    width: 70px;
    display: block;
    height: 70px;
    border: 4px solid #ffffff3b;
    /* 边框为4像素的实线，颜色为淡白色带透明度的3b */
    content: "";
    border-radius: 50%;
    bottom: -19px;
    left: 48px;
    box-shadow: -20px -50px 0px #ffffff1a;
    z-index: -1;
}

.card {
    overflow: hidden;
    border-radius: 1rem;
    z-index: 10;
}
</style>
