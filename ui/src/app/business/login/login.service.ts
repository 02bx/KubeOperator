import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {LoginCredential} from './login-credential';
import {Observable} from 'rxjs';
import {User} from '../user/user';

@Injectable({
    providedIn: 'root'
})
export class LoginService {

    loginUrl = '/api/v1/login';

    constructor(private http: HttpClient) {
    }

    login(item: LoginCredential): Observable<User> {
        return this.http.post<User>(this.loginUrl, item);
    }
}
