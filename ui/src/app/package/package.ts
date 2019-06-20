export class Package {
  id: string;
  name: string;
  meta: PackageMeta;
  date_created: string;

  getLogo() {
    let logo = null;
    const path = 'assets/images';
    switch (this.meta.resource) {
      case 'kubernetes':
        logo = path + 'logo-k8s.png';
        break;
      case 'okd':
        logo = path + 'logo-okd.png';
        break;
    }
    return logo;
  }
}

export class PackageMeta {
  resource: string;
  version: string;
  templates: Template[];
}

export class Config {
  name: string;
  alias: string;
  type: string;
  options: Option[];
  value: string;
  default: any;
  help_text: string;
  display: boolean;
}

export class Option {
  name: string;
  alias: string;
}

export class Role {
  name: string;
  meta: RoleMeta;
}

export class Requires {
  nodes_require: any[];
  volumes_require: Require[];
  device_require: Require[];
}

export class Require {
  name: string;
  verbose: string;
  minimal: number;
  excellent: number;
  unit: string;
  comment: string;
}

export class NodeVars {
  name: string;
  template: string;
  verbose: string;
  comment: string;
  type: string;
  options: any;
  placeholder: string;
  require: true;
}


export class RoleMeta {
  hidden: boolean;
  allow_os: Os[];
  requires: Requires;
  node_vars: NodeVars[];


}


export class Os {
  name: string;
  version: string[];
}

export class Template {
  name: string;
  roles: Role[] = [];
  private_config: Config[] = [];
  portals: Portal[];
  comment: string;
}

export class Portal {
  name: string;
  redirect: string;
}
