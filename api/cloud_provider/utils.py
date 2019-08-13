import os
import zipfile
from download import download
from jinja2 import FileSystemLoader, Environment

from fit2ansible.settings import TERRAFORM_DIR


def generate_terraform_file(target_path, cloud_path, vars):
    terraform_path = os.path.join(cloud_path, "terraform")
    lorder = FileSystemLoader(terraform_path)
    env = Environment(loader=lorder)
    _template = env.get_template("terraform.tf.j2")
    result = _template.render(vars)
    if not os.path.exists(target_path):
        os.makedirs(target_path)
    file = os.path.join(target_path, 'main.tf')
    with open(file, 'w') as f:
        f.write(result)
        return os.path.dirname(file)


def create_terrafrom_working_dir(cluster_name):
    if not os.path.exists(TERRAFORM_DIR):
        os.makedirs(TERRAFORM_DIR)
    cluster_dir = os.path.join(TERRAFORM_DIR, cluster_name)
    if not os.path.exists(cluster_dir):
        os.mkdir(cluster_dir)
    return cluster_dir


def download_plugins(url, target):
    f = download_file(url, target)
    unzip_plugin(f)


def download_file(url, target):
    basename = os.path.basename(url)
    target = os.path.join(target, basename)
    print(target)
    return download(url, target, progressbar=True)


def unzip_plugin(f):
    file_zip = zipfile.ZipFile(f, 'r')
    for file in file_zip.namelist():
        file_zip.extract(file, r'.')
    file_zip.close()
    os.remove(f)
