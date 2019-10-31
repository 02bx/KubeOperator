import os
import threading
import uuid

import yaml
from django.db import models
from common.models import JsonTextField
from django.utils.translation import ugettext_lazy as _
from fit2ansible.settings import PACKAGE_DIR, DEV, DEV_PACKAGE_DIR
from kubeops_api.package_manage import *

__all__ = ['Package']


class Package(models.Model):
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.CharField(max_length=20, unique=True, verbose_name=_('Name'))
    meta = JsonTextField(blank=True, null=True, verbose_name=_('Meta'))
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))

    def __str__(self):
        return self.name

    class Meta:
        verbose_name = _('Package')

    def repo_port(self):
        return self.meta['vars']['repo_port']

    @property
    def registry_port(self):
        return self.meta['vars']['registry_port']

    @classmethod
    def lookup(cls):
        packages_dir = PACKAGE_DIR
        if DEV:
            packages_dir = DEV_PACKAGE_DIR
        for d in os.listdir(packages_dir):
            full_path = os.path.join(packages_dir, d)
            meta_path = os.path.join(full_path, 'meta.yml')
            if not os.path.isdir(full_path) or not os.path.isfile(meta_path):
                continue
            with open(meta_path) as f:
                metadata = yaml.load(f)
            defaults = {'name': d, 'meta': metadata}
            instance = cls.objects.update_or_create(defaults=defaults, name=d)[0]
            if not DEV:
                thread = threading.Thread(target=cls.start_container(instance))
                thread.start()

    @classmethod
    def start_container(cls, package):
        if not is_package_container_exists(package.name):
            create_package_container(package)
            return
        if not is_package_container_start(package.name):
            start_package_container(package)
