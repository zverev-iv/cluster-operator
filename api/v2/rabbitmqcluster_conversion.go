package v2

import (
	"github.com/rabbitmq/cluster-operator/api/v1beta1"
	"github.com/rabbitmq/cluster-operator/internal/status"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts RabbitmqCluster v2 to the Hub version v1beta1.
func (src *RabbitmqCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.RabbitmqCluster)

	// copy object meta
	dst.ObjectMeta = src.ObjectMeta

	// conversion of Service to ClientService
	if src.Spec.ClientService.Type != "" {
		dst.Spec.Service.Type = src.Spec.ClientService.Type
	}
	if src.Spec.ClientService.Annotations != nil {
		dst.Spec.Service.Annotations = src.Spec.ClientService.Annotations
	}

	// copy the rest of the spec
	dst.Spec.Image = src.Spec.Image
	dst.Spec.Replicas = src.Spec.Replicas
	dst.Spec.Rabbitmq.AdvancedConfig = src.Spec.Rabbitmq.AdvancedConfig
	dst.Spec.Rabbitmq.AdditionalConfig = src.Spec.Rabbitmq.AdditionalConfig
	dst.Spec.Rabbitmq.EnvConfig = src.Spec.Rabbitmq.EnvConfig
	dst.Spec.ImagePullSecret = src.Spec.ImagePullSecret
	copyPluginsTov1beta1(dst, src.Spec.Rabbitmq.AdditionalPlugins)

	// copy status
	if src.Status.Admin != nil {
		dst.Status.Admin = &v1beta1.RabbitmqClusterAdmin{
			SecretReference: &v1beta1.RabbitmqClusterSecretReference{
				Name:      src.Status.Admin.SecretReference.Name,
				Namespace: src.Status.Admin.SecretReference.Namespace,
				Keys:      src.Status.Admin.SecretReference.Keys,
			},
			ServiceReference: &v1beta1.RabbitmqClusterServiceReference{
				Name:      src.Status.Admin.ServiceReference.Name,
				Namespace: src.Status.Admin.ServiceReference.Namespace,
			},
		}
	}

	if len(src.Status.Conditions) == 0 {
		dst.Status.Conditions = []status.RabbitmqClusterCondition{}
	} else {
		dst.Status.Conditions = src.Status.Conditions
	}

	dst.Status.ClusterStatus = src.Status.ClusterStatus

	return nil
}

// ConvertFrom converts RabbitmqCluster v1beta1 to this version version v2.
func (dst *RabbitmqCluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.RabbitmqCluster)

	// copy object meta
	dst.ObjectMeta = src.ObjectMeta

	// conversion of ClientService to Service
	if src.Spec.Service.Type != "" {
		dst.Spec.ClientService.Type = src.Spec.Service.Type
	}
	if src.Spec.Service.Annotations != nil {
		dst.Spec.ClientService.Annotations = src.Spec.Service.Annotations
	}

	// copy the rest of the spec
	dst.Spec.Image = src.Spec.Image
	dst.Spec.Replicas = src.Spec.Replicas
	dst.Spec.Rabbitmq.AdvancedConfig = src.Spec.Rabbitmq.AdvancedConfig
	dst.Spec.Rabbitmq.AdditionalConfig = src.Spec.Rabbitmq.AdditionalConfig
	dst.Spec.Rabbitmq.EnvConfig = src.Spec.Rabbitmq.EnvConfig
	dst.Spec.ImagePullSecret = src.Spec.ImagePullSecret
	copyPluginsTov2(dst, src.Spec.Rabbitmq.AdditionalPlugins)

	// copy status
	if src.Status.Admin != nil {
		dst.Status.Admin = &RabbitmqClusterAdmin{
			SecretReference: &RabbitmqClusterSecretReference{
				Name:      src.Status.Admin.SecretReference.Name,
				Namespace: src.Status.Admin.SecretReference.Namespace,
				Keys:      src.Status.Admin.SecretReference.Keys,
			},
			ServiceReference: &RabbitmqClusterServiceReference{
				Name:      src.Status.Admin.ServiceReference.Name,
				Namespace: src.Status.Admin.ServiceReference.Namespace,
			},
		}
	}

	if len(src.Status.Conditions) == 0 {
		dst.Status.Conditions = []status.RabbitmqClusterCondition{}
	} else {
		dst.Status.Conditions = src.Status.Conditions
	}

	dst.Status.ClusterStatus = src.Status.ClusterStatus

	return nil
}

func copyPluginsTov2(dst *RabbitmqCluster, src []v1beta1.Plugin) {
	if len(src) == 0 {
		return
	}
	dst.Spec.Rabbitmq.AdditionalPlugins = make([]Plugin, len(src))
	for i, plugin := range src {
		dst.Spec.Rabbitmq.AdditionalPlugins[i] = Plugin(string(plugin))
	}

}

func copyPluginsTov1beta1(dst *v1beta1.RabbitmqCluster, src []Plugin) {
	if len(src) == 0 {
		return
	}
	dst.Spec.Rabbitmq.AdditionalPlugins = make([]v1beta1.Plugin, len(src))
	for i, plugin := range src {
		dst.Spec.Rabbitmq.AdditionalPlugins[i] = v1beta1.Plugin(string(plugin))
	}

}
