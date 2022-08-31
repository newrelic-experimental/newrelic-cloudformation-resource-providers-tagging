# newrelic::cloudformation::tagging

CRUDL operations for New Relic Tags via the NerdGraph API

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "newrelic::cloudformation::tagging",
    "Properties" : {
        "<a href="#endpoint" title="Endpoint">Endpoint</a>" : <i>String</i>,
        "<a href="#apikey" title="APIKey">APIKey</a>" : <i>String</i>,
        "<a href="#guid" title="Guid">Guid</a>" : <i>String</i>,
        "<a href="#entityguid" title="EntityGuid">EntityGuid</a>" : <i>String</i>,
        "<a href="#listqueryfilter" title="ListQueryFilter">ListQueryFilter</a>" : <i>String</i>,
        "<a href="#variables" title="Variables">Variables</a>" : <i><a href="variables.md">Variables</a></i>,
        "<a href="#tags" title="Tags">Tags</a>" : <i>[ <a href="tagobject.md">TagObject</a>, ... ]</i>,
        "<a href="#semantics" title="Semantics">Semantics</a>" : <i>String</i>
    }
}
</pre>

### YAML

<pre>
Type: newrelic::cloudformation::tagging
Properties:
    <a href="#endpoint" title="Endpoint">Endpoint</a>: <i>String</i>
    <a href="#apikey" title="APIKey">APIKey</a>: <i>String</i>
    <a href="#guid" title="Guid">Guid</a>: <i>String</i>
    <a href="#entityguid" title="EntityGuid">EntityGuid</a>: <i>String</i>
    <a href="#listqueryfilter" title="ListQueryFilter">ListQueryFilter</a>: <i>String</i>
    <a href="#variables" title="Variables">Variables</a>: <i><a href="variables.md">Variables</a></i>
    <a href="#tags" title="Tags">Tags</a>: <i>
      - <a href="tagobject.md">TagObject</a></i>
    <a href="#semantics" title="Semantics">Semantics</a>: <i>String</i>
</pre>

## Properties

#### Endpoint

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### APIKey

_Required_: Yes

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Guid

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### EntityGuid

_Required_: Yes

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### ListQueryFilter

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Variables

_Required_: No

_Type_: <a href="variables.md">Variables</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Tags

_Required_: Yes

_Type_: List of <a href="tagobject.md">TagObject</a>

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

#### Semantics

_Required_: No

_Type_: String

_Update requires_: [No interruption](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-updating-stacks-update-behaviors.html#update-no-interrupt)

## Return Values

### Ref

When you pass the logical ID of this resource to the intrinsic `Ref` function, Ref returns the Guid.
