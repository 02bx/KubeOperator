import os
import uuid
import yaml
import logging
from django.conf import settings
from django.db import models
from django.utils.translation import ugettext_lazy as _
from ansible_api.models.inventory import BaseHost
from ansible_api.models.utils import name_validator
from ansible_api.tasks import run_im_adhoc
from common.models import JsonTextField
from ansible_api.models import Project, Group, Playbook, AdHoc
from ansible_api.models import Host as Ansible_Host
from ansible_api.models.mixins import (
    AbstractProjectResourceModel, AbstractExecutionModel
)
from .signals import pre_deploy_execution_start, post_deploy_execution_start

__all__ = ['Package', 'Cluster', 'Node', 'Role', 'DeployExecution', 'Host', 'HostInfo', 'Setting']
logger = logging.getLogger(__name__)


# 离线包的model
class Package(models.Model):
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.CharField(max_length=20, unique=True, verbose_name=_('Name'))
    meta = JsonTextField(blank=True, null=True, verbose_name=_('Meta'))
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))
    packages_dir = os.path.join(settings.BASE_DIR, 'resource', 'packages')

    def __str__(self):
        return self.name

    class Meta:
        verbose_name = _('Package')

    @property
    def path(self):
        return os.path.join(self.packages_dir, self.name)

    @classmethod
    def lookup(cls):
        for d in os.listdir(cls.packages_dir):
            full_path = os.path.join(cls.packages_dir, d)
            meta_path = os.path.join(full_path, 'meta.yml')
            if not os.path.isdir(full_path) or not os.path.isfile(meta_path):
                continue
            with open(meta_path) as f:
                metadata = yaml.load(f)
            defaults = {'name': d, 'meta': metadata}
            cls.objects.update_or_create(defaults=defaults, name=d)


