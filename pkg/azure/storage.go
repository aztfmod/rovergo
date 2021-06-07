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
	"github.com/spf13/cobra"
)

// FindStorageAccount returns resource id for a CAF tagged storage account
func FindStorageAccount(level string, environment string, subID string) (string, error) {
	query := fmt.Sprintf(`Resources 
		| where type == 'microsoft.storage/storageaccounts' 
		| where tags.level == '%s'
		| where tags.environment == '%s'
		| limit 1
		| project id`, level, environment)

	queryResults := RunQuery(query, subID)
	resSlice, ok := queryResults.([]interface{})
	if !ok {
		cobra.CheckErr("FindStorageAccount: Failed to parse query results")
	}

	if len(resSlice) <= 0 {
		return "", errors.New("No storage account found")
	}

	resMap, ok := resSlice[0].(map[string]interface{})
	if !ok {
		cobra.CheckErr("FindStorageAccount: Failed to parse query results")
	}

	return resMap["id"].(string), nil
}

// GetAccountKey fetches the access key for a storage account
func GetAccountKey(subID string, accountName string, resGrp string) string {
	client := storage.NewAccountsClient(subID)
	client.Authorizer = GetAuthorizer()

	keysRes, err := client.ListKeys(context.Background(), resGrp, accountName, storage.Kerb)
	cobra.CheckErr(err)

	return *(*keysRes.Keys)[0].Value
}

// UploadFileToBlob does what you might expect it to
func UploadFileToBlob(storageAcctID string, blobContainer string, blobName string, filePath string) {
	subID, resGrp, accountName := ParseResourceID(storageAcctID)
	console.Debugf("Uploading to storage account '%s' in res grp '%s' and subscription '%s'\n", accountName, resGrp, subID)
	console.Debugf("Will upload file '%s' to container '%s' to blob '%s'\n", filePath, blobContainer, blobName)

	accountKey := GetAccountKey(subID, accountName, resGrp)

	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	cobra.CheckErr(err)
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	containerURL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, blobContainer))

	blobContainerURL := azblob.NewContainerURL(*containerURL, pipeline)
	blobURL := blobContainerURL.NewBlockBlobURL(blobName)
	file, err := os.Open(filePath)
	cobra.CheckErr(err)

	uploadResp, err := azblob.UploadFileToBlockBlob(context.Background(), file, blobURL, azblob.UploadToBlockBlobOptions{})
	if uploadResp.Response().StatusCode > 201 {
		cobra.CheckErr(fmt.Sprintf("UploadFileToBlob failed with status %d to upload file '%s' to %s/%s", uploadResp.Response().StatusCode, filePath, blobContainer, blobName))
	}
	cobra.CheckErr(err)
}

// DownloadFileFromBlob does what you might expect it to
func DownloadFileFromBlob(storageAcctID string, blobContainer string, blobName string, filePath string) {
	subID, resGrp, accountName := ParseResourceID(storageAcctID)
	console.Debugf("Downloading from storage account '%s' in res grp '%s' and subscription '%s'\n", accountName, resGrp, subID)
	console.Debugf("Will download blob '%s' from container '%s' to file '%s'\n", blobName, blobContainer, filePath)

	accountKey := GetAccountKey(subID, accountName, resGrp)

	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	cobra.CheckErr(err)
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	containerURL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, blobContainer))

	blobContainerURL := azblob.NewContainerURL(*containerURL, pipeline)
	blobURL := blobContainerURL.NewBlockBlobURL(blobName)
	file, err := os.Create(filePath)
	cobra.CheckErr(err)

	err = azblob.DownloadBlobToFile(context.Background(), blobURL.BlobURL, 0, 0, file, azblob.DownloadFromBlobOptions{})
	cobra.CheckErr(err)
}
