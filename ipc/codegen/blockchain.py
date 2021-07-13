# WARNING: THIS CODE GENERATOR IS DEPRECATED
# TODO: Reimplement using go code generators if possible

from argparse import ArgumentParser


class HookSignature:
    def __init__(self, name, input, output, error=False, badReturn="nil"):
        super().__init__()
        self.name = name
        self.input = input
        self.output = output
        self.error = error
        self.badReturn = badReturn


signatures = [
    HookSignature(name="NewAddress", input=[("creatorAddress", "[]byte"), ("creatorNonce", "uint64"), ("vmType", "[]byte")], output=[("result", "[]byte")], error=True),
    HookSignature(name="GetStorageData", input=[("accountAddress", "[]byte"), ("index", "[]byte")], output=[("data", "[]byte")], error=True),
    HookSignature(name="GetBlockhash", input=[("nonce", "uint64")], output=[("result", "[]byte")], error=True),
    HookSignature(name="LastNonce", input=[], output=[("result", "uint64")], badReturn="0"),
    HookSignature(name="LastRound", input=[], output=[("result", "uint64")], badReturn="0"),    
    HookSignature(name="LastTimeStamp", input=[], output=[("result", "uint64")], badReturn="0"), 
    HookSignature(name="LastRandomSeed", input=[], output=[("result", "[]byte")], badReturn="nil"), 
    HookSignature(name="LastEpoch", input=[], output=[("result", "uint32")], badReturn="0"),   
    HookSignature(name="GetStateRootHash", input=[], output=[("result", "[]byte")], badReturn="nil"), 
    HookSignature(name="CurrentNonce", input=[], output=[("result", "uint64")], badReturn="0"),   
    HookSignature(name="CurrentRound", input=[], output=[("result", "uint64")], badReturn="0"),
    HookSignature(name="CurrentTimeStamp", input=[], output=[("result", "uint64")], badReturn="0"),  
    HookSignature(name="CurrentRandomSeed", input=[], output=[("result", "[]byte")], badReturn="nil"), 
    HookSignature(name="CurrentEpoch", input=[], output=[("result", "uint32")], badReturn="0"),
    HookSignature(name="ProcessBuiltInFunction", input=[("input", "*vmcommon.ContractCallInput")], output=[("vmOutput", "*vmcommon.VMOutput")], error=True),
    HookSignature(name="GetBuiltinFunctionNames", input=[], output=[("result", "vmcommon.FunctionNames")], badReturn="nil"),
    HookSignature(name="GetAllState", input=[("address", "[]byte")], output=[("result", "map[string][]byte")], error=True),             
    HookSignature(name="GetUserAccount", input=[("address", "[]byte")], output=[("result", "vmcommon.UserAccountHandler")], error=True),             
    HookSignature(name="GetCode", input=[("account", "vmcommon.UserAccountHandler")], output=[("code", "[]byte")], badReturn="nil"),
    HookSignature(name="GetShardOfAddress", input=[("address", "[]byte")], output=[("result", "uint32")], badReturn="0"),    
    HookSignature(name="IsSmartContract", input=[("address", "[]byte")], output=[("result", "bool")], badReturn="false"),
    HookSignature(name="IsPayable", input=[("address", "[]byte")], output=[("result", "bool")], error=True, badReturn="false"),
    HookSignature(name="SaveCompiledCode", input=[("codeHash", "[]byte"), ("code", "[]byte")], output=[], badReturn=""),
    HookSignature(name="GetCompiledCode", input=[("codeHash", "[]byte")], output=[("found", "bool"), ("code", "[]byte")], badReturn="false, nil"),
    HookSignature(name="ClearCompiledCodes", input=[], output=[], badReturn=""),
    HookSignature(name="GetESDTToken", input=[("address", "[]byte"), ("tokenID", "[]byte"), ("nonce", "uint64")], output=[("result", "*esdt.ESDigitalToken")], error=True),
    HookSignature(name="IsInterfaceNil", input=[], output=[("result", "bool")], badReturn="false"),
    HookSignature(name="GetSnapshot", input=[], output=[("result", "int")], badReturn="0"),
    HookSignature(name="RevertToSnapshot", input=[("snapshot", "int")], output=[], error=True)
]

interfaceTOImplementations = { "vmcommon.UserAccountHandler" : "common.Account" }
fieldsOfTO = { "common.Account": [ ("Nonce", None), ("Balance", None), ("CodeHash", None), ("RootHash", None), ("Address", "AddressBytes"), ("DeveloperReward", None), 
                                    ("OwnerAddress", None), ("UserName", None), ("CodeMetadata", None) ] }

