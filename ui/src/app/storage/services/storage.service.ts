import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {StorageTemplate} from '../models/storage-template';
import {Observable} from 'rxjs';
import {Storage} from '../models/storage';

@Injectable({
  providedIn: 'root'
})
export class StorageService {

  private storageUrl = '/api/v1/storage/';

  constructor(private http: HttpClient) {
  }

  listStorage(): Observable<Storage[]> {
    return this.http.get<Storage[]>(this.storageUrl);
  }

  getStorage(name: string): Observable<Storage> {
    return this.http.get<Storage>(this.storageUrl + name);
  }

  createStorage(storage: Storage): Observable<Storage> {
    return this.http.post<Storage>(this.storageUrl, storage);
  }

  deleteStorage(name: string): Observable<Storage> {
    return this.http.delete<Storage>(this.storageUrl + name);
  }
}

