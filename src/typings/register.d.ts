declare namespace System.Register{
    interface SendRegisterVcodeRquest {
        email     ?:string 
	    username  ?:string 
	    password    ?:string 
	    vcode     ?:string 
	    emailVCode ?:string 
		verification?:Common.VerificationRequest
		referralCode?:string
    }

	interface CommitRquest extends SendRegisterVcodeRquest {
		
    }
	
}