class SerializableType:
    def __init__(self, type, fromTypeConverterMethod, toTypeConverterFunction):
        super().__init__()
        self.type = type
        self.fromTypeConverterMethod = fromTypeConverterMethod
        self.toTypeConverterFunction = toTypeConverterFunction

serializableTypes = {
                        "map[string][]byte" :
                        SerializableType(type="SerializableMapStringBytes", fromTypeConverterMethod="ConvertToMap", toTypeConverterFunction="NewSerializableMapStringBytes") ,
                        "*vmcommon.VMOutput" :
                        SerializableType(type="SerializableVMOutput", fromTypeConverterMethod="ConvertToVMOutput", toTypeConverterFunction="NewSerializableVMOutput") ,
                    }

def main():
    parser = ArgumentParser()
    subparsers = parser.add_subparsers()

    messages_parser = subparsers.add_parser("messages")
    messages_parser.set_defaults(func=generate_messages)

    repliers_parser = subparsers.add_parser("repliers")
    repliers_parser.set_defaults(func=generate_repliers)

    reply_slots_parser = subparsers.add_parser("reply-slots")
    reply_slots_parser.set_defaults(func=generate_reply_slots)

    gateway_parser = subparsers.add_parser("gateway")
    gateway_parser.set_defaults(func=generate_gateway)

    factory_parser = subparsers.add_parser("factory")
    factory_parser.set_defaults(func=generate_factory)

    args = parser.parse_args()

    if not hasattr(args, "func"):
        parser.print_help()
    else:
        args.func(args)


def generate_messages(args):
    package = "common"
    print("package " + package)

    print("""
        import (        
        \"github.com/ElrondNetwork/elrond-vm-common\"
        \"github.com/ElrondNetwork/elrond-vm-common/data/esdt\"
        )
	""")

    for signature in signatures:
        request_kind = f"Blockchain{signature.name}Request"
        request_go = f"""
        // Message{request_kind} represents a request message
        type Message{request_kind} struct {{
            Message
            {get_struct_fields_go(signature.input, package)}
        }}

        // NewMessage{request_kind} creates a request message
        func NewMessage{request_kind}({get_ctor_args(signature.input, package)}) *Message{request_kind} {{
            message := &Message{request_kind}{{}}
            message.Kind = {request_kind}
            {get_field_assignments(signature.input)}
            return message
        }}
        """

        response_kind = f"Blockchain{signature.name}Response"
        response_go = f"""
        // Message{response_kind} represents a response message
        type Message{response_kind} struct {{
            Message
             {get_struct_fields_go(signature.output, package)}
        }}

        // NewMessage{response_kind} creates a response message
        func NewMessage{response_kind}({get_ctor_args(signature.output, package, error=signature.error)}) *Message{response_kind} {{
            message := &Message{response_kind}{{}}
            message.Kind = {response_kind}
            {get_field_assignments(signature.output, error=signature.error)}
            return message
        }}
        """

        print(request_go)
        print(response_go)


def get_struct_fields_go(input_output, package):
    fields = []
    for arg_name, arg_type in input_output:
        field_name = my_capitalize(arg_name)
        if arg_type in interfaceTOImplementations:
            arg_type = "*" + interfaceTOImplementations[arg_type].replace(package + ".", '')
        if arg_type in serializableTypes:
            arg_type = "*" + serializableTypes[arg_type].type
        fields.append(f"{field_name} {arg_type}")

    return "\n".join(fields)


def get_ctor_args(input_output, package, useInterfaces=False, error=False):
    args = []
    for arg_name, arg_type in input_output:        
        if not useInterfaces and (arg_type in interfaceTOImplementations):
            arg_type = "*" + interfaceTOImplementations[arg_type].replace(package+".", '')
        args.append(f"{arg_name} {arg_type}")

    if error:
        args.append(f"err error")

    return ", ".join(args)


def get_field_assignments(input_output, error=False):
    assignments = []
    for arg_name, arg_type in input_output:
        field_name = my_capitalize(arg_name)
        if arg_type in serializableTypes:
            arg_name = f"{serializableTypes[arg_type].toTypeConverterFunction}({arg_name})"
        assignments.append(f"message.{field_name} = {arg_name}")

    if error:
        assignments.append(f"message.SetError(err)")

    return "\n".join(assignments)


def my_capitalize(input):
    return input[0].upper() + input[1:]


