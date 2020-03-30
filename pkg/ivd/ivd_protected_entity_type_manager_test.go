/*
 * Copyright 2019 the Astrolabe contributors
 * SPDX-License-Identifier: Apache-2.0
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ivd

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi/cns"
	cnstypes "github.com/vmware/govmomi/cns/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/url"
	"os"
	"testing"
)

func TestProtectedEntityTypeManager(t *testing.T) {
	var vcUrl url.URL
	vcUrl.Scheme = "https"
	vcUrl.Host = "10.160.127.39"
	vcUrl.User = url.UserPassword("administrator@vsphere.local", "Admin!23")
	vcUrl.Path = "/sdk"

	t.Logf("%s\n", vcUrl.String())

	ivdPETM, err := NewIVDProtectedEntityTypeManagerFromURL(&vcUrl, "/ivd", true, logrus.New())
	ctx := context.Background()

	pes, err := ivdPETM.GetProtectedEntities(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("# of PEs returned = %d\n", len(pes))
}

func getVcConfigFromParams(params map[string]interface{}) (*url.URL, bool, error) {
	var vcUrl url.URL
	vcUrl.Scheme = "https"
	vcHostStr, ok := params["VirtualCenter"].(string)
	if !ok {
		return nil, false, errors.New("Missing vcHost param")
	}
	vcHostPortStr, ok := params["port"].(string)
	if !ok {
		return nil, false, errors.New("Missing port param")
	}

	vcUrl.Host = fmt.Sprintf("%s:%s", vcHostStr, vcHostPortStr)

	vcUser, ok := params["user"].(string)
	if !ok {
		return nil, false, errors.New("Missing vcUser param")
	}
	vcPassword, ok := params["password"].(string)
	if !ok {
		return nil, false, errors.New("Missing vcPassword param")
	}
	vcUrl.User = url.UserPassword(vcUser, vcPassword)
	vcUrl.Path = "/sdk"

	insecure := false
	insecureStr, ok := params["insecure-flag"].(string)
	if ok && (insecureStr == "TRUE" || insecureStr == "true") {
		insecure = true
	}

	return &vcUrl, insecure, nil
}

func getVcUrlFromConfig(ctx context.Context, config *rest.Config) (*url.URL, bool, error) {
	params := make(map[string]interface{})

	err := retrievePlatformInfoFromConfig(ctx, config, params, nil)
	if err != nil {
		return nil, false, errors.Errorf("Failed to retrieve VC config secret: %+v", err)
	}

	vcUrl, insecure, err := getVcConfigFromParams(params)
	if err != nil {
		return nil, false, errors.Errorf("Failed to get VC config from params: %+v", err)
	}

	return vcUrl, insecure, nil
}

func verifyMdIsRestoredAsExpected(md metadata) bool {
	reservedLabels := []string {
		"cns.containerCluster.clusterFlavor",
		"cns.containerCluster.clusterId",
		"cns.containerCluster.clusterType",
		"cns.containerCluster.vSphereUser",
		"cns.k8s.pv.name",
		"cns.tag",
		"cns.version",
	}

	extendedMdMap := make(map[string]string)

	for _, label := range md.ExtendedMetadata {
		extendedMdMap[label.Key] = label.Value
	}

	for _, key := range reservedLabels {
		_, ok := extendedMdMap[key]
		if !ok {
			return false
		}
	}

	return true
}

func TestCreateCnsVolume(t *testing.T) {
	path := os.Getenv("HOME") + "/.kube/config"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path/to/whatever does not exist
		t.Skipf("The KubeConfig file, %v, is not exist", path)
	}

	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		t.Fatalf("Failed to build k8s config from kubeconfig file: %+v ", err)
	}

	ctx := context.Background()

	// Step 1: To create the IVD PETM, get all PEs and select one as the reference.
	vcUrl, insecure, err := getVcUrlFromConfig(config)
	if err != nil {
		t.Fatalf("Failed to get VC config from params: %+v", err)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	ivdPETM, err := NewIVDProtectedEntityTypeManagerFromURL(vcUrl, "/ivd", insecure, logger)
	if err != nil {
		t.Fatalf("Failed to get a new ivd PETM: %+v", err)
	}

	peIDs, err := ivdPETM.GetProtectedEntities(ctx)
	if err != nil {
		t.Fatalf("Failed to get all PEs: %+v", err)
	}
	t.Logf("# of PEs returned = %d\n", len(peIDs))
	if len(peIDs) <= 0 {
		t.Skip("No enough number of PEs under the ivd PETM")
	}

	peID := peIDs[0]
	t.Logf("Selected PE ID: %v", peID.String())

	// Get general govmomi client and cns client
	// Step 2: Query the volume result for the selected protected entity/volume
	var queryFilter cnstypes.CnsQueryFilter
	var volumeIDList []cnstypes.CnsVolumeId
	volumeIDList = append(volumeIDList, cnstypes.CnsVolumeId{Id: peID.GetID()})

	queryFilter.VolumeIds = volumeIDList
	queryResult, err := ivdPETM.cnsClient.QueryVolume(ctx, queryFilter)
	if err != nil {
		t.Errorf("Failed to query volume. Error: %+v \n", err)
		t.Fatal(err)
	}
	logger.Debugf("Sucessfully Queried Volumes. queryResult: %+v", queryResult)

	// Step 3: Create a new volume with the same metadata as the selected one
	pe, err := newIVDProtectedEntity(ivdPETM, peID)
	if err != nil {
		t.Fatalf("Failed to get a new PE from the peID, %v: %v", peID.String(), err)
	}

	md, err := pe.getMetadata(ctx)
	if err != nil {
		t.Fatalf("Failed to get the metadata of the PE, %v: %v", pe.id.String(), err)
	}

	logger.Debugf("IVD md: %v", md.ExtendedMetadata)

	t.Logf("PE name, %v", md.VirtualStorageObject.Config.Name)
	md = FilterLabelsFromMetadataForCnsAPIs(md, "cns", logger)
	volumeId, err := createCnsVolumeWithClusterConfig(ctx, config, ivdPETM.client, ivdPETM.cnsClient, md, logger)
	if err != nil {
		t.Fatal("Fail to provision a new volume")
	}

	t.Logf("CNS volume, %v, created", volumeId)
	var volumeIDListToDelete []cnstypes.CnsVolumeId
	volumeIDList = append(volumeIDListToDelete, cnstypes.CnsVolumeId{Id: volumeId})

	defer func () {
		// Always delete the newly created volume at the end of test
		t.Logf("Deleting volume: %+v", volumeIDList)
		deleteTask, err := ivdPETM.cnsClient.DeleteVolume(ctx, volumeIDList, true)
		if err != nil {
			t.Errorf("Failed to delete volume. Error: %+v \n", err)
			t.Fatal(err)
		}
		deleteTaskInfo, err := cns.GetTaskInfo(ctx, deleteTask)
		if err != nil {
			t.Errorf("Failed to delete volume. Error: %+v \n", err)
			t.Fatal(err)
		}
		deleteTaskResult, err := cns.GetTaskResult(ctx, deleteTaskInfo)
		if err != nil {
			t.Errorf("Failed to detach volume. Error: %+v \n", err)
			t.Fatal(err)
		}
		if deleteTaskResult == nil {
			t.Fatalf("Empty delete task results")
		}
		deleteVolumeOperationRes := deleteTaskResult.GetCnsVolumeOperationResult()
		if deleteVolumeOperationRes.Fault != nil {
			t.Fatalf("Failed to delete volume: fault=%+v", deleteVolumeOperationRes.Fault)
		}
		t.Logf("Volume deleted sucessfully")
	} ()

	// Step 4: Query the volume result for the newly created protected entity/volume
	queryFilter.VolumeIds = volumeIDList
	queryResult, err = ivdPETM.cnsClient.QueryVolume(ctx, queryFilter)
	if err != nil {
		t.Errorf("Failed to query volume. Error: %+v \n", err)
		t.Fatal(err)
	}
	logger.Debugf("Sucessfully Queried Volumes. queryResult: %+v", queryResult)

	newPE, err := newIVDProtectedEntity(ivdPETM, newProtectedEntityID(NewIDFromString(volumeId)))
	if err != nil {
		t.Fatalf("Failed to get a new PE from the peID, %v: %v", peID.String(), err)
	}

	newMD, err := newPE.getMetadata(ctx)
	if err != nil {
		t.Fatalf("Failed to get the metadata of the PE, %v: %v", pe.id.String(), err)
	}

	logger.Debugf("IVD md: %v", newMD.ExtendedMetadata)

	// Verify the test result between the actual and expected
	if md.VirtualStorageObject.Config.Name != queryResult.Volumes[0].Name {
		t.Errorf("Volume names mismatch, src: %v, dst: %v", md.VirtualStorageObject.Config.Name, queryResult.Volumes[0].Name)
	} else {
		t.Logf("Volume names match, name: %v", md.VirtualStorageObject.Config.Name)
	}

	if verifyMdIsRestoredAsExpected(newMD) {
		t.Logf("Volume metadata is restored as expected")
	} else {
		t.Errorf("Volume metadata is NOT restored as expected")
	}
}
