declare namespace GlobalSetting {

    interface SiteSetting {
        favicoUrl: string
        title: string
        loginBackgroundUrl: string   
    }

    interface GlobalSettingCustomJsAndCss {
        jsContent: string;
        cssContent: string;
    }

    interface GlobalSettingGetCustomPathReq {
        customPath: string;
    }
    

}