def generate_repliers(args):
    print("package nodepart")
    print("""
    	import (
	    \"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/ipc/common\"
        \"github.com/ElrondNetwork/arwen-wasm-vm/v1_3/arwen\"
	    \"github.com/ElrondNetwork/elrond-vm-common/data/esdt\"
	)
	""")

    for signature in signatures:
        call_go, output_args_for_call, output_args_for_err, interface_results = get_call(signature)
        typedRequest = f"typedRequest := request.(*common.MessageBlockchain{signature.name}Request)\n" if signature.input else ""

        errCode = f"""
            if err != nil || arwen.IfNil(result) {{
                return common.NewMessageBlockchain{signature.name}Response({output_args_for_err})
            }}
        """

        func_go = f"""
        func (part *NodePart) replyToBlockchain{signature.name}(request common.MessageHandler) common.MessageHandler {{
            {typedRequest}{call_go}
            {errCode if interface_results else ""}
            response := common.NewMessageBlockchain{signature.name}Response({output_args_for_call})
            return response
        }}
        """
        print(func_go)


def get_call(signature):
    results = []
    interface_results = []
    output_args_for_call = []
    output_args_for_err = []
    call_args = []

    for arg_name, arg_type in signature.output:
        if arg_type in interfaceTOImplementations:
            output_args_for_call.append(generate_TO_for_interface(arg_name, arg_type))
            interface_results.append(arg_name)
        else:
            output_args_for_call.append(arg_name)
        results.append(arg_name)
        output_args_for_err.append("nil")

    if signature.error:
        results.append("err")
        output_args_for_call.append("err")
        output_args_for_err.append("err")

    for arg_name, _ in signature.input:
        call_args.append(f"typedRequest.{my_capitalize(arg_name)}")

    results = ", ".join(results)
    output_args_for_call = ", ".join(output_args_for_call)
    output_args_for_err = ", ".join(output_args_for_err)
    call_args = ", ".join(call_args)

    if signature.output:
        return f"{results} := part.blockchain.{signature.name}({call_args})", output_args_for_call, output_args_for_err, interface_results
    elif signature.error:
        return f"err := part.blockchain.{signature.name}({call_args})", output_args_for_call, output_args_for_err, interface_results
    else:
        return f"part.blockchain.{signature.name}({call_args})", output_args_for_call, output_args_for_err, interface_results

def generate_reply_slots(args):
    print("part.Repliers = common.CreateReplySlots(part.noopReplier)")

    for signature in signatures:
        print(f"part.Repliers[common.Blockchain{signature.name}Request] = part.replyToBlockchain{signature.name}")


def generate_gateway(args):
    package = "arwenpart"
    print("package " + package)
    print("""

import (
    "github.com/ElrondNetwork/elrond-vm-common/data/esdt"

    "github.com/ElrondNetwork/arwen-wasm-vm/v1_3/ipc/common"
    "github.com/ElrondNetwork/elrond-vm-common"
)

var _ vmcommon.BlockchainHook = (*BlockchainHookGateway)(nil)

// BlockchainHookGateway forwards requests to the actual hook
type BlockchainHookGateway struct {
    messenger *ArwenMessenger
}

// NewBlockchainHookGateway creates a new gateway
func NewBlockchainHookGateway(messenger *ArwenMessenger) *BlockchainHookGateway {
    return &BlockchainHookGateway{messenger: messenger}
}
""")

    for signature in signatures:
        errorSeparator = f", " if signature.error and signature.output else ""
        badReturn = f"{signature.badReturn} " if signature.output else ""
        func_go = f"""
        // {signature.name} forwards a message to the actual hook
        func (blockchain *BlockchainHookGateway) {signature.name}({get_ctor_args(signature.input, package, useInterfaces=True)}) {get_output_args(signature)} {{
            {generate_TOs_for_interfaces_gateway(signature)}
            request := common.NewMessageBlockchain{signature.name}Request({get_call_args(signature)})
            rawResponse, err := blockchain.messenger.SendHookCallRequest(request)
            if err != nil {{
                return {badReturn}{errorSeparator}{"err" if signature.error else ""}
            }}

            if rawResponse.GetKind() != common.Blockchain{signature.name}Response {{
                return {badReturn}{errorSeparator}{"common.ErrBadHookResponseFromNode" if signature.error else ""}
            }}

            """ 
        if signature.output :
                func_go = func_go + f"response := rawResponse.(*common.MessageBlockchain{signature.name}Response) "
        func_go = func_go + f"""
            {get_gateway_return(signature)}
        }}
        """

        print(func_go)

def generate_TOs_for_interfaces_gateway(signature):    
    generatedCode = ""
    for arg_name, arg_type in signature.input:
        if arg_type in interfaceTOImplementations:            
            buildCodeOfTO = "request" + my_capitalize(arg_name) + " := " + generate_TO_for_interface(arg_name, arg_type)
            generatedCode += "\n" + buildCodeOfTO

    return generatedCode

