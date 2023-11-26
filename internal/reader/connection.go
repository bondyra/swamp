package reader

type ConnectionPool interface {
	Init([]string) error
}
type Connection interface{}

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
