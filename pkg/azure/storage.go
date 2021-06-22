//
// Rover - Azure storage
// * For working with azure storage for accessing & storing TF state
// * Ben C, May 2021
//

package azure

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/aztfmod/rover/pkg/console"
)

// FindStorageAccount returns resource id for a CAF tagged storage account
func FindStorageAccount(level string, environment string, subID string) (string, error) {
	query := fmt.Sprintf(`Resources 
		| where type == 'microsoft.storage/storageaccounts' 
		| where tags.level == '%s'
		| where tags.environment == '%s'
		| limit 1
		| project id`, level, environment)

	queryResults, err := RunQuery(query, subID)
	if err != nil {
		return "", err
	}

	resSlice, ok := queryResults.([]interface{})
	if !ok {
		return "", errors.New("FindStorageAccount: Failed to parse query results")
	}

	if len(resSlice) <= 0 {
		return "", errors.New("No storage account found")
	}

	resMap, ok := resSlice[0].(map[string]interface{})
	if !ok {
		return "", errors.New("FindStorageAccount: Failed to parse query results")
	}

	return resMap["id"].(string), nil
}

// GetAccountKey fetches the access key for a storage account
func GetAccountKey(subID string, accountName string, resGrp string) (string, error) {
	client := storage.NewAccountsClient(subID)
	authorizer, err := GetAuthorizer()
	if err != nil {
		return "", err
	}
	client.Authorizer = authorizer

	keysRes, err := client.ListKeys(context.Background(), resGrp, accountName, storage.Kerb)
	if err != nil {
		return "", err
	}

	return *(*keysRes.Keys)[0].Value, nil
}

// UploadFileToBlob does what you might expect it to
func UploadFileToBlob(storageAcctID string, blobContainer string, blobName string, filePath string) error {
	subID, resGrp, accountName, err := ParseResourceID(storageAcctID)
	if err != nil {
		return err
	}
	console.Debugf("Uploading to storage account '%s' in res grp '%s' and subscription '%s'\n", accountName, resGrp, subID)
	console.Debugf("Will upload file '%s' to container '%s' to blob '%s'\n", filePath, blobContainer, blobName)

	accountKey, err := GetAccountKey(subID, accountName, resGrp)
	if err != nil {
		return err
	}

	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	endpoint, err := StorageEndpointForSubscription()
	if err != nil {
		return err
	}
	containerURL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.%s/%s", accountName, endpoint, blobContainer))

	blobContainerURL := azblob.NewContainerURL(*containerURL, pipeline)
	blobURL := blobContainerURL.NewBlockBlobURL(blobName)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	uploadResp, err := azblob.UploadFileToBlockBlob(context.Background(), file, blobURL, azblob.UploadToBlockBlobOptions{})
	if uploadResp.Response().StatusCode > 201 {
		return fmt.Errorf("UploadFileToBlob failed with status %d to upload file '%s' to %s/%s", uploadResp.Response().StatusCode, filePath, blobContainer, blobName)
	}
	if err != nil {
		return err
	}
	return nil
}

// DownloadFileFromBlob does what you might expect it to
func DownloadFileFromBlob(storageAcctID string, blobContainer string, blobName string, filePath string) error {
	subID, resGrp, accountName, err := ParseResourceID(storageAcctID)
	if err != nil {
		return err
	}
	console.Debugf("Downloading from storage account '%s' in res grp '%s' and subscription '%s'\n", accountName, resGrp, subID)
	console.Debugf("Will download blob '%s' from container '%s' to file '%s'\n", blobName, blobContainer, filePath)

	accountKey, err := GetAccountKey(subID, accountName, resGrp)
	if err != nil {
		return err
	}

	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return err
	}
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	endpoint, err := StorageEndpointForSubscription()
	if err != nil {
		return err
	}
	containerURL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.%s/%s", accountName, endpoint, blobContainer))

	blobContainerURL := azblob.NewContainerURL(*containerURL, pipeline)
	blobURL := blobContainerURL.NewBlockBlobURL(blobName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	err = azblob.DownloadBlobToFile(context.Background(), blobURL.BlobURL, 0, 0, file, azblob.DownloadFromBlobOptions{})
	if err != nil {
		return err
	}

	return nil
}

// ListBlobs does what you might expect it to
func ListBlobs(storageAcctID string, blobContainer string) (blobs []azblob.BlobItemInternal, err error) {
	subID, resGrp, accountName, err := ParseResourceID(storageAcctID)
	if err != nil {
		return nil, err
	}

	accountKey, err := GetAccountKey(subID, accountName, resGrp)
	if err != nil {
		return nil, err
	}

	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, err
	}
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	endpoint, err := StorageEndpointForSubscription()
	if err != nil {
		return nil, err
	}
	containerURL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.%s/%s", accountName, endpoint, blobContainer))

	blobContainerURL := azblob.NewContainerURL(*containerURL, pipeline)

	for marker := (azblob.Marker{}); marker.NotDone(); {
		listBlob, err := blobContainerURL.ListBlobsFlatSegment(context.Background(), marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			return nil, err
		}

		marker = listBlob.NextMarker
		blobs = append(blobs, listBlob.Segment.BlobItems...)
	}

	return blobs, nil
}