def generate_TO_for_interface(arg_name, arg_type):
    result = "&" + interfaceTOImplementations[arg_type] + "{\n"
    for field, function in fieldsOfTO[interfaceTOImplementations[arg_type]]:
        if function is None:
            function = "Get" + field
        result += f"{field}:{arg_name}.{function}(),\n"
    result += "}"    
    return result

def get_call_args(signature):
    call_args = []
    for arg_name, arg_type in signature.input:
        if arg_type in interfaceTOImplementations: 
            call_args.append("request" + my_capitalize(arg_name))
        else:
            call_args.append(arg_name)

    return ", ".join(call_args)


def get_output_args(signature):
    output_args = []
    for arg_name, arg_type in signature.output:
        output_args.append(arg_type)
    if signature.error:
        output_args.append(f"error")

    result = ", ".join(output_args)
    if len(output_args) > 1:
        result = f"({result})"

    return result


def get_gateway_return(signature):
    if not signature.output:
        if not signature.error:
            return f"return"
        else:
            return f"return err"
    result_field, _ = signature.output[0]
    result_field = my_capitalize(result_field)

    return_args = []
    for arg_name, arg_type in signature.output:
        fromTypeConverterMethod = ""
        if arg_type in serializableTypes:
            fromTypeConverterMethod = f".{serializableTypes[arg_type].fromTypeConverterMethod}()"
        return_args.append("response." + my_capitalize(arg_name) + fromTypeConverterMethod)
    returnResult = ", ".join(return_args)        

    if signature.error:
        return f"return {returnResult}, response.GetError()"

    return f"return {returnResult}"


def generate_factory(args):
    assignments = ""

    for signature in signatures:
        assignments += f"messageCreators[Blockchain{signature.name}Request] = createMessageBlockchain{signature.name}Request"
        assignments += "\n"
        assignments += f"messageCreators[Blockchain{signature.name}Response] = createMessageBlockchain{signature.name}Response"
        assignments += "\n"

    print(f"""
package common


// CreateMessage creates a message given its kind
func CreateMessage(kind MessageKind) MessageHandler {{
    kindIndex := uint32(kind)
    length := uint32(len(messageCreators))
    if kindIndex < length {{
        message := messageCreators[kindIndex]()
        message.SetKind(kind)
        return message
    }}

    return createUndefinedMessage()
}}

type messageCreator func() MessageHandler

var messageCreators = make([]messageCreator, LastKind)

func init() {{
    for i := 0; i < len(messageCreators); i++ {{
        messageCreators[i] = createUndefinedMessage
    }}

    messageCreators[Initialize] = createMessageInitialize
    messageCreators[Stop] = createMessageStop
    messageCreators[ContractDeployRequest] = createMessageContractDeployRequest
    messageCreators[ContractCallRequest] = createMessageContractCallRequest
    messageCreators[ContractResponse] = createMessageContractResponse
    messageCreators[DiagnoseWaitRequest] = createMessageDiagnoseWaitRequest
    messageCreators[DiagnoseWaitResponse] = createMessageDiagnoseWaitResponse
    messageCreators[VersionRequest] = createMessageVersionRequest
	messageCreators[VersionResponse] = createMessageVersionResponse

    {assignments}
}}

func createMessageInitialize() MessageHandler {{
    return &MessageInitialize{{}}
}}

func createMessageStop() MessageHandler {{
    return &MessageStop{{}}
}}

func createMessageContractDeployRequest() MessageHandler {{
    return &MessageContractDeployRequest{{}}
}}

func createMessageContractCallRequest() MessageHandler {{
    return &MessageContractCallRequest{{}}
}}

func createMessageContractResponse() MessageHandler {{
    return &MessageContractResponse{{}}
}}

func createMessageDiagnoseWaitRequest() MessageHandler {{
    return &MessageDiagnoseWaitRequest{{}}
}}

func createMessageDiagnoseWaitResponse() MessageHandler {{
    return &MessageDiagnoseWaitResponse{{}}
}}

func createMessageVersionRequest() MessageHandler {{
	return &MessageVersionRequest{{}}
}}

func createMessageVersionResponse() MessageHandler {{
	return &MessageVersionResponse{{}}
}}

func createUndefinedMessage() MessageHandler {{
    return NewUndefinedMessage()
}}
""")

    for signature in signatures:
        print(f"""
        func createMessageBlockchain{signature.name}Request() MessageHandler {{
            return &MessageBlockchain{signature.name}Request{{}}
        }}
        """)

        print(f"""
        func createMessageBlockchain{signature.name}Response() MessageHandler {{
            return &MessageBlockchain{signature.name}Response{{}}
        }}
        """)


if __name__ == "__main__":
    main()
