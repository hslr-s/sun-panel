export enum PanelStateNetworkModeEnum {
  // 局域网
  'lan' = 0,
  // 互联网
  'wan' = 1,
}

export enum PanelPanelConfigStyleEnum {
  'icon' = 1, // 图标风格
  'info' = 0, // 详情风格
  'small' = 1, // 同icon
}

export enum PanelItemIconTypeEnum {
  'text' = 1, // 文字
  'url' = 2, // 图标
  'online' = 3, // 在线
}

export enum PanelItemOpenMethodEnum {
  'currentTab' = 1,
  'newTab' = 2,
  'modal' = 3,
}

// 卡片类型
export enum PanelCardType {
  Normal = 1, // 标准卡片
  Docker = 2, // docker
  BuiltInApp = 3, // 内置应用

}

export enum FileType {
  Unknown = 0,
  Wallpaper = 1,
  CardIcon = 2,
}
