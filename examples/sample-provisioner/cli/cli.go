/*
Copyright 2020 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/container-object-storage-interface/cosi-provisioner-sidecar/pkg/grpcclient"
	cosispec "github.com/container-object-storage-interface/spec"
	"google.golang.org/grpc"
	"k8s.io/klog"
)

const (
	// Default timeout of short COSI calls like CreateBucket, etc
	cosiTimeout = time.Second
	// Verify (and update, if needed) the node ID at this freqeuency.
	sleepDuration = 2 * time.Minute
	// Interval of logging connection errors
	connectionLoggingInterval = 10 * time.Second
)

// Command line flags
var (
	connectionTimeout = flag.Duration("connection-timeout", 0, "The --connection-timeout flag is deprecated")
	cosiAddress       = flag.String("cosi-address", "tcp://0.0.0.0:9000", "Path of the COSI driver socket that the provisioner  will connect to.")
	// List of supported versions
	supportedVersions = []string{"1.0.0"}
)

func main() {
	grpcClient, err := grpcclient.NewGRPCClient(*cosiAddress, []grpc.DialOption{}, nil)
	if err != nil {
		klog.Errorf("error creating GRPC Client: %v", err)
		os.Exit(1)
	}
	cosiConn, err := grpcClient.ConnectWithLogging(connectionLoggingInterval)
	if err != nil {
		klog.Errorf("error connecting to COSI driver: %v", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), cosiTimeout)
	defer cancel()
	klog.V(1).Infof("Creating COSI client")
	client := cosispec.NewProvisionerClient(cosiConn)
	klog.Infof("Calling COSI driver to discover driver name")
	// getPovisonerInfo(ctx, client)
	// createBucket(ctx, client)
	crateBucketAcces(ctx, client)
	// revokeBucketAccess(ctx, client)
	// deleteBucket(ctx, client)
}
func getPovisonerInfo(ctx context.Context, client cosispec.ProvisionerClient) {
	log.Println("Inside Provisioner info")
	req := cosispec.ProvisionerGetInfoRequest{}
	rsp, err := client.ProvisionerGetInfo(ctx, &req)
	if err != nil {
		klog.Errorf("error calling COSI cosispec.ProvisionerGetInfoRequest: %v", err)
		os.Exit(1)
	}
	provisionerName := rsp.ProvisionerIdentity
	// TODO: Register provisioner using internal type
	klog.Info("This provsioner is working with the driver identified as: ", provisionerName)
}
func createBucket(ctx context.Context, client cosispec.ProvisionerClient) {
	log.Println("Inside Bucket Create")
	req := cosispec.ProvisionerCreateBucketRequest{
		BucketName: "bucket10000000",
		Region:     "us-west",
		Zone:       "south",
		// BucketContext: ["endPoint"][""]
	}
	rsp, err := client.ProvisionerCreateBucket(ctx, &req)
	if err != nil {
		klog.Errorf("error calling COSI cosispec.ProvisionerGetInfoRequest: %v", err)
		os.Exit(1)
	}
	fmt.Println(rsp)
}
func crateBucketAcces(ctx context.Context, client cosispec.ProvisionerClient) {
	req := cosispec.ProvisionerGrantBucketAccessRequest{
		BucketName:   "bucket10000000",
		Region:       "unknown",
		Zone:         "us-south",
		Principal:    "username",
		AccessPolicy: "write",
	}
	rsp, err := client.ProvisionerGrantBucketAccess(ctx, &req)
	if err != nil {
		klog.Errorf("error calling COSI cosispec.ProvisionerGrantBucketAccess: %v", err)
		os.Exit(1)
	}
	fmt.Println("Printing bucket access\n", rsp.CredentialsFileContents)
}

func revokeBucketAccess(ctx context.Context, client cosispec.ProvisionerClient) {

	req := cosispec.ProvisionerRevokeBucketAccessRequest{
		BucketName: "bucket10011111",
		Region:     "unknown",
		Zone:       "us-south",
		Principal:  "SSLEBNV7T8QBHLUO06ZH",
	}
	rsp, err := client.ProvisionerRevokeBucketAccess(ctx, &req)
	if err != nil {
		klog.Errorf("error calling COSI cosispec.ProvisionerRevokeBucketAccess: %v", err)
		os.Exit(1)
	}
	fmt.Println("Printing bucket access", rsp)

}

func deleteBucket(ctx context.Context, client cosispec.ProvisionerClient) {

	req := cosispec.ProvisionerDeleteBucketRequest{
		BucketName: "bucket10000000",
		Region:     "unknown",
		Zone:       "us-south",
	}
	rsp, err := client.ProvisionerDeleteBucket(ctx, &req)
	if err != nil {
		klog.Errorf("error calling COSI cosispec.ProvisionerDeleteBucket: %v", err)
		os.Exit(1)
	}
	fmt.Println("Printing bucket delete", rsp)

}

// ~ docker run -p 9001:9000 \
//   -e "MINIO_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE" \
//   -e "MINIO_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" \
//   minio/minio server /data
