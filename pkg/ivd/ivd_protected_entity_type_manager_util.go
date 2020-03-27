package ivd

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/cns"
	cnstypes "github.com/vmware/govmomi/cns/types"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25"
	vim25types "github.com/vmware/govmomi/vim25/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"
)

func findAllDatastores(ctx context.Context, client *vim25.Client) ([]vim25types.ManagedObjectReference, error) {
	finder := find.NewFinder(client)
	dss, err := finder.DatastoreList(ctx, "*")
	if err != nil {
		return nil, err
	}

	var dsList []vim25types.ManagedObjectReference
	for _, ds := range dss {
		dsList = append(dsList, ds.Reference())
	}
	if len(dsList) == 0 {
		return nil, errors.New("No datastore can be found")
	}

	return dsList, nil
}

func retrievePlatformInfoFromConfig(ctx context.Context, config *rest.Config, params map[string]interface{}, logger logrus.FieldLogger) error {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.WithError(err).Errorf("Failed to get k8s clientset from the given config: %v", config)
		return err
	}

	ns := "kube-system"
	secretApis := clientset.CoreV1().Secrets(ns)
	vsphere_secret := "vsphere-config-secret"
	// v0.18.0 requires context secret, err := secretApis.Get(ctx, vsphere_secret, metav1.GetOptions{})
	secret, err := secretApis.Get(vsphere_secret, metav1.GetOptions{})

	if err != nil {
		logger.WithError(err).Errorf("Failed to get k8s secret, %s", vsphere_secret)
		return err
	}
	sEnc := string(secret.Data["csi-vsphere.conf"])
	lines := strings.Split(sEnc, "\n")

	for _, line := range lines {
		if strings.Contains(line, "VirtualCenter") {
			parts := strings.Split(line, "\"")
			params["VirtualCenter"] = parts[1]
		} else if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			params[key] = value[1 : len(value)-1]
		}
	}

	return nil
}

func createCnsVolumeWithClusterConfig(ctx context.Context, config *rest.Config, client *govmomi.Client, cnsClient *cns.Client, md metadata, logger logrus.FieldLogger) (string, error) {
	logger.Debugf("createCnsVolumeWithClusterConfig called with args, metadata: %v", md)

	reservedLabelsMap, err := fillInClusterSpecificParams(ctx, config, logger)
	if err != nil {
		logger.WithError(err).Error("Failed at calling fillInClusterSpecificParams")
		return "", err
	}

	// Preparing for the VolumeCreateSpec for the volume provisioning
	logger.Debug("Preparing for the VolumeCreateSpec for the volume provisioning")
	dsList, err := findAllDatastores(ctx, client.Client)
	if err != nil {
		logger.WithError(err).Error("Failed to find any datastore in the underlying vSphere")
		return "", err
	}

	var metadataList []cnstypes.BaseCnsEntityMetadata
	metadata := &cnstypes.CnsKubernetesEntityMetadata{
		CnsEntityMetadata: cnstypes.CnsEntityMetadata{
			EntityName:  md.VirtualStorageObject.Config.Name,
			Labels:      md.ExtendedMetadata,
		},
		EntityType: string(cnstypes.CnsKubernetesEntityTypePV),
	}
	metadataList = append(metadataList, cnstypes.BaseCnsEntityMetadata(metadata))

	var cnsVolumeCreateSpecList []cnstypes.CnsVolumeCreateSpec
	cnsVolumeCreateSpec := cnstypes.CnsVolumeCreateSpec{
		Name:        md.VirtualStorageObject.Config.Name,
		VolumeType: string(cnstypes.CnsVolumeTypeBlock),
		Datastores: dsList,
		Metadata: cnstypes.CnsVolumeMetadata{
			ContainerCluster: cnstypes.CnsContainerCluster{
				ClusterType: string(cnstypes.CnsClusterTypeKubernetes), // hard coded for the moment
				ClusterId:   reservedLabelsMap["cns.containerCluster.clusterId"],
				VSphereUser: reservedLabelsMap["cns.containerCluster.vSphereUser"],
			},
			EntityMetadata: metadataList,
		},
		BackingObjectDetails: &cnstypes.CnsBlockBackingDetails{
			CnsBackingObjectDetails: cnstypes.CnsBackingObjectDetails{
				CapacityInMb: md.VirtualStorageObject.Config.CapacityInMB,
			},
		},
	}

	cnsVolumeCreateSpecList = append(cnsVolumeCreateSpecList, cnsVolumeCreateSpec)
	logger.Debugf("Provisioning volume using the spec: %v", cnsVolumeCreateSpec)

	// provision volume using CNS API
	createTask, err := cnsClient.CreateVolume(ctx, cnsVolumeCreateSpecList)
	if err != nil {
		logger.WithError(err).Errorf("Failed to create volume. Error: %+v", err)
		return "", err
	}
	createTaskInfo, err := cns.GetTaskInfo(ctx, createTask)
	if err != nil {
		logger.WithError(err).Errorf("Failed to create volume. Error: %+v", err)
		return "", err
	}
	createTaskResult, err := cns.GetTaskResult(ctx, createTaskInfo)
	if err != nil {
		logger.WithError(err).Errorf("Failed to create volume. Error: %+v", err)
		return "", err
	}
	if createTaskResult == nil {
		err := errors.New("Empty create task results")
		logger.Error(err.Error())
		return "", err
	}
	createVolumeOperationRes := createTaskResult.GetCnsVolumeOperationResult()
	if createVolumeOperationRes.Fault != nil {
		logger.Errorf("Failed to create volume: fault=%+v", createVolumeOperationRes.Fault)
		return "", errors.New(createVolumeOperationRes.Fault.LocalizedMessage)
	}

	volumeId := createVolumeOperationRes.VolumeId.Id
	logger.Infof("CNS volume, %v, created", volumeId)
	return volumeId, nil
}

