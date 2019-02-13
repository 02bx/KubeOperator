import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {Cluster, ExtraConfig} from '../cluster';
import {TipService} from '../../tip/tip.service';
import {ClrWizard} from '@clr/angular';
import {Config, Package, Template} from '../../package/package';
import {PackageService} from '../../package/package.service';
import {TipLevels} from '../../tip/tipLevels';
import {Node} from '../../node/node';
import {ClusterService} from '../cluster.service';
import {NodeService} from '../../node/node.service';
import {RelationService} from '../relation.service';

@Component({
  selector: 'app-cluster-create',
  templateUrl: './cluster-create.component.html',
  styleUrls: ['./cluster-create.component.css']
})
export class ClusterCreateComponent implements OnInit {


  @ViewChild('wizard') wizard: ClrWizard;
  createClusterOpened: boolean;
  isSubmitGoing = false;
  cluster: Cluster = new Cluster();
  template: Template;
  configs: Config[] = [];
  packages: Package[] = [];
  templates: Template[] = [];
  nodes: Node[] = [];
  options = {};
  @Output() create = new EventEmitter<boolean>();
  loadingFlag = false;

  constructor(private tipService: TipService, private nodeService: NodeService, private clusterService: ClusterService,
              private packageService: PackageService, private relationService: RelationService) {
  }

  ngOnInit() {
    this.listPackages();
    this.generateChars();
  }

  newCluster() {
    // 清空对象
    this.reset();
    this.createClusterOpened = true;
  }

  reset() {
    this.wizard.reset();
    this.cluster = new Cluster();
    this.template = null;
    this.templates = null;
    this.nodes = null;
    this.configs = null;
  }

  packgeOnChange() {
    this.packages.forEach((pak) => {
      if (pak.name === this.cluster.package) {
        this.configs = pak.meta.configs;
        this.templates = pak.meta.templates;
      }
    });
  }

  listPackages() {
    this.packageService.listPackage().subscribe(data => {
      this.packages = data;
    }, error => {
      this.tipService.showTip('加载离线包错误!: \n' + error, TipLevels.ERROR);
    });
  }

  templateOnChange() {
    this.nodes = [];
    this.templates.forEach(tmp => {
      if (tmp.name === this.cluster.template) {
        tmp.roles.forEach(role => {
          if (!role.meta.hidden) {
            const name = role.name;
            const roleNumber = role.meta.nodes_require[1];
            for (let i = 0; i < roleNumber; i++) {
              const node: Node = new Node();
              node.name = role.name + '-' + i;
              node.roles.push(role.name);
              this.nodes.push(node);
            }
          }
        });
      }
    });
  }


  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.clusterService.createCluster(this.cluster).subscribe(data => {
      this.createNodes(this.cluster.name);
      this.configCluster(this.cluster.name);
      this.isSubmitGoing = false;
      this.createClusterOpened = false;
      this.create.emit(true);
    });
  }

  configCluster(clusterName) {
    this.configs.forEach(config => {
      const extraConfig: ExtraConfig = new ExtraConfig();
      extraConfig.key = config.name;
      extraConfig.value = config.value;
      this.clusterService.configCluster(clusterName, extraConfig).subscribe();
    });
  }


  createNodes(clusterName) {
    this.isSubmitGoing = true;
    this.nodes.forEach(node => {
      this.nodeService.createNode(clusterName, node).subscribe();
    });
  }


  generateChars() {
    this.options = this.relationService.genOptions(this.nodes);
  }


  onCancel() {
    this.reset();
    this.createClusterOpened = false;
  }

}
