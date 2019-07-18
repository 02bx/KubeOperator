import logging
import os
import yaml
from django.db import models
from ansible_api.models import Project, Playbook
from common import models as common_models
from fit2ansible import settings
from django.utils.translation import ugettext_lazy as _
from ansible_api.models import Host as Ansible_Host

from kubeops_api.models.role import Role

logger = logging.getLogger(__name__)


class StorageTemplate(models.Model):
    name = models.CharField(max_length=128, unique=True, verbose_name='名称')
    meta = common_models.JsonDictTextField(blank=True, null=True)
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))
    templates_dir = os.path.join(settings.BASE_DIR, 'resource', 'storage')

    def __str__(self):
        return self.name

    @property
    def path(self):
        return os.path.join(self.templates_dir, self.name)

    @classmethod
    def lookup(cls):
        for d in os.listdir(cls.templates_dir):
            full_path = os.path.join(cls.templates_dir, d)
            meta_path = os.path.join(full_path, 'meta.yml')
            if not os.path.isdir(full_path) or not os.path.isfile(meta_path):
                continue
            with open(meta_path) as f:
                metadata = yaml.load(f)
            defaults = {'name': d, 'meta': metadata}
            cls.objects.update_or_create(defaults=defaults, name=d)


class StorageNode(Ansible_Host):
    STORATE_NODE_STATUS_UNKNOWN = 'UNKNOWN'
    STORATE_NODE_STATUS_RUNNING = 'RUNNING'
    STORATE_NODE_STATUS_ERROR = 'ERROR'
    STORATE_NODE_STATUS_CHOICES = (
        (STORATE_NODE_STATUS_RUNNING, 'running'),
        (STORATE_NODE_STATUS_UNKNOWN, 'unknown'),
        (STORATE_NODE_STATUS_ERROR, 'error'),
    )
    status = models.CharField(max_length=128, choices=STORATE_NODE_STATUS_CHOICES, default=STORATE_NODE_STATUS_UNKNOWN)
    message = models.TextField(default=None, null=True)

    @property
    def roles(self):
        return self.groups

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


class Storage(Project):
    STORATE_STATUS_UNKNOWN = 'UNKNOWN'
    STORATE_STATUS_RUNNING = 'RUNNING'
    STORATE_STATUS_ERROR = 'ERROR'
    STORATE_STATUS_CHECKING = 'CHECKING'

    STORATE_STATUS_CHOICES = (
        (STORATE_STATUS_RUNNING, 'running'),
        (STORATE_STATUS_UNKNOWN, 'unknown'),
        (STORATE_STATUS_ERROR, 'error'),
        (STORATE_STATUS_CHECKING, 'checking'),
    )
    template = models.ForeignKey("StorageTemplate", null=True, on_delete=models.SET_NULL)
    vars = common_models.JsonDictTextField(default={}, blank=True, null=True, verbose_name=_('Vars'))
    status = models.CharField(max_length=128, choices=STORATE_STATUS_CHOICES, default=STORATE_STATUS_UNKNOWN)
    nodes = models.ManyToManyField('StorageNode')

    def create_playbooks(self):
        for playbook in self.template.meta.get('playbooks', []):
            url = 'file:///{}'.format(os.path.join(self.template.path))
            Playbook.objects.create(
                name=playbook['name'], alias=playbook['alias'],
                type=Playbook.TYPE_LOCAL, url=url, project=self
            )

    def create_roles(self):
        _roles = {}
        for role in self.template.meta.get('roles', []):
            _roles[role['name']] = role
        roles_data = [role for role in _roles.values()]
        for data in roles_data:
            Role.objects.update_or_create(defaults=data, name=data['name'])

    def set_vars(self):
        self.vars = self.template.meta.get('vars', {})

    def config(self, k, v):
        if isinstance(v, str):
            v = v.strip()
        self.vars[k] = v
        self.save()

    def on_storage_create(self):
        self.change_to()
        # self.create_roles()
        # self.create_playbooks()
        self.set_vars()