class Setting(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    key = models.CharField(max_length=128, blank=False)
    value = models.CharField(max_length=255, blank=True, default=None, null=True)
    name = models.CharField(max_length=128, blank=False)
    helper = models.CharField(max_length=255, blank=True)
    order = models.IntegerField(default=0)


class Cluster(Project):
    package = models.ForeignKey("Package", null=True, on_delete=models.SET_NULL)
    template = models.CharField(max_length=64, blank=True, default='')
    current_task_id = models.CharField(max_length=128, blank=True, default='')
    is_super = models.BooleanField(default=False)

    @property
    def state(self):
        if not self.current_task_id is "":
            c = DeployExecution.objects.filter(id=self.current_task_id).first()
            return c.state
        return None

    def create_roles(self):
        _roles = {}
        for role in self.package.meta.get('roles', []):
            _roles[role['name']] = role
        template = None
        for tmp in self.package.meta.get('templates', []):
            if tmp['name'] == self.template:
                template = tmp
                break

        for role in template.get('roles', []):
            _roles[role['name']] = role
        roles_data = [role for role in _roles.values()]

        children_data = {}
        for data in roles_data:
            children_data[data['name']] = data.pop('children', [])
            Role.objects.update_or_create(defaults=data, name=data['name'])

        for name, children_name in children_data.items():
            try:
                role = Role.objects.get(name=name)
                children = Role.objects.filter(name__in=children_name)
                role.children.set(children)
            except Role.DoesNotExist:
                pass

        ose_role = Role.objects.get(name='OSEv3')
        private_var = template['private_vars']
        role_vars = ose_role.vars
        role_vars.update(private_var)
        ose_role.vars = role_vars
        ose_role.save()

    def create_node_localhost(self):
        Node.objects.create(
            name="localhost", vars={"ansible_connection": "local"},
            project=self, meta={"hidden": True}
        )

    def create_install_playbooks(self):
        for data in self.package.meta.get('install_playbooks', []):
            url = 'file:///{}'.format(os.path.join(self.package.path))
            Playbook.objects.create(
                name=data['name'], alias=data['alias'],
                type=Playbook.TYPE_LOCAL, url=url, project=self,
            )

    def create_upgrade_playbooks(self):
        for data in self.package.meta.get('upgrade_playbooks', []):
            url = 'file:///{}'.format(os.path.join(self.package.path))
            Playbook.objects.create(
                name=data['name'], alias=data['alias'],
                type=Playbook.TYPE_LOCAL, url=url, project=self,
            )

    def create_uninstall_playbooks(self):
        for data in self.package.meta.get('uninstall_playbooks', []):
            url = 'file:///{}'.format(os.path.join(self.package.path))
            Playbook.objects.create(
                name=data['name'], alias=data['alias'],
                type=Playbook.TYPE_LOCAL, url=url, project=self,
            )

    def on_cluster_create(self):
        self.change_to()
        self.create_roles()
        self.create_node_localhost()
        self.create_install_playbooks()
        self.create_upgrade_playbooks()
        self.create_uninstall_playbooks()

    def configs(self, tp='list'):
        self.change_to()
        role = Role.objects.get(name='OSEv3')
        configs = role.vars
        if tp == 'list':
            configs = [{'key': k, 'value': v} for k, v in configs.items()]
        return configs

    def set_config(self, k, v):
        self.change_to()
        role = Role.objects.select_for_update().get(name='OSEv3')
        _vars = role.vars
        if isinstance(v, str):
            v = v.strip()
        _vars[k] = v
        role.vars = _vars
        role.save()

    def get_config(self, k):
        v = self.configs(tp='dict').get(k)
        return {'key': k, 'value': v}

    def del_config(self, k):
        self.change_to()
        role = Role.objects.get(name='OSEv3')
        _vars = role.vars
        _vars.pop(k, None)
        role.vars = _vars
        role.save()


class ClusterConfig(models.Model):
    key = models.CharField(max_length=1024)
    value = JsonTextField()

    class Meta:
        abstract = True

    target = models.ForeignKey('ansible_api.Project', related_name="target", on_delete=models.CASCADE)


class Host(BaseHost):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    node = models.ForeignKey('Node', default=None, null=True, related_name='node',
                             on_delete=models.SET_NULL)
    name = models.CharField(max_length=128, validators=[name_validator], unique=True)

    @property
    def cluster(self):
        if self.node is not None:
            return self.node.project.name
        else:
            return '无'

    @property
    def info(self):
        return self.infos.all().latest()

    def gather_info(self):
        info = HostInfo.objects.create(host_id=self.id)
        info.gather_info()

    class Meta:
        ordering = ('name',)


class HostInfo(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    memory = models.fields.BigIntegerField(default=0)
    os = models.fields.CharField(max_length=128, default="")
    os_version = models.fields.CharField(max_length=128, default="")
    cpu_core = models.fields.IntegerField(default=0)
    host = models.ForeignKey('Host', on_delete=models.CASCADE, null=True, related_name='infos')
    volumes = models.ManyToManyField('Volume')
    date_created = models.DateTimeField(auto_now_add=True)

    class Meta:
        get_latest_by = 'date_created'

    def gather_info(self):
        host = self.host
        hosts = [self.host.__dict__]
        result = run_im_adhoc(adhoc_data={'pattern': host.name, 'module': 'setup'},
                              inventory_data={'hosts': hosts, 'vars': {}})
        if not result.get('summary', {}).get('success', False):
            raise Exception("get os info failed!")
        else:
            facts = result["raw"]["ok"][host.name]["setup"]["ansible_facts"]
            self.memory = facts["ansible_memtotal_mb"]
            self.cpu_core = facts["ansible_processor_count"]
            self.os = facts["ansible_distribution"]
            self.os_version = facts["ansible_distribution_version"]
            self.save()
            devices = facts["ansible_devices"]
            volumes = []
            for name in devices:
                if not name.startswith('dm'):
                    volume = Volume(name='/dev/' + name)
                    volume.size = devices[name]['size']
                    volume.save()
                    volumes.append(volume)
            self.volumes.set(volumes)


class Volume(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4)
    name = models.CharField(max_length=128)
    size = models.CharField(max_length=16)
    blank = models.BooleanField(default=False)

    class Meta:
        ordering = ('size',)


class Node(Ansible_Host):
    host = models.ForeignKey('Host', related_name='host', default=None, null=True, on_delete=models.CASCADE)

    @property
    def roles(self):
        return self.groups

    @property
    def host_memory(self):
        return self.host.info.memory

    @property
    def host_cpu_core(self):
        return self.host.info.cpu_core

    @property
    def host_os(self):
        return self.host.info.os

    @property
    def host_os_version(self):
        return self.host.info.os_version

    @property
    def host_volumes(self):
        return self.host.info.volumes

    @roles.setter
    def roles(self, value):
        self.groups.set(value)

    def on_node_save(self):
        self.ip = self.host.ip
        self.username = self.host.username
        self.password = self.host.password
        self.private_key = self.host.private_key
        self.host.node_id = self.id
        self.host.save()
        self.save()

    def add_vars(self, _vars):
        __vars = {k: v for k, v in self.vars.items()}
        __vars.update(_vars)
        if self.vars != __vars:
            self.vars = __vars
            self.save()

    def remove_var(self, key):
        __vars = self.vars
        if key in __vars:
            del __vars[key]
            self.vars = __vars
            self.save()

    def get_var(self, key, default):
        return self.vars.get(key, default)


class Role(Group):
    class Meta:
        proxy = True

    @property
    def nodes(self):
        return self.hosts

    @nodes.setter
    def nodes(self, value):
        self.hosts.set(value)

    def __str__(self):
        return "%s %s" % (self.project, self.name)


class DeployExecution(AbstractProjectResourceModel, AbstractExecutionModel):
    OPERATION_INSTALL = 'install'
    OPERATION_UPGRADE = 'upgrade'
    OPERATION_UNINSTALL = 'uninstall'
    OPERATION_UPDATE = 'update'

    OPERATION_CHOICES = (
        (OPERATION_INSTALL, _('install')),
        (OPERATION_UPGRADE, _('upgrade')),
        (OPERATION_UNINSTALL, _('uninstall')),
        (OPERATION_UPDATE, _('update')),
    )

    project = models.ForeignKey('ansible_api.Project', on_delete=models.CASCADE)
    operation = models.CharField(max_length=128, choices=OPERATION_CHOICES, blank=True, default=OPERATION_INSTALL)
    current_task = models.CharField(max_length=128, null=True, blank=True, default=None)
    progress = models.FloatField(max_length=64, null=True, blank=True, default=0.0)

    def save(self, force_insert=False, force_update=False, using=None, update_fields=None):
        super().save(force_insert, force_update, using, update_fields)
        Cluster.objects.filter(id=self.project.id).update(current_task_id=self.id)

    def start(self):
        local_hostname = Setting.objects.filter(key="local_hostname").first()
        registry_hostname = Setting.objects.filter(key="registry_hostname").first()

        result = {"raw": {}, "summary": {}}
        pre_deploy_execution_start.send(self.__class__, execution=self)
        playbooks = self.project.playbook_set.filter(name__endswith='-' + self.operation).order_by('name')
        try:
            for index, playbook in enumerate(playbooks):
                print("\n>>> Start run {} ".format(playbook.name))
                self.update_task(playbook.name)
                _result = playbook.execute(extra_vars={
                    "cluster_name": self.project.name,
                    "registry_hostname": registry_hostname.value,
                    "local_hostname": local_hostname.value
                })
                result["summary"].update(_result["summary"])
                if not _result.get('summary', {}).get('success', False):
                    break
                else:
                    self.update_progress((index + 1) / len(playbooks))
                if len(playbooks) == index + 1:
                    self.update_task('Finish')
        except Exception as e:
            logger.error(e, exc_info=True)
            result['summary'] = {'error': 'Unexpect error occur: {}'.format(e)}
        post_deploy_execution_start.send(self.__class__, execution=self, result=result)
        return result

    def update_task(self, task):

        self.current_task = task
        self.save()

    def update_progress(self, progress):
        self.progress = progress
        self.save()

    def to_json(self):
        return {
            'id': self.id.__str__(),
            'progress': self.progress,
            'current_task': self.current_task,
            'operation': self.operation,
            'state': self.state
        }

    class Meta:
        get_latest_by = 'date_created'
        ordering = ('-date_created',)
