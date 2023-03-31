import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ChatComponent } from './components/chat/chat.component';
import { MessageBoxComponent } from './components/message-box/message-box.component';
import { LoginRegisterComponent } from './components/login-register/login-register.component';
import { ChatBoxComponent } from './components/chat-box/chat-box.component';
@NgModule({
  declarations: [
    AppComponent,
    ChatComponent,
    MessageBoxComponent,
    LoginRegisterComponent,
    ChatBoxComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    HttpClientModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
