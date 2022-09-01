[sh![New Relic Experimental header](https://github.com/newrelic/opensource-website/raw/master/src/images/categories/Experimental.png)](https://opensource.newrelic.com/oss-category/#new-relic-experimental)

# newrelic::cloudformation::tagging

![GitHub forks](https://img.shields.io/github/forks/newrelic-experimental/newrelic-experimental-FIT-template?style=social)
![GitHub stars](https://img.shields.io/github/stars/newrelic-experimental/newrelic-experimental-FIT-template?style=social)
![GitHub watchers](https://img.shields.io/github/watchers/newrelic-experimental/newrelic-experimental-FIT-template?style=social)

![GitHub all releases](https://img.shields.io/github/downloads/newrelic-experimental/newrelic-experimental-FIT-template/total)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/newrelic-experimental/newrelic-experimental-FIT-template)
![GitHub last commit](https://img.shields.io/github/last-commit/newrelic-experimental/newrelic-experimental-FIT-template)
![GitHub Release Date](https://img.shields.io/github/release-date/newrelic-experimental/newrelic-experimental-FIT-template)


![GitHub issues](https://img.shields.io/github/issues/newrelic-experimental/newrelic-experimental-FIT-template)
![GitHub issues closed](https://img.shields.io/github/issues-closed/newrelic-experimental/newrelic-experimental-FIT-template)
![GitHub pull requests](https://img.shields.io/github/issues-pr/newrelic-experimental/newrelic-experimental-FIT-template)
![GitHub pull requests closed](https://img.shields.io/github/issues-pr-closed/newrelic-experimental/newrelic-experimental-FIT-template)

## Description
This Cloud Formation Custom Resource provides a CRUD interface to the New Relic [NerdGraph (GraphQL) Tagging API](https://docs.newrelic.com/docs/apis/nerdgraph/examples/nerdgraph-tagging-api-tutorial) for Cloud Formation stacks.

(For performance reasons CloudFormation _List_ is not a supported operation.)

## Use

## Model
| Field      | Type        | Default                          | Create | Update | Delete | Read | Notes                                                                                                                       |
|------------|-------------|----------------------------------|:------:|:------:|:------:|:----:|-----------------------------------------------------------------------------------------------------------------------------|
| AccountID  | string      | none                             |   R    |        |        |  R   | [New Relic Account ID](https://docs.newrelic.com/docs/accounts/accounts-billing/account-structure/account-id/)              |
| APIKey     | string      | none                             |   R    |   R    |   R    |  R   | [New Relic User Key](https://docs.newrelic.com/docs/apis/intro-apis/new-relic-api-keys/#overview-keys)                      |
| Endpoint   | string      | https://api.newrelic.com/graphql |   O    |   O    |   O    |  O   | [API endpoints](https://docs.newrelic.com/docs/apis/nerdgraph/get-started/introduction-new-relic-nerdgraph/#authentication) |
| EntityGuid | string      | none                             |   R    |   R    |   R    |  R   |                                                                                                                             |
| Guid       | string      | none                             |        |   R    |   R    |  R   |                                                                                                                             |
| Tags       | []TagObject | none                             |   R    |   R    |   R    |  R   |                                                                                                                             |
| Semantics  | string      | AWSSemantics                     |   O    |   O    |   O    |  R   |                                                                                                                             |                                                                                                                             |
| Variables  | Object      | none                             |        |        |        |      |                                                                                                                             |


Key:
- R- Requird
- O- Optional
- Blank- unused

### EntityGuid
The New Relic Guid of the Entity that is manipulated. Must be provided for all operations.

### Guid
`Guid` this is to satisfy CloudWatch, ignore it

### Tags
An array of Tag Objects

### Tag Object
A JSON or YAML array of objects consisting of a 
- Key
- Values an array of string values to associate with the Key

See the example below

### Semantics
This flag controls how the Resource interacts with New Relic:
- `AWSSemantics` (default) follows typical CloudWatch semantics and is fairly close to 'normal' Map semantics
  - Create: remove the Key if it already exists and write the new value(s) in its place
  - Update:
    - Key present: delete the Key and replace it (same as Create)
    - Key not present: fail with NOT_FOUND
  - Delete
    - Key present: delete the Key and value(s)
    - Key not present: fail with NOT_FOUND
- `NewRelicSemantics` 
  - Create: if the Key already exists _add_ the new values to the existing key, otherwise a normal 'create'
  - Update: replace the _entire_ existing tag set with these tags
  - Delete: delete any matching key

### Variables
`Variables` is a JSON object of key/value pairs (string/string) that are substituted in the `DashboardInput` string with [Moustache](#Moustache) allowing for parameterized dashboards at the CloudFormation level.

## Moustache
All text substitution is done using a [Go implementation](https://github.com/cbroglie/mustache) of the [Moustache specification](https://github.com/mustache/spec)`. [The manual is here](http://mustache.github.io/mustache.5.html), in
general all you need to know is use triple curly braces.

## Example
An example of a YAML configuration, this is valid and works:
```yaml
AWSTemplateFormatVersion: 2010-09-09
Description: Sample New Relic Tagging Template
Resources:
  Dashboard1:
    Type: 'newrelic::cloudformation::tagging'
    Properties:
      APIKey: "YOUR_API_KEY_HERE"
      EntityGuid: "YOUR_ENTITY_GUID_HERE"
      Tags:
        - Key: "Template-Key1"
          Values: [ "Template-Key1-Value1-Create", "Template-Key1-Value2-Create" ]
```

## Important Notes
- Key values are _case sensitive_, `ThisIsAKey` *is not* the same as `thisisakey`.

## Troubleshooting
- Error log
- `Debug` log level for mustache substitution
- Validate mutation using the [Explorer](https://api.newrelic.com/graphiql)

## Building
- Install Docker
- Install Golang
- [Install the CloudFormation Command Line Interface (CFN-CLI)](https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/what-is-cloudformation-cli.html)
- `make clean build `

## Testing
- Start Docker
- Activate `cfn-cli`, usually `source ~/.virtenv/aws-cfn-cli/bin/activate`
- Build and start the container `make clean build ; sam local start-lambda --warm-containers eager`
- Run the [CloudFormation Contract Tests](https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/contract-tests.html) `cfn test`

## Publishing
```bash
# Double check the resulting zip file for security leaks!
cfn submit --dry-run
# Send the resource to AWS, the result is a private resource
cfn submit --set-default --no-role --region us-east-1
# Test the private resource with the sample template to ensure it works
aws cloudformation deploy --region us-east-1 --template-file template-examples-live/sample.yml --stack-name tagging-test-stack
# Tell AWS to run the Contract Tests, required for going public
aws cloudformation test-type --region us-east-1 --log-delivery-bucket newrelic--cloudformation--custom--resources --arn arn:aws:cloudformation:us-east-1:830139413159:type/resource/newrelic-cloudformation-tagging
# Check the result
aws cloudformation describe-type --region us-east-1  --arn arn:aws:cloudformation:us-east-1:830139413159:type/resource/newrelic-cloudformation-tagging
# Also the logs are in CloudWatch. They end in .zip but are really gunzip so
# gunzip -S .zip <file>
# Publish the extension publicly AFTER pushing the final version to GitHub AND generating/tagging a release
# IMPORTANT!
#   --public-version-number (string)
#   The version number to assign to this version of the extension.
#   Use the following format, and adhere to semantic versioning when assigning a version number to your extension:
#     MAJOR.MINOR.PATCH
#   For more information, see Semantic Versioning 2.0.0 .
#   If you don’t specify a version number, CloudFormation increments the version number by one minor version release.
#   You cannot specify a version number the first time you publish a type. CloudFormation automatically sets the first version number to be 1.0.0 .
#
# It's a good idea to not publish until everything is ready AND version 1.0.0 is release in Git!
# KEEP IT ALL IN-SYNC!
#
# Git
#
aws cloudformation publish-type --region us-east-1 --arn arn:aws:cloudformation:us-east-1:830139413159:type/resource/newrelic-cloudformation-dashboards  
aws cloudformation describe-type --region us-east-1 --arn arn:aws:cloudformation:us-east-1:830139413159:type/resource/newrelic-cloudformation-dashboards
```
### Notes
- If you see an `Exception: Could not assume specified role...` error message in the test log then you probably have a `cfn` generated role and they're not correct. As this resource uses no other AWS resources a Role is not required.

## Helpful links
- [CloudFormation CLI User Guide](https://docs.aws.amazon.com/cloudformation-cli/latest/userguide/what-is-cloudformation-cli.html)
- [New Relic GraphQL Explorer](https://api.newrelic.com/graphiql) 


## Support
New Relic has open-sourced this project. This project is provided AS-IS WITHOUT WARRANTY OR DEDICATED SUPPORT. Issues and contributions should be reported to the project here on GitHub.

We encourage you to bring your experiences and questions to the [Explorers Hub](https://discuss.newrelic.com) where our community members collaborate on solutions and new ideas.

## Contributing
We encourage your contributions to improve New Relic NerdGraph CloudFormation Custom Resource! Keep in mind when you submit your pull request, you'll need to sign the CLA via the click-through using CLA-Assistant. You only have to sign the CLA one time per project. If you have any questions, or to execute our corporate CLA, required if your contribution is on behalf of a company, please drop us an email at opensource@newrelic.com.

**A note about vulnerabilities**

As noted in our [security policy](../../security/policy), New Relic is committed to the privacy and security of our customers and their data. We believe that providing coordinated disclosure by security researchers and engaging with the security community are important means to achieve our security goals.

If you believe you have found a security vulnerability in this project or any of New Relic's products or websites, we welcome and greatly appreciate you reporting it to New Relic through [HackerOne](https://hackerone.com/newrelic).

## License
New Relic NerdGraph CloudFormation Custom Resource is licensed under the [Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) License.
