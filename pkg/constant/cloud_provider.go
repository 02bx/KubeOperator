package constant

const (
	OpenStack                = "OpenStack"
	OpenStackImageName       = "kubeoperator_centos_7.6.1810"
	OpenStackImageDiskFormat = "qcow2"
	OpenStackImageVMDkPath   = "/terraform/images/openstack/kubeoperator_centos_7.6.1810-1.qcow2"
	VSphere                  = "vSphere"
	VSphereImageName         = "kubeoperator_centos_7.6.1810"
	VSphereImageVMDkPath     = "/data/iso/vsphere/kubeoperator_centos_7.6.1810-1.vmdk"
	VSphereImageOvfPath      = "/data/iso/vsphere/kubeoperator_centos_7.6.1810.ovf"
	VSphereFolder            = "kubeoperator"
	ImageDefaultPassword     = "KubeOperator@2019"
	ImageCredentialName      = "kubeoperator"
	ImageUserName            = "root"
	ImagePasswordType        = "password"
)
