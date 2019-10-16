import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {tap} from 'rxjs/operators';
import {Observable, Observer, Subject} from 'rxjs';
import {Execution} from './execution';

@Injectable({
  providedIn: 'root'
})
export class OperaterService {
  private baseUrl = '/api/v1/clusters/{clusterName}/executions/';

  executionQueue = new Subject<Execution>();
  $executionQueue = this.executionQueue.asObservable();

  constructor(private http: HttpClient) {
  }

  executeOperate(clusterName, opt: string, params?): Observable<any> {
    if (!params) {
      params = {};
    }
    return this.http.post(`${this.baseUrl.replace('{clusterName}', clusterName)}`, {'operation': opt, 'params': params}).pipe(
      tap(data => {
        this.executionQueue.next(data);
      })
    );
  }


}
