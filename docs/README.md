# NewRelic::CloudFormation::Tagging

CRUD operations for New Relic Tags via the NerdGraph API

## Syntax

To declare this entity in your AWS CloudFormation template, use the following syntax:

### JSON

<pre>
{
    "Type" : "NewRelic::CloudFormation::Tagging",
    "Properties" : {
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
Type: NewRelic::CloudFormation::Tagging
Properties:
    <a href="#entityguid" title="EntityGuid">EntityGuid</a>: <i>String</i>
    <a href="#listqueryfilter" title="ListQueryFilter">ListQueryFilter</a>: <i>String</i>
    <a href="#variables" title="Variables">Variables</a>: <i><a href="variables.md">Variables</a></i>
    <a href="#tags" title="Tags">Tags</a>: <i>
      - <a href="tagobject.md">TagObject</a></i>
    <a href="#semantics" title="Semantics">Semantics</a>: <i>String</i>
</pre>

## Properties

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

### Fn::GetAtt

The `Fn::GetAtt` intrinsic function returns a value for a specified attribute of this type. The following are the available attributes and sample return values.

For more information about using the `Fn::GetAtt` intrinsic function, see [Fn::GetAtt](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-getatt.html).

#### Guid

Returns the <code>Guid</code> value.

