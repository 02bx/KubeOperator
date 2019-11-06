import {Component, OnInit} from '@angular/core';
import {ClusterService} from '../cluster/cluster.service';
import {Cluster} from '../cluster/cluster';
import {Router} from '@angular/router';
import {DashboardSearch} from './dashboardSearch';
import {DashboardService} from './dashboard.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

  loading = true;
  clusters: Cluster[] = [];
  selectClusters: Cluster[] = [];
  dashboardSearch: DashboardSearch = new DashboardSearch();
  clusterData = [];
  podCount = 0;
  nodeCount = 0;
  namespaceCount = 0;
  deploymentCount = 0;
  containerCount = 0;
  restartPods = [];
  warnContainers = [];
  cpu_usage = 0;
  mem_usage = 0;
  cpu_total = 0;
  mem_total = 0;
  show_pod_detail = false;
  show_container_detail = false;

  constructor(private clusterService: ClusterService, private router: Router, private dashboardService: DashboardService) {
  }

  ngOnInit() {
    this.dashboardSearch.cluster = 'all';
    this.dashboardSearch.dateLimit = 1;
    this.search();
  }

  data_init() {
    this.clusterData = [];
    this.podCount = 0;
    this.nodeCount = 0;
    this.namespaceCount = 0;
    this.deploymentCount = 0;
    this.containerCount = 0;
    this.restartPods = [];
    this.warnContainers = [];
    this.cpu_usage = 0;
    this.mem_usage = 0;
    this.cpu_total = 0;
    this.mem_total = 0;
    this.show_pod_detail = false;
    this.show_container_detail = false;
  }

  listCluster() {
    this.clusterService.listCluster().subscribe(data => {
      this.clusters = data;
      this.selectClusters = data;
      this.getClusterData();
    }, error => {
      this.loading = false;
    });
  }

  getCluster() {
    this.clusterService.getCluster(this.dashboardSearch.cluster).subscribe(data => {
      this.clusters = [];
      this.clusters.push(data);
      this.getClusterData();
    });
  }

  getClusterData() {
    this.data_init();
    this.dashboardService.getDashboard(this.dashboardSearch.cluster).subscribe(data => {
      this.clusterData = data.data;
      this.restartPods = data.restartPods;
      this.warnContainers = data.warnContainers;
      for (const cd of this.clusterData) {
        const d = JSON.parse(cd);
        this.podCount = this.podCount + d['pods'].length;
        this.namespaceCount = this.namespaceCount + d['namespaces'].length;
        this.deploymentCount = this.deploymentCount + d['deployments'].length;
        this.nodeCount = this.nodeCount + d['nodes'].length;
        this.cpu_total = this.cpu_total + d['cpu_total'];
        this.mem_total = this.mem_total + d['mem_total'];
        this.cpu_usage = this.cpu_usage + d['cpu_usage'];
        this.mem_usage = this.mem_usage + d['mem_usage'];
      }
      if (this.clusterData.length > 0) {
        this.cpu_usage = this.cpu_usage / this.clusterData.length * 100;
        this.mem_usage = this.mem_usage / this.clusterData.length * 100;
      }
      this.loading = false;
    });
  }

  search() {
    this.loading = true;
    if (this.dashboardSearch.cluster === 'all') {
      this.listCluster();
    } else {
      this.getCluster();
    }
  }

  refresh() {
    this.search();
  }

  toPage(url) {
    this.redirect(url);
  }

  redirect(url: string) {
    if (url) {
      const linkUrl = ['kubeOperator', url];
      this.router.navigate(linkUrl);
    }
  }
}
