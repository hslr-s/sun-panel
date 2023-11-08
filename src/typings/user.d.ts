declare namespace User{

	interface Info{
		id?:number
		name ?:string
		createTime?:string
		username?:string
		password?:string
		headImage?:string
		status?:number
		role?:number
		mail?:string
		// userId?:string // id代替
		token?:string
		isAdmin?:number
	}

	interface GetReferralCodeResponse{
		referralCode:string
	}


	
}