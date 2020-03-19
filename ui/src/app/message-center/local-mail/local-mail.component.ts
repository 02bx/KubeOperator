import {Component, OnInit} from '@angular/core';
import {MessageCenterService} from '../message-center.service';
import {AlertLevels} from '../../base/header/components/common-alert/alert';
import {CommonAlertService} from '../../base/header/common-alert.service';

@Component({
  selector: 'app-local-mail',
  templateUrl: './local-mail.component.html',
  styleUrls: ['./local-mail.component.css']
})
export class LocalMailComponent implements OnInit {

  loading = false;
  messages = [];
  selectedMessages = [];
  limit = 10;
  page = 1;
  type = 'ALL';
  readStatus = 'ALL';
  level = 'ALL';

  constructor(private messageCenterService: MessageCenterService, private alertService: CommonAlertService) {
  }

  ngOnInit() {
    this.listMessage(this.limit, this.type, this.readStatus, this.level);
  }

  listMessage(limit, type, readStatus, level) {
    this.limit = limit;
    this.type = type;
    this.readStatus = readStatus;
    this.level = level;
    this.loading = true;
    this.messageCenterService.listUserMessageByPage(this.limit, this.page, this.type, this.readStatus, this.level).subscribe(data => {
      this.messages = data;
      this.loading = false;
    });
  }

  updateMessage() {
    const promises: Promise<{}>[] = [];
    this.selectedMessages.forEach(msg => {
      promises.push(this.messageCenterService.updateUserMessage(msg).toPromise());
    });
    Promise.all(promises).then(() => {
      this.listMessage(this.limit, this.type, this.readStatus, this.level);
      this.alertService.showAlert('更新成功', AlertLevels.SUCCESS);
    }, res => {
      this.alertService.showAlert('更新失败' + res.error.msg, AlertLevels.ERROR);
    }).finally(() => {
      this.selectedMessages = [];
    });
  }

  updateAllMessage() {
    this.messageCenterService.updateAllUserMessage().subscribe(data => {
      this.alertService.showAlert('更新成功', AlertLevels.SUCCESS);
      this.listMessage(this.limit, this.type, this.readStatus, this.level);
    });
  }


}
