<script setup lang="ts">
import {
  NButton, NCard, NCheckbox, NCheckboxGroup,
  NFormItem, NGrid, NGridItem, NInput,
  NInputGroup, NSelect, NSpace, NSwitch, useMessage,
} from 'naive-ui'

import type { FormRules } from 'naive-ui'
import { onMounted, ref, watch } from 'vue'
import IconEditor from '../IconEditor.vue'
import { verifyUrlHttp } from './common'
import type { SystemApp } from '@/utils/cmn/builtInApps'
import { t } from '@/locales'
import { AppLoader, ProAuthOverlay, ProTag, RoundCardModal } from '@/components/common/'
import { getList as getContainerListApi } from '@/api/docker/container'
import { getSiteFavicon } from '@/api/panel/itemIcon'
import { DockerContainerInfoState } from '@/enums/docker'
import { useProAuthorizeStore } from '@/store'

const props = defineProps<{
  itemInfo: Panel.Info
  rules: FormRules
}>()

const emit = defineEmits<{

  (e: 'update:itemInfo', info?: Panel.Info): void
  (e: 'update:rules', info?: FormRules): void
}>()

const getIconLoading = ref([false, false])
const containerListLoading = ref(false)
const ms = useMessage()
const proAuthStore = useProAuthorizeStore()
const containers = ref<Docker.Container.Info[]>([])
async function getContainerList() {
  containerListLoading.value = true
  await getContainerListApi<Common.ListResponse<Docker.Container.Info[]>>().then((res) => {
    containers.value = res.data.list
  }).catch((res) => {
    if (res.code === -1) {
      ms.error(`${t('iconItem.dockerGetContainerListFailed')}: ${res.msg}`, {
        closable: true,
        duration: 10000,
        keepAliveOnHover: true,
      })
    }
  })
  containerListLoading.value = false
}

const item = ref(props.itemInfo)
const builtinAppModalShow = ref(false)

// 打开方式
const options = [
  {
    default: true,
    label: t('iconItem.currentPageOpen'),
    value: 1,
  },
  {
    label: t('iconItem.newWindowOpen'),
    value: 2,
  },
  {
    label: t('iconItem.currentPageLayerOpen'),
    value: 3,
  },
]
const rulesExtends: FormRules = {
  url: [
    {
      required: true,
      trigger: 'blur',
      type: 'string',
      message: t('form.required'),
    },
  ],
  expandParam: {
    containerId: {
      required: true,
      trigger: 'blur',
      type: 'string',
      message: t('form.required'),
    },
  },
}

const extendParam: Panel.ItemCardExpandParam.Docker = {
  containerId: '',
  containerName: '',
  openPageAndStart: false,
  cardType: '',
  showBadgeStatus: [],
  openPageAndStartDelaySecond: 0,
}

const model = ref(props.itemInfo)

watch(model.value, () => {
  emit('update:itemInfo', model.value)
})

onMounted(() => {
  const r = { ...props.rules, ...rulesExtends }
  emit('update:rules', r)

  getContainerList()
  // 设置默认值
  if (!props.itemInfo.expandParam) {
    model.value.expandParam = { ...extendParam }
  }
  else {
    const ep = { ...extendParam }
    // 给子字段赋值，防止无用的字段
    ep.containerId = props.itemInfo.expandParam.containerId
    ep.openPageAndStart = props.itemInfo.expandParam.openPageAndStart
    ep.cardType = props.itemInfo.expandParam.cardType
    ep.showBadgeStatus = props.itemInfo.expandParam.showBadgeStatus
    ep.openPageAndStartDelaySecond = props.itemInfo.expandParam.openPageAndStartDelaySecond
    ep.containerName = props.itemInfo.expandParam.containerName

    model.value.expandParam = ep
  }
})

async function getIconByUrl(url: string, loadingIndex: number) {
  getIconLoading.value[loadingIndex] = true
  try {
    const { code, data } = await getSiteFavicon<{ iconUrl: string }>(url)
    if (code === 0) {
      model.value.icon = {
        itemType: 2,
        src: data.iconUrl,
      }
    }
    else {
      ms.error(t('iconItem.geticonFail'))
    }
  }
  catch (error) {
    ms.error(t('iconItem.geticonFail'))
  }
  getIconLoading.value[loadingIndex] = false
}

function handleUpdateshowBadgeStatusValue(value: (string | number)[]) {
  model.value.expandParam.showBadgeStatus = value
  // message.info(JSON.stringify(value))
  // console.log(JSON.stringify(value))
}

function handleChange(value: string, appInfo: SystemApp) {
  const containerName = appInfo.name
  if (model.value.title === '')
    model.value.title = containerName

  model.value.expandParam.containerName = containerName
}
</script>

