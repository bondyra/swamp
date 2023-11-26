package aws

type AwsConnection struct {
	profile string
	client  AwsClientInterface
}

func (ac AwsConnection) Init(factory AwsFactory) error {
	var err error
	ac.client, err = factory.NewClient(ac.profile)
	if err != nil {
		return err
	}
	return nil
}

// "GetResourceInput":{
//     "TypeName":string*,
//     "TypeVersionId":string,
//     "RoleArn":string,
//     "Identifier":string*
//   }

//   "GetResourceOutput":{
//     "TypeName":string,
//     "ResourceDescription":{
//       "Identifier":string,
//       "Properties":string json
//     }
//   }

//   "ListResourcesInput":{
//     "TypeName":string*,
//     "TypeVersionId":string,
//     "RoleArn":string,
//     "NextToken":string,
//     "MaxResults":string,
//     "ResourceModel":string json
//   },
//   "ListResourcesOutput":{
//     "TypeName":string,
//     "ResourceDescriptions":[
//       {
//         "Identifier":string,
//         "Properties":string json
//       }
//     ],
//     "NextToken":string
//   }
