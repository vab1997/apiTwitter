package secretManager

import (
	"encoding/json"
	"fmt"

	"apiTwitter/awsConfig"
	"apiTwitter/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetSecret(secretName string) (models.SecretManager, error) {
	var dataSecret models.SecretManager
	fmt.Println("> Secret " + secretName)

	svc := secretsmanager.NewFromConfig(awsConfig.Cfg)
	key, err := svc.GetSecretValue(awsConfig.Ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})

	if err != nil {
		fmt.Println("Error to get secret " + err.Error())
		return dataSecret, err
	}

	json.Unmarshal([]byte(*key.SecretString), &dataSecret)
	fmt.Println("> Read secret OK")

	return dataSecret, nil
}
