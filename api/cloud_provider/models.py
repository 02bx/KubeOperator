import os
import uuid

import yaml
from django.db import models

# Create your models here.
from common import models as common_models
from fit2ansible import settings
from django.utils.translation import ugettext_lazy as _


class CloudProviderTemplate(models.Model):
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.CharField(max_length=20, unique=True, verbose_name=_('Name'))
    meta = common_models.JsonTextField(blank=True, null=True, verbose_name=_('Meta'))
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))
    template_dir = os.path.join(settings.BASE_DIR, 'resource', 'clouds')

    @property
    def path(self):
        return os.path.join(self.template_dir, self.name)

    @classmethod
    def lookup(cls):
        for d in os.listdir(cls.template_dir):
            full_path = os.path.join(cls.template_dir, d)
            meta_path = os.path.join(full_path, 'meta.yml')
            if not os.path.isdir(full_path) or not os.path.isfile(meta_path):
                continue
            with open(meta_path) as f:
                metadata = yaml.load(f)
            defaults = {'name': d, 'meta': metadata}
            cls.objects.update_or_create(defaults=defaults, name=d)


class Region(models.Model):
    id = models.UUIDField(default=uuid.uuid4, primary_key=True)
    name = models.CharField(max_length=20, unique=True, verbose_name=_('Name'))
    date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))
    template = models.ForeignKey('CloudProviderTemplate', on_delete=models.SET_NULL, null=True)
    cloud_region = models.CharField(max_length=128, null=True, default=None)
    vars = common_models.JsonDictTextField(default={})
    comment = models.CharField(max_length=128, blank=True, null=True, verbose_name=_("Comment"))

    def set_vars(self):
        meta = self.template.meta.get('meta', None)
        if meta:
            _vars = meta.get('vars', {})
            self.vars.update(_vars)
            self.save()

    def on_region_create(self):
        self.set_vars()


# class Zone(models.Model):
#     id = models.UUIDField(default=uuid.uuid4, primary_key=True)
#     name = models.CharField(max_length=20, unique=True, verbose_name=_('Name'))
#     date_created = models.DateTimeField(auto_now_add=True, verbose_name=_('Date created'))
#     comment = models.CharField(max_length=128, blank=True, null=True, verbose_name=_("Comment"))
#     vars = common_models.JsonDictTextField(default={})
#     region = models.ForeignKey('Region', on_delete=models.CASCADE, null=True)
#     cloud_zone = models.CharField(max_length=128, null=True, default=None)
