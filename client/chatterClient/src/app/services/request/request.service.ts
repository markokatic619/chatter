import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http'
import { Observable } from 'rxjs';

interface LoginResponse {
  loginCookie: string;
}
interface RegisterResponce {
  loginCookie: string;
  responseMessage:string;
}
@Injectable({
  providedIn: 'root'
})

export class RequestService {

  private loginUrl = "http://localhost:8080/login";
  private registerUrl = "http://localhost:8080/register";

  constructor( private http:HttpClient) { }

  login(data: any): Observable<LoginResponse> 
  {
    return this.http.post<LoginResponse>(this.loginUrl,data)
  }
  register(data:any)
  {
    return this.http.post<RegisterResponce>(this.registerUrl,data)
  }
}
