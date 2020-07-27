import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {ProjectMember, ProjectMemberRequest} from '../project-member';
import {ProjectMemberService} from '../project-member.service';
import {NgForm} from '@angular/forms';
import {ResourceTypes} from '../../../../constant/shared.const';
import {ActivatedRoute} from '@angular/router';
import {Project} from '../../project';
import {ModalAlertService} from '../../../../shared/common-component/modal-alert/modal-alert.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';

@Component({
    selector: 'app-project-member-create',
    templateUrl: './project-member-create.component.html',
    styleUrls: ['./project-member-create.component.css']
})
export class ProjectMemberCreateComponent extends BaseModelComponent<ProjectMember> implements OnInit {

    opened = false;
    item: ProjectMemberRequest = new ProjectMemberRequest();
    selectUsers: string[] = [];
    roles: string[] = [];
    currentProject: Project = new Project();
    @Output() created = new EventEmitter();
    @ViewChild('memberForm') memberForm: NgForm;


    constructor(private projectMemberService: ProjectMemberService,
                private route: ActivatedRoute,
                private modalAlertService: ModalAlertService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
        super(projectMemberService);
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentProject = data.project;
        });
    }

    open() {
        this.opened = true;
        this.item = new ProjectMemberRequest();
        this.getRoles();
    }

    onCancel() {
        this.opened = false;
        // this.item = new ProjectMemberRequest();
        // this.memberForm.resetForm(this.item);
    }

    onSubmit() {
        this.item.projectId = this.currentProject.id;
        this.projectMemberService.create(this.item).subscribe(res => {
            this.opened = false;
            this.created.emit();
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    leaveInput() {
        this.projectMemberService.getUsers(this.item.name).subscribe(res => {
            this.selectUsers = res.items;
        });
    }

    handleValidation() {

    }

    selectedName(name) {
        this.item.name = name;
        this.selectUsers = [];
    }

    getRoles() {
        this.projectMemberService.getRoles().subscribe(res => {
            this.roles = res;
        });
    }
}
