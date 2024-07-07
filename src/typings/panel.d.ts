enum DockerContainerInfoState {
    created = 'created',
    restarting = 'restarting',
    removing = 'removing',
    running = 'running',
    paused = 'paused',
    exited = 'exited',
    dead = 'dead',
}


declare namespace Panel {

    interface Info extends ItemInfo {

    }

    interface ItemInfo extends Common.InfoBase {
        icon: ItemIcon | null
        title: string
        url: string
        sort?: number
        lanUrl?: string
        description?: string
        openMethod: number
        itemIconGroupId?: number

        // 以下1.4.0新增
        onlyName?: string

        // 以下1.5.0新增
        backgroundColor?: string
        cardType: number
        expandParam?: any
    }

    interface ItemIconGroup extends Common.InfoBase {
        icon?: string
        title?: string
        onlyName?:string
        sort?: number
        publicVisitHide?:boolean
    }

    interface ItemIcon {
        itemType: number
        src?: string
        text?: string
        backgroundColor?: string // 从1.5.0开始逐渐废弃
    }

    interface State {
        rightSiderCollapsed: boolean
        leftSiderCollapsed: boolean
        networkMode: PanelStateNetworkModeEnum | null
        panelConfig: panelConfig
    }

    interface panelConfig {
        backgroundImageSrc?: string
        backgroundBlur?: number
        backgroundMaskNumber?: number
        iconStyle?: PanelPanelConfigStyleEnum    // 1.5.0+ 弃用、迁移至单独分组设置
        iconTextColor?: string   // 1.5.0+ 弃用、迁移至单独分组设置
        iconTextInfoHideDescription?: boolean    // 1.5.0+ 弃用、迁移至单独分组设置
        iconTextIconHideTitle?: boolean  // 1.5.0+ 弃用、迁移至单独分组设置
        logoText?: string
        logoImageSrc?: string
        clockShowSecond?: boolean
        clockColor?: string
        searchBoxShow?: boolean
        searchBoxSearchIcon?: boolean
        searchBoxBackground?:string
        searchBoxBorderColor?:string
        // searchBoxPlaceholderColor?:string
        searchBoxTextColor?:string
        marginTop?: number
        marginBottom?: number
        maxWidth?: number
        maxWidthUnit: string
        marginX?: number
        footerHtml?: string
        systemMonitorShow?: boolean
        systemMonitorShowTitle?: boolean
        systemMonitorPublicVisitModeShow?: boolean
        netModeChangeButtonShow?: boolean
    }

    interface userConfig {
        panel: panelConfig
        searchEngine?: any
    }

    interface ItemIconSortRequest {
        sortItems: Common.SortItemRequest[]
        itemIconGroupId: number
    }

    interface ItemGroupInfoCardStyle {
        style: number;
        textColor: string;
        textInfoHideDescription: boolean;
        textIconHideTitle: boolean;
    }

    interface ItemGroupInfo extends ItemIconGroup {
        cardStyle: ItemGroupInfoCardStyle
    }
}

declare namespace Panel.ItemCardExpandParam {


    // interface CardInfo extends Panel.ItemInfo {

    // }

    // // 普通
    // interface Standard extends Panel.ItemInfo {

    // }

    // 系统/内置应用
    interface SystemApp {
        systemAppName: string

    }

    // docker卡片
    interface Docker {
        containerId: string
        containerName: string
        openPageAndStart: bool
        openPageAndStartDelaySecond:number // 启动容器后延时n秒打开页面
        cardType: string
        showBadgeStatus?: DockerContainerInfoState[]
    }

    // // 网络模式切换按钮
    // interface SwitchNetMode {
    //     wanIcon: ItemIcon
    //     lanIcon: ItemIcon
    // }


    // docker
    // interface Docker extends Panel.ItemInfo {
    //     containerName:string // 容器名称
    //     imageName:string // 镜像名称

    // }

    // // iframe
    // interface iframe extends Panel.ItemInfo {
    //     url:string // iframe 地址

    // }

}