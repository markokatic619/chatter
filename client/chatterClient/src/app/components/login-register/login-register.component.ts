import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { RequestService } from '../../../app/services/request/request.service'
import { CookieService } from 'ngx-cookie-service'
@Component({
  selector: 'login-register',
  templateUrl: './login-register.component.html',
  styleUrls: ['./login-register.component.css']
})


export class LoginRegisterComponent {

  aplicationName = "Chatter"
  selectedButton = "login"
  
  constructor( private http: HttpClient,private requestService:RequestService, private cookies : CookieService){}

  getUserData(data: any) {
    if(data.email != "" || data.password != "")
    {
      this.requestService.login(data).subscribe((result) => {
      const dateNow = new Date();
      dateNow.setDate(dateNow.getDate() + 14);
      this.cookies.set("loginCookie", result.loginCookie, dateNow);
    });
    }
  }
  getUserDataForRegistration(data:any){
    if(data.name != "" && data.name.split(' ').length == 2 && data.email != "" && data.birthday != "" && data.password != "") 
    {
      var firstName = data.name.split(' ')[0]
      var lastName = data.name.split(' ')[1]
      data.firstName = firstName
      data.lastName = lastName
      delete data.name
      this.requestService.register(data).subscribe((result) => {
        if(result.responseMessage=="success")
        {
          alert("success")
          const dateNow = new Date();
          dateNow.setDate(dateNow.getDate() + 14);
          this.cookies.set("loginCookie", result.loginCookie, dateNow);
        }
        else{
          this.errorMessage(result.responseMessage)
        }
      })
    }
    else 
    {
      this.errorMessage("Invalid registration form")
    }
   
  }
  getLoginForm(){this.selectedButton="login"}

  getRegisterForm(){this.selectedButton="register"}

  errorMessage(message:string){
    alert(message)
  }
}