<template>
  <div>
    <NGrid cols="2" :x-gap="10" item-responsive>
      <NGridItem span="2 500:1">
        <NFormItem
          path="expandParam.containerId"
          :label="$t('iconItem.containerName')"
        >
          <NSelect
            v-model:value="model.expandParam.containerId"
            :loading="containerListLoading"
            label-field="name"
            value-field="id"
            :options="containers"
            @update:value="handleChange"
          />
        </NFormItem>
      </NGridItem>
      <NGridItem span="2 500:1">
        <NFormItem path="title" :label="$t('common.title')">
          <NInput v-model:value="model.title" type="text" show-count :maxlength="20" clearable />
        </NFormItem>
      </NGridItem>
    </NGrid>

    <NFormItem :label="t('iconItem.showBadgeStatus')">
      <NCheckboxGroup :value="model.expandParam.showBadgeStatus" @update:value="handleUpdateshowBadgeStatusValue">
        <NSpace item-style="display: flex;" align="center">
          <NCheckbox :value="DockerContainerInfoState.running" :label="t('iconItem.showBadgeStatusRunning')" />
          <NCheckbox :value="DockerContainerInfoState.restarting" :label="t('iconItem.showBadgeStatusRestarting')" />
          <NCheckbox :value="DockerContainerInfoState.exited" :label="t('iconItem.showBadgeStatusExited')" />
        </NSpace>
      </NCheckboxGroup>
    </NFormItem>

    <NFormItem path="icon" :label="$t('iconItem.icon.iconStyle')">
      <NCard size="small">
        <IconEditor v-model:item-icon="model.icon" />
      </NCard>
    </NFormItem>

    <NFormItem path="url" :label="$t('iconItem.url')">
      <div class="w-full">
        <NInputGroup>
          <NInput v-model:value="model.url" type="text" :maxlength="1000" placeholder="http(s)://" />
          <NButton :disabled="!model.url" :loading="getIconLoading[0]" @click="getIconByUrl(model.url, 0)">
            {{ $t('iconItem.getIcon') }}
          </NButton>
        </NInputGroup>
        <span v-if="!verifyUrlHttp(model.url)" class="text-xs text-orange-400">
          {{ t("iconItem.urlNoHttpStartWarn") }}
        </span>
      </div>
    </NFormItem>

    <NFormItem path="lanUrl" :label="$t('iconItem.lanUrl')">
      <div class="w-full">
        <NInputGroup>
          <NInput v-model:value="model.lanUrl" type="text" :maxlength="1000" :placeholder="$t('iconItem.lanUrlInputPlaceholder')" />
          <NButton :disabled="!model.lanUrl" :loading="getIconLoading[1]" @click="getIconByUrl(model.lanUrl || '', 1)">
            {{ $t('iconItem.getIcon') }}
          </NButton>
        </NInputGroup>
        <span v-if="!verifyUrlHttp(model.lanUrl || '')" class="text-xs text-orange-400">
          {{ t("iconItem.urlNoHttpStartWarn") }}
        </span>
      </div>
    </NFormItem>

    <!-- <NFormItem path="description" :label="$t('common.description')">
      <NInput v-model:value="model.description" type="text" show-count :maxlength="100" clearable />
    </NFormItem> -->
    <NFormItem :label="$t('iconItem.openPageAndStartContainer')">
      <ProAuthOverlay :show="proAuthStore.isExpired">
        <NSwitch v-model:value="model.expandParam.openPageAndStart" />
        <!-- <div v-if="model.expandParam.openPageAndStart" class="flex items-center">
        <span class="mr-2 ml-5">容器启动后延时打开页面</span>
        <NInputNumber
          v-model:value="model.expandParam.openPageAndStartDelaySecond"
          style="width: 80px;"
          :min="0"
          :max="60"
        />
        <span class="ml-2">s</span>
      </div> -->
        <ProTag class="ml-2" />
      </ProAuthOverlay>
    </NFormItem>

    <NFormItem path="openMethod" :label="$t('iconItem.openMethod')">
      <NSelect v-model:value="model.openMethod" :options="options" />
    </NFormItem>
    <!-- <NGrid cols="2" :x-gap="10" item-responsive>
      <NGridItem span="2 500:1">

      </NGridItem>
      <NGridItem span="2 500:1">
        <NFormItem path="openMethod" :label="$t('iconItem.openMethod')">
          <NSelect v-model:value="model.openMethod" :options="options" />
        </NFormItem>
      </NGridItem>
    </NGrid> -->

    <RoundCardModal v-model:show="builtinAppModalShow">
      <AppLoader :component-name="item.expandParam.systemAppName" />
      <!--  -->
    </RoundCardModal>
  </div>
</template>
