declare namespace HomePage{


    interface State{
        active:string
        spaceId:number // 空间的id
        notesList:Info[]
    }

}

declare namespace HomePage.quest{

    interface LoginReqest{
        username:string 
        password:string
        vcode?:string
    }

	interface LoginResponse extends User.Info{
		token :string
	}

    interface ResetPasswordByVCodeReqest extends System.Register.SendRegisterVcodeRquest{
    }

}