func fillInClusterSpecificParams(ctx context.Context, config *rest.Config, logger logrus.FieldLogger) (map[string]string, error) {
	params := make(map[string]interface{})
	err := retrievePlatformInfoFromConfig(ctx, config, params, logger)
	if err != nil {
		logger.WithError(err).Errorf("Failed to retrieve VC config secret: %+v", err)
		return map[string]string{}, err
	}

	clusterId, ok := params["cluster-id"].(string)
	if !ok {
		logger.WithError(err).Errorf("Failed to retrieve cluster id")
		return map[string]string{}, err
	}

	user, ok := params["user"].(string)
	if !ok {
		logger.WithError(err).Errorf("Failed to retrieve vsphere user")
		return map[string]string{}, err
	}
	logger.Debugf("Retrieved cluster id, %v, and vSphere user, %v", clusterId, user)

	// currently, we only pick up two cluster specific labels, cluster-id and vsphere-user.
	// For the following labels,
	//    cns.containerCluster.clusterType -- always "KUBERNETES", and no other type available for the moment
	//    cns.containerCluster.clusterFlavor -- the most recent govmomi version doesn't provide field to set the cluster flavor
	//    others are not cluster specfic, but cns specific
	reservedLabelsMap := map[string]string {
		//"cns.containerCluster.clusterFlavor",
		//"cns.containerCluster.clusterType",
		//"cns.k8s.pv.name",
		//"cns.tag",
		//"cns.version",
		"cns.containerCluster.clusterId": clusterId,
		"cns.containerCluster.vSphereUser": user,
	}

	return reservedLabelsMap, nil
}

func FilterLabelsFromMetadataForVslmAPIs(ctx context.Context, md metadata, logger logrus.FieldLogger) (metadata, error) {
	var kvsList []vim25types.KeyValue

	logger.Debugf("labels of CNS volume before filtering: %v", md.ExtendedMetadata)

	// Retrieving cluster id and vSphere user
	logger.Debug("Retrieving cluster id and vSphere user required by provisioning volume")
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.WithError(err).Error("Failed to get k8s inClusterConfig")
		return metadata{}, err
	}

	reservedLabelsMap, err := fillInClusterSpecificParams(ctx, config, logger)
	if err != nil {
		logger.WithError(err).Error("Failed at calling fillInClusterSpecificParams")
		return metadata{}, err
	}

	for key, value := range reservedLabelsMap {
		kvsList = append(kvsList, vim25types.KeyValue {
			Key: key,
			Value: value,
		})
	}

	for _, label := range md.ExtendedMetadata {
		value, ok := reservedLabelsMap[label.Key]
		if !ok {
			value = label.Value
		}
		kvsList = append(kvsList, vim25types.KeyValue {
			Key: label.Key,
			Value: value,
		})
	}
	md.ExtendedMetadata = kvsList

	logger.Debugf("labels of CNS volume after filtering: %v", md.ExtendedMetadata)

	return md, nil
}

func FilterLabelsFromMetadataForCnsAPIs(md metadata, prefix string, logger logrus.FieldLogger) metadata {
	// prefix: cns.containerCluster
	var kvsList []vim25types.KeyValue

	logger.Debugf("labels of CNS volume before filtering ones with certain prefix, %v: %v", prefix, md.ExtendedMetadata)

	for _, label := range md.ExtendedMetadata {
		if !strings.HasPrefix(label.Key, prefix) {
			kvsList = append(kvsList, vim25types.KeyValue {
				Key: label.Key,
				Value: label.Value,
			})
		}
	}
	md.ExtendedMetadata = kvsList

	logger.Debugf("labels of CNS volume after filtering ones with certain prefix, %v: %v", prefix, md.ExtendedMetadata)

	return md
}

func CreateCnsVolumeInCluster(ctx context.Context, client *govmomi.Client, cnsClient *cns.Client, md metadata, logger logrus.FieldLogger) (vim25types.ID, error) {
	logger.Infof("CreateCnsVolumeInCluster called with args, metadata: %v", md)

	// Retrieving cluster id and vSphere user
	logger.Debug("Retrieving cluster id and vSphere user required by provisioning volume")
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.WithError(err).Error("Failed to get k8s inClusterConfig")
		return vim25types.ID{}, err
	}

	volumeId, err := createCnsVolumeWithClusterConfig(ctx, config, client, cnsClient, md, logger)
	if err != nil {
		logger.WithError(err).Error("Failed to call createCnsVolumeWithClusterConfig")
		return vim25types.ID{}, err
	}

	return NewIDFromString(volumeId), nil
}