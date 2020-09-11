# Multiversion conversion with RabbitmqCluster

## Experiment

Operator supports 2 versions: RabbitmqCluster `v1beta1` and a new version `v2` which has `spec.service` renamed to `spec.clientService`.

### What I did

Steps before writing version conversion.
I created a placeholder mutating and validating webhook to make sure that the wiring of webhooks and cert manager is correct before writing conversion.

1. Install cert manager `kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.1/cert-manager.yaml`
1. Generated a brand new kubebuilder project with api RabbitmqCluster `v1beta1`. There are manifests required for implementing a webhook that's not in our project. The new project was used to copy over some manifests files. Everything under `config/webhook` and `config/certmanager` was copied over. There are also kustomize variables needed to be copied in file `config/default/base/kustomization.yaml`.
1. Created a validating and a mutating webhooks for version `v1beta1`. `kubebuilder create webhook --group rabbitmq.com --version v1beta1 --kind RabbitmqCluster --defaulting --programmatic-validation`
1. Enable all webhook and cert manager related kustomize files. It's mostly uncommenting sections in `config/default/base/kustomization.yaml`
1. Add webbook as controller-gen options so that it creates webbook manifests to register both the mutating and the validating webhook.
   `controller-gen $(CRD_OPTIONS) rbac:roleName=operator-role webhook paths="./api/...;./controllers/..." output:crd:artifacts:config=config/crd/bases`
1. `makenv deploy-dev` succeeds. Operator failed to come up because it does not have permission to bind on port `443`, which is the default port, and also a privilege port. So I changed the webhook service that `targetPort` is set to `9443` and customized the webhook port in `api/v1beta1/rabbitmqcluster_webhook.go`.
1. `makenv deploy-dev` now succeeds. However, webhooks are all failing and I can't create RabbitmqCluster. It's failing before kubebuilder generated manifests on webhooks and cert manager have errors. To fix the issue, I had to: made sure that cert name and cert namespace were using the correct variables defined in vars section of `config/default/base/kustomization.yaml` (correct values are `CERTICIFATE_NAME` and `CERTIFICATE_NAMESPACE`).

After making sure that the wiring between certmanager and webhooks are working. I moved on to writing the conversion webhooks.

1. Generate a new api version `kubebuilder create api --group rabbitmq.com version v2 --kind RabbitmqCluster`
1. Update group name in `api/v2/groupversion_info.go` to `rabbitmq.com`. By default, the group name is constructed as "group name + domain name", which is `rabbitmq.com.rabbitmq.com` for our apis.
1. Copied over RabbitmqCluster definitions from `v1beta1` to the new version. Rename `spec.service` to `spec.clientService`
1. Set `v1beta1` as the storage version by adding annotations `+kubebuilder:storageversion`.
1. Update makefile to support multi version CRD generation by removing `trivialVersions=true` from crd options.
1. `make manifests generate` and then `make install`. `make install` fails because it runs `k apply` which adds the entire last applied manifests to the annotations, and our crd definitions exceeds the 1MB file limit. This is fix by deleting the crd manually, and use `k create -f` in the makefile. Not an idea solution, need to look into alternative in the future.
1. Create conversion file in `api/v2` and in `api/v1beta1`. Where `v1beta1` is hub version, and `v2`is the spoke version (manually writing conversion funtions, there will be later documented steps on using conversion-gen)
1. Create RabbitmqCluster now fails with webhook error "certificate signed by unknown authority". This is because ca injection are not done correctly. I needed to fix annotations about ca injection in all files to `cert-manager.io/inject-ca-from: rabbitmq-system/rabbitmq-cluster-serving-cert`. Kustomize vars were not templating values..
1. Destroy and `makenv deploy-dev` again. I can now create RabbitmqCluster `v1beta1` and `v2`.

Steps for using conversion-gen

1. Delete all conversion logic in `api/v2/rabbitmqcluster_conversion.go`.
1. Install the tool by adding it to tools/tools.go `_ "k8s.io/code-generator/cmd/conversion-gen"`.
1. Create doc.go in `api/v2/`.
1. Add conversion-gen annotation `// +k8s:conversion-gen=github.com/rabbitmq/cluster-operator/api/v1beta1` to `api/v2/doc.go`.
1. Add conversion-gen annotation `// +k8s:conversion-gen=false` to skip fields that need manually conversion.
1. Use conversion-gen to generate conversion file in `make generate`. This command can take several minutes to finish.

```
conversion-gen \
	--input-dirs=./api/v2 \ # directory for it to scan for annotations
	--output-file-base=zz_generated.conversion \ # name of the output generated conversion file
	--output-base "." \ # it does not output the file without a output-base specified
	--go-header-file=./hack/boilerplate.go.txt # without specified go-header-file, conversion-gen tries to find an non-existing file in its repo
```

1. Use the generated conversion methods in `api/v2/rabbitmqcluster_conversion.go`.
1. Manually adds conversion between `spec.service` and `spec.clientService` in both `convertTo` and `convertFrom`.
1. Initialize `status.conditions` in both `convertTo` and `convertFrom` for creation. `status.conditions` cannot be nil.
1. PROFIT

### Reference documentations

1. [conversion-gen](https://godoc.org/k8s.io/code-generator/cmd/conversion-gen)
1. [install certmanager](https://cert-manager.io/docs/installation/kubernetes/)
1. [kubebuilder book](https://book.kubebuilder.io/multiversion-tutorial/api-changes.html)
1. [certmanager webhook issue annotations](https://github.com/jetstack/cert-manager/issues/2920#issuecomment-658779302)
1. [kubebuilder issue on supporting conversion-gen](https://github.com/kubernetes-sigs/kubebuilder/issues/1529#issuecomment-656359330)

## Documentation

Other operators with multi version conversions that I used as references.

- [cluster-api](https://github.com/kubernetes-sigs/cluster-api/tree/master/api/v1alpha3)
- [kubeflow](https://github.com/kubeflow/kubeflow/tree/master/components/notebook-controller/api)
- [certmanager](https://github.com/jetstack/cert-manager/tree/66d45afcdb3d7b3eb06a445916fd48b045d9e218/pkg/internal/apis/meta/v1)
