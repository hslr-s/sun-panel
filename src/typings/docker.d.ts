declare namespace Docker {


    interface State {
        actionMenuBar: string | 'container'
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
    

}

