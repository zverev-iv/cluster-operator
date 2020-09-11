package v2

import (
	"github.com/rabbitmq/cluster-operator/api/v1beta1"
	"github.com/rabbitmq/cluster-operator/internal/status"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts RabbitmqCluster v2 to the Hub version v1beta1.
func (src *RabbitmqCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.RabbitmqCluster)

	if err := Convert_v2_RabbitmqCluster_To_v1beta1_RabbitmqCluster(src, dst, nil); err != nil {
		return err
	}

	// manual conversion of Service to ClientService
	if src.Spec.ClientService.Type != "" {
		dst.Spec.Service.Type = src.Spec.ClientService.Type
	}
	if src.Spec.ClientService.Annotations != nil {
		dst.Spec.Service.Annotations = src.Spec.ClientService.Annotations
	}

	// initialize empty status.Conditions
	// status.Conditions cannot be nil
	if len(dst.Status.Conditions) == 0 {
		dst.Status.Conditions = []status.RabbitmqClusterCondition{}
	}

	return nil
}

// ConvertFrom converts RabbitmqCluster v1beta1 to this version version v2.
func (dst *RabbitmqCluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.RabbitmqCluster)

	if err := Convert_v1beta1_RabbitmqCluster_To_v2_RabbitmqCluster(src, dst, nil); err != nil {
		return err
	}

	// manual conversion of ClientService to Service
	if src.Spec.Service.Type != "" {
		dst.Spec.ClientService.Type = src.Spec.Service.Type
	}
	if src.Spec.Service.Annotations != nil {
		dst.Spec.ClientService.Annotations = src.Spec.Service.Annotations
	}

	// initialize empty status.Conditions
	// status.Conditions cannot be nil
	if len(dst.Status.Conditions) == 0 {
		dst.Status.Conditions = []status.RabbitmqClusterCondition{}
	}

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
