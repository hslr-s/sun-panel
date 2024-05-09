declare namespace Docker {


    interface State {
        actionMenuBar: string | 'container'
        containerList:Info[]
    }

    

}

declare namespace Docker.Container {
    interface Port {
        ip: string;
        privatePort: number;
        publicPort: number;
        type: "tcp" | "udp" | "sctp"
    }
    
    interface Info {
        id :string
        name: string;
        image: string;
        imageID?: string;
        command?: string;
        createTime: string;
        ports: Port[];
        state: string;
        status: string;
    }

    interface ExecSwitchActionReq {
        id :string
        action:boolean
    }

    interface Stats {
        cpuPercent: number;
        memory: number;
        cache: number;
        ioRead: number;
        ioWrite: number;
        networkRX: number;
        networkTX: number;
        shotTime?: string; // 时间戳字符串
        runStatus:string
        errorMessage:string
    }
    
    

}

